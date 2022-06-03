package deploystack

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/diff"
	domainspb "google.golang.org/genproto/googleapis/cloud/domains/v1beta1"
	"google.golang.org/genproto/googleapis/type/postaladdress"
)

func TestDomainRegistrarContactYAML(t *testing.T) {
	tests := map[string]struct {
		file    string
		contact ContactData
		err     error
	}{
		"simple": {
			file: "test_files/contact_sample.yaml",
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
			dat, err := os.ReadFile(tc.file)
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

func TestDomainRegistrarContactReadYAML(t *testing.T) {
	tests := map[string]struct {
		file string
		want ContactData
		err  error
	}{
		"simple": {
			file: "test_files/contact_sample.yaml",
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
			got, err := newContactDataFromFile(tc.file)

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

func TestDomainIsAvailable(t *testing.T) {
	_, rescueStdout := blockOutput()
	defer func() { os.Stdout = rescueStdout }()
	tests := map[string]struct {
		domain    string
		wantAvail string
		wantCost  string
		err       error
	}{
		// TODO: Get this test to work with testing service account.
		// "example.com": {
		// 	domain:    "example.com",
		// 	wantAvail: "UNAVAILABLE",
		// 	wantCost:  "",
		// 	err:       nil,
		// },
		"dsadsahcashfhfdsh.com": {
			domain:    "dsadsahcashfhfdsh.com",
			wantAvail: "AVAILABLE",
			wantCost:  "12USD",
			err:       nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := domainIsAvailable(projectID, tc.domain)
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
		// TODO: Get this test to work with testing service account.
		// "yesornositetester.com": {
		// 	domain:  "yesornositetester.com",
		// 	project: "ds-tester-yesornosite",
		// 	want:    true,
		// 	err:     nil,
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := domainsIsVerified(tc.project, tc.domain)
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

func TestRegistrarContactManage(t *testing.T) {
	tests := map[string]struct {
		input string
		want  ContactData
		err   error
	}{
		"simple": {
			input: "testing.yaml",
			want: ContactData{DomainRegistrarContact{
				"person@example.com",
				"+1.4155551234",
				PostalAddress{
					"US",
					"94502",
					"CA",
					"San Francisco",
					[]string{"345 Spear Street"},
					[]string{"Googler"},
				},
			}},
			err: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, rescueStdout := blockOutput()
			defer func() { os.Stdout = rescueStdout }()
			content := []byte("")

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("error setting up environment for testing %v", err)
			}

			_, err = w.Write(content)
			if err != nil {
				t.Error(err)
			}
			w.Close()

			stdin := os.Stdin
			// Restore stdin right after the test.
			defer func() { os.Stdin = stdin }()
			os.Stdin = r

			got, err := RegistrarContactManage(tc.input)

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
