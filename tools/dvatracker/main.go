package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/gcloudtf"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
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

	for _, repo := range repos {
		DSMeta, err := NewDSMeta(repo)
		if err != nil {
			log.Fatalf("error loading project %s", err)
		}

		dvas := GetDVAs(DSMeta.GetShortName(), DSMeta.Blocks, config)

		dvasGlobal = append(dvasGlobal, dvas...)
	}
	if err := dvasGlobal.ToCSV("out/ref.csv"); err != nil {
		log.Fatalf("error writing csv: %s", err)
	}
}

func GetDVAs(shortname string, b gcloudtf.Blocks, config gcloudtf.GCPResources) []DVA {
	result := []DVA{}
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
					result = append(result, d)
				}
			}

		}
	}

	return result
}

type DVAs []DVA

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

type DSMeta struct {
	DeployStack deploystack.Config
	Blocks      gcloudtf.Blocks
	GitRepo     string
	GitBranch   string
}

func NewDSMeta(repo string) (DSMeta, error) {
	d := DSMeta{}
	d.GitRepo = repo
	d.GitBranch = "main"

	if strings.Contains(repo, "/tree/") {
		end := strings.Index(repo, "/tree/")
		d.GitRepo = repo[:end]
		d.GitBranch = repo[end+6:]
	}

	repoPath := filepath.Base(d.GitRepo)
	repoPath = strings.ReplaceAll(repoPath, "deploystack-", "")
	repoPath = fmt.Sprintf("./repo/%s", repoPath)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fname := filepath.Join(os.TempDir(), "stdout")
		old := os.Stdout            // keep backup of the real stdout
		temp, _ := os.Create(fname) // create temp file
		defer temp.Close()
		os.Stdout = temp
		_, err = git.PlainClone(
			repoPath,
			false,
			&git.CloneOptions{
				URL:           d.GitRepo,
				ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", d.GitBranch)),
				Progress:      temp,
			})

		if err != nil {
			os.Stdout = old
			out, _ := ioutil.ReadFile(fname)
			fmt.Printf("git response: \n%s\n", string(out))
			return d, fmt.Errorf("cannot get repo: %s", err)
		}

		os.Stdout = old
	}

	orgpwd, err := os.Getwd()
	if err != nil {
		return d, fmt.Errorf("could not get the wd: %s", err)
	}
	if err := os.Chdir(repoPath); err != nil {
		return d, fmt.Errorf("could not change the wd: %s", err)
	}

	s := deploystack.NewStack()

	if err := s.FindAndReadRequired(); err != nil {
		return d, fmt.Errorf("could not read config file: %s", err)
	}

	b, err := gcloudtf.Extract(s.Config.PathTerraform)
	if err != nil {
		log.Fatalf("couldn't extract from TF file: %s", err)
	}

	d.Blocks = *b
	d.DeployStack = s.Config

	if err := os.Chdir(orgpwd); err != nil {
		return d, fmt.Errorf("could not change the wd back: %s", err)
	}

	return d, nil
}

func (d DSMeta) GetShortName() string {
	r := filepath.Base(d.GitRepo)
	r = strings.ReplaceAll(r, "deploystack-", "")
	return r
}

func (d DSMeta) GetShortNameUnderScore() string {
	r := d.GetShortName()
	r = strings.ReplaceAll(r, "-", "_")
	return r
}
