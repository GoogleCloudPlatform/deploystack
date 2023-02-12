package gcloud

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func TestGetProjectNumbers(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)

	tests := map[string]struct {
		input string
		want  string
	}{
		"1": {input: creds["project_id"], want: creds["project_number"]},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.ProjectNumberGet(tc.input)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestCheckProject(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)

	tests := map[string]struct {
		input string
		want  bool
	}{
		"Does Exists":     {input: creds["project_id"], want: true},
		"Does Not Exists": {input: "ds-does-not-exst", want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := c.ProjectExists(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestGetProjectParent(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		input string
		want  *cloudresourcemanager.ResourceId
	}{
		"1": {
			input: creds["project_id"],
			want: &cloudresourcemanager.ResourceId{
				Id:   creds["parent"],
				Type: creds["parent_type"],
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.ProjectParentGet(tc.input)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestGetProjects(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		filepath string
	}{
		"1": {filepath: "../test_files/gcloudout/projects.txt"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.ProjectList()

			raw := readTestFile(tc.filepath)
			want := strings.Split(strings.TrimSpace(raw), "\n")

			gotfiltered := []string{}

			for _, v := range got {
				gotfiltered = append(gotfiltered, v.Name)
			}

			sort.Strings(want)
			sort.Strings(gotfiltered)

			failures := 0
			for i, v := range gotfiltered {
				if v != want[i] {
					failures++
					t.Logf("failing pair %s:%s ", v, want[i])
				}
			}

			if failures > 0 {
				t.Logf("Expected:%s\n", want)
				t.Logf("Got     :%s", gotfiltered)
				t.Fatalf("expected: %v got: %v", len(want), len(gotfiltered))
			}

			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		input string
		err   error
	}{
		"Too long": {
			input: "zprojectnamedeletethisprojectnamehastoomanycharacters",
			err:   ErrorProjectCreateTooLong,
		},
		"Bad Chars": {
			input: "ALLUPERCASEDONESTWORK",
			err:   ErrorProjectInvalidCharacters,
		},
		"Spaces": {
			input: "spaces in name",
			err:   ErrorProjectInvalidCharacters,
		},
		// "Duplicate": {input: projectID, err: ErrorProjectAlreadyExists},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			name := tc.input + randSeq(5)
			err := c.ProjectCreate(name, "", "")
			c.ProjectDelete(name)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v project: %s", tc.err, err, name)
			}
		})
	}
}

func TestGetProject(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	expected := projectID

	old, err := c.ProjectIDGet()
	if err != nil {
		t.Fatalf("retrieving old project: expected: no error, got: %v", err)
	}

	if err := c.ProjectIDSet(expected); err != nil {
		t.Fatalf("setting expecgted project: expected: no error, got: %v", err)
	}

	got, err := c.ProjectIDGet()
	if err != nil {
		t.Fatalf("expected: no error, got: %v", err)
	}

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("expected: %v, got: %v", expected, got)
	}

	if err := c.ProjectIDSet(old); err != nil {
		t.Fatalf("resetting old project: expected: no error, got: %v", err)
	}
}
