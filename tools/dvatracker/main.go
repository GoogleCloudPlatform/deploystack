package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/deploystack/dsgithub"
	"github.com/GoogleCloudPlatform/deploystack/gcloudtf"
)

var versionStrings = []string{"v1", "alpha", "beta", "v1beta1", "v1beta2", "v1beta3", "v1beta4"}

var repos = []string{}

func main() {
	dvasGlobal := DVAs{}

	config, err := gcloudtf.NewGCPResources("../../config/resources.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}

	repos, err := gcloudtf.NewRepos("../../config/repos.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}

	repopath := "."
	outfolder := "out"

	for _, repo := range repos {
		Meta, err := dsgithub.NewMeta(repo, repopath)
		if err != nil {
			log.Fatalf("error loading project %s", err)
		}

		dvas := GetDVAs(Meta.ShortName(), Meta.Terraform, config)

		dvasGlobal = append(dvasGlobal, dvas...)
	}

	if err := os.MkdirAll(outfolder, os.ModePerm); err != nil {
		log.Fatalf("cannot create out folder: %s", err)
	}

	if err := dvasGlobal.ToCSV(fmt.Sprintf("%s/ref.csv", outfolder)); err != nil {
		log.Fatalf("error writing csv: %s", err)
	}
}

func GetDVAs(shortname string, b gcloudtf.Blocks, config gcloudtf.GCPResources) DVAs {
	result := DVAs{}
	for _, block := range b {
		if block.Kind == "managed" || block.Kind == "module" {
			cfg, ok := config[block.Type]
			if !ok {
				fmt.Printf("Needs an entry: %s\n", block.Type)
				continue
			}
			if len(cfg.APICalls) == 0 {
				continue
			}
			for _, v := range cfg.APICalls {
				for _, version := range versionStrings {
					d := DVA{shortname, strings.Replace(v, "[version]", version, 1)}
					result.Add(d)
				}
			}

		}
	}

	return result
}

type DVAs []DVA

func (ds *DVAs) Add(d DVA) {
	if !ds.Exists(d) {
		*ds = append(*ds, d)
	}
}

func (ds DVAs) Exists(d DVA) bool {
	for _, v := range ds {
		if d.API == v.API && d.Stack == v.Stack {
			return true
		}
	}

	return false
}

func (d DVAs) ToCSV(file string) error {
	f, err := os.Create(file)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("could not create file %s: %s", file, err)
	}
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"appname", "api_method"}); err != nil {
		log.Fatalln("error writing header to file", err)
	}

	for _, record := range d {
		if err := w.Write(record.toSlice()); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

	return nil
}

type DVA struct {
	Stack string
	API   string
}

func (d DVA) toSlice() []string {
	return []string{d.Stack, d.API}
}
