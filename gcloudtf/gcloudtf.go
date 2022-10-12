// Package gcloudtf handles parsing Terraform files to extract structured
// information out of it, to support a few tools to speed up the DeployStack
// authoring process
package gcloudtf

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"gopkg.in/yaml.v2"
)

// Extract points to a path that includes Terraform files and extracts all of
// the information out of it for use with DeployStack Tools
func Extract(path string) (*Blocks, error) {
	mod, dia := tfconfig.LoadModule(path)
	if dia.Err() != nil {
		return nil, fmt.Errorf("terraform config problem %+v", dia)
	}

	b, err := NewBlocks(mod)
	if dia.Err() != nil {
		return nil, fmt.Errorf("could not properly parse blocks %+v", err)
	}

	return b, nil
}

// Block represents one of several kinds of Terraform constructs: resources,
// variables, module
type Block struct {
	Name  string            `json:"name" yaml:"name"`
	Text  string            `json:"text" yaml:"text"`
	Kind  string            `json:"kind" yaml:"kind"`
	Type  string            `json:"type" yaml:"type"`
	Attr  map[string]string `json:"attr" yaml:"attr"`
	file  string
	start int
}

// NewResourceBlock converts a parsed Terraform Resource to a Block
func NewResourceBlock(t *tfconfig.Resource) (Block, error) {
	b := Block{}
	var err error
	b.Name = t.Name
	b.Type = t.Type
	b.Kind = t.Mode.String()
	b.start = t.Pos.Line
	b.file = t.Pos.Filename
	b.Text, err = getResourceText(t.Pos.Filename, t.Pos.Line)
	if err != nil {
		return b, fmt.Errorf("could not extract text from Resource: %s", err)
	}

	return b, nil
}

// NewVariableBlock converts a parsed Terraform Variable to a Block
func NewVariableBlock(t *tfconfig.Variable) (Block, error) {
	b := Block{}
	var err error
	b.Name = t.Name
	b.Type = t.Type
	b.Kind = "variable"
	b.start = t.Pos.Line
	b.file = t.Pos.Filename
	b.Text, err = getResourceText(t.Pos.Filename, t.Pos.Line)
	if err != nil {
		return b, fmt.Errorf("could not extract text from Variable: %s", err)
	}
	return b, nil
}

// NewModuleBlock converts a parsed Terraform Module to a Block
func NewModuleBlock(t *tfconfig.ModuleCall) (Block, error) {
	b := Block{}
	var err error
	b.Name = t.Name
	b.Type = t.Source
	b.Kind = "module"
	b.start = t.Pos.Line
	b.file = t.Pos.Filename
	b.Text, err = getResourceText(t.Pos.Filename, t.Pos.Line)
	if err != nil {
		return b, fmt.Errorf("could not extract text from Module: %s", err)
	}
	return b, nil
}

// IsResource returns true if block is a Terraform resource
func (b Block) IsResource() bool {
	return b.Kind == "managed"
}

// IsModule returns true if block is a Terraform module
func (b Block) IsModule() bool {
	return b.Kind == "module"
}

// IsVariable returns true if block is a Terraform variable
func (b Block) IsVariable() bool {
	return b.Kind == "variable"
}

// NoDefault returns true if block does not contain a default value
func (b Block) NoDefault() bool {
	return !strings.Contains(b.Text, "default")
}

func getResourceText(file string, start int) (string, error) {
	dat, _ := os.ReadFile(file)
	sl := strings.Split(string(dat), "\n")

	resultSl := []string{}

	startpos := start - 1

	end := len(sl) - 1

	end = findClosingBracket(start, sl) + 1

	if startpos == 0 {
		startpos = 1
	}

	for i := startpos - 1; i < end; i++ {
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

// Blocks is a slice of type Block
type Blocks []Block

// NewBlocks converts the results from a Terraform parse operation to Blocks.
func NewBlocks(mod *tfconfig.Module) (*Blocks, error) {
	result := Blocks{}

	for _, v := range mod.ModuleCalls {
		b, err := NewModuleBlock(v)
		if err != nil {
			return nil, fmt.Errorf("could not parse Module Calls: %s", err)
		}
		result = append(result, b)
	}

	for _, v := range mod.ManagedResources {
		b, err := NewResourceBlock(v)
		if err != nil {
			return nil, fmt.Errorf("could not parse ManagedResources: %s", err)
		}
		result = append(result, b)
	}

	for _, v := range mod.DataResources {
		b, err := NewResourceBlock(v)
		if err != nil {
			return nil, fmt.Errorf("could not parse DataResources: %s", err)
		}
		result = append(result, b)
	}

	for _, v := range mod.Variables {
		b, err := NewVariableBlock(v)
		if err != nil {
			return nil, fmt.Errorf("could not parse Variables: %s", err)
		}
		result = append(result, b)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].start < result[j].start
	})

	return &result, nil
}

// List is a slice of strings that we add extra functionality to
type List []string

// Matches determines if a given string is an entry in the list.
func (l List) Matches(s string) bool {
	tmp := strings.ToLower(s)
	for _, v := range l {
		if strings.Contains(tmp, strings.ToLower(v)) {
			return true
		}
	}
	return false
}

// GCPResources is a collection of GCPResource
type GCPResources map[string]GCPResource

// GetProduct returns the prouct name assoicated with the terraform resource
func (g GCPResources) GetProduct(key string) string {
	v, ok := g[key]
	if !ok {
		return ""
	}

	return v.Product
}

// GCPResource is a Terraform resource that matches up with a GCP product. This
// is used to automate the generation of tests and documentation
type GCPResource struct {
	Label      string     `json:"label" yaml:"label"`
	Product    string     `json:"product" yaml:"product"`
	APICalls   []string   `json:"api_calls" yaml:"api_calls"`
	TestConfig TestConfig `json:"test_config" yaml:"test_config"`
}

// TestConfig is the information needed to automate test creation.
type TestConfig struct {
	TestType    string `json:"test_type" yaml:"test_type"`
	TestCommand string `json:"test_command" yaml:"test_command"`
	Suffix      string `json:"suffix" yaml:"suffix"`
	Region      bool   `json:"region" yaml:"region"`
	Zone        bool   `json:"zone" yaml:"zone"`
	LabelField  string `json:"label_field" yaml:"label_field"`
	Expected    string `json:"expected" yaml:"expected"`
	Todo        string `json:"todo" yaml:"todo"`
}

// HasTest returns true if test config actually has a test in it.
func (t TestConfig) HasTest() bool {
	return len(t.TestCommand) > 0
}

// HasTodo returns true if test config actually has a todo in it.
func (t TestConfig) HasTodo() bool {
	return len(t.Todo) > 0
}

// NewGCPResources reads in a yaml file as a config
func NewGCPResources(path string) (GCPResources, error) {
	result := GCPResources{}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return result, fmt.Errorf("unable to find or read config file: %s", err)
	}

	if err := yaml.Unmarshal(content, &result); err != nil {
		return result, fmt.Errorf("unable to convert content to GCPResources: %s", err)
	}

	return result, nil
}

// Repos is a slice of strings containing github urls
type Repos []string

// NewRepos retrieves the list of repos from a ymal file.
func NewRepos(path string) (Repos, error) {
	result := Repos{}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return result, fmt.Errorf("unable to find or read config file: %s", err)
	}

	if err := yaml.Unmarshal(content, &result); err != nil {
		return result, fmt.Errorf("unable to convert content to a list of repos: %s", err)
	}

	return result, nil
}
