package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func TestPreprocessors(t *testing.T) {
	tests := map[string]struct {
		f        func(q *Queue) tea.Cmd
		count    int
		label1st string
		value1st string
		settings map[string]string
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
			label1st: "asia-east1-b",
			value1st: "asia-east1-b",
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")

			if tc.settings != nil {
				for i, v := range tc.settings {
					q.stack.AddSetting(i, v)
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
