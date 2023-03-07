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

package gcloud

import (
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func TestGetProjectNumbers(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)

	tests := map[string]struct {
		input string
		want  string
	}{
		"1": {input: creds["project_id"], want: creds["project_number"]},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.ProjectNumberGet(tc.input)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestCheckProject(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)

	tests := map[string]struct {
		input string
		want  bool
	}{
		"Does Exists":     {input: creds["project_id"], want: true},
		"Does Not Exists": {input: "ds-does-not-exst", want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := c.ProjectExists(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestGetProjectParent(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		input string
		want  *cloudresourcemanager.ResourceId
	}{
		"1": {
			input: creds["project_id"],
			want: &cloudresourcemanager.ResourceId{
				Id:   creds["parent"],
				Type: creds["parent_type"],
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.ProjectParentGet(tc.input)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestGetProjects(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		filepath string
	}{
		"1": {filepath: "gcloudout/projects.txt"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.ProjectList()

			f := filepath.Join(testFilesDir, tc.filepath)
			raw := readTestFile(f)
			want := strings.Split(strings.TrimSpace(raw), "\n")

			gotfiltered := []string{}

			for _, v := range got {
				gotfiltered = append(gotfiltered, v.Name)
			}

			sort.Strings(want)
			sort.Strings(gotfiltered)

			extraGots := []string{}
			for _, gotItem := range gotfiltered {
				found := false
				for _, wantItem := range want {
					if wantItem == gotItem {
						found = true
						break
					}
				}

				if !found {
					extraGots = append(extraGots, gotItem)
				}

			}

			if len(extraGots) > 0 {
				for _, v := range extraGots {
					if !strings.Contains(v, "ds-unittest") && !strings.Contains(v, "ds-test-") {
						t.Logf("extra gots: %v ", extraGots)
						t.Fatalf("expected: %v got: %v", len(want), len(gotfiltered))
					}
				}

			}

			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		input   string
		err     error
		noRando bool
	}{
		"Too long": {
			input: "zprojectnamedeletethisprojectnamehastoomanycharacters",
			err:   ErrorProjectCreateTooLong,
		},
		"Bad Chars": {
			input: "ALLUPERCASEDONESTWORK",
			err:   ErrorProjectInvalidCharacters,
		},
		"Spaces": {
			input: "spaces in name",
			err:   ErrorProjectInvalidCharacters,
		},
		// TODO: Figure out why this isn't working for test account
		// "Duplicate": {
		// 	input:   projectID,
		// 	err:     ErrorProjectAlreadyExists,
		// 	noRando: true,
		// },
		"Too short": {
			input: "",
			err:   ErrorProjectCreateTooShort,
		},
		// TODO: Figure out why this isn't working for test account
		// "Should work": {
		// 	input: "ds-unittest",
		// 	err:   nil,
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			name := tc.input + randSeq(5)
			if tc.noRando {
				name = tc.input
			}

			err := c.ProjectCreate(name, creds["parent"], creds["parent_type"])

			// Don't accidently delete the project that you are using to run
			// these tests. Yes I found out the hard way
			if name != projectID {
				c.ProjectDelete(name)
			}

			if err != tc.err {
				t.Fatalf("expected: %v, got: %v project: %s", tc.err, err, name)
			}
		})
	}
}

func TestGetProject(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	expected := projectID

	old, err := c.ProjectIDGet()
	if err != nil {
		t.Fatalf("retrieving old project: expected: no error, got: %v", err)
	}

	if err := c.ProjectIDSet(expected); err != nil {
		t.Fatalf("setting expecgted project: expected: no error, got: %v", err)
	}

	got, err := c.ProjectIDGet()
	if err != nil {
		t.Fatalf("expected: no error, got: %v", err)
	}

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("expected: %v, got: %v", expected, got)
	}

	if err := c.ProjectIDSet(old); err != nil {
		t.Fatalf("resetting old project: expected: no error, got: %v", err)
	}
}
