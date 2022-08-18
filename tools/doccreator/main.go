package main

import (
	"html"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

func main() {
	mod, dia := tfconfig.LoadModule("./repo")
	if dia.Err() != nil {
		log.Fatalf("error in the diagnostics %+v\n", dia)
	}

	s := deploystack.NewStack()
	if err := s.ReadConfig("./repo/deploystack.json", "./repo/messages/description.txt"); err != nil {
		log.Fatalf("could not read config file: %s", err)
	}

	f, err := os.Create("./out/results.html")
	if err != nil {
		log.Fatalf("create file: %s", err)
		return
	}

	DSMeta := DSMeta{}
	DSMeta.Entities = NewEntities(mod)
	DSMeta.DeployStack = s.Config
	DSMeta.GitRepo = "https://github.com/GoogleCloudPlatform/deploystack-nosql-client-server"

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
		"Title":   strings.Title,
	}

	t, err := template.New("dsdev.html").Funcs(funcMap).ParseFiles("./tmpl/dsdev.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(f, DSMeta)
	if err != nil {
		log.Fatalf("execute: %s", err)
	}
}

type DSMeta struct {
	DeployStack deploystack.Config
	Entities    entities
	GitRepo     string
}

type entity struct {
	Name  string
	File  string
	Start int
	Text  string
	Kind  string
	Type  string
}

type entities []entity

func NewEntities(mod *tfconfig.Module) entities {
	result := entities{}

	for _, v := range mod.ManagedResources {
		e := entity{}
		e.Name = v.Name
		e.File = v.Pos.Filename
		e.Start = v.Pos.Line
		e.Type = v.Type
		e.Kind = "resource"
		e.Text, _ = e.GetResourceText()
		result = append(result, e)
	}

	for _, v := range mod.DataResources {
		e := entity{}
		e.Name = v.Name
		e.File = v.Pos.Filename
		e.Start = v.Pos.Line
		e.Type = v.Type
		e.Kind = "data"
		e.Text, _ = e.GetResourceText()
		result = append(result, e)
	}

	for _, v := range mod.Variables {
		e := entity{}
		e.Name = v.Name
		e.File = v.Pos.Filename
		e.Start = v.Pos.Line
		e.Type = v.Type
		e.Kind = "data"
		e.Text, _ = e.GetResourceText()
		result = append(result, e)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Start < result[j].Start
	})

	return result
}

func (e entity) GetResourceText() (string, error) {
	dat, _ := os.ReadFile(e.File)
	sl := strings.Split(string(dat), "\n")

	resultSl := []string{}

	start := e.Start - 1

	end := len(sl) - 1

	end = findClosingBracket(e, sl) + 1

	for i := start - 1; i < end; i++ {
		resultSl = append(resultSl, sl[i])
	}

	result := strings.Join(resultSl, "\n")

	result = html.EscapeString(result)

	return result, nil
}

func findClosingBracket(e entity, sl []string) int {
	count := 0
	for i := e.Start - 1; i < len(sl); i++ {
		if strings.Contains(sl[i], "{") {
			count++
		}
		if strings.Contains(sl[i], "}") {
			count--
		}
		if count == 0 {
			return i
		}
	}

	return len(sl)
}
