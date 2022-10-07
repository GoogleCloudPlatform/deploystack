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

	"github.com/GoogleCloudPlatform/deploystack/gcloudtf"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

var skipLsit = list{
	"google_secret_manager_secret_version",
	"google_project_service",
	"google_service_networking_connection",
	"google_cloud_run_service_iam_policy",
	"random_id",
	"google_project_iam_member",
	"null_resource",
	"google_storage_bucket_iam_binding",
}

func main() {
	config, err := gcloudtf.NewGCPResources("../../config/resources.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}

	folderPtr := flag.String("folder", "", "the file path of a local deploystack repo in progress")
	flag.Parse()

	if len(*folderPtr) == 0 {
		log.Fatalf("a folder is required, please pass a repo using -repo flag ")
	}

	DSMeta, err := NewDSMeta(*folderPtr, config)
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

func NewDSMeta(folder string, config gcloudtf.GCPResources) (DSMeta, error) {
	d := DSMeta{}

	mod, dia := tfconfig.LoadModule(folder)
	if dia.Err() != nil {
		return d, fmt.Errorf("terraform config problem %+v", dia)
	}

	d.Entities = NewEntities(mod, config)

	return d, nil
}

type list []string

func (l list) Matches(s string) bool {
	tmp := strings.ToLower(s)
	for _, v := range l {
		if strings.Contains(tmp, strings.ToLower(v)) {
			return true
		}
	}
	return false
}

type entity struct {
	Name    string
	File    string
	Start   int
	Text    string
	Kind    string
	Type    string
	Attr    map[string]string
	Product gcloudtf.TestConfig
}

type entities []entity

func NewEntities(mod *tfconfig.Module, config gcloudtf.GCPResources) entities {
	result := entities{}

	for _, v := range mod.ManagedResources {

		product, ok := config[v.Type]
		if !ok {
			if !skipLsit.Matches(v.Type) {
				fmt.Printf("cannot find proper test for %s\n", v.Type)
			}

			continue
		}

		e := entity{}
		e.Name = v.Name
		e.File = v.Pos.Filename
		e.Start = v.Pos.Line
		e.Type = v.Type
		e.Kind = "resource"
		e.Text, _ = e.GetResourceText()
		e.Product = product.TestConfig
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
	terms := list{"zone", "region"}

	label := e.Product.LabelField
	if e.Product.LabelField == "" {
		label = "name"
	}
	terms = append(terms, label)

	sl := strings.Split(e.Text, "\n")
	m := map[string]string{}

	for _, t := range terms {
		for _, v := range sl {

			if strings.Contains(v, "#") {
				continue
			}

			lsl := strings.Split(v, "=")

			if strings.Contains(lsl[0], t) {
				m[t] = strings.TrimSpace(lsl[1])
				break
			}
		}
	}

	if val, ok := m[label]; ok {
		m["label"] = val
	}

	m["label"] = strings.Trim(m["label"], "\"")

	e.Attr = m
}

func (e entity) IsResource() bool {
	return e.Kind == "resource"
}

func (e entity) IsVariable() bool {
	return e.Kind == "variable"
}

func (e entity) HasLabel() bool {
	_, ok := e.Attr["label"]
	return ok
}

func (e entity) HasTest() bool {
	return e.Product.TestCommand != ""
}

func (e entity) HasTodo() bool {
	return len(e.Product.Todo) > 0
}

func (e entity) NoDefault() bool {
	return !strings.Contains(e.Text, "default")
}

func (e entity) GenerateTestCommand() string {
	if e.Product.TestCommand == "" {
		return ""
	}

	cmdsl := []string{}
	cmdsl = append(cmdsl, e.Product.TestCommand)
	cmdsl = append(cmdsl, e.FormatLabel())

	if e.Product.Zone {
		cmdsl = append(cmdsl, "--zone $ZONE")
	}
	if e.Product.Region {
		cmdsl = append(cmdsl, "--region $REGION")
	}

	if len(e.Product.Suffix) > 0 {
		cmdsl = append(cmdsl, e.Product.Suffix)
	}

	result := strings.Join(cmdsl, " ")
	result = strings.ReplaceAll(result, "gs:// ", "gs://")

	return result
}

func (e entity) Expected() string {
	if e.Product.Expected != "" {
		return e.Product.Expected
	}
	return e.FormatLabel()
}

func (e entity) FormatLabel() string {
	name := e.Attr["label"]

	name = strings.ReplaceAll(name, "\"", "")
	name = strings.ReplaceAll(name, "${var.basename}", "$BASENAME")
	name = strings.ReplaceAll(name, "${var.project_id}", "$PROJECT")

	return name
}

func (e entity) GetResourceText() (string, error) {
	dat, _ := os.ReadFile(e.File)
	sl := strings.Split(string(dat), "\n")
	resultSl := []string{}

	start := e.Start - 1

	end := len(sl) - 1

	end = findClosingBracket(e.Start, sl) + 1

	for i := start - 1; i < end; i++ {
		if i < 0 {
			i = 0
		}
		resultSl = append(resultSl, sl[i])
	}

	result := strings.Join(resultSl, "\n")

	return result, nil
}

func findClosingBracket(start int, sl []string) int {
	count := 0
	for i := start - 1; i < len(sl); i++ {
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
