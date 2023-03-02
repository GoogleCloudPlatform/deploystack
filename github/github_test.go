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
	"strings"
	"testing"
)

func TestClone(t *testing.T) {
	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}

	testdata := fmt.Sprintf("%s/test_files", wd)

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
			path: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
		},
		"error": {
			in: Repo{
				Name:   "deploystack-nosql-client-server-not-exist",
				Owner:  "GoogleCloudPlatform",
				Branch: "main",
			},
			path: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
			err:  fmt.Errorf("cannot get repo"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.in.Clone(tc.path)

			if tc.err == nil && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if tc.err != nil {
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("expected: %v, got: %v", tc.err, err)
				}
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
