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

			extraGots := []string{}
			for _, gotItem := range gotfiltered {
				found := false
				for _, wantItem := range want {
					if wantItem == gotItem {
						found = true
						break
					}
				}

				if !found {
					extraGots = append(extraGots, gotItem)
				}

			}

			// extraWants := []string{}
			// for _, wantItem := range want {
			// 	found := false
			// 	for _, gotItem := range gotfiltered {
			// 		if wantItem == gotItem {
			// 			found = true
			// 			break
			// 		}
			// 	}

			// 	if !found {
			// 		extraWants = append(extraWants, wantItem)
			// 	}

			// }

			if len(extraGots) > 0 {
				// if len(extraGots) > 0 || len(extraWants) > 0 {
				// t.Logf("extra wants: %v ", extraWants)
				t.Logf("extra gots: %v ", extraGots)
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
