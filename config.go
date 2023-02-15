// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploystack

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
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
	CustomSettings       []Custom          `json:"custom_settings" yaml:"custom_settings"`
	ConfigureGCEInstance bool              `json:"configure_gce_instance" yaml:"configure_gce_instance"`
	DocumentationLink    string            `json:"documentation_link" yaml:"documentation_link"`
	PathTerraform        string            `json:"path_terraform" yaml:"path_terraform"`
	PathMessages         string            `json:"path_messages" yaml:"path_messages"`
	PathScripts          string            `json:"path_scripts" yaml:"path_scripts"`
	Projects             Projects          `json:"projects" yaml:"projects"`
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
