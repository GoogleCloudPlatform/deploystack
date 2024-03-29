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
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"google.golang.org/api/cloudbilling/v1"
)

func TestGetBillingAccounts(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)

	buildtestfile := filepath.Join(testFilesDir, "gcloudout/billing_accounts.json")
	localtestfile := filepath.Join(testFilesDir, "gcloudout/billing_accounts_local.json")
	testfile := localtestfile

	if _, err := os.Stat(localtestfile); errors.Is(err, os.ErrNotExist) {
		testfile = buildtestfile
	}

	if os.Getenv("USESA") != "" {
		testfile = buildtestfile
	}

	if os.Getenv("BUILD") != "" {
		testfile = buildtestfile
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
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project string
		account string
		err     error
	}{
		"BadProject":  {project: "stackinaboxstackinabox", account: billingAccount, err: ErrorBillingNoPermission},
		"BaddAccount": {project: projectID, account: "AAAAAA-BBBBBB-CCCCCC", err: ErrorBillingInvalidAccount},
		// TODO: get this working properly again
		// "ShouldWork":  {project: "ds-deleteme-exp2", account: billingAccount, err: nil},
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
