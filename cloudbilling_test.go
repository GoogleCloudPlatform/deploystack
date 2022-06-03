package deploystack

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"

	"google.golang.org/api/cloudbilling/v1"
)

func TestGetBillingAccounts(t *testing.T) {
	dat, err := os.ReadFile("test_files/gcloudout/billing_accounts.json")
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
			got, err := billingAccounts()

			sort.Slice(got[:], func(i, j int) bool {
				return got[i].DisplayName < got[j].DisplayName
			})

			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				fmt.Printf("%+v\n", got[0])
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestLinkProjectToBillingAccount(t *testing.T) {
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
			err := BillingAccountProjectAttach(tc.project, tc.account)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v", tc.err, err)
			}
		})
	}
}
