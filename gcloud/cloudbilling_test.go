package gcloud

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"sort"
	"testing"

	"google.golang.org/api/cloudbilling/v1"
)

func TestGetBillingAccounts(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)

	buildtestfile := "test_files/gcloudout/billing_accounts.json"
	localtestfile := "test_files/gcloudout/billing_accounts_local.json"
	testfile := localtestfile

	if _, err := os.Stat(localtestfile); errors.Is(err, os.ErrNotExist) {
		testfile = buildtestfile
	}

	if os.Getenv("BUILD") != "" {
		testfile = "test_files/gcloudout/billing_accounts.json"
	}

	dat, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	bas := []*cloudbilling.BillingAccount{}
	err = json.Unmarshal(dat, &bas)
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	tests := map[string]struct {
		want []*cloudbilling.BillingAccount
	}{
		"NoErrorNoAccounts": {want: bas},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.BillingAccountList()

			sort.Slice(got[:], func(i, j int) bool {
				return got[i].DisplayName < got[j].DisplayName
			})

			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			for i, v := range got {
				if reflect.DeepEqual(tc.want[i].DisplayName, v.DisplayName) {
					break
				}

				if !reflect.DeepEqual(tc.want[i].DisplayName, v.DisplayName) {
					t.Fatalf("expected: %v, got: %v", tc.want[i].DisplayName, v.DisplayName)
				}
			}
		})
	}
}

func TestLinkProjectToBillingAccount(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project string
		account string
		err     error
	}{
		"BadProject":  {project: "stackinaboxstackinabox", account: "0145C0-557C58-C970F3", err: ErrorBillingNoPermission},
		"BaddAccount": {project: projectID, account: "AAAAAA-BBBBBB-CCCCCC", err: ErrorBillingInvalidAccount},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := c.BillingAccountAttach(tc.project, tc.account)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v", tc.err, err)
			}
		})
	}
}
