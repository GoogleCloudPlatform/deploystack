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

package tui

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/deploystack/config"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func TestPreprocessors(t *testing.T) {
	testdata := ""
	tests := map[string]struct {
		f        func(q *Queue) tea.Cmd
		count    int
		label1st string
		value1st string
		settings map[string]string
		cache    map[string]interface{}
	}{
		"getDiskTypes": {
			f:        getDiskTypes,
			count:    3,
			label1st: "Standard",
			value1st: "pd-standard",
		},
		"getYesOrNo": {
			f:        getYesOrNo,
			count:    2,
			label1st: "Yes",
			value1st: "y",
		},
		"getNoOrYes": {
			f:        getNoOrYes,
			count:    2,
			label1st: "No",
			value1st: "n",
		},
		"getProjects": {
			f:        getProjects,
			count:    86,
			label1st: "aiab-test-project",
			value1st: "aiab-test-project",
		},

		"getRegions": {
			f:        getRegions,
			count:    35,
			label1st: "asia-east1",
			value1st: "asia-east1",
		},

		"getZones": {
			f:        getZones,
			count:    3,
			label1st: "asia-east1-a",
			value1st: "asia-east1-a",
			settings: map[string]string{"region": "asia-east1"},
		},

		"getMachineTypeFamilies": {
			f:        getMachineTypeFamilies,
			count:    34,
			label1st: "a2 highgpu",
			value1st: "a2-highgpu",
			settings: map[string]string{"zone": "asia-east1-b"},
		},

		"getMachineTypes": {
			f:        getMachineTypes,
			count:    4,
			label1st: "a2-highgpu-1g",
			value1st: "a2-highgpu-1g",
			settings: map[string]string{
				"zone":                         "asia-east1-b",
				"instance-machine-type-family": "a2-highgpu",
			},
		},

		"getDiskProjects": {
			f:        getDiskProjects,
			count:    14,
			label1st: "CentOS",
			value1st: "centos-cloud",
		},

		"getImageFamilies": {
			f:        getImageFamilies,
			count:    3,
			label1st: "centos-7",
			value1st: "centos-7",
			settings: map[string]string{
				"instance-image-project": "centos-cloud",
			},
		},

		"getImageDisks": {
			f:        getImageDisks,
			count:    1,
			label1st: "centos-7-v20230203  (Latest)",
			value1st: "centos-cloud/centos-7-v20230203",
			settings: map[string]string{
				"instance-image-project": "centos-cloud",
				"instance-image-family":  "centos-7",
			},
		},

		"handleReports": {
			f:        handleReports,
			count:    2,
			settings: map[string]string{},
			label1st: "Minimal JSON",
			value1st: "/minimaljson",
			cache: map[string]interface{}{
				"reports": []config.Report{
					{
						WD:     fmt.Sprintf("%s/minimaljson", testdata),
						Path:   fmt.Sprintf("%s/minimaljson/.deploystack/deploystack.json", testdata),
						Config: config.Config{Title: "Minimal JSON"},
					},
					{
						WD:     fmt.Sprintf("%s/minimalyaml", testdata),
						Path:   fmt.Sprintf("%s/minimalyaml/.deploystack/deploystack.yaml", testdata),
						Config: config.Config{Title: "Minimal YAML"},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")

			if tc.settings != nil {
				for i, v := range tc.settings {
					q.stack.AddSetting(i, v)
				}
			}

			if tc.cache != nil {
				for i, v := range tc.cache {
					q.Save(i, v)
				}
			}

			cmd := tc.f(&q)
			raw := cmd()

			got := raw.([]list.Item)

			if tc.count != len(got) {
				t.Fatalf("count - want '%d' got '%d'", tc.count, len(got))
			}

			i := got[0].(item)

			if tc.label1st != i.label {
				t.Fatalf("label - want '%s' got '%s'", tc.label1st, i.label)
			}

			if tc.value1st != i.value {
				t.Fatalf("value - want '%s' got '%s'", tc.value1st, i.value)
			}
		})
	}
}
