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
	"os"
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/deploystack/config"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
)

func compareValues(label string, want interface{}, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("%s: expected: \n|%v|\ngot: \n|%v|", label, want, got)
	}
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

func TestBasic(t *testing.T) {
	tests := map[string]struct {
		in   string
		want string
	}{
		"basic": {
			in:   "test",
			want: "test",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
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
