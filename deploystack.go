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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/deploystack/config"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	"github.com/GoogleCloudPlatform/deploystack/github"
	"github.com/GoogleCloudPlatform/deploystack/terraform"
	"github.com/GoogleCloudPlatform/deploystack/tui"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

var (
	opts             = option.WithCredentialsFile("")
	credspath        = ""
	defaultUserAgent = "deploystack"
	contactfile      = "contact.yaml"
)

// Init initializes a Deploystack stack by looking on the local file system
func Init(path string) (*config.Stack, error) {
	s := config.NewStack()

	if err := s.FindAndReadRequired(path); err != nil {
		return &s, fmt.Errorf("could not read config file: %s", err)
	}

	if s.Config.Name == "" {
		if err := s.Config.ComputeName(path); err != nil {
			return &s, fmt.Errorf("could not retrieve name of stack: %s \nDeployStack author: fix this by adding a 'name' key and value to the deploystack config", err)
		}
		s.AddSetting("stack_name", s.Config.Name)
	}
	s.Config.Setwd(path)

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

// ContactCheck checks the local file system for a file containing domain
// registar contact info
func ContactCheck() gcloud.ContactData {
	contact := gcloud.ContactData{}
	if _, err := os.Stat(contactfile); err == nil {
		f, err := os.Open(contactfile)
		if err != nil {
			return contact
		}

		if _, err = contact.ReadFrom(f); err != nil {
			return contact
		}
	}
	return contact
}

// ContactSave writes a file containing domain registar contact info to disk
// if it exists
func ContactSave(i interface{}) {
	// We can ignore errors - this is an convenience to the user
	// not a necessity
	switch v := i.(type) {
	case gcloud.ContactData:
		if v.AllContacts.Email == "" {
			return
		}

		f, err := os.Create(contactfile)
		if err != nil {
			return
		}

		if _, err := v.WriteTo(f); err != nil {
			return
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

// NewMeta allows project to point at local directories for info
// as well as pulling down from github
func NewMeta(path string) (Meta, error) {
	d := Meta{}

	s := config.NewStack()

	if err := s.FindAndRead(path, false); err == nil {
		d.DeployStack = s.Config
	}

	tfPath := filepath.Join(path, s.Config.PathTerraform)
	if b, err := terraform.Extract(tfPath); err == nil {

		if b != nil {
			d.Terraform = *b
		}
	}

	return d, nil
}

// ShortName retrieves the shortname of whatever we are calling this stack
func (m Meta) ShortName() string {
	r := filepath.Base(m.Github.Name)
	r = strings.ReplaceAll(r, "deploystack-", "")
	return r
}

// ShortNameUnderscore retrieves the shortname of whatever we are calling
// this stack replacing hyphens with underscores
func (m Meta) ShortNameUnderscore() string {
	r := m.ShortName()
	r = strings.ReplaceAll(r, "-", "_")
	return r
}

// Suggest will provide it's best guess of what the deploystack config should
// be based on the contents of the repo, including an existing deploystack config
func (m Meta) Suggest() (config.Config, error) {
	out := m.DeployStack.Copy()

	name := filepath.Base(m.Github.URL())
	name = strings.ReplaceAll(name, "deploystack-", "")
	title := strings.ReplaceAll(name, "-", " ")
	caser := cases.Title(language.AmericanEnglish)
	title = caser.String(title)

	if m.DeployStack.Name == "" {
		out.Name = name
	}

	if m.DeployStack.Title == "" {
		out.Title = title
	}

	if len(m.Terraform) == 0 {
		return out, errors.New("suggest: terraform was empty")
	}

	if m.DeployStack.PathTerraform == "" {
		out.PathTerraform = filepath.Dir(m.Terraform[0].File)
	}

	resources, err := terraform.NewGCPResources()
	if err != nil {
		return out, fmt.Errorf("could not get terraform resource meta data: %w", err)
	}

	for _, v := range m.Terraform {
		switch v.Kind {
		case "variable":
			// For now if there are default values, don't bother capturing
			if !v.NoDefault() {
				continue
			}

			switch v.Name {
			case "project_id":
				out.Project = true
			case "project_number":
				out.ProjectNumber = true
			case "billing_account":
				out.BillingAccount = true
			case "region":
				out.RegionDefault = "us-central1"
				out.Region = true
				out.RegionType = "compute"

				if r := m.Terraform.Search("google_cloud_run", "type"); len(r) > 0 {
					out.RegionType = "run"
				}

				if r := m.Terraform.Search("google_cloudfunctions", "type"); len(r) > 0 {
					out.RegionType = "functions"
				}

			case "zone":
				out.Zone = true
			default:
				checkCustom := out.CustomSettings.Get(v.Name)
				checkAuthor := out.AuthorSettings.Find(v.Name)

				if checkCustom.Name == "" && checkAuthor == nil {
					cust := config.Custom{}
					cust.Name = v.Name
					cust.Type = v.Type
					out.CustomSettings = append(out.CustomSettings, cust)
				}

			}
		case "managed":
			product := resources.GetProduct(v.Type)

			if product == "" {
				continue
			}

			add := true
			for _, v := range out.Products {
				if v.Product == product {
					add = false
					break
				}
			}

			if add {
				p := config.Product{Product: product}
				out.Products = append(out.Products, p)
			}

		}
	}

	return out, nil
}

// DownloadRepo takes a name of a GoogleCloudPlatform repo or a
// GoogleCloudPlatform/deploystack-[name] repo, and downloads it into a unique
// folder name, and outputs that name
func DownloadRepo(repo github.Repo, path string) (string, error) {
	candidate := fmt.Sprintf("%s/%s", path, repo.Name)
	dir := UniquePath(candidate)
	return dir, repo.Clone(dir)

}

// UniquePath returns either the input candidate path if it does not exist,
// or a path like the input candidate with increasing numbers appended to it
// until the output name is a path that does not exist
func UniquePath(candidate string) string {
	if _, err := os.Stat(candidate); os.IsNotExist(err) {
		return candidate
	}
	i := 1
	for {
		attempted := fmt.Sprintf("%s_%d", candidate, i)
		if _, err := os.Stat(attempted); os.IsNotExist(err) {
			return attempted
		}
		i++
	}
}

// AttemptRepo will try to download a repo - only from GoogleCloudPlatform. If
// it fails, it will append "deploystack-" to the front of the requested name
// and try again.
func AttemptRepo(name, wd string) (string, github.Repo, error) {

	gh := github.New(name)

	dir, err := DownloadRepo(gh, wd)
	if err != nil {
		// This allows using a shortened name of the repo as the label here.
		if !strings.Contains(gh.Name, "deploystack-") {
			gh = github.New(fmt.Sprintf("deploystack-%s", name))
			dir, err = DownloadRepo(gh, wd)
		}
	}
	return dir, gh, err
}

// WriteConfig will drop a .deploystack folder with deploystack.yaml file for
// repos that do not have one.
func WriteConfig(dir string, gh github.Repo) error {
	m, err := NewMeta(dir)
	if err != nil {
		tui.Fatal(err)
	}

	m.Github = gh

	sb := strings.Builder{}
	sb.WriteString("This DeployStack is running automatically based on our best guess\n")
	sb.WriteString("of what the Terraform files present in the github repo you chose \n")
	sb.WriteString("need in terms of input \n")
	sb.WriteString("\n\n")
	sb.WriteString("If you would like to see proper information, please file an issue at \n")
	sb.WriteString(fmt.Sprintf("%s/issues", gh.URL()))
	m.DeployStack.Description = sb.String()

	config, err := m.Suggest()
	if err != nil {
		return fmt.Errorf("could not make suggestion based on repo: %s", err)
	}

	configyaml, err := config.Marshal("yaml")
	if err != nil {
		return fmt.Errorf("could turn suggestions into yaml: %s", err)
	}

	target := fmt.Sprintf("%s/.deploystack", dir)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		os.MkdirAll(target+"/messages", 0700)
		err := os.WriteFile(target+"/deploystack.yaml", []byte(configyaml), 0644)
		if err != nil {
			return fmt.Errorf("could not write suggestions down as config: %s", err)
		}
	}

	return nil
}
