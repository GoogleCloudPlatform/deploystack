package deploystack

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func TestGetProjectNumbers(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"1": {input: creds["project_id"], want: creds["project_number"]},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ProjectNumber(tc.input)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestGetProjectParent(t *testing.T) {
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
			got, err := ProjectParent(tc.input)
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
	tests := map[string]struct {
		want []string
	}{
		"1": {want: []string{
			creds["project_id"],
		}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ListProjects()

			gotfiltered := []string{}

			for _, v := range got {
				if !strings.Contains(v.Name, "zprojectnamedelete") {
					gotfiltered = append(gotfiltered, v.Name)
				}
			}

			sort.Strings(tc.want)
			sort.Strings(gotfiltered)

			pass := false
			for _, v := range gotfiltered {
				if v == tc.want[0] {
					pass = true
				}
			}

			if !pass {
				t.Logf("Expected:%s\n", tc.want)
				t.Logf("Got     :%s", gotfiltered)
				t.Fatalf("expected: %v got: %v", len(tc.want), len(gotfiltered))
			}

			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
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
			err := CreateProject(name, "", "")
			DeleteProject(name)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v project: %s", tc.err, err, name)
			}
		})
	}
}

func TestGetProject(t *testing.T) {
	expected := projectID

	old, err := ProjectID()
	if err != nil {
		t.Fatalf("retrieving old project: expected: no error, got: %v", err)
	}

	if err := ProjectIDSet(expected); err != nil {
		t.Fatalf("setting expecgted project: expected: no error, got: %v", err)
	}

	got, err := ProjectID()
	if err != nil {
		t.Fatalf("expected: no error, got: %v", err)
	}

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("expected: %v, got: %v", expected, got)
	}

	if err := ProjectIDSet(old); err != nil {
		t.Fatalf("resetting old project: expected: no error, got: %v", err)
	}
}
