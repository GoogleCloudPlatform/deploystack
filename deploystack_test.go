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

package deploystack

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/deploystack/config"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	"github.com/GoogleCloudPlatform/deploystack/github"
)

func compareValues(label string, want interface{}, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("%s: expected: \n|%v|\ngot: \n|%v|", label, want, got)
	}
}

func TestPrecheck(t *testing.T) {
	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}

	testdata := fmt.Sprintf("%s/test_files/configs", wd)
	tests := map[string]struct {
		wd   string
		want string
	}{
		"single": {
			wd:   fmt.Sprintf("%s/preferred", testdata),
			want: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			oldWD, _ := os.Getwd()
			err := os.Chdir(tc.wd)
			if err != nil {
				t.Fatalf("error changing wd: %s", err)
			}

			out := captureOutput(func() {
				Precheck()
			})

			if !strings.Contains(tc.want, string(out)) {
				t.Fatalf("expected to contain: %+v, got: %+v", tc.want, string(out))
			}

			os.Chdir(oldWD)
		})
	}
}
func TestPrecheckMulti(t *testing.T) {
	// Precheck exits if it is called in testing with mutliple stacks
	// so make throwing the exit the test

	if os.Getenv("BE_CRASHER") == "1" {
		wd, err := filepath.Abs(".")
		if err != nil {
			t.Fatalf("error setting up environment for testing %v", err)
		}

		testdata := fmt.Sprintf("%s/test_files/configs", wd)
		path := fmt.Sprintf("%s/multi", testdata)
		oldWD, _ := os.Getwd()
		if err := os.Chdir(path); err != nil {
			t.Fatalf("error changing wd: %s", err)
		}

		Precheck()
		os.Chdir(oldWD)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestPrecheckMulti")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)

}

func captureOutput(f func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}

func TestCacheContact(t *testing.T) {
	tests := map[string]struct {
		in  gcloud.ContactData
		err error
	}{
		"basic": {
			in: gcloud.ContactData{
				AllContacts: gcloud.DomainRegistrarContact{
					Email: "test@example.com",
					Phone: "+155555551212",
					PostalAddress: gcloud.PostalAddress{
						RegionCode:         "US",
						PostalCode:         "94502",
						AdministrativeArea: "CA",
						Locality:           "San Francisco",
						AddressLines:       []string{"345 Spear Street"},
						Recipients:         []string{"Googler"},
					},
				},
			},
			err: nil,
		},
		"err": {
			in:  gcloud.ContactData{},
			err: fmt.Errorf("stat contact.yaml: no such file or directory"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			CacheContact(tc.in)

			if tc.err == nil {
				if _, err := os.Stat(contactfile); errors.Is(err, os.ErrNotExist) {
					t.Fatalf("expected no error,  got: %+v", err)
				}
			} else {
				if _, err := os.Stat(contactfile); err.Error() != tc.err.Error() {
					t.Fatalf("expected %+v, got: %+v", tc.err, err)
				}

			}

			os.Remove(contactfile)

		})
	}
}

func TestNewContactDataFromFile(t *testing.T) {
	tests := map[string]struct {
		in   string
		want gcloud.ContactData
		err  error
	}{
		"basic": {
			in: "test_files/contact/contact.yaml",
			want: gcloud.ContactData{
				AllContacts: gcloud.DomainRegistrarContact{
					Email: "test@example.com",
					Phone: "+155555551212",
					PostalAddress: gcloud.PostalAddress{
						RegionCode:         "US",
						PostalCode:         "94502",
						AdministrativeArea: "CA",
						Locality:           "San Francisco",
						AddressLines:       []string{"345 Spear Street"},
						Recipients:         []string{"Googler"},
					},
				},
			},
			err: nil,
		},
		"error": {
			in:   "test_files/contact/noexists.yaml",
			want: gcloud.ContactData{},
			err:  fmt.Errorf("open test_files/contact/noexists.yaml: no such file or directory"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewContactDataFromFile(tc.in)

			if tc.err == nil {

				if err != nil {
					t.Fatalf("expected no error,  got: %+v", err)
				}

				if !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("expected: %+v, got: %+v", tc.want, got)
				}

			} else {
				if err.Error() != tc.err.Error() {
					t.Fatalf("expected %+v, got: %+v", tc.err, err)
				}
			}

		})
	}
}

func TestCheckForContact(t *testing.T) {
	tests := map[string]struct {
		in   string
		want gcloud.ContactData
	}{
		"basic": {
			in: "test_files/contact/contact.yaml",
			want: gcloud.ContactData{
				AllContacts: gcloud.DomainRegistrarContact{
					Email: "test@example.com",
					Phone: "+155555551212",
					PostalAddress: gcloud.PostalAddress{
						RegionCode:         "US",
						PostalCode:         "94502",
						AdministrativeArea: "CA",
						Locality:           "San Francisco",
						AddressLines:       []string{"345 Spear Street"},
						Recipients:         []string{"Googler"},
					},
				},
			},
		},

		"empty": {
			in:   contactfile,
			want: gcloud.ContactData{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			oldContactFile := contactfile
			contactfile = tc.in

			got := CheckForContact()
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}

			contactfile = oldContactFile
		})
	}
}

func TestInit(t *testing.T) {
	errUnableToRead := errors.New("unable to read config file: ")
	tests := map[string]struct {
		path string
		want config.Stack
		err  error
	}{
		"error": {
			path: "sadasd",
			want: config.Stack{},
			err:  errUnableToRead,
		},
		"no_custom": {
			path: "test_files/dsfolders/no_customs",
			want: config.Stack{
				Config: config.Config{
					Title:         "TESTCONFIG",
					Name:          "test",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       true,
					Region:        true,
					RegionType:    "functions",
					RegionDefault: "us-central1",
				},
			},
			err: nil,
		},
		"no_name": {
			path: "test_files/dsfolders/no_name",
			want: config.Stack{
				Config: config.Config{
					Title:         "NONAME",
					Name:          "",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       true,
					Region:        true,
					RegionType:    "functions",
					RegionDefault: "us-central1",
				},
			},
			err: fmt.Errorf("could retrieve name of stack: could not open local git directory: repository does not exist \nDeployStack author: fix this by adding a 'name' key and value to the deploystack config"),
		},
		"custom": {
			path: "test_files/dsfolders/customs",
			want: config.Stack{
				Config: config.Config{
					Title:         "TESTCONFIG",
					Name:          "test",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       false,
					Region:        false,
					RegionType:    "",
					RegionDefault: "",
					CustomSettings: []config.Custom{
						{Name: "nodes", Description: "Nodes", Default: "3"},
					},
				},
			},
			err: nil,
		},
		"custom_options": {
			path: "test_files/dsfolders/customs_options",
			want: config.Stack{
				Config: config.Config{
					Title:         "TESTCONFIG",
					Name:          "test",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       false,
					Region:        false,
					RegionType:    "",
					RegionDefault: "",

					CustomSettings: []config.Custom{
						{
							Name:        "nodes",
							Description: "Nodes",
							Default:     "3",
							Options:     []string{"1", "2", "3"},
						},
					},
				},
			},
			err: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			oldWD, _ := os.Getwd()
			os.Chdir(tc.path)

			s, err := Init()

			if tc.err == nil {
				if err != nil {
					t.Fatalf("expected: no error got: %+v", err)
				}
			}

			if errors.Is(err, tc.err) {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

			os.Chdir(oldWD)

			compareValues("Name", tc.want.Config.Name, s.Config.Name, t)
			compareValues("Title", tc.want.Config.Title, s.Config.Title, t)
			compareValues("Description", tc.want.Config.Description, s.Config.Description, t)
			compareValues("Duration", tc.want.Config.Duration, s.Config.Duration, t)
			compareValues("Project", tc.want.Config.Project, s.Config.Project, t)
			compareValues("Region", tc.want.Config.Region, s.Config.Region, t)
			compareValues("RegionType", tc.want.Config.RegionType, s.Config.RegionType, t)
			compareValues("RegionDefault", tc.want.Config.RegionDefault, s.Config.RegionDefault, t)
			for i, v := range s.Config.CustomSettings {
				compareValues(v.Name, tc.want.Config.CustomSettings[i], v, t)
			}
		})
	}
}

func TestShortName(t *testing.T) {
	tests := map[string]struct {
		in   string
		want string
	}{
		"deploystack-repo":     {in: "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry", want: "cost-sentry"},
		"non-deploystack-repo": {in: "https://github.com/tpryan/microservices-demo", want: "microservices-demo"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := Meta{}
			m.Github.Name = tc.in

			got := m.ShortName()
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestShortNameUnderscore(t *testing.T) {
	tests := map[string]struct {
		in   string
		want string
	}{
		"deploystack-repo":     {in: "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry", want: "cost_sentry"},
		"non-deploystack-repo": {in: "https://github.com/tpryan/microservices-demo", want: "microservices_demo"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := Meta{}
			m.Github.Name = tc.in

			got := m.ShortNameUnderscore()
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestNewMeta(t *testing.T) {
	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}

	testdata := fmt.Sprintf("%s/test_files/repos", wd)
	tests := map[string]struct {
		repo string
		path string
		want Meta
	}{
		"defaultbranch": {
			repo: "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry",
			path: testdata,
			want: Meta{
				Github: github.Repo{
					Name:   "deploystack-cost-sentry",
					Owner:  "GoogleCloudPlatform",
					Branch: "main",
				},
				LocalPath: fmt.Sprintf("%s/deploystack-cost-sentry", testdata)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewMeta(tc.repo, tc.path, "")
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want.Github.Name, got.Github.Name) {
				t.Fatalf("expected: %v, got: %v", tc.want.Github.Name, got.Github.Name)
			}

			if !reflect.DeepEqual(tc.want.Github.Branch, got.Github.Branch) {
				t.Fatalf("expected: %v, got: %v", tc.want.Github.Branch, got.Github.Branch)
			}

			if !reflect.DeepEqual(tc.want.LocalPath, got.LocalPath) {
				t.Fatalf("expected: %v, got: %v", tc.want.LocalPath, got.LocalPath)
			}

			os.RemoveAll(got.LocalPath)
		})
	}
}

func TestGetRepo(t *testing.T) {
	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}

	testdata := fmt.Sprintf("%s/test_files/repoforgithub", wd)
	tests := map[string]struct {
		repo string
		path string
		want string
		err  error
	}{
		"deploystack-nosql-client-server": {
			repo: "deploystack-nosql-client-server",
			path: testdata,
			want: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
		},
		"nosql-client-server": {
			repo: "nosql-client-server",
			path: testdata,
			want: fmt.Sprintf("%s/nosql-client-server", testdata),
		},

		"deploystack-cost-sentry": {
			repo: "deploystack-cost-sentry",
			path: testdata,
			want: fmt.Sprintf("%s/deploystack-cost-sentry_1", testdata),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			got, err := DownloadRepo(tc.repo, tc.path)

			if tc.err == nil && err != nil {
				t.Fatalf("expected: no error got: %+v", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}

			if _, err := os.Stat(tc.want); os.IsNotExist(err) {
				t.Fatalf("expected: %s to exist it does not", err)
			}

			err = os.RemoveAll(tc.want)
			if err != nil {
				t.Logf(err.Error())
			}

		})
	}
}

func TestGetAcceptableDir(t *testing.T) {
	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}
	testdata := fmt.Sprintf("%s/test_files/repoforgithub", wd)

	tests := map[string]struct {
		in   string
		want string
	}{
		"doesnotexist": {
			in:   fmt.Sprintf("%s/testfolder", testdata),
			want: fmt.Sprintf("%s/testfolder", testdata),
		},
		"exists": {
			in:   fmt.Sprintf("%s/alreadyexists", testdata),
			want: fmt.Sprintf("%s/alreadyexists_2", testdata),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := UniquePath(tc.in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestCloneFromRepo(t *testing.T) {

	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}

	testdata := fmt.Sprintf("%s/test_files/repoforgithub", wd)
	tests := map[string]struct {
		repo string
		path string
		want string
		err  error
	}{
		"deploystack-nosql-client-server": {
			repo: "deploystack-nosql-client-server",
			path: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
			want: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
		},
		"nosql-client-server": {
			repo: "nosql-client-server",
			path: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
			want: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			err := CloneByName(tc.repo, tc.path)

			if tc.err == nil && err != nil {
				t.Fatalf("expected: no error got: %+v", err)
			}

			if _, err := os.Stat(tc.want); os.IsNotExist(err) {
				t.Fatalf("expected: %s to exist it does not", err)
			}

			err = os.RemoveAll(tc.path)
			if err != nil {
				t.Logf(err.Error())
			}

		})
	}
}
