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
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/kylelemons/godebug/diff"
	domainspb "google.golang.org/genproto/googleapis/cloud/domains/v1beta1"
	"google.golang.org/genproto/googleapis/type/postaladdress"
)

func TestContactDataYAML(t *testing.T) {
	tests := map[string]struct {
		file    string
		contact ContactData
		err     error
	}{
		"simple": {
			file: "contact/contact_sample.yaml",
			contact: ContactData{DomainRegistrarContact{
				"you@example.com",
				"+1 555 555 1234",
				PostalAddress{
					"US",
					"94105",
					"CA",
					"San Francisco",
					[]string{"345 Spear Street"},
					[]string{"Your Name"},
				},
			}},
			err: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := filepath.Join(testFilesDir, tc.file)
			dat, err := os.ReadFile(f)
			if err != nil {
				t.Fatalf("err could not get file for testing: (%s)", err)
			}

			want := string(dat)
			got, err := tc.contact.YAML()

			if err != tc.err {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

			if !reflect.DeepEqual(want, got) {
				fmt.Println(diff.Diff(want, got))
				t.Fatalf("expected: \n|%v|\ngot: \n|%v|", want, got)
			}
		})
	}
}

func TestContactDataReadFrom(t *testing.T) {
	tests := map[string]struct {
		file string
		want ContactData
		err  error
	}{
		"simple": {
			file: "contact/contact_sample.yaml",
			want: ContactData{DomainRegistrarContact{
				"you@example.com",
				"+1 555 555 1234",
				PostalAddress{
					"US",
					"94105",
					"CA",
					"San Francisco",
					[]string{"345 Spear Street"},
					[]string{"Your Name"},
				},
			}},
			err: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cf := filepath.Join(testFilesDir, tc.file)
			f, _ := os.Open(cf)
			got := newContactData()
			_, err := got.ReadFrom(f)

			if err != tc.err {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v got: %v", tc.want, got)
			}
		})
	}
}

func TestContactDataWriteTo(t *testing.T) {
	tests := map[string]struct {
		contact ContactData
	}{
		"basic": {
			contact: ContactData{DomainRegistrarContact{
				"you@example.com",
				"+1 555 555 1234",
				PostalAddress{
					"US",
					"94105",
					"CA",
					"San Francisco",
					[]string{"345 Spear Street"},
					[]string{"Your Name"},
				},
			}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f, err := os.CreateTemp("", "contact")
			defer os.Remove(f.Name())
			if err != nil {
				t.Fatalf("creating tmp: expected no error, got %s", err)
			}

			if _, err = tc.contact.WriteTo(f); err != nil {
				t.Fatalf("writing tmp: expected no error, got %s", err)
			}

			f2, err := os.Open(f.Name())
			got := newContactData()
			if _, err = got.ReadFrom(f2); err != nil {
				t.Fatalf("reading tmp: expected no error, got %s", err)
			}

			if !reflect.DeepEqual(tc.contact, got) {
				diff := deep.Equal(tc.contact, got)
				t.Errorf("compare failed: %v", diff)
			}

		})
	}
}

func TestDomainIsAvailable(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		domain    string
		wantAvail string
		wantCost  string
		err       error
	}{
		"example.com": {
			domain:    "example.com",
			wantAvail: "UNAVAILABLE",
			wantCost:  "",
			err:       nil,
		},
		"dsadsahcashfhfdsh.com": {
			domain:    "dsadsahcashfhfdsh.com",
			wantAvail: "AVAILABLE",
			wantCost:  "12USD",
			err:       nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.DomainIsAvailable(projectID, tc.domain)
			if err != tc.err {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

			if got != nil {

				if !reflect.DeepEqual(tc.wantAvail, got.Availability.String()) {
					t.Fatalf("expected: %v got: %v", tc.wantAvail, got)
				}
				if got.Availability.String() == "AVAILABLE" {
					cost := fmt.Sprintf("%d%s", got.YearlyPrice.Units, got.YearlyPrice.CurrencyCode)
					if !reflect.DeepEqual(tc.wantCost, cost) {
						t.Fatalf("expected: %v got: %v", tc.wantCost, cost)
					}
				}

			}
		})
	}
}

func TestDomainIsVerified(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		domain  string
		project string
		want    bool
		err     error
	}{
		"example.com": {
			domain:  "example.com",
			project: projectID,
			want:    false,
			err:     nil,
		},
		// TODO: fix this broken test
		// "yesornositetester.com": {
		// 	domain:  "yesornositetester.com",
		// 	project: "ds-tester-yesornosite",
		// 	want:    true,
		// 	err:     nil,
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.DomainIsVerified(tc.project, tc.domain)
			if err != tc.err {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v got: %v", tc.want, got)
			}
		})
	}
}

func TestDomainContact(t *testing.T) {
	contact := &domainspb.ContactSettings_Contact{
		PostalAddress: &postaladdress.PostalAddress{
			RegionCode:         "US",
			PostalCode:         "94105",
			AdministrativeArea: "CA",
			Locality:           "San Francisco",
			AddressLines:       []string{"345 Spear Street"},
			Recipients:         []string{"Your Name"},
		},
		Email:       "you@example.com",
		PhoneNumber: "+1 555 555 1234",
	}

	tests := map[string]struct {
		input ContactData
		want  domainspb.ContactSettings
		err   error
	}{
		"simple": {
			input: ContactData{DomainRegistrarContact{
				"you@example.com",
				"+1 555 555 1234",
				PostalAddress{
					"US",
					"94105",
					"CA",
					"San Francisco",
					[]string{"345 Spear Street"},
					[]string{"Your Name"},
				},
			}},
			want: domainspb.ContactSettings{
				Privacy:           domainspb.ContactPrivacy_PRIVATE_CONTACT_DATA,
				RegistrantContact: contact,
				AdminContact:      contact,
				TechnicalContact:  contact,
			},
			err: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.input.DomainContact()

			if err != tc.err {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v got: %+v", tc.want, got)
			}
		})
	}
}
