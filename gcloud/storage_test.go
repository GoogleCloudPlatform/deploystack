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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		"basic": {projectID, projectID + "-testing-object", "../README.md", nil},
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

func TestStorageBadProject(t *testing.T) {
	t.Parallel()
	bad := "notavalidprojectnameanditshouldfaildasdas"
	tests := map[string]struct {
		servicefunc func() error
		err         error
	}{
		"StorageBucketCreate": {
			servicefunc: func() error {
				c := NewClient(context.Background(), "testing")
				return c.StorageBucketCreate(bad, "")
			},
			err: fmt.Errorf("error activating service for polling"),
		},
		"StorageBucketDelete": {
			servicefunc: func() error {
				c := NewClient(context.Background(), "testing")
				return c.StorageBucketDelete(bad, "")
			},
			err: fmt.Errorf("error activating service for polling"),
		},
		"StorageObjectCreate": {
			servicefunc: func() error {
				c := NewClient(context.Background(), "testing")
				_, err := c.StorageObjectCreate(bad, "", "")
				return err
			},
			err: fmt.Errorf("error activating service for polling"),
		},
		"StorageObjectDelete": {
			servicefunc: func() error {
				c := NewClient(context.Background(), "testing")
				return c.StorageObjectDelete(bad, "", "")
			},
			err: fmt.Errorf("error activating service for polling"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.servicefunc()
			assert.ErrorContains(t, err, tc.err.Error())
		})
	}
}

func TestStorageErrors(t *testing.T) {
	tests := map[string]struct {
		servicefunc func() error
		err         error
	}{
		"StorageBucketCreate": {
			servicefunc: func() error {
				c := NewClient(context.Background(), "testing")
				return c.StorageBucketCreate(projectID, "ALLCAPSSHOULDERR")
			},
			err: fmt.Errorf("Invalid bucket name"),
		},
		"StorageBucketDelete": {
			servicefunc: func() error {
				c := NewClient(context.Background(), "testing")
				return c.StorageBucketDelete(projectID, "ALLCAPSSHOULDERR")
			},
			err: fmt.Errorf("Invalid bucket name"),
		},
		"StorageObjectCreate": {
			servicefunc: func() error {
				c := NewClient(context.Background(), "testing")
				_, err := c.StorageObjectCreate(projectID, "ALLCAPSSHOULDERR", "")
				return err
			},
			err: fmt.Errorf("no such file or directory"),
		},
		"StorageObjectDelete": {
			servicefunc: func() error {
				c := NewClient(context.Background(), "testing")
				return c.StorageObjectDelete(projectID, "ALLCAPSSHOULDERR", "")
			},
			err: fmt.Errorf("object doesn't exist"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			err := tc.servicefunc()
			assert.ErrorContains(t, err, tc.err.Error())
		})
	}
}
