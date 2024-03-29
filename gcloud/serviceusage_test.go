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
	"errors"
	"reflect"
	"testing"
)

const FAKESERVICE Service = 1000004

func TestServiceEnable(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		service Service
		project string
		err     error
		want    bool
		disable bool
	}{
		"vault":        {Vault, projectID, nil, true, true},
		"compute":      {Compute, projectID, nil, true, false},
		"emptyproject": {Compute, "", ErrorProjectRequired, false, false},
		"fakeservice":  {FAKESERVICE, projectID, ErrorServiceNotExistOrNotAllowed, false, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			err := c.ServiceEnable(tc.project, tc.service)

			if tc.err == nil && err != nil {
				t.Fatalf("expected: no error got: %v", err)
			}

			if tc.err != nil && err != nil {
				if errors.Is(tc.err, err) {
					t.Fatalf("expected: %v got: %v", tc.err, err)
				}
			}

			got, err := c.ServiceIsEnabled(tc.project, tc.service)
			if err != tc.err {
				t.Fatalf("expected: %v got: %v", tc.err, err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}

			if tc.disable {
				c.ServiceDisable(tc.project, tc.service)
			}
		})
	}
}

func TestServiceDisable(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		service Service
		project string
		err     error
		want    bool
	}{
		"vault":       {Vault, projectID, nil, false},
		"fakeservice": {FAKESERVICE, projectID, ErrorServiceNotExistOrNotAllowed, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := c.ServiceDisable(tc.project, tc.service)
			if err != tc.err {
				t.Fatalf("expected: %v got: %v", tc.err, err)
			}

			got, err := c.ServiceIsEnabled(tc.project, tc.service)
			if err != tc.err {
				t.Fatalf("expected: %v got: %v", tc.err, err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}
