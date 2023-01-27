package deploystack

import (
	"os"
	"reflect"
	"testing"
)

func TestServiceEnable(t *testing.T) {
	_, rescueStdout := blockOutput()
	defer func() { os.Stdout = rescueStdout }()
	tests := map[string]struct {
		service string
		project string
		err     error
		want    bool
		disable bool
	}{
		"vault":       {"vault.googleapis.com", projectID, nil, true, true},
		"compute":     {"compute.googleapis.com", projectID, nil, true, false},
		"fakeservice": {"fakeservice.googleapis.com", projectID, ErrorServiceNotExistOrNotAllowed, false, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ServiceEnable(tc.project, tc.service)
			if err != tc.err {
				t.Fatalf("expected: %v got: %v", tc.err, err)
			}

			got, err := ServiceIsEnabled(tc.project, tc.service)
			if err != tc.err {
				t.Fatalf("expected: %v got: %v", tc.err, err)
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
		"vault":       {"vault.googleapis.com", projectID, nil, false},
		"fakeservice": {"fakeservice.googleapis.com", projectID, ErrorServiceNotExistOrNotAllowed, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ServiceDisable(tc.project, tc.service)
			if err != tc.err {
				t.Fatalf("expected: %v got: %v", tc.err, err)
			}

			got, err := ServiceIsEnabled(tc.project, tc.service)
			if err != tc.err {
				t.Fatalf("expected: %v got: %v", tc.err, err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}
