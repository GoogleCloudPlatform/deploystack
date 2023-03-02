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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Repo contains the details of a github repo for the purpose of downloading
type Repo struct {
	Name   string `json:"name" yaml:"name"`
	Owner  string `json:"owner" yaml:"owner"`
	Branch string `json:"branch" yaml:"branch"`
}

func (r Repo) URL() string {
	return fmt.Sprintf("https://github.com/%s/%s", r.Owner, r.Name)
}

func (r Repo) ReferenceName() string {
	return fmt.Sprintf("refs/heads/%s", r.Branch)
}

// NewRepo generates Github from a url that might contain branch information
func NewRepo(repo string) Repo {
	result := Repo{}

	result.Branch = "main"

	input := strings.ReplaceAll(repo, "https://github.com/", "")

	sl := strings.Split(input, "/")
	result.Owner = sl[0]
	result.Name = sl[1]

	if strings.Contains(repo, "/tree/") {
		end := strings.Index(repo, "/tree/")
		result.Branch = repo[end+6:]
	}

	return result
}

// Path returns where this repo should exist locally given the input path
func (r Repo) Path(path string) string {
	result := filepath.Base(r.Name)
	result = fmt.Sprintf("%s/%s", path, result)
	return result
}

// Clone performs a git clone to the directory of our choosing
func (r Repo) Clone(path string) error {
	log.Printf("Clone called %+v %s", r, path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err = git.PlainClone(
			path,
			false,
			&git.CloneOptions{
				URL:           r.URL(),
				ReferenceName: plumbing.ReferenceName(r.ReferenceName()),
				Progress:      nil,
			})

		if err != nil {
			return fmt.Errorf("cannot get repo (%s) : %s", r.URL(), err)
		}

	}

	return nil
}
