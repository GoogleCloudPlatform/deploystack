package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func compareValues(label string, want interface{}, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("%s: expected: \n|%v|\ngot: \n|%v|", label, want, got)
	}
}

func TestConfig(t *testing.T) {
	testdata := "../test_files/configs"
	tests := map[string]struct {
		pwd      string
		want     Config
		descPath string
	}{
		"Original": {
			pwd: "original",
			want: Config{
				Title:             "Three Tier App (TODO)",
				Duration:          9,
				DocumentationLink: "https://cloud.google.com/shell/docs/cloud-shell-tutorials/deploystack/three-tier-app",
				Project:           true,
				ProjectNumber:     true,
				Region:            true,
				BillingAccount:    false,
				RegionType:        "run",
				RegionDefault:     "us-central1",
				Zone:              true,
				AuthorSettings:    Settings{Setting{Name: "basename", Value: "three-tier-app", Type: "string"}},
				PathTerraform:     ".",
				PathMessages:      "messages",
				PathScripts:       "scripts",
			},
			descPath: "messages/description.txt",
		},
		"YAML": {
			pwd: "preferredyaml",
			want: Config{
				Title:             "Three Tier App (TODO)",
				Duration:          9,
				DocumentationLink: "https://cloud.google.com/shell/docs/cloud-shell-tutorials/deploystack/three-tier-app",
				Project:           true,
				ProjectNumber:     true,
				Region:            true,
				BillingAccount:    false,
				RegionType:        "run",
				RegionDefault:     "us-central1",
				Zone:              true,
				AuthorSettings:    Settings{Setting{Name: "basename", Value: "three-tier-app", Type: "string"}},
				PathTerraform:     "terraform",
				PathMessages:      ".deploystack/messages",
				PathScripts:       ".deploystack/scripts",
				CustomSettings: []Custom{
					{
						Name:        "nodes",
						Description: "Please enter the number of nodes",
						Default:     "roles/owner|Project Owner",
						Options: []string{
							"roles/reviewer|Project Reviewer",
							"roles/owner|Project Owner",
							"roles/vison.reader|Cloud Vision Reader",
						},
					},
				},
			},
			descPath: ".deploystack/messages/description.txt",
		},
		"withAuthorSettings": {
			pwd: "withauthorsettings",
			want: Config{
				Title:          "Three Tier App (TODO)",
				Duration:       9,
				AuthorSettings: Settings{Setting{Name: "basename", Value: "three-tier-app", Type: "string"}},
				PathTerraform:  "terraform",
				PathMessages:   ".deploystack/messages",
				PathScripts:    ".deploystack/scripts",
			},
			descPath: ".deploystack/messages/description.txt",
		},
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if err := os.Chdir(fmt.Sprintf("%s/%s", testdata, tc.pwd)); err != nil {
				t.Fatalf("failed to set the wd: %v", err)
			}

			s := NewStack()

			if err := s.FindAndReadRequired(); err != nil {
				t.Fatalf("could not read config file: %s", err)
			}

			dat, err := os.ReadFile(tc.descPath)
			if err != nil {
				t.Fatalf("could not read description file: %s", err)
			}
			tc.want.Description = string(dat)

			if !reflect.DeepEqual(tc.want, s.Config) {
				diff := deep.Equal(s.Config, tc.want)
				t.Errorf("compare failed: %v", diff)
				// t.Fatalf("expected: \n%+v, \ngot: \n%+v", tc.want, s.Config)
			}
		})
		if err := os.Chdir(wd); err != nil {
			t.Errorf("failed to reset the wd: %v", err)
		}
	}
}

func TestConfigSetAuthorSettings(t *testing.T) {
	tests := map[string]struct {
		in   Config
		want Settings
	}{
		"Original": {
			in: Config{
				HardSet: map[string]string{"basename": "three-tier-app"},
			},
			want: Settings{
				{Name: "basename", Value: "three-tier-app", Type: "string"},
			},
		},

		"Mix": {
			in: Config{
				HardSet: map[string]string{"basename": "three-tier-app"},
				AuthorSettings: Settings{
					{Name: "nodes", Value: "3", Type: "numeric"},
				},
			},
			want: Settings{
				{Name: "basename", Value: "three-tier-app", Type: "string"},
				{Name: "nodes", Value: "3", Type: "numeric"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			got := tc.in.GetAuthorSettings()

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestComputeNames(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
		err   error
	}{
		"http": {
			"../test_files/computenames_repos/deploystack-single-vm",
			"single-vm",
			nil,
		},
		"ssh": {
			"../test_files/computenames_repos/deploystack-gcs-to-bq-with-least-privileges",
			"gcs-to-bq-with-least-privileges",
			nil,
		},
		"nogit": {
			"../test_files/computenames_repos/folder-no-git",
			"",
			fmt.Errorf("could not open local git directory: repository does not exist"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			oldWD, _ := os.Getwd()
			os.Chdir(tc.input)
			defer os.Chdir(oldWD)

			s := NewStack()
			s.FindAndReadRequired()
			err := s.Config.ComputeName()

			os.Chdir(oldWD)

			if !(tc.err == nil && err == nil) {
				if errors.Is(tc.err, err) {
					t.Fatalf("error expected: %v, got: %v", tc.err, err)
				}
			}

			if !reflect.DeepEqual(tc.want, s.Config.Name) {
				t.Fatalf("expected: %v, got: %v", tc.want, s.Config.Name)
			}
		})
	}
}

func TestReadConfig(t *testing.T) {
	errUnableToRead := errors.New("unable to read config file: ")
	tests := map[string]struct {
		path string
		want Stack
		err  error
	}{
		"error": {
			path: "sadasd",
			want: Stack{},
			err:  errUnableToRead,
		},
		"no_custom": {
			path: "../test_files/dsfolders/no_customs",
			want: Stack{
				Config: Config{
					Title:         "TESTCONFIG",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       true,
					Region:        true,
					RegionType:    "functions",
					RegionDefault: "us-central1",
				},
			},
			err: nil,
		},
		"custom": {
			path: "../test_files/dsfolders/customs",
			want: Stack{
				Config: Config{
					Title:         "TESTCONFIG",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       false,
					Region:        false,
					RegionType:    "",
					RegionDefault: "",
					CustomSettings: []Custom{
						{Name: "nodes", Description: "Nodes", Default: "3"},
					},
				},
			},
			err: nil,
		},
		"custom_options": {
			path: "../test_files/dsfolders/customs_options",
			want: Stack{
				Config: Config{
					Title:         "TESTCONFIG",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       false,
					Region:        false,
					RegionType:    "",
					RegionDefault: "",

					CustomSettings: []Custom{
						{
							Name:        "nodes",
							Description: "Nodes",
							Default:     "3",
							Options:     []string{"1", "2", "3"},
						},
					},
				},
			},
			err: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewStack()
			oldWD, _ := os.Getwd()
			os.Chdir(tc.path)

			err := s.FindAndReadRequired()

			if errors.Is(err, tc.err) {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

			os.Chdir(oldWD)

			compareValues("Title", tc.want.Config.Title, s.Config.Title, t)
			compareValues("Description", tc.want.Config.Description, s.Config.Description, t)
			compareValues("Duration", tc.want.Config.Duration, s.Config.Duration, t)
			compareValues("Project", tc.want.Config.Project, s.Config.Project, t)
			compareValues("Region", tc.want.Config.Region, s.Config.Region, t)
			compareValues("RegionType", tc.want.Config.RegionType, s.Config.RegionType, t)
			compareValues("RegionDefault", tc.want.Config.RegionDefault, s.Config.RegionDefault, t)
			for i, v := range s.Config.CustomSettings {
				compareValues(v.Name, tc.want.Config.CustomSettings[i], v, t)
			}
		})
	}
}

func TestSettingSort(t *testing.T) {
	tests := map[string]struct {
		in         Settings
		want       Settings
		deletekeys []string
	}{
		"basic": {
			in: Settings{
				Setting{Name: "test1", Value: "value1"},
				Setting{Name: "test_project", Value: "project_name"},
				Setting{Name: "another", Value: "thing"},
				Setting{Name: "once", Value: "more"},
			},
			want: Settings{
				Setting{Name: "another", Value: "thing"},
				Setting{Name: "once", Value: "more"},
				Setting{Name: "test1", Value: "value1"},
				Setting{Name: "test_project", Value: "project_name"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.in.Sort()
			if !reflect.DeepEqual(tc.want, tc.in) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, tc.in)
			}
		})
	}
}

func TestSettingsAdd(t *testing.T) {
	tests := map[string]struct {
		in    Settings
		key   string
		value string
		want  *Setting
	}{
		"not set yet": {
			in: Settings{
				Setting{Name: "test1", Value: "value1", Type: "string"},
				Setting{Name: "test_project", Value: "project_name", Type: "string"},
				Setting{Name: "another", Value: "thing", Type: "string"},
			},
			key:   "once",
			value: "with feeling",
			want:  &Setting{Name: "once", Value: "with feeling", Type: "string"},
		},
		"already set": {
			in: Settings{
				Setting{Name: "test1", Value: "value1", Type: "string"},
				Setting{Name: "test_project", Value: "project_name", Type: "string"},
				Setting{Name: "another", Value: "thing", Type: "string"},
				Setting{Name: "once", Value: "more", Type: "string"},
			},
			key:   "once",
			value: "with feeling",
			want:  &Setting{Name: "once", Value: "with feeling", Type: "string"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.in.Add(tc.key, tc.value)

			got := tc.in.Find(tc.key)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestSettingsReplace(t *testing.T) {
	tests := map[string]struct {
		in    Settings
		want  Settings
		key   string
		value string
	}{
		"basic": {
			in: Settings{
				Setting{Name: "test1", Value: "value1"},
				Setting{Name: "test_project", Value: "project_name"},
				Setting{Name: "another", Value: "thing"},
				Setting{Name: "once", Value: "more"},
			},
			key:   "once",
			value: "withFeeling",
			want: Settings{
				Setting{Name: "another", Value: "thing"},
				Setting{Name: "once", Value: "withFeeling"},
				Setting{Name: "test1", Value: "value1"},
				Setting{Name: "test_project", Value: "project_name"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := Setting{Name: tc.key, Value: tc.value}
			tc.in.Replace(s)

			tc.in.Sort()
			tc.want.Sort()

			if !reflect.DeepEqual(tc.want, tc.in) {
				t.Fatalf("expected: \n%+v, \ngot: \n%+v", tc.want, tc.in)
			}
		})
	}
}

func TestBasic(t *testing.T) {
	tests := map[string]struct {
		in   Settings
		q    string
		want Settings
	}{
		"basic": {
			in: Settings{
				Setting{Name: "test1", Value: "value1"},
				Setting{Name: "test_project", Value: "project_name"},
				Setting{Name: "another", Value: "thing"},
				Setting{Name: "once", Value: "more"},
			},
			q: "test",
			want: Settings{
				Setting{Name: "test1", Value: "value1"},
				Setting{Name: "test_project", Value: "project_name"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Search(tc.q)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestCustomsGet(t *testing.T) {
	tests := map[string]struct {
		in   Customs
		key  string
		want Custom
	}{
		"basic": {
			in: Customs{
				Custom{Name: "nodes", Description: "test", Default: "3"},
				Custom{Name: "role", Description: "test", Default: "Viewer"},
			},
			key:  "role",
			want: Custom{Name: "role", Description: "test", Default: "Viewer"},
		},
		"nil": {
			in: Customs{
				Custom{Name: "nodes", Description: "test", Default: "3"},
				Custom{Name: "role", Description: "test", Default: "Viewer"},
			},
			key:  "role_not_here",
			want: Custom{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Get(tc.key)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}
