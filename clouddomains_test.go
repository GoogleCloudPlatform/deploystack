package deploystack

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/diff"
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
