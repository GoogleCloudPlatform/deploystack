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

// Package github includes all of the client calls to github for building
// automation tools for deploystack
package github

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// GITHUBHOST is the hostname of where to find github on the web
const GITHUBHOST = "github.com"

// Repo contains the details of a github repo for the purpose of downloading
type Repo struct {
	Name        string `json:"name" yaml:"name"`
	Owner       string `json:"owner" yaml:"owner"`
	Branch      string `json:"branch" yaml:"branch"`
	Description string `json:"description" yaml:"description"`
}

// URL returns the github url of this repo
func (r Repo) URL() string {
	return fmt.Sprintf("https://%s/%s/%s", GITHUBHOST, r.Owner, r.Name)
}

// ReferenceName is a shortcut to the reference for the current branch
func (r Repo) ReferenceName() string {
	return fmt.Sprintf("refs/heads/%s", r.Branch)
}

// NewRepo generates Github from a url that might contain branch information
//
// Deprecated: Use New() Instead
func NewRepo(repo string) Repo {
	return New("", SiteURL(repo))
}

// Populate gets metadata from github
func (r *Repo) Populate() error {

	client := github.NewClient(nil)

	repo, _, err := client.Repositories.Get(context.Background(), r.Owner, r.Name)
	if err != nil {
		return err
	}

	r.Description = repo.GetDescription()

	return nil
}

// Path returns where this repo should exist locally given the input path
func (r Repo) Path(path string) string {
	result := filepath.Base(r.Name)
	result = fmt.Sprintf("%s/%s", path, result)
	return result
}

// Clone performs a git clone to the directory of our choosing
func (r Repo) Clone(path string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("directory (%s) already exists", path)
	}

	options := &git.CloneOptions{
		URL:           r.URL(),
		ReferenceName: plumbing.ReferenceName(r.ReferenceName()),
	}

	if _, err := git.PlainClone(path, false, options); err != nil {
		return fmt.Errorf("cannot clone repo (%s) : %s", r.URL(), err)
	}

	return nil
}

// Option defines option functions that can be passed in to modify repos on New
type Option func(*Repo)

// New generates Github from a series of options
func New(name string, options ...Option) Repo {
	r := Repo{
		Name:   name,
		Owner:  "GoogleCloudPlatform",
		Branch: "main",
	}

	for _, f := range options {
		f(&r)
	}

	return r
}

// Owner sets the owner on a repo
func Owner(o string) Option {
	return func(r *Repo) {
		r.Owner = o
	}
}

// Branch sets the branch on a repo
func Branch(b string) Option {
	return func(r *Repo) {
		r.Branch = b
	}
}

// SiteURL sets the owner, name, and branch of a repo based on a URL
func SiteURL(u string) Option {
	return func(r *Repo) {

		input := strings.ReplaceAll(u, fmt.Sprintf("https://%s/", GITHUBHOST), "")
		sl := strings.Split(input, "/")
		r.Owner = sl[0]
		r.Name = sl[1]
		if strings.Contains(u, "/tree/") {
			end := strings.Index(u, "/tree/")
			r.Branch = u[end+6:]
		}

	}
}
