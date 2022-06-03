package deploystack

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestGetRunRegions(t *testing.T) {
	_, rescueStdout := blockOutput()
	defer func() { os.Stdout = rescueStdout }()
	rRegions, err := regionsListHelper("test_files/gcloudout/regions_run.txt")
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
			got, err := regionsRun(tc.project)
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
