package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/deploystack/dsgithub"
	"github.com/GoogleCloudPlatform/deploystack/gcloudtf"
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
		log.Fatalf("a folder is required, please pass a folder using -folder flag ")
	}

	meta, err := dsgithub.NewMetaFromLocal(*folderPtr)
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

	entities := NewEntities(meta.Terraform, config)

	err = t.Execute(f, entities)
	if err != nil {
		log.Fatalf("execute: %s", err)
	}
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
	Block      gcloudtf.Block
	TestConfig gcloudtf.TestConfig
	Attr       map[string]string
}

type entities []entity

func NewEntities(blocks gcloudtf.Blocks, config gcloudtf.GCPResources) entities {
	result := entities{}

	for _, block := range blocks {

		e := entity{}
		e.Block = block

		product, ok := config[block.Type]
		if !ok {
			if !skipLsit.Matches(block.Type) && block.Kind != "variable" && block.Kind != "data" {
				fmt.Printf("cannot find proper test for %s kind(%s) \n", block.Type, block.Kind)
			}
		}
		if ok {
			e.TestConfig = product.TestConfig
			if e.TestConfig.HasTest() {
				e.generateMap()
			}

		}

		result = append(result, e)

	}

	return result
}

func (e entity) GenerateTestCommand() string {
	if e.TestConfig.TestCommand == "" {
		return ""
	}

	cmdsl := []string{}
	cmdsl = append(cmdsl, e.TestConfig.TestCommand)
	cmdsl = append(cmdsl, e.FormatLabel())

	if e.TestConfig.Zone {
		cmdsl = append(cmdsl, "--zone $ZONE")
	}
	if e.TestConfig.Region {
		cmdsl = append(cmdsl, "--region $REGION")
	}

	if len(e.TestConfig.Suffix) > 0 {
		cmdsl = append(cmdsl, e.TestConfig.Suffix)
	}

	result := strings.Join(cmdsl, " ")
	result = strings.ReplaceAll(result, "gs:// ", "gs://")

	return result
}

func (e entity) FormatLabel() string {
	name := e.Attr["label"]
	name = strings.ReplaceAll(name, "\"", "")
	name = strings.ReplaceAll(name, "${var.basename}", "$BASENAME")
	name = strings.ReplaceAll(name, "${var.project_id}", "$PROJECT")

	return name
}

func (e *entity) generateMap() {
	terms := list{"zone", "region"}

	label := e.TestConfig.LabelField
	if e.TestConfig.LabelField == "" {
		label = "name"
	}
	terms = append(terms, label)

	sl := strings.Split(e.Block.Text, "\n")
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
