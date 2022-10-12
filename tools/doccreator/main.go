package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/deploystack/dsgithub"
	"github.com/GoogleCloudPlatform/deploystack/gcloudtf"
)

func main() {
	config, err := gcloudtf.NewGCPResources("../../config/resources.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}

	repoPtr := flag.String("repo", "", "the url of a publicly available github deploystack repo ")
	flag.Parse()

	if len(*repoPtr) == 0 {
		log.Fatalf("a repo is required, please pass a repo using -repo flag ")
	}

	repopath := "."
	Meta, err := dsgithub.NewMeta(*repoPtr, repopath)
	if err != nil {
		log.Fatalf("error loading project %s", err)
	}

	outfile := fmt.Sprintf("%s.md", Meta.ShortName())
	tmpl := "./tmpl/devsite.md"

	outfolder := fmt.Sprintf("./out/%s", filepath.Base(*repoPtr))

	if err := os.MkdirAll(outfolder, os.ModePerm); err != nil {
		log.Fatalf("cannot create out folder: %s", err)
	}

	fmt.Printf("%+v\n", outfolder)

	f, err := os.Create(fmt.Sprintf("%s/%s", outfolder, outfile))
	if err != nil {
		log.Fatalf("cannot create out file: %s", err)
		return
	}

	docMeta := NewDocMeta(Meta, config)

	funcMap := template.FuncMap{
		"ToLower":             strings.ToLower,
		"Title":               strings.Title,
		"Escape":              url.QueryEscape,
		"TrimSpace":           strings.TrimSpace,
		"ShortNameUnderscore": ShortNameUnderscore,
	}

	t, err := template.New(filepath.Base(tmpl)).Funcs(funcMap).ParseFiles(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(f, docMeta)
	if err != nil {
		log.Fatalf("execute: %s", err)
	}
}

func ShortNameUnderscore(name string) string {
	r := strings.ToLower(name)
	r = strings.ReplaceAll(r, " ", "_")
	return r
}

type DocMeta struct {
	DSMeta   dsgithub.Meta
	Products map[string]string
}

func NewDocMeta(dsMeta dsgithub.Meta, config gcloudtf.GCPResources) DocMeta {
	result := DocMeta{}
	result.DSMeta = dsMeta
	result.Products = map[string]string{}

	for _, v := range dsMeta.Terraform {
		val := config.GetProduct(v.Type)

		if val != "" {
			result.Products[val] = val
		}
	}

	return result
}
