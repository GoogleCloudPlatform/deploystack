// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package deploystack provides a series of interfaces for getting Google Cloud
// settings and configurations for use with DeplyStack
package deploystack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	"google.golang.org/api/option"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/yaml.v2"
)

var (
	opts             = option.WithCredentialsFile("")
	credspath        = ""
	defaultUserAgent = "deploystack"
	contactfile      = "contact.yaml"
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
	CustomSettings       []Custom          `json:"custom_settings" yaml:"custom_settings"`
	ConfigureGCEInstance bool              `json:"configure_gce_instance" yaml:"configure_gce_instance"`
	DocumentationLink    string            `json:"documentation_link" yaml:"documentation_link"`
	PathTerraform        string            `json:"path_terraform" yaml:"path_terraform"`
	PathMessages         string            `json:"path_messages" yaml:"path_messages"`
	PathScripts          string            `json:"path_scripts" yaml:"path_scripts"`
	Projects             Projects          `json:"projects" yaml:"projects"`
}

// Project represets a GCP project for use in a stack
type Project struct {
	Name         string `json:"variable_name"  yaml:"variable_name"`
	UserPrompt   string `json:"user_prompt"  yaml:"user_prompt"`
	SetAsDefault bool   `json:"set_as_default"  yaml:"set_as_default"`
	value        string
}

// Projects is a list of projects that we will collect info for
type Projects struct {
	Items           []Project `json:"items"  yaml:"items"`
	AllowDuplicates bool      `json:"allow_duplicates"  yaml:"allow_duplicates"`
}

// Custom represents a custom setting that we would like to collect from a user
// We will collect these settings from the user before continuing.
type Custom struct {
	Name           string   `json:"name"  yaml:"name"`
	Description    string   `json:"description"  yaml:"description"`
	Default        string   `json:"default"  yaml:"default"`
	Value          string   `json:"value"  yaml:"value"`
	Options        []string `json:"options"  yaml:"options"`
	PrependProject bool     `json:"prepend_project"  yaml:"prepend_project"`
	Validation     string   `json:"validation,omitempty"  yaml:"validation,omitempty"`
	project        string
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

// ComputeName uses the git repo in the working directory to compute the
// shortname for the application.
func (c *Config) ComputeName() error {
	repo, err := git.PlainOpen(".")
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

// func handleProcessError(err error) {
// 	fmt.Printf("\n\n%sThere was an issue collecting the information it takes to run this application.                             %s\n\n", TERMREDREV, TERMCLEAR)
// 	fmt.Printf("%sYou can try again by typing %sdeploystack install%s at the command prompt  %s\n\n", TERMREDB, TERMREDREV, TERMCLEAR+TERMREDB, TERMCLEAR)
// 	fmt.Printf("%sIf the issue persists, please report at https://github.com/GoogleCloudPlatform/deploystack/issues %s\n\n", TERMREDB, TERMCLEAR)

// 	fmt.Printf("Extra diagnostic information:\n")

// 	if strings.Contains(err.Error(), "invalid token JSON from metadata") {
// 		fmt.Printf("timed out waiting for API activation, you must authorize API use to continue \n")
// 	}

// 	fmt.Println(err)
// 	os.Exit(1)
// }

// func handleEarlyShutdown(err error) {
// 	fmt.Printf("\n\n%sYou've chosen to stop moving forward through Deploystack.                             %s\n\n", TERMCYANB, TERMCLEAR)
// 	fmt.Printf("If this was an error, you can try again by typing %sdeploystack install%s at the command prompt. \n\n", TERMCYANB, TERMCLEAR)

// 	fmt.Printf("Reason: %s\n", err)
// 	os.Exit(1)
// }

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

// Stack represents the input config and output settings for this DeployStack
type Stack struct {
	Settings map[string]string
	Config   Config
}

// NewStack returns an initialized Stack
func NewStack() Stack {
	s := Stack{}
	s.Settings = make(map[string]string)
	return s
}

func (s *Stack) findAndReadConfig() (Config, error) {
	flavor := "json"
	config := Config{}

	configPath := ".deploystack/deploystack.json"
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		flavor = "yaml"
		configPath = ".deploystack/deploystack.yaml"
	}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		configPath = "deploystack.json"
	}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		wd, _ := os.Getwd()
		return config, fmt.Errorf("config file not present, looking for deploystack.json or .deploystack/deploystack.json in (%s)", wd)
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("unable to find or read config file: %s", err)
	}

	if flavor == "yaml" {
		config, err = NewConfigYAML(content)
		if err != nil {
			return config, fmt.Errorf("unable to parse config file: %s", err)
		}
		return config, nil
	}

	config, err = NewConfigJSON(content)
	if err != nil {
		return config, fmt.Errorf("unable to parse config file: %s", err)
	}

	return config, nil
}

func (s *Stack) findDSFolder(c Config, folder string) (string, error) {
	switch folder {
	case "messages":
		if c.PathMessages != "" {
			return c.PathMessages, nil
		}
	case "scripts":
		if c.PathScripts != "" {
			return c.PathScripts, nil
		}
	}

	path := fmt.Sprintf(".deploystack/%s", folder)

	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	path = folder

	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	return fmt.Sprintf("./%s", folder), fmt.Errorf("requirement (%s) was not found either in the root, or in .deploystack folder nor was it set in deploystack.json", folder)
}

func (s *Stack) findTFFolder(c Config) (string, error) {
	if c.PathTerraform != "" {
		return c.PathTerraform, nil
	}

	path := "terraform"

	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	return ".", nil
}

// FindAndReadRequired finds and reads in a Config from a json file.
func (s *Stack) FindAndReadRequired() error {
	config, err := s.findAndReadConfig()
	if err != nil {
		return fmt.Errorf("unable to parse config file: %s", err)
	}

	tfPath, err := s.findTFFolder(config)
	if err != nil {
		return fmt.Errorf("unable to locate terraform folder: %s", err)
	}
	config.PathTerraform = tfPath

	scriptPath, _ := s.findDSFolder(config, "scripts")
	if err != nil {
		log.Printf("WARNING - unable to locate scripts folder, folder not required, : %s", err)
	}
	config.PathScripts = scriptPath

	messagePath, err := s.findDSFolder(config, "messages")
	if err != nil {
		wd, _ := os.Getwd()
		return fmt.Errorf("unable to locate messages folder in (%s): %s", wd, err)
	}
	config.PathMessages = messagePath

	descText := fmt.Sprintf("%s/description.txt", messagePath)
	if _, err := os.Stat(descText); err == nil {
		description, err := ioutil.ReadFile(descText)
		if err != nil {
			return fmt.Errorf("unable to read description file: %s", err)
		}

		config.Description = string(description)
	}

	s.Config = config

	return nil
}

// AddSetting stores a setting key/value pair.
func (s Stack) AddSetting(key, value string) {
	k := strings.ToLower(key)
	s.Settings[k] = value
}

// GetSetting returns a setting value.
func (s Stack) GetSetting(key string) string {
	return s.Settings[key]
}

// DeleteSetting removes a setting value.
func (s Stack) DeleteSetting(key string) {
	delete(s.Settings, key)
}

// Terraform returns all of the settings as a Terraform variables format.
func (s Stack) Terraform() string {
	result := strings.Builder{}

	keys := []string{}
	for i := range s.Settings {
		keys = append(keys, i)
	}

	sort.Strings(keys)

	for _, v := range keys {
		if len(v) < 1 {
			continue
		}
		label := strings.ToLower(strings.ReplaceAll(v, " ", "_"))
		val := s.Settings[v]

		if label == "project_name" {
			continue
		}

		if label == "stack_name" {
			continue
		}

		if len(val) < 1 {
			continue
		}

		if val[0:1] == "[" {
			sb := strings.Builder{}
			sb.WriteString("[")
			tmp := strings.ReplaceAll(val, "[", "")
			tmp = strings.ReplaceAll(tmp, "]", "")
			sl := strings.Split(tmp, ",")

			for i, v := range sl {
				sl[i] = fmt.Sprintf("\"%s\"", v)
			}

			delimtext := strings.Join(sl, ",")

			sb.WriteString(delimtext)
			sb.WriteString("]")
			set := sb.String()
			set = strings.ReplaceAll(set, "\"\"", "")

			result.WriteString(fmt.Sprintf("%s=%s\n", label, set))
			continue
		}

		result.WriteString(fmt.Sprintf("%s=\"%s\"\n", label, val))

	}

	return result.String()
}

// TerraformFile exports TFVars format to input file.
func (s Stack) TerraformFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(s.Terraform()); err != nil {
		return err
	}

	return nil
}

// NewContactDataFromFile generates a new ContactData from a cached yaml file
func NewContactDataFromFile(file string) (gcloud.ContactData, error) {
	c := gcloud.NewContactData()

	dat, err := os.ReadFile(file)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(dat, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

// CheckForContact checks the local file system for a file containg domain
// registar contact info
func CheckForContact() gcloud.ContactData {
	contact := gcloud.ContactData{}
	if _, err := os.Stat(contactfile); err == nil {
		contact, err = NewContactDataFromFile(contactfile)
		if err != nil {
			log.Printf("domain registrar contact not cached")
		}
	}
	return contact
}

// CacheContact writes a file containg domain registar contact info to disk
// if it exists
func CacheContact(i interface{}) {
	switch v := i.(type) {
	case gcloud.ContactData:
		if v.AllContacts.Email != "" {
			yaml, err := v.YAML()
			if err != nil {
				log.Printf("could not convert contact to yaml: %s", err)
			}

			if err := os.WriteFile(contactfile, []byte(yaml), 0o644); err != nil {
				log.Printf("could not write contact to file: %s", err)
			}
		}
	}
}
