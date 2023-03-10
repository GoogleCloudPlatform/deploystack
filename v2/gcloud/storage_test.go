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

func TestBucketCreate(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project string
		bucket  string
		err     error
	}{
		"basic": {projectID, projectID + "-testing-bucket", nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := c.StorageBucketCreate(tc.project, tc.bucket)
			if err != tc.err {
				t.Fatalf("create: expected: %+v, got: %+v", tc.err, err)
			}

			err = c.StorageBucketDelete(tc.project, tc.bucket)
			if err != tc.err {
				t.Fatalf("delete: expected: no error got: %+v", err)
			}
		})
	}
}

func TestObjectCreate(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		project string
		bucket  string
		path    string
		err     error
	}{
		"basic": {
			project: projectID,
			bucket:  projectID + "-testing-object",
			path:    "../../README.md",
			err:     nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := c.StorageBucketCreate(tc.project, tc.bucket)
			if err != tc.err {
				t.Fatalf("create bucket: expected: %+v, got: %+v", tc.err, err)
			}

			gspath, err := c.StorageObjectCreate(tc.project, tc.bucket, tc.path)
			if err != tc.err {
				t.Fatalf("create object : expected: %+v, got: %+v", tc.err, err)
			}

			err = c.StorageObjectDelete(tc.project, tc.bucket, gspath)
			if err != tc.err {
				t.Fatalf("delete object : expected: %+v, got: %+v", tc.err, err)
			}

			err = c.StorageBucketDelete(tc.project, tc.bucket)
			if err != tc.err {
				t.Fatalf("delete bucket: expected: no error got: %+v", err)
			}
		})
	}
}
