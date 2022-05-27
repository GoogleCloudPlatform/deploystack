package deploystack

import (
	"reflect"
	"testing"
)

func TestServiceEnable(t *testing.T) {
	tests := map[string]struct {
		service string
		project string
		err     error
		want    bool
		disable bool
	}{
		"vault":   {"vault.googleapis.com", projectID, nil, true, true},
		"compute": {"compute.googleapis.com", projectID, nil, true, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ServiceEnable(tc.project, tc.service)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}

			got, err := ServiceIsEnabled(tc.project, tc.service)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}

			if tc.disable {
				ServiceDisable(tc.project, tc.service)
			}
		})
	}
}

func TestServiceDisable(t *testing.T) {
	tests := map[string]struct {
		service string
		project string
		err     error
		want    bool
	}{
		"vault": {"vault.googleapis.com", projectID, nil, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ServiceDisable(tc.project, tc.service)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}

			got, err := ServiceIsEnabled(tc.project, tc.service)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}
