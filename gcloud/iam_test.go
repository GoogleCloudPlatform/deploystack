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
	"testing"
	"time"
)

func TestServiceAccountCreate(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project  string
		username string
		err      error
	}{
		"basic": {projectID, "testSA-" + randSeq(5), nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			email, err := c.ServiceAccountCreate(tc.project, tc.username, tc.username)
			if err != tc.err {
				t.Fatalf("create: expected: %+v, got: %+v", tc.err, err)
			}
			// Immediately deleting caused intermittent failures
			time.Sleep(time.Second * 2)

			err = c.ServiceAccountDelete(tc.project, email)
			if err != tc.err {
				t.Logf("delete: trying to delete: %s", email)
				t.Fatalf("delete: expected: no error got: %+v", err)
			}
		})
	}
}
