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
	"reflect"
	"regexp"
	"sort"
	"testing"
	"time"

	"google.golang.org/api/cloudfunctions/v1"
)

func TestGenerateFunctionSignedURL(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project string
		region  string
		want    string
		err     error
	}{
		"basic": {projectID, "us-central1", "", nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			retext := fmt.Sprintf("https://storage.googleapis.com/uploads-[0-9]+.%s.cloudfunctions.appspot.com/[0-9a-fA-F-]+.zip\\?GoogleAccessId=service-[0-9]+@gcf-admin-robot.iam.gserviceaccount.com&Expires=[0-9]+&Signature=[0-9a-zA-Z%%]+", tc.region)
			reSignedURL := regexp.MustCompile(retext)

			got, err := c.FunctionGenerateSignedURL(tc.project, tc.region)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}

			if len(reSignedURL.Find([]byte(got))) <= 0 {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestGetFunctionRegions(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	fRegions, err := regionsListHelper("test_files/gcloudout/regions_functions.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	tests := map[string]struct {
		project string
		want    []string
	}{
		"functionsRegions": {projectID, fRegions},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.FunctionRegionList(tc.project)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}

			sort.Strings(got)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

// TODO: Spend time fixing this.
func TestCloudFunctionCreate(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project  string
		region   string
		function cloudfunctions.CloudFunction
		err      error
	}{
		"basic": {
			projectID,
			"us-central1",
			cloudfunctions.CloudFunction{
				Name:              fmt.Sprintf("projects/%s/locations/%s/functions/testFunctionName", projectID, "us-central1"),
				AvailableMemoryMb: 256,
				DockerRegistry:    "CONTAINER_REGISTRY",
				EntryPoint:        "RecordTest",
				EventTrigger: &cloudfunctions.EventTrigger{
					EventType: "google.pubsub.topic.publish",
					Resource:  fmt.Sprintf("projects/%s/topics/cloud-builds", projectID),
				},
				SourceArchiveUrl: "gs://ds-tester-helper-testing-artifacts/func.zip",
				Runtime:          "go116",
				Timeout:          "60s",
			},
			nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := c.FunctionDeploy(tc.project, tc.region, tc.function)
			if err != tc.err {
				t.Fatalf("expected: %+v, got: %+v", tc.err, err)
			}

			functionDeletable := false
			limit := 36
			count := 0
			for !functionDeletable {
				f, err := c.FunctionGet(tc.project, tc.region, "testFunctionName")
				if err != nil {
					t.Fatalf("polling function: expected: no error got: %+v", err)
				}
				if f.Status != "DEPLOY_IN_PROGRESS" {
					functionDeletable = true
				}

				count++

				if count > limit {
					t.Fatalf("polling function: took too long")
				}
				time.Sleep(5 * time.Second)

			}

			err = c.FunctionDelete(tc.project, tc.region, "testFunctionName")
			if err != tc.err {
				t.Fatalf("deleting function: expected: no error got: %+v", err)
			}
		})
	}
}
