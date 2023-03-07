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
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestGetRunRegions(t *testing.T) {
	t.Parallel()
	c := NewClient(ctx, defaultUserAgent)
	f := filepath.Join(testFilesDir, "gcloudout/regions_run.txt")

	rRegions, err := regionsListHelper(f)
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	tests := map[string]struct {
		project string
		want    []string
	}{
		"runRegions": {projectID, rRegions},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.RunRegionList(tc.project)
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
