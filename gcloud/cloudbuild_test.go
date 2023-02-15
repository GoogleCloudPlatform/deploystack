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
	"testing"

	"google.golang.org/api/cloudbuild/v1"
)

func TestTriggerCreate(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project string
		trigger cloudbuild.BuildTrigger
		err     error
	}{
		"basic": {
			projectID,
			cloudbuild.BuildTrigger{
				Name:      "testtrigger",
				EventType: "MANUAL",
				GitFileSource: &cloudbuild.GitFileSource{
					RepoType: "GITHUB",
					Uri:      "http://github.com/googlecloudplatform/deploystack-single-vm",
					Path:     "test.yaml",
					Revision: fmt.Sprintf("refs/heads/%s", "main"),
				},
				SourceToBuild: &cloudbuild.GitRepoSource{
					Uri:      "http://github.com/googlecloudplatform/deploystack-single-vm",
					Ref:      fmt.Sprintf("refs/heads/%s", "main"),
					RepoType: "GITHUB",
				},
			},
			nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.CloudBuildTriggerCreate(tc.project, tc.trigger)
			if err != tc.err {
				t.Fatalf("expected: %+v, got: %+v", tc.err, err)
			}

			err = c.CloudBuildTriggerDelete(tc.project, got.Id)
			if err != tc.err {
				t.Fatalf("expected: no error got: %+v", err)
			}
		})
	}
}
