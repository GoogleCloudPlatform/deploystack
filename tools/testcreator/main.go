package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

var required = []string{
	"./repo/deploystack.json",
	"./repo/messages/description.txt",
	"./repo/main.tf",
}

func main() {
	folderPtr := flag.String("folder", "", "the file path of a local deploystack repo in progress")
	flag.Parse()

	if len(*folderPtr) == 0 {
		log.Fatalf("a folder is required, please pass a repo using -repo flag ")
	}

	DSMeta, err := NewDSMeta(*folderPtr)
	if err != nil {
		log.Fatalf("error loading project %s", err)
	}

	outfolder := fmt.Sprintf("./out/")

	if err := os.MkdirAll(outfolder, os.ModePerm); err != nil {
		log.Fatalf("cannot create out folder: %s", err)
	}

	f, err := os.Create(fmt.Sprintf("%s/test", outfolder))
	if err != nil {
		log.Fatalf("cannot create out file: %s", err)
		return
	}

	funcMap := template.FuncMap{
		"ToLower":   strings.ToLower,
		"ToUpper":   strings.ToUpper,
		"Title":     strings.Title,
		"Escape":    url.QueryEscape,
		"TrimSpace": strings.TrimSpace,
	}

	t, err := template.New("test").Funcs(funcMap).ParseFiles("./tmpl/test")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(f, DSMeta)
	if err != nil {
		log.Fatalf("execute: %s", err)
	}
}

type DSMeta struct {
	Entities entities
}

func NewDSMeta(folder string) (DSMeta, error) {
	d := DSMeta{}

	mod, dia := tfconfig.LoadModule(folder)
	if dia.Err() != nil {
		return d, fmt.Errorf("terraform config problem %+v", dia)
	}

	d.Entities = NewEntities(mod)

	return d, nil
}

type entity struct {
	Name   string
	File   string
	Start  int
	Text   string
	Kind   string
	Type   string
	Attr   map[string]string
	Gcloud string
}

type entities []entity

func NewEntities(mod *tfconfig.Module) entities {
	result := entities{}

	for _, v := range mod.ManagedResources {

		gcloud, ok := prods[v.Type]
		if !ok {
			continue
		}

		e := entity{}
		e.Name = v.Name
		e.File = v.Pos.Filename
		e.Start = v.Pos.Line
		e.Type = v.Type
		e.Kind = "resource"
		e.Text, _ = e.GetResourceText()
		e.Gcloud = gcloud
		e.generateMap()
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

func (e *entity) generateMap() {
	terms := []string{"name", "zone", "region", "secret_id"}
	sl := strings.Split(e.Text, "\n")
	m := map[string]string{}

	for _, t := range terms {
		for _, v := range sl {

			lsl := strings.Split(v, "=")

			if strings.Contains(lsl[0], t) {
				m[t] = strings.TrimSpace(lsl[1])
				break
			}
		}
		if _, ok := m[t]; ok {
			break
		}
	}

	e.Attr = m
}

func (e entity) IsResource() bool {
	return e.Kind == "resource"
}

func (e entity) IsVariable() bool {
	return e.Kind == "variable"
}

func (e entity) HasName() bool {
	_, ok := e.Attr["name"]
	return ok
}

func (e entity) NoDefault() bool {
	return !strings.Contains(e.Text, "default")
}

func (e entity) GenerateGcloudCommand() string {
	zone := e.Attr["zone"]
	region := e.Attr["region"]

	cmdsl := []string{}
	cmdsl = append(cmdsl, "gcloud")
	cmdsl = append(cmdsl, e.Gcloud)
	cmdsl = append(cmdsl, "describe")
	cmdsl = append(cmdsl, e.GenerateGcloudName())

	if len(zone) > 0 {
		cmdsl = append(cmdsl, "--zone $ZONE")
	}
	if len(zone) == 0 && len(region) > 0 {
		cmdsl = append(cmdsl, "--region $REGION")
	}

	if strings.Contains(e.Type, "google_cloud_run") {
		cmdsl = append(cmdsl, "--region $REGION")
	}

	cmdsl = append(cmdsl, `--format="value(name)"`)

	return strings.Join(cmdsl, " ")
}

func (e entity) GenerateGcloudName() string {
	name := e.Attr["name"]
	name = strings.ReplaceAll(name, "\"", "")
	name = strings.ReplaceAll(name, "${var.basename}", "$BASENAME")

	if _, ok := e.Attr["secret_id"]; ok {
		name = e.Attr["secret_id"]
	}

	return name
}

func (e entity) GetResourceText() (string, error) {
	dat, _ := os.ReadFile(e.File)
	sl := strings.Split(string(dat), "\n")
	resultSl := []string{}

	start := e.Start - 1

	end := len(sl) - 1

	end = findClosingBracket(e, sl) + 1

	for i := start - 1; i < end; i++ {
		if i < 0 {
			i = 0
		}
		resultSl = append(resultSl, sl[i])
	}

	result := strings.Join(resultSl, "\n")

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
