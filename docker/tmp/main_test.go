package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetRepo(t *testing.T) {
	wd, err := filepath.Abs("../")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}

	testdata := fmt.Sprintf("%s/test_files", wd)
	tests := map[string]struct {
		repo string
		path string
		want string
		err  error
	}{
		"deploystack-nosql-client-server": {
			repo: "deploystack-nosql-client-server",
			path: testdata,
			want: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
		},
		"nosql-client-server": {
			repo: "nosql-client-server",
			path: testdata,
			want: fmt.Sprintf("%s/nosql-client-server", testdata),
		},

		"deploystack-cost-sentry": {
			repo: "deploystack-cost-sentry",
			path: testdata,
			want: fmt.Sprintf("%s/deploystack-cost-sentry_1", testdata),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			err := getRepo(tc.repo, tc.path)

			if tc.err == nil && err != nil {
				t.Fatalf("expected: no error got: %+v", err)
			}

			if _, err := os.Stat(tc.want); os.IsNotExist(err) {
				t.Fatalf("expected: %s to exist it does not", err)
			}

			err = os.RemoveAll(tc.want)
			if err != nil {
				t.Logf(err.Error())
			}

		})
	}
}

func TestGetAcceptableDir(t *testing.T) {
	wd, err := filepath.Abs("../")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}
	testdata := fmt.Sprintf("%s/test_files", wd)

	tests := map[string]struct {
		in   string
		want string
	}{
		"doesnotexist": {
			in:   fmt.Sprintf("%s/testfolder", testdata),
			want: fmt.Sprintf("%s/testfolder", testdata),
		},
		"exists": {
			in:   fmt.Sprintf("%s/alreadyexists", testdata),
			want: fmt.Sprintf("%s/alreadyexists_2", testdata),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := getAcceptableDir(tc.in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestCloneFromRepo(t *testing.T) {

	wd, err := filepath.Abs("../")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}

	testdata := fmt.Sprintf("%s/test_files", wd)
	tests := map[string]struct {
		repo string
		path string
		want string
		err  error
	}{
		"deploystack-nosql-client-server": {
			repo: "deploystack-nosql-client-server",
			path: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
			want: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
		},
		"nosql-client-server": {
			repo: "nosql-client-server",
			path: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
			want: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			err := cloneFromRepo(tc.repo, tc.path)

			if tc.err == nil && err != nil {
				t.Fatalf("expected: no error got: %+v", err)
			}

			if _, err := os.Stat(tc.want); os.IsNotExist(err) {
				t.Fatalf("expected: %s to exist it does not", err)
			}

			err = os.RemoveAll(tc.path)
			if err != nil {
				t.Logf(err.Error())
			}

		})
	}
}
