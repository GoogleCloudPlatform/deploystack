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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/deploystack/config"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	"github.com/GoogleCloudPlatform/deploystack/github"
	"github.com/GoogleCloudPlatform/deploystack/terraform"
	"github.com/GoogleCloudPlatform/deploystack/tui"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v2"
)

var (
	opts             = option.WithCredentialsFile("")
	credspath        = ""
	defaultUserAgent = "deploystack"
	contactfile      = "contact.yaml"
)

// Init initializes a Deploystack stack by looking on teh local file system
func Init() (*config.Stack, error) {
	s := config.NewStack()

	if err := s.FindAndReadRequired(); err != nil {
		return &s, fmt.Errorf("could not read config file: %s", err)
	}

	if s.Config.Name == "" {
		if err := s.Config.ComputeName(); err != nil {
			return &s, fmt.Errorf("could not retrieve name of stack: %s \nDeployStack author: fix this by adding a 'name' key and value to the deploystack config", err)
		}
		s.AddSetting("stack_name", s.Config.Name)
	}

	return &s, nil
}

// Precheck handles the logic around switching working directories for multiple
// stacks in one repo
func Precheck() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	reports, err := config.FindConfigReports(wd)
	if err != nil {
		return err
	}

	if len(reports) > 1 {
		stackPath := tui.PreCheck(reports)
		if err := os.Chdir(stackPath); err != nil {
			return err
		}
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

	// If anything goes wrong with this, it's fine, this is a convenience
	// for the next time someone runs this.
	case gcloud.ContactData:
		if v.AllContacts.Email == "" {
			return
		}

		if v.AllContacts.Email != "" {
			yaml, err := v.YAML()
			if err != nil {
				return
			}

			if err := os.WriteFile(contactfile, []byte(yaml), 0o644); err != nil {
				return
			}
		}
	}
}

// Meta is a datastructure that combines the Deploystack, github and Terraform
// bits of metadata about a stack.
type Meta struct {
	DeployStack config.Config
	Terraform   terraform.Blocks `json:"terraform" yaml:"terraform"`
	Github      github.Repo      `json:"github" yaml:"github"`
	LocalPath   string           `json:"localpath" yaml:"localpath"`
}

// NewMeta downloads a github repo and parses the DeployStack and Terraform
// information from the stack.
func NewMeta(repo, path, dspath string) (Meta, error) {
	g := github.NewRepo(repo)

	log.Printf("cloning to path: %s", path)
	if err := g.Clone(g.Path(path)); err != nil {
		return Meta{}, fmt.Errorf("cannot clone repo: %s", err)
	}

	d, err := NewMetaFromLocal(g.Path(path) + dspath)
	if err != nil {
		return Meta{}, fmt.Errorf("cannot parse deploystack into: %s", err)
	}
	d.Github = g
	d.LocalPath = g.Path(path)

	return d, nil
}

// NewMetaFromLocal allows project to point at local directories for info
// as well as pulling down from github
func NewMetaFromLocal(path string) (Meta, error) {
	d := Meta{}
	orgpwd, err := os.Getwd()
	if err != nil {
		return d, fmt.Errorf("could not get the wd: %s", err)
	}
	if err := os.Chdir(path); err != nil {
		return d, fmt.Errorf("could not change the wd: %s", err)
	}

	s := config.NewStack()

	if err := s.FindAndReadRequired(); err != nil {
		log.Printf("could not read config file: %s", err)
	}

	b, err := terraform.Extract(s.Config.PathTerraform)
	if err != nil {
		log.Printf("couldn't extract from TF file: %s", err)
	}

	if b != nil {
		d.Terraform = *b
	}

	d.DeployStack = s.Config

	if err := os.Chdir(orgpwd); err != nil {
		return d, fmt.Errorf("could not change the wd back: %s", err)
	}
	return d, nil
}

// ShortName retrieves the shortname of whatever we are calling this stack
func (d Meta) ShortName() string {
	r := filepath.Base(d.Github.Name)
	r = strings.ReplaceAll(r, "deploystack-", "")
	return r
}

// ShortNameUnderscore retrieves the shortname of whatever we are calling
// this stack replacing hyphens with underscores
func (d Meta) ShortNameUnderscore() string {
	r := d.ShortName()
	r = strings.ReplaceAll(r, "-", "_")
	return r
}
