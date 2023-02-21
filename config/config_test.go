package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func compareValues(label string, want interface{}, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("%s: expected: \n|%v|\ngot: \n|%v|", label, want, got)
	}
}

func TestConfig(t *testing.T) {
	testdata := "test_files/configs"
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
				HardSet:           map[string]string{"basename": "three-tier-app"},
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
				HardSet:           map[string]string{"basename": "three-tier-app"},
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
				t.Fatalf("expected: %+v, got: %+v", tc.want, s.Config)
			}
		})
		if err := os.Chdir(wd); err != nil {
			t.Errorf("failed to reset the wd: %v", err)
		}
	}
}

func TestComputeNames(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
		err   error
	}{
		"http": {
			"test_files/computenames_repos/deploystack-single-vm",
			"single-vm",
			nil,
		},
		"ssh": {
			"test_files/computenames_repos/deploystack-gcs-to-bq-with-least-privileges",
			"gcs-to-bq-with-least-privileges",
			nil,
		},
		"nogit": {
			"test_files/computenames_repos/folder-no-git",
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
			path: "test_files/no_customs",
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
			path: "test_files/customs",
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
			path: "test_files/customs_options",
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
