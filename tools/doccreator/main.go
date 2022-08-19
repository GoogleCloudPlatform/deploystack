package main

import (
	"fmt"
	"html"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"gopkg.in/src-d/go-git.v4"
)

func main() {
	DSMeta, err := NewDSMeta("https://github.com/GoogleCloudPlatform/deploystack-nosql-client-server")
	if err != nil {
		log.Fatalf("error in the diagnostics %s", err)
	}

	f, err := os.Create("./out/results.html")
	if err != nil {
		log.Fatalf("create file: %s", err)
		return
	}

	funcMap := template.FuncMap{
		"ToLower":   strings.ToLower,
		"Title":     strings.Title,
		"Escape":    url.QueryEscape,
		"TrimSpace": strings.TrimSpace,
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

type product struct {
	Title         string
	Youtube       string
	Walkthrough   string
	Documentation string
	Description   string
	Logo          string
}

type DSMeta struct {
	DeployStack deploystack.Config
	Entities    entities
	GitRepo     string
	Products    map[string]product
}

func NewDSMeta(repo string) (DSMeta, error) {
	var err error
	d := DSMeta{}
	d.GitRepo = repo
	d.Products = make(map[string]product)

	if _, err := os.Stat("./repo"); os.IsNotExist(err) {
		_, err = git.PlainClone("./repo", false, &git.CloneOptions{
			URL:      repo,
			Progress: os.Stdout,
		})
	}

	if err != nil {
		return d, err
	}

	mod, dia := tfconfig.LoadModule("./repo")
	if dia.Err() != nil {
		return d, fmt.Errorf("error in the diagnostics %+v", dia)
	}

	s := deploystack.NewStack()
	if err := s.ReadConfig("./repo/deploystack.json", "./repo/messages/description.txt"); err != nil {
		return d, fmt.Errorf("could not read config file: %s", err)
	}

	d.Entities = NewEntities(mod)
	d.DeployStack = s.Config

	for _, v := range d.Entities {
		if v.Kind == "resource" {
			if val, ok := prods[v.Type]; ok {
				if p, pok := pinfo[val]; pok {
					d.Products[val] = p
				}
			}
		}
	}

	return d, nil
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
		e.Kind = "variable"
		e.Text, _ = e.GetResourceText()
		result = append(result, e)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Start < result[j].Start
	})

	return result
}

func (e entity) IsResource() bool {
	return e.Kind == "resource"
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
