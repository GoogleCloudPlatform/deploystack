package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/gcloudtf"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"gopkg.in/src-d/go-git.v4"
)

func main() {
	config, err := gcloudtf.NewGCPResources("../../config/resources.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}

	repoPtr := flag.String("repo", "", "the url of a publicly available github deploystack repo ")
	devsitePtr := flag.Bool("devsite", true, "whether or not not make this for the internal doc site")
	flag.Parse()

	if len(*repoPtr) == 0 {
		log.Fatalf("a repo is required, please pass a repo using -repo flag ")
	}

	outfile := "index.html"
	tmpl := "./tmpl/dsdev.html"

	DSMeta, err := NewDSMeta(*repoPtr, config)
	if err != nil {
		log.Fatalf("error loading project %s", err)
	}

	if *devsitePtr {
		outfile = fmt.Sprintf("%s.md", DSMeta.GetShortName())
		tmpl = "./tmpl/devsite.md"
	}

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

	funcMap := template.FuncMap{
		"ToLower":   strings.ToLower,
		"Title":     strings.Title,
		"Escape":    url.QueryEscape,
		"TrimSpace": strings.TrimSpace,
	}

	t, err := template.New(filepath.Base(tmpl)).Funcs(funcMap).ParseFiles(tmpl)
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

func (p product) GetShortNameUnderScore() string {
	r := strings.ToLower(p.Title)

	r = strings.ReplaceAll(r, " ", "_")

	return r
}

type DSMeta struct {
	DeployStack deploystack.Config
	Entities    entities
	GitRepo     string
	Products    map[string]product
}

func NewDSMeta(repo string, config gcloudtf.GCPResources) (DSMeta, error) {
	d := DSMeta{}
	d.GitRepo = repo
	d.Products = make(map[string]product)

	if _, err := os.Stat("./repo"); os.IsNotExist(err) {
		fname := filepath.Join(os.TempDir(), "stdout")
		old := os.Stdout            // keep backup of the real stdout
		temp, _ := os.Create(fname) // create temp file
		defer temp.Close()
		os.Stdout = temp
		_, err = git.PlainClone("./repo", false, &git.CloneOptions{
			URL:      repo,
			Progress: temp,
		})
		if err != nil {
			os.Stdout = old
			out, _ := ioutil.ReadFile(fname)
			fmt.Printf("git response: \n%s\n", string(out))
			return d, fmt.Errorf("cannot get repo: %s", err)
		}

		os.Stdout = old
	}

	// for _, v := range required {
	// 	if _, err := os.Stat(v); os.IsNotExist(err) {
	// 		return d, fmt.Errorf("required file does not exist: %s", v)
	// 	}
	// }

	orgpwd, err := os.Getwd()
	if err != nil {
		return d, fmt.Errorf("could not get the wd: %s", err)
	}

	path := fmt.Sprintf("%s/repo", orgpwd)

	mod, dia := tfconfig.LoadModule(path)
	if dia.Err() != nil {
		return d, fmt.Errorf("terraform config problem %+v", dia)
	}

	if err := os.Chdir(path); err != nil {
		return d, fmt.Errorf("could not change the wd: %s", err)
	}

	s := deploystack.NewStack()
	if err := s.FindAndReadRequired(); err != nil {
		return d, fmt.Errorf("could not read config file (%s): %s", path, err)
	}

	if err := os.Chdir(orgpwd); err != nil {
		return d, fmt.Errorf("could not change the wd back: %s", err)
	}

	d.Entities = NewEntities(mod)
	d.DeployStack = s.Config

	for _, v := range d.Entities {
		if v.Kind == "resource" {
			val := config.GetProduct(v.Type)

			if val != "" {
				p, ok := pinfo[val]
				if ok {
					d.Products[val] = p
				}
			}
		}
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

	// result = html.EscapeString(result)

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
