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

package github

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testFilesDir = filepath.Join(os.Getenv("DEPLOYSTACK_PATH"), "testdata")

func TestClone(t *testing.T) {

	testdata := filepath.Join(testFilesDir, "repoforgithub")

	tests := map[string]struct {
		in   Repo
		path string
		err  error
	}{
		"basic": {
			in: Repo{
				Name:   "deploystack-nosql-client-server",
				Owner:  "GoogleCloudPlatform",
				Branch: "main",
			},
			path: filepath.Join(testdata, "deploystack-nosql-client-server"),
		},
		"error": {
			in: Repo{
				Name:   "deploystack-nosql-client-server-not-exist",
				Owner:  "GoogleCloudPlatform",
				Branch: "main",
			},
			path: filepath.Join(testdata, "deploystack-nosql-client-server-not-exist"),
			err:  fmt.Errorf("cannot clone repo"),
		},
		"overwrite": {
			in: Repo{
				Name:   "deploystack-nosql-client-server",
				Owner:  "GoogleCloudPlatform",
				Branch: "main",
			},
			path: filepath.Join(testdata, "alreadyexists"),
			err:  fmt.Errorf("already exists"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.in.Clone(tc.path)

			if tc.err == nil && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
				t.Skip()
			}

			if _, err := os.Stat(tc.path + "/.git"); os.IsNotExist(err) {
				t.Errorf("expected: %s to exist it does not", err)
			}

			err = os.RemoveAll(tc.path)
			if err != nil {
				t.Logf(err.Error())
			}
		})
	}
}

func TestNewRepo(t *testing.T) {
	tests := map[string]struct {
		in   string
		want Repo
	}{
		"defaultbranch": {
			in: "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry",
			want: Repo{
				Name:   "deploystack-cost-sentry",
				Owner:  "GoogleCloudPlatform",
				Branch: "main",
			}},
		"otherbranch": {
			in: "https://github.com/tpryan/microservices-demo/tree/deploystack-enable",
			want: Repo{
				Name:   "microservices-demo",
				Owner:  "tpryan",
				Branch: "deploystack-enable",
			}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := NewRepo(tc.in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestRepoPath(t *testing.T) {
	tests := map[string]struct {
		in   Repo
		path string
		want string
	}{
		"defaultbranch": {
			in: Repo{
				Name:   "deploystack-cost-sentry",
				Owner:  "GoogleCloudPlatform",
				Branch: "main",
			},
			path: ".",
			want: "./deploystack-cost-sentry",
		},
		"otherbranch": {
			in: Repo{
				Name:   "microservices-demo",
				Owner:  "tpryan",
				Branch: "deploystack-enable",
			},
			path: ".",
			want: "./microservices-demo",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Path(tc.path)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := map[string]struct {
		name    string
		options []Option
		want    Repo
	}{
		"basic": {
			name: "deploystack-nosql-client-server",
			want: Repo{
				Owner:  "GoogleCloudPlatform",
				Name:   "deploystack-nosql-client-server",
				Branch: "main",
			},
		},
		"owner": {
			name: "deploystack-nosql-client-server",
			want: Repo{
				Owner:  "tpryan",
				Name:   "deploystack-nosql-client-server",
				Branch: "main",
			},
			options: []Option{
				Owner("tpryan"),
			},
		},
		"branch": {
			name: "deploystack-nosql-client-server",
			want: Repo{
				Owner:  "GoogleCloudPlatform",
				Name:   "deploystack-nosql-client-server",
				Branch: "experimental",
			},
			options: []Option{
				Branch("experimental"),
			},
		},
		"owner and branch": {
			name: "deploystack-nosql-client-server",
			want: Repo{
				Owner:  "tpryan",
				Name:   "deploystack-nosql-client-server",
				Branch: "experimental",
			},
			options: []Option{
				Owner("tpryan"),
				Branch("experimental"),
			},
		},
		"siteurl": {
			name: "deploystack-cost-sentry",
			want: Repo{
				Owner:  "GoogleCloudPlatform",
				Name:   "deploystack-cost-sentry",
				Branch: "main",
			},
			options: []Option{
				SiteURL("https://github.com/GoogleCloudPlatform/deploystack-cost-sentry"),
			},
		},
		"siteurl branch": {
			name: "deploystack-nosql-client-server",
			want: Repo{
				Owner:  "tpryan",
				Name:   "microservices-demo",
				Branch: "deploystack-enable",
			},
			options: []Option{
				SiteURL("https://github.com/tpryan/microservices-demo/tree/deploystack-enable"),
			},
		},
		"siteurl name empty": {
			name: "",
			want: Repo{
				Owner:  "GoogleCloudPlatform",
				Name:   "deploystack-cost-sentry",
				Branch: "main",
			},
			options: []Option{
				SiteURL("https://github.com/GoogleCloudPlatform/deploystack-cost-sentry"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := New(tc.name, tc.options...)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPopulate(t *testing.T) {
	tests := map[string]struct {
		in   Repo
		want Repo
		err  error
	}{
		"basic": {
			in: Repo{
				Owner:  "GoogleCloudPlatform",
				Name:   "deploystack-nosql-client-server",
				Branch: "main",
			},
			want: Repo{
				Owner:       "GoogleCloudPlatform",
				Name:        "deploystack-nosql-client-server",
				Branch:      "main",
				Description: "A terraform solution that will create 2 VMS and a firewall rule, connect them all, to serve up a API powered by mongo.",
			},
		},
		"no description": {
			in: Repo{
				Owner:  "tpryan",
				Name:   "deploystack",
				Branch: "main",
			},
			want: Repo{
				Owner:       "tpryan",
				Name:        "deploystack",
				Branch:      "main",
				Description: "",
			},
		},
		"doesnt exist": {
			in: Repo{
				Owner:  "tpryan",
				Name:   "deploystack-dontexist",
				Branch: "main",
			},
			want: Repo{
				Owner:       "tpryan",
				Name:        "deploystack-dontexist",
				Branch:      "main",
				Description: "",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in
			got.Populate()
			assert.Equal(t, tc.want, got)

		})
	}
}
