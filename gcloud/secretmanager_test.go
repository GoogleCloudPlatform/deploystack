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

import "testing"

func TestSecretCreate(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project string
		name    string
		payload string
		err     error
	}{
		"basic": {projectID, "testsecret", "secretshhhhhhhhhh", nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := c.SecretCreate(tc.project, tc.name, tc.payload)
			if err != tc.err {
				t.Fatalf("expected: %+v, got: %+v", tc.err, err)
			}

			err = c.SecretDelete(tc.project, tc.name)
			if err != tc.err {
				t.Fatalf("expected: no error got: %+v", err)
			}
		})
	}
}
