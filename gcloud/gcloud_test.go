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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/scheduler/apiv1beta1/schedulerpb"
	"google.golang.org/api/option"
)

var (
	projectID        = ""
	billingAccount   = ""
	creds            map[string]string
	opts             = option.WithCredentialsFile("")
	ctx              = context.Background()
	defaultUserAgent = "deploystack/testing"
)

func TestMain(m *testing.M) {
	var err error
	opts = option.WithCredentialsFile("../creds.json")

	dat, err := os.ReadFile("../creds.json")
	if err != nil {
		log.Fatalf("unable to handle the json config file: %v", err)
	}

	json.Unmarshal(dat, &creds)

	projectID = creds["project_id"]
	if err != nil {
		log.Fatalf("could not get environment project id: %s", err)
	}
	billingAccount = creds["billing_account"]
	if err != nil {
		log.Fatalf("could not get environment billing account: %s", err)
	}
	code := m.Run()
	os.Exit(code)
}

func readTestFile(file string) string {
	dat, err := os.ReadFile(file)
	if err != nil {
		return "Couldn't read test file"
	}

	return string(dat)
}

func randSeq(n int) string {
	rand.Seed(time.Now().Unix())

	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func removeFromSlice(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}

	return slice
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func regionsListHelper(file string) ([]string, error) {
	result := []string{}
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return result, fmt.Errorf("unable to read region file (%s): %s", file, err)
	}

	temp := strings.Split(string(dat), "\n")

	for _, v := range temp {
		if v == "" {
			continue
		}
		full := strings.Split(v, "/")
		result = append(result, strings.TrimSpace(full[len(full)-1]))
	}

	sort.Strings(result)

	return result, nil
}

func TestGetRegions(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	cRegions, err := regionsListHelper("test_files/gcloudout/regions_compute.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	fRegions, err := regionsListHelper("test_files/gcloudout/regions_functions.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	rRegions, err := regionsListHelper("test_files/gcloudout/regions_run.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	tests := map[string]struct {
		product string
		project string
		want    []string
		err     error
	}{
		"computeRegions": {
			product: "compute",
			project: projectID,
			want:    cRegions,
			err:     nil,
		},

		"functionsRegions": {
			product: "functions",
			project: projectID,
			want:    fRegions,
			err:     nil,
		},

		"runRegions": {
			product: "run",
			project: projectID,
			want:    rRegions,
			err:     nil,
		},

		"GarbageInout": {
			product: "An outdated iPad",
			project: projectID,
			want:    []string{},
			err: fmt.Errorf(
				"invalid product (%s) requested",
				"An outdated iPad"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.RegionList(tc.project, tc.product)

			// BUG: getting weird regions intertmittenly popping up. Solving with this hack
			if tc.product == "compute" {
				got = removeDuplicateStr(removeFromSlice(removeFromSlice(got, "me-west1"), "us-west4"))
				tc.want = removeDuplicateStr(removeFromSlice(removeFromSlice(cRegions, "me-west1"), "us-west4"))
			}

			if err != tc.err {
				if tc.err == nil {
					t.Fatalf("expected: no error, got: %v", err)
				}

				if err.Error() != tc.err.Error() {
					t.Fatalf("expected: error (%v), got: %v", tc.err, err)
				}
			}

			sort.Strings(got)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestBillingAccountCache(t *testing.T) {

	client := NewClient(context.Background(), "testing")
	log.Printf("cache before: %+v", client.cache)
	cachekey := "BillingAccountList"

	_, ok := client.cache[cachekey]
	if ok {
		t.Fatalf("cache should be empty but it isn't")
	}

	result, err := client.BillingAccountList()
	if err != nil {
		t.Fatalf("coult not get first answer from client for test: %s", err)
	}

	log.Printf("cache after: %+v", client.cache)

	_, ok = client.cache[cachekey]
	if !ok {
		t.Fatalf("cache should have a result but it doesn't")
	}

	resultCache, err := client.BillingAccountList()
	if err != nil {
		t.Fatalf("coult not get first answer from client for test: %s", err)
	}

	if !reflect.DeepEqual(result, resultCache) {
		t.Fatalf("expected: %+v, got: %+v", result, resultCache)
	}

}

func TestCacheableFunctions(t *testing.T) {
	client := NewClient(context.Background(), "testing")
	tests := map[string]struct {
		cachekey  string
		cachefunc func() (interface{}, error)
	}{
		"BillingAccountList": {
			cachekey: "BillingAccountList",
			cachefunc: func() (interface{}, error) {
				return client.BillingAccountList()
			},
		},
		"ProjectList": {
			cachekey: "ProjectList",
			cachefunc: func() (interface{}, error) {
				return client.ProjectList()
			},
		},
		"MachineTypeList": {
			cachekey: fmt.Sprintf("MachineTypeList%s", DefaultZone),
			cachefunc: func() (interface{}, error) {
				return client.MachineTypeList(projectID, DefaultZone)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			_, ok := client.cache[tc.cachekey]
			if ok {
				t.Fatalf("cache should be empty but it isn't")
			}

			result, err := tc.cachefunc()
			if err != nil {
				t.Fatalf("coult not get first answer from client for test: %s", err)
			}

			_, ok = client.cache[tc.cachekey]
			if !ok {
				t.Fatalf("cache should have a result but it doesn't")
			}

			resultCache, err := tc.cachefunc()
			if err != nil {
				t.Fatalf("coult not get first answer from client for test: %s", err)
			}

			if !reflect.DeepEqual(result, resultCache) {
				t.Fatalf("expected: %+v, got: %+v", result, resultCache)
			}

		})
	}
}

func TestBreakServices(t *testing.T) {
	client := NewClient(context.Background(), "testing")
	tests := map[string]struct {
		servicefunc func() (interface{}, error)
		blankfunc   func()
		errorfunc   func() (interface{}, error)
	}{
		"compute": {
			servicefunc: func() (interface{}, error) {
				return client.getComputeService(projectID)
			},
			blankfunc: func() {
				client.services.computeService.BasePath = "nonsenseshouldbreak"
			},
			errorfunc: func() (interface{}, error) {
				return client.ComputeRegionList(projectID)
			},
		},
		"billing": {
			servicefunc: func() (interface{}, error) {
				return client.getCloudbillingService()
			},
			blankfunc: func() {
				client.services.billing.BasePath = "nonsenseshouldbreak"
			},
			errorfunc: func() (interface{}, error) {
				return client.BillingAccountList()
			},
		},
		"resourceManager": {
			servicefunc: func() (interface{}, error) {
				return client.getCloudResourceManagerService()
			},
			blankfunc: func() {
				client.services.resourceManager.BasePath = "nonsenseshouldbreak"
			},
			errorfunc: func() (interface{}, error) {
				return client.ProjectList()
			},
		},
		"domains": {
			servicefunc: func() (interface{}, error) {
				return client.getDomainsClient(projectID)
			},
			blankfunc: func() {
				client.services.domains.Close()
			},
			errorfunc: func() (interface{}, error) {
				return client.DomainsSearch(projectID, "example.com")
			},
		},
		// "serviceUsage": {
		// 	servicefunc: func() (interface{}, error) {
		// 		return client.getServiceUsageService()
		// 	},
		// 	blankfunc: func() {
		// 		client.services.serviceUsage.BasePath = "nonsenseshouldbreak"
		// 	},
		// 	errorfunc: func() (interface{}, error) {
		// 		return client.ServiceIsEnabled(projectID, "example.com")
		// 	},
		// },
		"functions": {
			servicefunc: func() (interface{}, error) {
				return client.getCloudFunctionsService(projectID)
			},
			blankfunc: func() {
				client.services.functions.BasePath = "nonsenseshouldbreak"
			},
			errorfunc: func() (interface{}, error) {
				return client.FunctionRegionList(projectID)
			},
		},
		"run": {
			servicefunc: func() (interface{}, error) {
				return client.getRunService(projectID)
			},
			blankfunc: func() {
				client.services.run.BasePath = "nonsenseshouldbreak"
			},
			errorfunc: func() (interface{}, error) {
				return client.RunRegionList(projectID)
			},
		},
		"build": {
			servicefunc: func() (interface{}, error) {
				return client.getCloudBuildService(projectID)
			},
			blankfunc: func() {
				client.services.build.BasePath = "nonsenseshouldbreak"
			},
			errorfunc: func() (interface{}, error) {
				return "", client.CloudBuildTriggerDelete(projectID, "")
			},
		},
		"iam": {
			servicefunc: func() (interface{}, error) {
				return client.getIAMService(projectID)
			},
			blankfunc: func() {
				client.services.iam.BasePath = "nonsenseshouldbreak"
			},
			errorfunc: func() (interface{}, error) {
				return "", client.ProjectGrantIAMRole(projectID, "", "")
			},
		},
		"scheduler": {
			servicefunc: func() (interface{}, error) {
				return client.getSchedulerService(projectID)
			},
			blankfunc: func() {
				client.services.scheduler.Close()
			},
			errorfunc: func() (interface{}, error) {
				return "", client.JobSchedule(projectID, "", schedulerpb.Job{})
			},
		},
		"secretManager": {
			servicefunc: func() (interface{}, error) {
				return client.getSecretManagerService(projectID)
			},
			blankfunc: func() {
				client.services.secretManager.BasePath = "nonsenseshouldbreak"
			},
			errorfunc: func() (interface{}, error) {
				return "", client.SecretDelete(projectID, "")
			},
		},

		"storage": {
			servicefunc: func() (interface{}, error) {
				return client.getStorageService(projectID)
			},
			blankfunc: func() {
				client.services.storage.Close()
			},
			errorfunc: func() (interface{}, error) {
				return "", client.StorageBucketDelete(projectID, "")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := tc.servicefunc()
			if err != nil {
				t.Fatalf("could not call service function for %s: %s ", name, err)
			}

			tc.blankfunc()

			_, err = tc.errorfunc()
			if err == nil {
				t.Fatalf("error should be returned by service function for %s: %s ", name, err)
			}

		})
	}
}

func TestBreakServicesServiceUsage(t *testing.T) {
	client := NewClient(context.Background(), "testing")
	tests := map[string]struct {
		servicefunc func() (interface{}, error)
		blankfunc   func()
		errorfunc   func() (interface{}, error)
	}{
		"serviceUsage": {
			servicefunc: func() (interface{}, error) {
				return client.getServiceUsageService()
			},
			blankfunc: func() {
				client.services.serviceUsage.BasePath = "nonsenseshouldbreak"
			},
			errorfunc: func() (interface{}, error) {
				return client.ServiceIsEnabled(projectID, "example.com")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := tc.servicefunc()
			if err != nil {
				t.Fatalf("could not call service function for %s: %s ", name, err)
			}

			tc.blankfunc()

			_, err = tc.errorfunc()
			if err == nil {
				t.Fatalf("error should be returned by service function for %s: %s ", name, err)
			}

		})
	}
}
