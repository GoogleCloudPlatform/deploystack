package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/kylelemons/godebug/diff"
)

var testFilesDir = filepath.Join(os.Getenv("DEPLOYSTACK_PATH"), "test_files")

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

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := fmt.Sprintf("%s/%s", testdata, tc.pwd)
			descPath := filepath.Join(path, tc.descPath)

			s := NewStack()

			if err := s.FindAndReadRequired(path); err != nil {
				t.Fatalf("could not read config file: %s", err)
			}

			dat, err := os.ReadFile(descPath)
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
			"computenames_repos/deploystack-single-vm",
			"single-vm",
			nil,
		},
		"ssh": {
			"computenames_repos/deploystack-gcs-to-bq-with-least-privileges",
			"gcs-to-bq-with-least-privileges",
			nil,
		},
		"nogit": {
			"computenames_repos/folder-no-git",
			"",
			fmt.Errorf("could not open local git directory: repository does not exist"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testdata := filepath.Join(testFilesDir, tc.input)
			s := NewStack()
			s.FindAndReadRequired(testdata)
			err := s.Config.ComputeName(testdata)

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

			err := s.FindAndReadRequired(tc.path)

			if errors.Is(err, tc.err) {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

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

func TestSettingsAddComplete(t *testing.T) {
	tests := map[string]struct {
		in   Settings
		set  Setting
		want *Setting
	}{
		"not set yet": {
			in: Settings{
				Setting{Name: "test1", Value: "value1", Type: "string"},
				Setting{Name: "test_project", Value: "project_name", Type: "string"},
				Setting{Name: "another", Value: "thing", Type: "string"},
			},
			set:  Setting{Name: "once", Value: "with feeling", Type: "string"},
			want: &Setting{Name: "once", Value: "with feeling", Type: "string"},
		},
		"already set": {
			in: Settings{
				Setting{Name: "test1", Value: "value1", Type: "string"},
				Setting{Name: "test_project", Value: "project_name", Type: "string"},
				Setting{Name: "another", Value: "thing", Type: "string"},
				Setting{Name: "once", Value: "more", Type: "string"},
			},
			set:  Setting{Name: "once", Value: "with feeling", Type: "string"},
			want: &Setting{Name: "once", Value: "with feeling", Type: "string"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.in.AddComplete(tc.set)

			got := tc.in.Find(tc.set.Name)

			if !reflect.DeepEqual(tc.want, got) {
				diff := deep.Equal(tc.want, got)
				t.Errorf("compare failed: %v", diff)
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

func TestSettingsSearch(t *testing.T) {
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

func TestNewConfigReport(t *testing.T) {
	wd, err := filepath.Abs("../")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}
	testdata := fmt.Sprintf("%s/test_files/configs/multi", wd)

	tests := map[string]struct {
		in   string
		want Report
		err  error
	}{
		"basic-yaml": {
			in: fmt.Sprintf("%s/minimalyaml/.deploystack/deploystack.yaml", testdata),
			want: Report{
				WD:     fmt.Sprintf("%s/minimalyaml", testdata),
				Path:   fmt.Sprintf("%s/minimalyaml/.deploystack/deploystack.yaml", testdata),
				Config: Config{Title: "Minimal YAML"},
			},
			err: nil,
		},
		"basic-json": {
			in: fmt.Sprintf("%s/minimaljson/.deploystack/deploystack.json", testdata),
			want: Report{
				WD:     fmt.Sprintf("%s/minimaljson", testdata),
				Path:   fmt.Sprintf("%s/minimaljson/.deploystack/deploystack.json", testdata),
				Config: Config{Title: "Minimal JSON"},
			},
			err: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewReport(tc.in)

			if tc.err == nil {
				if err != nil {
					t.Fatalf("expected no error, got %s", err)
				}
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestFindConfigReports(t *testing.T) {

	wd, err := filepath.Abs("../")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}
	testdata := fmt.Sprintf("%s/test_files/configs/multi", wd)

	tests := map[string]struct {
		in   string
		want []Report
		err  error
	}{
		"basic": {
			in: testdata,
			want: []Report{
				{
					WD:     fmt.Sprintf("%s/minimaljson", testdata),
					Path:   fmt.Sprintf("%s/minimaljson/.deploystack/deploystack.json", testdata),
					Config: Config{Title: "Minimal JSON"},
				},
				{
					WD:     fmt.Sprintf("%s/minimalyaml", testdata),
					Path:   fmt.Sprintf("%s/minimalyaml/.deploystack/deploystack.yaml", testdata),
					Config: Config{Title: "Minimal YAML"},
				},
			},
			err: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := FindConfigReports(testdata)

			if tc.err == nil {
				if err != nil {
					t.Fatalf("expected no error, got %s", err)
				}
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestConfigCopy(t *testing.T) {
	tests := map[string]struct {
		in   Config
		want Config
	}{
		"empty": {
			in:   Config{},
			want: Config{},
		},
		"full": {
			in: Config{
				Title:          "TESTCONFIG",
				Description:    "A test string for usage with this stuff.",
				Duration:       5,
				Project:        true,
				ProjectNumber:  true,
				Region:         true,
				BillingAccount: false,
				RegionType:     "run",
				RegionDefault:  "us-central1",
				Zone:           true,
				PathTerraform:  "terraform",
				PathMessages:   ".deploystack/messages",
				PathScripts:    ".deploystack/scripts",
				CustomSettings: []Custom{
					{
						Name:        "nodes",
						Description: "Nodes",
						Default:     "3"},
					{
						Name:        "nodes2",
						Description: "Nodes",
						Default:     "3",
						Options:     []string{"1", "2", "3"},
					},
				},
				AuthorSettings: Settings{
					{Name: "basename", Value: "basename", Type: "string"},
				},
				Products: []Product{
					{Info: "A VM", Product: "Compute Engine"},
				},
			},

			want: Config{
				Title:          "TESTCONFIG",
				Description:    "A test string for usage with this stuff.",
				Duration:       5,
				Project:        true,
				ProjectNumber:  true,
				Region:         true,
				BillingAccount: false,
				RegionType:     "run",
				RegionDefault:  "us-central1",
				Zone:           true,
				PathTerraform:  "terraform",
				PathMessages:   ".deploystack/messages",
				PathScripts:    ".deploystack/scripts",
				CustomSettings: []Custom{
					{
						Name:        "nodes",
						Description: "Nodes",
						Default:     "3"},
					{
						Name:        "nodes2",
						Description: "Nodes",
						Default:     "3",
						Options:     []string{"1", "2", "3"},
					},
				},
				AuthorSettings: Settings{
					{Name: "basename", Value: "basename", Type: "string"},
				},
				Products: []Product{
					{Info: "A VM", Product: "Compute Engine"},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Copy()
			if !reflect.DeepEqual(tc.want, got) {
				wantYAML, err := tc.want.Marshal("yaml")
				if err != nil {
					t.Fatalf("couldn't even marshall the wanted result: %s", err)
				}

				gotYaml, err := got.Marshal("yaml")
				if err != nil {
					t.Fatalf("couldn't even marshall the gotten result: %s", err)
				}

				fmt.Println(diff.Diff(string(wantYAML), string(gotYaml)))
				t.Fatalf("objects didn't match")
			}
		})
	}
}

func TestConfigMarshall(t *testing.T) {
	tests := map[string]struct {
		in     Config
		format string
		want   string
	}{
		"yaml": {
			format: "yaml",
			in: Config{
				Title:          "TESTCONFIG",
				Description:    "A test string for usage with this stuff.",
				Duration:       5,
				Project:        true,
				ProjectNumber:  true,
				Region:         true,
				BillingAccount: false,
				RegionType:     "run",
				RegionDefault:  "us-central1",
				Zone:           true,
				PathTerraform:  "terraform",
				PathMessages:   ".deploystack/messages",
				PathScripts:    ".deploystack/scripts",
				CustomSettings: []Custom{
					{
						Name:        "nodes",
						Description: "Nodes",
						Default:     "3"},
					{
						Name:        "nodes2",
						Description: "Nodes",
						Default:     "3",
						Options:     []string{"1", "2", "3"},
					},
				},
				AuthorSettings: Settings{
					{Name: "basename", Value: "basename", Type: "string"},
				},
				Products: []Product{
					{Info: "A VM", Product: "Compute Engine"},
				},
			},
			want: `title: TESTCONFIG
name: ""
description: A test string for usage with this stuff.
duration: 5
collect_project: true
collect_project_number: true
collect_billing_account: false
register_domain: false
collect_region: true
region_type: run
region_default: us-central1
collect_zone: true
hard_settings: {}
custom_settings:
- name: nodes
  description: Nodes
  default: "3"
  options: []
  prepend_project: false
- name: nodes2
  description: Nodes
  default: "3"
  options:
  - "1"
  - "2"
  - "3"
  prepend_project: false
author_settings:
- name: basename
  value: basename
  type: string
  list: []
  map: {}
configure_gce_instance: false
documentation_link: ""
path_terraform: terraform
path_messages: .deploystack/messages
path_scripts: .deploystack/scripts
projects:
  items: []
  allow_duplicates: false
products:
- info: A VM
  product: Compute Engine
`,
		},
		"json": {
			format: "json",
			in: Config{
				Title:          "TESTCONFIG",
				Description:    "A test string for usage with this stuff.",
				Duration:       5,
				Project:        true,
				ProjectNumber:  true,
				Region:         true,
				BillingAccount: false,
				RegionType:     "run",
				RegionDefault:  "us-central1",
				Zone:           true,
				PathTerraform:  "terraform",
				PathMessages:   ".deploystack/messages",
				PathScripts:    ".deploystack/scripts",
				CustomSettings: []Custom{
					{
						Name:        "nodes",
						Description: "Nodes",
						Default:     "3"},
					{
						Name:        "nodes2",
						Description: "Nodes",
						Default:     "3",
						Options:     []string{"1", "2", "3"},
					},
				},
				AuthorSettings: Settings{
					{Name: "basename", Value: "basename", Type: "string"},
				},
				Products: []Product{
					{Info: "A VM", Product: "Compute Engine"},
				},
			},
			want: `{
	"title": "TESTCONFIG",
	"name": "",
	"description": "A test string for usage with this stuff.",
	"duration": 5,
	"collect_project": true,
	"collect_project_number": true,
	"collect_billing_account": false,
	"register_domain": false,
	"collect_region": true,
	"region_type": "run",
	"region_default": "us-central1",
	"collect_zone": true,
	"hard_settings": null,
	"custom_settings": [
		{
			"name": "nodes",
			"description": "Nodes",
			"default": "3",
			"options": null,
			"prepend_project": false
		},
		{
			"name": "nodes2",
			"description": "Nodes",
			"default": "3",
			"options": [
				"1",
				"2",
				"3"
			],
			"prepend_project": false
		}
	],
	"author_settings": [
		{
			"name": "basename",
			"value": "basename",
			"type": "string",
			"list": null,
			"map": null
		}
	],
	"configure_gce_instance": false,
	"documentation_link": "",
	"path_terraform": "terraform",
	"path_messages": ".deploystack/messages",
	"path_scripts": ".deploystack/scripts",
	"projects": {
		"items": null,
		"allow_duplicates": false
	},
	"products": [
		{
			"info": "A VM",
			"product": "Compute Engine"
		}
	]
}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, _ := tc.in.Marshal(tc.format)
			textdiff := diff.Diff(string(tc.want), string(got))
			if textdiff != "" {
				fmt.Println(textdiff)
				t.Fatalf("object didn't match")
			}
		})
	}
}
