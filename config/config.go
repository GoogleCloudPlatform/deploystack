// Package config holds all of the data structures for DeployStack.
// Having them in main package caused circular dependecy issues.
package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/yaml.v2"
)

// Config represents the settings this app will collect from a user. It should
// be in a json file. The idea is minimal programming has to be done to setup
// a DeployStack and export out a tfvars file for terraform part of solution.
type Config struct {
	Title                string            `json:"title" yaml:"title"`
	Name                 string            `json:"name" yaml:"name"`
	Description          string            `json:"description" yaml:"description"`
	Duration             int               `json:"duration" yaml:"duration"`
	Project              bool              `json:"collect_project" yaml:"collect_project"`
	ProjectNumber        bool              `json:"collect_project_number" yaml:"collect_project_number"`
	BillingAccount       bool              `json:"collect_billing_account" yaml:"collect_billing_account"`
	Domain               bool              `json:"register_domain" yaml:"register_domain"`
	Region               bool              `json:"collect_region" yaml:"collect_region"`
	RegionType           string            `json:"region_type" yaml:"region_type"`
	RegionDefault        string            `json:"region_default" yaml:"region_default"`
	Zone                 bool              `json:"collect_zone" yaml:"collect_zone"`
	HardSet              map[string]string `json:"hard_settings" yaml:"hard_settings"`
	CustomSettings       Customs           `json:"custom_settings" yaml:"custom_settings"`
	AuthorSettings       Settings          `json:"author_settings" yaml:"author_settings"`
	ConfigureGCEInstance bool              `json:"configure_gce_instance" yaml:"configure_gce_instance"`
	DocumentationLink    string            `json:"documentation_link" yaml:"documentation_link"`
	PathTerraform        string            `json:"path_terraform" yaml:"path_terraform"`
	PathMessages         string            `json:"path_messages" yaml:"path_messages"`
	PathScripts          string            `json:"path_scripts" yaml:"path_scripts"`
	Projects             Projects          `json:"projects" yaml:"projects"`
	Products             []Product         `json:"products" yaml:"products"`
	WD                   string            `json:"-" yaml:"-"`
}

func (c *Config) convertHardset() {
	for i, v := range c.HardSet {
		c.AuthorSettings.AddComplete(Setting{Name: i, Value: v, Type: "string"})
	}
	// Blow hardset away so that if anywhere is looking for them, it fails.
	c.HardSet = nil
}

// wd used to be unexported, but not it is not. Left the getter and setter to
// not break anything

// Getwd gets the working directory for the config.
func (c *Config) Getwd() string {
	return c.WD
}

// Setwd sets the working directory for the config.
func (c *Config) Setwd(wd string) {
	c.WD = wd
}

// Copy produces a copy of a config file for manipulating it without changing
// the original
func (c Config) Copy() Config {
	out := Config{}
	out.WD = c.WD
	out.Name = c.Name
	out.Title = c.Title
	out.Project = c.Project
	out.ProjectNumber = c.ProjectNumber
	out.Region = c.Region
	out.RegionType = c.RegionType
	out.RegionDefault = c.RegionDefault
	out.Zone = c.Zone
	out.Description = c.Description
	out.Duration = c.Duration
	out.DocumentationLink = c.DocumentationLink
	out.Domain = c.Domain
	out.ConfigureGCEInstance = c.ConfigureGCEInstance
	out.PathTerraform = c.PathTerraform
	out.PathMessages = c.PathMessages
	out.PathScripts = c.PathScripts

	for _, v := range c.AuthorSettings {
		out.AuthorSettings.AddComplete(v)
	}

	for _, v := range c.CustomSettings {
		out.CustomSettings = append(out.CustomSettings, v)
	}

	for _, v := range c.Products {
		out.Products = append(out.Products, v)
	}

	return out
}

// Marshal returns a string representation in format `json` or `yaml`
func (c Config) Marshal(format string) ([]byte, error) {

	if format == "yaml" {
		out, err := yaml.Marshal(&c)
		if err != nil {
			return nil, fmt.Errorf("cannot export test: %s", err)
		}

		return out, nil
	}

	out, err := json.MarshalIndent(&c, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("cannot export test: %s", err)
	}

	return out, nil
}

func (c *Config) defaultAuthorSettings() {
	for i, v := range c.AuthorSettings {
		if v.Type == "" {
			v.Type = "string"
			c.AuthorSettings[i] = v
		}

	}
}

// GetAuthorSettings delivers the combined Hardset and AuthorSettings variables
func (c *Config) GetAuthorSettings() Settings {
	c.convertHardset()
	c.AuthorSettings.Sort()
	return c.AuthorSettings
}

// ComputeName uses the git repo in the working directory to compute the
// shortname for the application.
func (c *Config) ComputeName(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("could not open local git directory: %s", err)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return err
	}

	remote := ""
	for _, v := range remotes {
		for _, url := range v.Config().URLs {
			if strings.Contains(strings.ToLower(url), "googlecloudplatform") {
				remote = strings.ToLower(url)
			}
		}
	}

	// Fixes bug where ssh called repos have issues. Super edge case, but
	// now its all testable
	remote = strings.ReplaceAll(remote, "git@github.com:", "https://github.com/")

	u, err := url.Parse(remote)
	if err != nil {
		return fmt.Errorf("could not parse git url: %s", err)
	}

	shortname := filepath.Base(u.Path)
	shortname = strings.ReplaceAll(shortname, ".git", "")
	shortname = strings.ReplaceAll(shortname, "deploystack-", "")
	c.Name = shortname

	return nil
}

// NewConfigJSON returns a Config object from a file read.
func NewConfigJSON(content []byte) (Config, error) {
	result := Config{}
	if err := json.Unmarshal(content, &result); err != nil {
		return result, fmt.Errorf("unable to convert content to Config: %s", err)
	}

	return result, nil
}

// NewConfigYAML returns a Config object from a file read.
func NewConfigYAML(content []byte) (Config, error) {
	result := Config{}

	if err := yaml.Unmarshal(content, &result); err != nil {
		return result, fmt.Errorf("unable to convert content to Config: %s", err)
	}

	return result, nil
}

// Product is some info about a GCP product
type Product struct {
	Info    string `json:"info" yaml:"info"`
	Product string `json:"product" yaml:"product"`
}

// Project represets a GCP project for use in a stack
type Project struct {
	Name         string `json:"variable_name"  yaml:"variable_name"`
	UserPrompt   string `json:"user_prompt"  yaml:"user_prompt"`
	SetAsDefault bool   `json:"set_as_default"  yaml:"set_as_default"`
	Value        string `json:"value"  yaml:"value"`
}

// Projects is a list of projects that we will collect info for
type Projects struct {
	Items           []Project `json:"items"  yaml:"items"`
	AllowDuplicates bool      `json:"allow_duplicates"  yaml:"allow_duplicates"`
}

// Setting is a item that will be translated to a varaible in a terraform file
type Setting struct {
	Name  string            `json:"name"  yaml:"name"`
	Value string            `json:"value"  yaml:"value"`
	Type  string            `json:"type"  yaml:"type"`
	List  []string          `json:"list"  yaml:"list"`
	Map   map[string]string `json:"map"  yaml:"map"`
}

// TFVars emits the name value combination here in away that terraform excepts
// in a tfvars file
func (s *Setting) TFVars() string {
	return fmt.Sprintf("%s=%s\n", s.TFvarsName(), s.TFvarsValue())
}

// TFvarsName formats the name for the tfvars format
func (s Setting) TFvarsName() string {
	name := strings.ToLower(strings.ReplaceAll(s.Name, " ", "_"))
	return name
}

// TFvarsValue formats the value for the tfvars format
func (s Setting) TFvarsValue() string {
	result := ""
	// If we used the workaround for lists in strings, convert it to a list
	// under the covers
	if s.Value != "" && s.Value[0:1] == "[" {
		replacer := strings.NewReplacer("[", "", "]", "")
		s.List = strings.Split(replacer.Replace(s.Value), ",")
		s.Type = "list"
		s.Value = ""
	}

	switch s.Type {
	case "string", "":
		result = fmt.Sprintf("\"%s\"", s.Value)
	case "list":
		tmp := []string{}
		for _, v := range s.List {
			tmp = append(tmp, fmt.Sprintf("\"%s\"", v))
		}
		str := strings.Join(tmp, ",")

		result = fmt.Sprintf("[%s]", str)
	case "map":
		tmp := []string{}

		for i, v := range s.Map {
			tmp = append(tmp, fmt.Sprintf("%s=\"%s\"", i, v))
		}

		sort.Strings(tmp)
		str := strings.Join(tmp, ",")
		result = fmt.Sprintf("{%s}", str)
	default:
		result = s.Value
	}

	return result
}

// Settings are a collection of setting
type Settings []Setting

// AddComplete adds an whole setting to the settings control
func (s *Settings) AddComplete(set Setting) {

	setting := s.Find(set.Name)
	if setting != nil {
		s.Replace(set)
		return
	}

	(*s) = append((*s), set)
	return
}

// Add either creates a new setting or updates the existing one
func (s *Settings) Add(key, value string) {
	k := strings.ToLower(key)

	set := s.Find(key)
	if set != nil {
		set.Name = key
		set.Value = value
		set.Type = "string"
		s.Replace(*set)
		return
	}

	set = &Setting{Name: k, Value: value, Type: "string"}
	(*s) = append((*s), *set)
	return
}

// Sort sorts the slice according to Setting.Name ascendings
func (s *Settings) Sort() {
	sort.Slice(*s, func(i, j int) bool {
		return (*s)[i].Name < (*s)[j].Name
	})
}

// Replace will look for a setting with the same name, and overwrite the value
func (s *Settings) Replace(set Setting) {
	for i, v := range *s {
		if v.Name == set.Name {
			(*s)[i] = set
		}
	}

}

// Search returns all settings whose names contain a particular string
func (s *Settings) Search(q string) Settings {
	result := Settings{}

	for _, v := range *s {
		if strings.Contains(v.Name, q) {
			result = append(result, v)
		}
	}

	return result
}

// Find locates a setting in the slice
func (s *Settings) Find(key string) *Setting {
	k := strings.ToLower(key)

	for _, v := range *s {
		if v.Name == k {
			return &v
		}
	}

	return nil
}

// Custom represents a custom setting that we would like to collect from a user
// We will collect these settings from the user before continuing.
type Custom struct {
	Setting        `json:"-"  yaml:"-"`
	Name           string   `json:"name"  yaml:"name"`
	Description    string   `json:"description"  yaml:"description"`
	Default        string   `json:"default"  yaml:"default"`
	Options        []string `json:"options"  yaml:"options"`
	PrependProject bool     `json:"prepend_project"  yaml:"prepend_project"`
	Validation     string   `json:"validation,omitempty"  yaml:"validation,omitempty"`
	Project        string   `json:"-"  yaml:"-"`
}

// Customs are a slice of Custom variables.
type Customs []Custom

// Get returns one Custom Variable
func (cs Customs) Get(name string) Custom {
	for _, v := range cs {
		if v.Name == name {
			return v
		}
	}

	return Custom{}
}

// Report is collection of data about multiple configs in the same root
// used for multi stack repos
type Report struct {
	Path   string
	WD     string
	Config Config
}

// NewReport Generates a new config report for a given file
func NewReport(file string) (Report, error) {
	result := Report{Path: file}

	result.WD = strings.ReplaceAll(filepath.Dir(file), "/.deploystack", "")

	dat, err := os.ReadFile(file)
	if err != nil {
		return result, err
	}

	switch filepath.Ext(file) {
	case ".json":
		result.Config, err = NewConfigJSON(dat)
		if err != nil {
			return result, err
		}
	case ".yaml":
		result.Config, err = NewConfigYAML(dat)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

// FindConfigReports walks through a directory and finds all of the configs in
// the folder
func FindConfigReports(dir string) ([]Report, error) {

	var result []Report
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if info.Name() == "deploystack.json" || info.Name() == "deploystack.yaml" {
			cr, err := NewReport(path)
			if err != nil {
				return err
			}

			result = append(result, cr)
		}
		return nil
	})
	return result, err

}
