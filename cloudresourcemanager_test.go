package deploystack

import (
	"reflect"
	"sort"
	"strings"
	"testing"
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
			got, err := projectNumber(tc.input)
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
			got, err := projects()

			gotfiltered := []string{}

			for _, v := range got {
				if !strings.Contains(v.Name, "zprojectnamedelete") {
					gotfiltered = append(gotfiltered, v.Name)
				}
			}

			sort.Strings(tc.want)
			sort.Strings(gotfiltered)

			if len(gotfiltered) != len(tc.want) {

				t.Logf("Expected:%s\n", tc.want)
				t.Logf("Got     :%s", gotfiltered)
				t.Fatalf("expected: %v, got: %v", len(tc.want), len(gotfiltered))
			}

			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, gotfiltered) {
				t.Fatalf("expected: %v, got: %v", tc.want, gotfiltered)
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	tests := map[string]struct {
		input string
		err   error
	}{
		"Too long":  {input: "zprojectnamedeletethisprojectnamehastoomanycharacters", err: ErrorProjectCreateTooLong},
		"Bad Chars": {input: "ALLUPERCASEDONESTWORK", err: ErrorProjectInvalidCharacters},
		"Spaces":    {input: "spaces in name", err: ErrorProjectInvalidCharacters},
		// "Duplicate": {input: projectID, err: ErrorProjectAlreadyExists},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			name := tc.input + randSeq(5)
			err := projectCreate(name)
			projectDelete(name)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v project: %s", tc.err, err, name)
			}
		})
	}
}
