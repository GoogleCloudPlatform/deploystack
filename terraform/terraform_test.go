// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package terraform

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/kylelemons/godebug/diff"
)

var testFilesDir = filepath.Join(os.Getenv("DEPLOYSTACK_PATH"), "test_files")

func TestExtract(t *testing.T) {
	wd, err := filepath.Abs("../")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}
	testdata := fmt.Sprintf("%s/terraform/testdata/extracttest", wd)
	tests := map[string]struct {
		in   string
		want Blocks
		err  error
	}{
		"basic": {
			in: testdata,
			want: Blocks{
				Block{
					Name: "snapshot",
					Text: `resource "google_compute_snapshot" "snapshot" {
  project           = var.project_id
  name              = "${var.basename}-snapshot"
  source_disk       = google_compute_instance.exemplar.boot_disk[0].source
  zone              = var.zone
  storage_locations = ["${var.region}"]
  depends_on        = [time_sleep.startup_completion]
}`,
					Kind:  "managed",
					Type:  "google_compute_snapshot",
					File:  fmt.Sprintf("%s/main.tf", testdata),
					Start: 15,
				},
				Block{
					Name: "project_id",
					Text: `variable "project_id" {
  type = string
}`,
					Kind:  "variable",
					Type:  "string",
					File:  fmt.Sprintf("%s/variables.tf", testdata),
					Start: 15,
				},
				Block{
					Name: "project-services",
					Text: `module "project-services" {
  source                      = "terraform-google-modules/project-factory/google//modules/project_services"
  version                     = "~> 13.0"
  disable_services_on_destroy = false

  project_id  = var.project_id
  enable_apis = var.enable_apis

  activate_apis = [
    "compute.googleapis.com"
  ]
}`,
					Kind:  "module",
					Type:  "terraform-google-modules/project-factory/google//modules/project_services",
					File:  fmt.Sprintf("%s/main.tf", testdata),
					Start: 24,
				},
				Block{
					Name:  "project",
					Type:  "google_project",
					Kind:  "data",
					Start: 37,
					File:  fmt.Sprintf("%s/main.tf", testdata),
					Text: `data "google_project" "project" {
}`,
				},
			},
		},
		"nofolder": {
			in:   testdata + "nofolder",
			want: Blocks{},
			err:  fmt.Errorf("terraform config problem"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := Extract(tc.in)

			if err != nil {
				if tc.err == nil {
					t.Fatalf("expected no error, got: %+v", err)
				}
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("expected %s, got: %s", tc.err, err)
				}
				t.Skip()
			}

			for i := 0; i < len(*got); i++ {

				if tc.want[i].Name != (*got)[i].Name {
					t.Fatalf("name expected: %+v, got: %+v", tc.want[i].Name, (*got)[i].Name)
				}

				if tc.want[i].Text != strings.TrimSpace((*got)[i].Text) {
					fmt.Println(diff.Diff(tc.want[i].Text, strings.TrimSpace((*got)[i].Text)))
					t.Fatalf("text expected: \n%+v, got: \n%+v", tc.want[i].Text, strings.TrimSpace((*got)[i].Text))
				}

				if tc.want[i].Kind != (*got)[i].Kind {
					t.Fatalf("kindexpected: %+v, got: %+v", tc.want[i].Kind, (*got)[i].Kind)
				}

				if tc.want[i].Type != (*got)[i].Type {
					t.Fatalf("type expected: %+v, got: %+v", tc.want[i].Type, (*got)[i].Type)
				}

				if tc.want[i].File != (*got)[i].File {
					t.Fatalf("file expected: %+v, got: %+v", tc.want[i].File, (*got)[i].File)
				}

				if tc.want[i].Start != (*got)[i].Start {
					t.Fatalf("start expected: %+v, got: %+v", tc.want[i].Start, (*got)[i].Start)
				}

			}
		})
	}
}

func TestNewBlocks(t *testing.T) {
	mod, dia := tfconfig.LoadModule("testdata/extracttest")
	if dia.Err() != nil {
		t.Fatalf("coult not initiate testdata: %v", dia.Err())
	}

	got, err := NewBlocks(mod)
	if dia.Err() != nil {
		t.Fatalf("coult not turn testdata into structured data: %v", err)
	}

	want := Blocks{
		Block{
			Name: "snapshot",
			Text: `resource "google_compute_snapshot" "snapshot" {
  project           = var.project_id
  name              = "${var.basename}-snapshot"
  source_disk       = google_compute_instance.exemplar.boot_disk[0].source
  zone              = var.zone
  storage_locations = ["${var.region}"]
  depends_on        = [time_sleep.startup_completion]
}`,
			Kind:  "managed",
			Type:  "google_compute_snapshot",
			File:  "testdata/extracttest/main.tf",
			Start: 15,
		},
		Block{
			Name: "project_id",
			Text: `variable "project_id" {
  type = string
}`,
			Kind:  "variable",
			Type:  "string",
			File:  "testdata/extracttest/variables.tf",
			Start: 15,
		},
		Block{
			Name: "project-services",
			Text: `module "project-services" {
  source                      = "terraform-google-modules/project-factory/google//modules/project_services"
  version                     = "~> 13.0"
  disable_services_on_destroy = false

  project_id  = var.project_id
  enable_apis = var.enable_apis

  activate_apis = [
    "compute.googleapis.com"
  ]
}`,
			Kind:  "module",
			Type:  "terraform-google-modules/project-factory/google//modules/project_services",
			File:  "testdata/extracttest/main.tf",
			Start: 24,
		},
		Block{
			Name:  "project",
			Type:  "google_project",
			Kind:  "data",
			Start: 37,
			File:  "testdata/extracttest/main.tf",
			Text: `data "google_project" "project" {
}`,
		},
		Block{
			Name:  "project",
			Type:  "google_project",
			Kind:  "data",
			Start: 37,
			File:  "testdata/extracttest/main.tf",
			Text: `data "google_project" "project" {
}`,
		},
	}

	for i := 0; i < len(*got); i++ {

		if want[i].Name != (*got)[i].Name {
			t.Fatalf("name expected: %+v, got: %+v", want[i].Name, (*got)[i].Name)
		}

		if want[i].Text != strings.TrimSpace((*got)[i].Text) {
			fmt.Println(diff.Diff((*got)[i].Text, strings.TrimSpace((*got)[i].Text)))
			t.Fatalf("text expected: \n%+v, got: \n%+v", want[i].Text, strings.TrimSpace((*got)[i].Text))
		}

		if want[i].Kind != (*got)[i].Kind {
			t.Fatalf("kindexpected: %+v, got: %+v", want[i].Kind, (*got)[i].Kind)
		}

		if want[i].Type != (*got)[i].Type {
			t.Fatalf("type expected: %+v, got: %+v", want[i].Type, (*got)[i].Type)
		}

		if want[i].File != (*got)[i].File {
			t.Fatalf("file expected: %+v, got: %+v", want[i].File, (*got)[i].File)
		}

		if want[i].Start != (*got)[i].Start {
			t.Fatalf("start expected: %+v, got: %+v", want[i].Start, (*got)[i].Start)
		}

	}
}

func TestVariableExtract(t *testing.T) {
	mod, dia := tfconfig.LoadModule("testdata/variables")
	if dia.Err() != nil {
		t.Fatalf("coult not initiate testdata: %v", dia.Err())
	}

	got, err := NewBlocks(mod)
	if dia.Err() != nil {
		t.Fatalf("coult not turn testdata into structured data: %v", err)
	}

	want := Blocks{
		Block{
			Name: "project_id",
			Text: `
variable "project_id" {
  type = string
}`,
			Kind:  "variable",
			Type:  "string",
			File:  "testdata/variables/variables.tf",
			Start: 15,
		},
	}

	if want[0].Name != (*got)[0].Name {
		t.Fatalf("expected: %+v, got: %+v", want[0].Name, (*got)[0].Name)
	}

	if want[0].Text != (*got)[0].Text {
		t.Fatalf("expected: %+v, got: %+v", want[0].Text, (*got)[0].Text)
	}

	if want[0].Kind != (*got)[0].Kind {
		t.Fatalf("expected: %+v, got: %+v", want[0].Kind, (*got)[0].Kind)
	}

	if want[0].Type != (*got)[0].Type {
		t.Fatalf("expected: %+v, got: %+v", want[0].Type, (*got)[0].Type)
	}

	if want[0].File != (*got)[0].File {
		t.Fatalf("expected: %+v, got: %+v", want[0].File, (*got)[0].File)
	}

	if want[0].Start != (*got)[0].Start {
		t.Fatalf("expected: %+v, got: %+v", want[0].Start, (*got)[0].Start)
	}
}

func TestResourceExtract(t *testing.T) {
	mod, dia := tfconfig.LoadModule("testdata/resources")
	if dia.Err() != nil {
		t.Fatalf("coult not initiate testdata: %v", dia.Err())
	}

	got, err := NewBlocks(mod)
	if dia.Err() != nil {
		t.Fatalf("coult not turn testdata into structured data: %v", err)
	}

	want := Blocks{
		Block{
			Name: "snapshot",
			Text: `
resource "google_compute_snapshot" "snapshot" {
  project           = var.project_id
  name              = "${var.basename}-snapshot"
  source_disk       = google_compute_instance.exemplar.boot_disk[0].source
  zone              = var.zone
  storage_locations = ["${var.region}"]
  depends_on        = [time_sleep.startup_completion]
}`,
			Kind:  "managed",
			Type:  "google_compute_snapshot",
			File:  "testdata/resources/main.tf",
			Start: 15,
		},
	}

	if want[0].Name != (*got)[0].Name {
		t.Fatalf("Name expected: %+v, got: %+v", want[0].Name, (*got)[0].Name)
	}

	if want[0].Text != (*got)[0].Text {
		fmt.Println(diff.Diff(want[0].Text, (*got)[0].Text))
		t.Fatalf("Text expected: %+v, got: %+v", want[0].Text, (*got)[0].Text)
	}

	if want[0].Kind != (*got)[0].Kind {
		t.Fatalf("Kind expected: %+v, got: %+v", want[0].Kind, (*got)[0].Kind)
	}

	if want[0].Type != (*got)[0].Type {
		t.Fatalf("Type expected: %+v, got: %+v", want[0].Type, (*got)[0].Type)
	}

	if want[0].File != (*got)[0].File {
		t.Fatalf("file expected: %+v, got: %+v", want[0].File, (*got)[0].File)
	}

	if want[0].Start != (*got)[0].Start {
		t.Fatalf("start expected: %+v, got: %+v", want[0].Start, (*got)[0].Start)
	}
}

func TestModuleExtract(t *testing.T) {
	mod, dia := tfconfig.LoadModule("testdata/modules")
	if dia.Err() != nil {
		t.Fatalf("coult not initiate testdata: %v", dia.Err())
	}

	got, err := NewBlocks(mod)
	if dia.Err() != nil {
		t.Fatalf("coult not turn testdata into structured data: %v", err)
	}

	want := Blocks{
		Block{
			Name: "project-services",
			Text: `
module "project-services" {
  source                      = "terraform-google-modules/project-factory/google//modules/project_services"
  version                     = "~> 13.0"
  disable_services_on_destroy = false

  project_id  = var.project_id
  enable_apis = var.enable_apis

  activate_apis = [
    "compute.googleapis.com"
  ]
}`,
			Kind:  "module",
			Type:  "terraform-google-modules/project-factory/google//modules/project_services",
			File:  "testdata/modules/main.tf",
			Start: 15,
		},
	}

	if want[0].Name != (*got)[0].Name {
		t.Fatalf("expected: %+v, got: %+v", want[0].Name, (*got)[0].Name)
	}

	if want[0].Text != (*got)[0].Text {
		t.Fatalf("expected: %+v, got: %+v", want[0].Text, (*got)[0].Text)
	}

	if want[0].Kind != (*got)[0].Kind {
		t.Fatalf("expected: %+v, got: %+v", want[0].Kind, (*got)[0].Kind)
	}

	if want[0].Type != (*got)[0].Type {
		t.Fatalf("expected: %+v, got: %+v", want[0].Type, (*got)[0].Type)
	}

	if want[0].File != (*got)[0].File {
		t.Fatalf("expected: %+v, got: %+v", want[0].File, (*got)[0].File)
	}

	if want[0].Start != (*got)[0].Start {
		t.Fatalf("expected: %+v, got: %+v", want[0].Start, (*got)[0].Start)
	}
}

func TestFindClosingBracket(t *testing.T) {
	tests := map[string]struct {
		start   int
		content string
		want    int
	}{
		"none": {start: 1, content: "", want: 0},
		"regular usage": {start: 4, content: `

		# Enabling services in your GCP project
		variable "gcp_service_list" {
		  description = "The list of apis necessary for the project"
		  type        = list(string)
		  default = [
			"compute.googleapis.com",
		  ]
		}`, want: 9},
		"broken": {start: 4, content: `

		# Enabling services in your GCP project
		variable "gcp_service_list" {
		  description = "The list of apis necessary for the project"
		  type        = list(string)
		  default = [
			"compute.googleapis.com",
		  ]
		`, want: 10},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := findClosingBracket(tc.start, strings.Split(tc.content, "\n"))
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestTestConfig(t *testing.T) {
	tests := map[string]struct {
		in      GCPResource
		hasTest bool
		hasTodo bool
	}{
		"both": {
			in: GCPResource{
				Label:   "google_bigquery_dataset",
				Product: "BigQuery",
				APICalls: []string{
					"google.cloud.bigquery.[version].DatasetService.InsertDataset",
				},
				TestConfig: TestConfig{
					TestType:    "bq",
					TestCommand: "bq ls | grep -c",
					LabelField:  "table_id",
					Expected:    "1",
					Todo:        "Double check this set of options for test",
				},
			},
			hasTest: true,
			hasTodo: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			if !reflect.DeepEqual(tc.hasTest, tc.in.TestConfig.HasTest()) {
				t.Fatalf("HasTest: expected: %+v, got: %+v", tc.hasTest, tc.in.TestConfig.HasTest())
			}

			if !reflect.DeepEqual(tc.hasTodo, tc.in.TestConfig.HasTodo()) {
				t.Fatalf("HasTodo: expected: %+v, got: %+v", tc.hasTodo, tc.in.TestConfig.HasTodo())
			}
		})
	}
}

func TestNewRepos(t *testing.T) {
	wd, err := filepath.Abs("../")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}
	testdata := fmt.Sprintf("%s/terraform/testdata/yaml", wd)

	tests := map[string]struct {
		in   string
		want Repos
		err  error
	}{
		"basic": {
			in: fmt.Sprintf("%s/repos.yaml", testdata),
			want: Repos{
				"https://github.com/GoogleCloudPlatform/deploystack-cost-sentry",
				"https://github.com/GoogleCloudPlatform/deploystack-etl-pipeline",
			},
		},
		"nofile": {
			in:  fmt.Sprintf("%s/noexist.yaml", testdata),
			err: fmt.Errorf("unable to find or read config file"),
		},
		"badfile": {
			in:  fmt.Sprintf("%s/bad.yaml", testdata),
			err: fmt.Errorf("unable to convert content to a list of repos"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewRepos(tc.in)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("expected no error, got: %+v", err)
				}
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("expected %s, got: %s", tc.err, err)
				}

			} else {
				if !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("expected: %+v, got: %+v", tc.want, got)
				}
			}

		})
	}
}

func TestListMatches(t *testing.T) {
	tests := map[string]struct {
		list List
		in   string
		want bool
	}{
		"test-false": {
			in:   "test",
			list: List{"compute", "sql", "run", "functions"},
			want: false,
		},
		"cloudfunctions-true": {
			in:   "cloudfunctions",
			list: List{"compute", "sql", "run", "functions"},
			want: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.list.Matches(tc.in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestBlockBools(t *testing.T) {
	tests := map[string]struct {
		in         Block
		IsResource bool
		IsModule   bool
		IsVariable bool
		NoDefault  bool
	}{
		"resource": {
			in:         Block{Name: "test", Kind: "managed", Text: "default "},
			IsResource: true,
		},
		"module": {
			in:       Block{Name: "test", Kind: "module", Text: "default "},
			IsModule: true,
		},
		"variable": {
			in:         Block{Name: "test", Kind: "variable", Text: "default "},
			IsVariable: true,
		},
		"nodefault": {
			in:        Block{Name: "test", Kind: "resource"},
			NoDefault: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			if !reflect.DeepEqual(tc.IsModule, tc.in.IsModule()) {
				t.Fatalf("IsModule expected: %+v, got: %+v", tc.IsModule, tc.in.IsModule())
			}
			if !reflect.DeepEqual(tc.IsVariable, tc.in.IsVariable()) {
				t.Fatalf("IsVariable expected: %+v, got: %+v", tc.IsVariable, tc.in.IsVariable())
			}
			if !reflect.DeepEqual(tc.IsResource, tc.in.IsResource()) {
				t.Fatalf("IsResource expected: %+v, got: %+v", tc.IsResource, tc.in.IsResource())
			}
			if !reflect.DeepEqual(tc.NoDefault, tc.in.NoDefault()) {
				t.Fatalf("NoDefault expected: %+v, got: %+v", tc.NoDefault, tc.in.NoDefault())
			}
		})
	}
}

func TestGCPResourceGetProduct(t *testing.T) {
	tests := map[string]struct {
		resources GCPResources
		in        string
		want      string
	}{
		"find": {
			resources: GCPResources{
				"google_bigquery_dataset": GCPResource{
					Label:   "google_bigquery_dataset",
					Product: "BigQuery",
					APICalls: []string{
						"google.cloud.bigquery.[version].DatasetService.InsertDataset",
					},
					TestConfig: TestConfig{
						TestType:    "bq",
						TestCommand: "bq ls | grep -c",
						LabelField:  "table_id",
						Expected:    "1",
						Todo:        "Double check this set of options for test",
					},
				},
				"google_bigquery_table": GCPResource{
					Label:   "google_bigquery_table",
					Product: "BigQuery",
					APICalls: []string{
						"google.cloud.bigquery.[version].TableService.InsertTable",
						"google.cloud.bigquery.[version].TableService.UpdateTable",
						"google.cloud.bigquery.[version].TableService.PatchTable",
					},
					TestConfig: TestConfig{
						TestType:    "bq",
						TestCommand: "bq ls | grep -c",
						LabelField:  "dataset_id",
						Todo:        "Double check this set of options for test",
					},
				},
			},
			in:   "google_bigquery_table",
			want: "BigQuery",
		},
		"no find": {
			resources: GCPResources{
				"google_bigquery_dataset": GCPResource{
					Label:   "google_bigquery_dataset",
					Product: "BigQuery",
					APICalls: []string{
						"google.cloud.bigquery.[version].DatasetService.InsertDataset",
					},
					TestConfig: TestConfig{
						TestType:    "bq",
						TestCommand: "bq ls | grep -c",
						LabelField:  "table_id",
						Expected:    "1",
						Todo:        "Double check this set of options for test",
					},
				},
				"google_bigquery_table": GCPResource{
					Label:   "google_bigquery_table",
					Product: "BigQuery",
					APICalls: []string{
						"google.cloud.bigquery.[version].TableService.InsertTable",
						"google.cloud.bigquery.[version].TableService.UpdateTable",
						"google.cloud.bigquery.[version].TableService.PatchTable",
					},
					TestConfig: TestConfig{
						TestType:    "bq",
						TestCommand: "bq ls | grep -c",
						LabelField:  "dataset_id",
						Todo:        "Double check this set of options for test",
					},
				},
			},
			in:   "google_compute_instace",
			want: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.resources.GetProduct(tc.in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestGetResourceText(t *testing.T) {
	tests := map[string]struct {
		in     string
		want   string
		err    error
		target string
	}{
		"basic": {
			in:     "testdata/resources",
			target: "google_compute_snapshot.snapshot",
			want: `
resource "google_compute_snapshot" "snapshot" {
  project           = var.project_id
  name              = "${var.basename}-snapshot"
  source_disk       = google_compute_instance.exemplar.boot_disk[0].source
  zone              = var.zone
  storage_locations = ["${var.region}"]
  depends_on        = [time_sleep.startup_completion]
}`,
		},
		"begin at zero": {
			in:     "testdata/resources_begin_at_zero",
			target: "google_compute_snapshot.snapshot",
			want: `resource "google_compute_snapshot" "snapshot" {
  project           = var.project_id
  name              = "${var.basename}-snapshot"
  source_disk       = google_compute_instance.exemplar.boot_disk[0].source
  zone              = var.zone
  storage_locations = ["${var.region}"]
  depends_on        = [time_sleep.startup_completion]
}`,
		},
		"filenotfound": {
			in:     "testdata/resources_not_exist",
			target: "google_compute_snapshot.snapshot",
			want:   ``,
			err:    fmt.Errorf("could not get terraform file"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			mod, _ := tfconfig.LoadModule(tc.in)

			r := mod.ManagedResources[tc.target]

			var got string
			var err error
			if r == nil {
				got, err = getResourceText(tc.in, 0)
			} else {
				got, err = getResourceText(r.Pos.Filename, r.Pos.Line)
			}

			if tc.err == nil && err != nil {
				t.Fatalf("expected:no error, got: %+v", err)
			}

			if tc.err != nil && !strings.Contains(err.Error(), tc.err.Error()) {
				t.Fatalf("expected error: %s, got: %s", tc.err, err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				fmt.Println(diff.Diff(tc.want, got))
				t.Fatalf("text doesn't match")
			}
		})
	}
}

func TestBlocks(t *testing.T) {
	tests := map[string]struct {
		in   interface{}
		want Block
		err  error
	}{
		"resource-good": {
			in: tfconfig.Resource{
				Name: "snapshot",
				Type: "google_compute_snapshot",
				Mode: tfconfig.ManagedResourceMode,
				Pos: tfconfig.SourcePos{
					Filename: "testdata/resources/main.tf",
					Line:     15,
				},
			},
			want: Block{
				Name:  "snapshot",
				Type:  "google_compute_snapshot",
				Kind:  "managed",
				Start: 15,
				File:  "testdata/resources/main.tf",
				Text: `
resource "google_compute_snapshot" "snapshot" {
  project           = var.project_id
  name              = "${var.basename}-snapshot"
  source_disk       = google_compute_instance.exemplar.boot_disk[0].source
  zone              = var.zone
  storage_locations = ["${var.region}"]
  depends_on        = [time_sleep.startup_completion]
}`,
			},
		},
		"resource-bad": {
			in: tfconfig.Resource{
				Name: "snapshot",
				Type: "google_compute_snapshot",
				Mode: tfconfig.ManagedResourceMode,
				Pos: tfconfig.SourcePos{
					Filename: "testdata/resources_notexist/main.tf",
					Line:     15,
				},
			},
			want: Block{
				Name:  "snapshot",
				Type:  "google_compute_snapshot",
				Kind:  "managed",
				Start: 15,
				File:  "testdata/resources_notexist/main.tf"},
			err: fmt.Errorf("could not extract text from Resource"),
		},
		"variable-good": {
			in: tfconfig.Variable{
				Name: "project_id",
				Type: "string",
				Pos: tfconfig.SourcePos{
					Filename: "testdata/variables/variables.tf",
					Line:     15,
				},
			},
			want: Block{
				Name:  "project_id",
				Type:  "string",
				Kind:  "variable",
				Start: 15,
				File:  "testdata/variables/variables.tf",
				Text: `
variable "project_id" {
  type = string
}`,
			},
		},
		"variable-bad": {
			in: tfconfig.Variable{
				Name: "project_id",
				Type: "string",
				Pos: tfconfig.SourcePos{
					Filename: "testdata/variables_not_exist/variables.tf",
					Line:     15,
				},
			},
			want: Block{
				Name:  "project_id",
				Type:  "string",
				Kind:  "variable",
				Start: 15,
				File:  "testdata/variables_not_exist/variables.tf",
			},
			err: fmt.Errorf("could not extract text from Variable"),
		},
		"module-good": {
			in: tfconfig.ModuleCall{
				Name:   "project-services",
				Source: "terraform-google-modules/project-factory/google//modules/project_services",
				Pos: tfconfig.SourcePos{
					Filename: "testdata/modules/main.tf",
					Line:     15,
				},
			},
			want: Block{
				Name:  "project-services",
				Type:  "terraform-google-modules/project-factory/google//modules/project_services",
				Kind:  "module",
				Start: 15,
				File:  "testdata/modules/main.tf",
				Text: `
module "project-services" {
  source                      = "terraform-google-modules/project-factory/google//modules/project_services"
  version                     = "~> 13.0"
  disable_services_on_destroy = false

  project_id  = var.project_id
  enable_apis = var.enable_apis

  activate_apis = [
    "compute.googleapis.com"
  ]
}`,
			},
		},
		"module-bad": {
			in: tfconfig.ModuleCall{
				Name:   "project-services",
				Source: "terraform-google-modules/project-factory/google//modules/project_services",
				Pos: tfconfig.SourcePos{
					Filename: "testdata/modules-not-exists/main.tf",
					Line:     15,
				},
			},
			want: Block{
				Name:  "project-services",
				Type:  "terraform-google-modules/project-factory/google//modules/project_services",
				Kind:  "module",
				Start: 15,
				File:  "testdata/modules-not-exists/main.tf",
			},
			err: fmt.Errorf("could not extract text from Module"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got Block
			var err error

			switch v := tc.in.(type) {
			case tfconfig.Resource:
				got, err = NewResourceBlock(&v)
			case tfconfig.Variable:
				got, err = NewVariableBlock(&v)
			case tfconfig.ModuleCall:
				got, err = NewModuleBlock(&v)
			}

			if tc.err == nil && err != nil {
				t.Fatalf("expected:no error, got: %+v", err)
			}

			if tc.err != nil && !strings.Contains(err.Error(), tc.err.Error()) {
				t.Fatalf("expected error: %s, got: %s", tc.err, err)
			}

			if !reflect.DeepEqual(tc.want.Type, got.Type) {
				t.Fatalf("Type expected: %+v, got: %+v", tc.want.Type, got.Type)
			}

			if !reflect.DeepEqual(tc.want.Kind, got.Kind) {
				t.Fatalf("Kind expected: %+v, got: %+v", tc.want.Kind, got.Kind)
			}

			if !reflect.DeepEqual(tc.want.Start, got.Start) {
				t.Fatalf("Start expected: %+v, got: %+v", tc.want.Start, got.Start)
			}

			if !reflect.DeepEqual(tc.want.File, got.File) {
				t.Fatalf("File expected: %+v, got: %+v", tc.want.File, got.File)
			}

			if !reflect.DeepEqual(tc.want.Text, got.Text) {
				fmt.Println(diff.Diff(tc.want.Text, got.Text))
				t.Fatalf("text doesn't match")
			}

		})
	}
}

func TestBadNewBlocks(t *testing.T) {
	tests := map[string]struct {
		in  tfconfig.Module
		err error
	}{
		"module": {
			in: tfconfig.Module{
				ModuleCalls: map[string]*tfconfig.ModuleCall{
					"t": {
						Pos: tfconfig.SourcePos{
							Filename: "testdata/modules-not-exists/main.tf",
						},
					},
				},
			},
			err: fmt.Errorf("could not parse Module Calls"),
		},
		"resource": {
			in: tfconfig.Module{
				ManagedResources: map[string]*tfconfig.Resource{
					"t": {
						Pos: tfconfig.SourcePos{
							Filename: "testdata/resources_notexist/main.tf",
						},
					},
				},
			},
			err: fmt.Errorf("could not parse ManagedResources"),
		},
		"variable": {
			in: tfconfig.Module{
				Variables: map[string]*tfconfig.Variable{
					"t": {
						Pos: tfconfig.SourcePos{
							Filename: "testdata/resources_notexist/main.tf",
						},
					},
				},
			},
			err: fmt.Errorf("could not parse Variables"),
		},
		"data": {
			in: tfconfig.Module{
				DataResources: map[string]*tfconfig.Resource{
					"t": {
						Pos: tfconfig.SourcePos{
							Filename: "testdata/resources_notexist/main.tf",
						},
					},
				},
			},
			err: fmt.Errorf("could not parse DataResources"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewBlocks(&tc.in)
			if tc.err != nil && !strings.Contains(err.Error(), tc.err.Error()) {
				t.Fatalf("expected error: %s, got: %s", tc.err, err)
			}

		})
	}
}

func TestBlocksSort(t *testing.T) {
	tests := map[string]struct {
		in   Blocks
		want Blocks
	}{
		"basic": {
			in: Blocks{
				{Start: 100, File: "variable.tf"},
				{Start: 100, File: "main.tf"},
				{Start: 56, File: "variable.tf"},
				{Start: 19, File: "main.tf"},
				{Start: 1, File: "variable.tf"},
				{Start: 36, File: "main.tf"},
			},
			want: Blocks{
				{Start: 19, File: "main.tf"},
				{Start: 36, File: "main.tf"},
				{Start: 100, File: "main.tf"},
				{Start: 1, File: "variable.tf"},
				{Start: 56, File: "variable.tf"},
				{Start: 100, File: "variable.tf"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.in.Sort()
			got := tc.in
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tests := map[string]struct {
		in    string
		field string
		want  Blocks
	}{
		"type": {
			in:    "project_service",
			field: "type",
			want: Blocks{
				{
					Name: "all",
					Kind: "managed",
					Type: "google_project_service",
					File: filepath.Join(
						nosqltestdata,
						"deploystack-nosql-client-server",
						"terraform",
						"main.tf",
					),
					Start: 21,
				},
			},
		},
		"name": {
			in:    "allow-http",
			field: "name",
			want: Blocks{
				{
					Name: "default-allow-http",
					Kind: "managed",
					Type: "google_compute_firewall",
					File: filepath.Join(
						nosqltestdata,
						"deploystack-nosql-client-server",
						"terraform",
						"main.tf",
					),
					Start: 41,
				},
			},
		},
		"kind": {
			in:    "variable",
			field: "kind",
			want: Blocks{
				{
					Name:  "project_id",
					Kind:  "variable",
					Type:  "string",
					File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
					Start: 17,
				},

				{
					Name:  "project_number",
					Kind:  "variable",
					Type:  "string",
					File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
					Start: 21,
				},

				{
					Name:  "zone",
					Kind:  "variable",
					Type:  "string",
					File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
					Start: 25,
				},

				{
					Name:  "region",
					Kind:  "variable",
					Type:  "string",
					File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
					Start: 29,
				},

				{
					Name:  "basename",
					Kind:  "variable",
					Type:  "string",
					File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
					Start: 33,
				},

				{
					Name:  "gcp_service_list",
					Kind:  "variable",
					Type:  "list(string)",
					File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
					Start: 37,
				},
			},
		},
		"file": {
			in:    "variables.tf",
			field: "file",
			want: Blocks{
				{
					Name:  "project_id",
					Kind:  "variable",
					Type:  "string",
					File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
					Start: 17,
				},

				{
					Name:  "project_number",
					Kind:  "variable",
					Type:  "string",
					File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
					Start: 21,
				},

				{
					Name:  "zone",
					Kind:  "variable",
					Type:  "string",
					File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
					Start: 25,
				},

				{
					Name:  "region",
					Kind:  "variable",
					Type:  "string",
					File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
					Start: 29,
				},

				{
					Name:  "basename",
					Kind:  "variable",
					Type:  "string",
					File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
					Start: 33,
				},

				{
					Name:  "gcp_service_list",
					Kind:  "variable",
					Type:  "list(string)",
					File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
					Start: 37,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := searchBlocks.Search(tc.in, tc.field)
			if !reflect.DeepEqual(tc.want, got) {
				diff := deep.Equal(tc.want, got)
				t.Errorf("compare failed: %v", diff)
			}
		})
	}
}

var nosqltestdata = filepath.Join(testFilesDir, "reposformeta")
var searchBlocks = Blocks{
	{
		Name: "all",
		Kind: "managed",
		Type: "google_project_service",
		File: filepath.Join(
			nosqltestdata,
			"deploystack-nosql-client-server",
			"terraform",
			"main.tf",
		),
		Start: 21,
	},
	{
		Name: "default",
		Kind: "data",
		Type: "google_compute_network",
		File: filepath.Join(
			nosqltestdata,
			"deploystack-nosql-client-server",
			"terraform",
			"main.tf",
		),
		Start: 29,
	},

	{
		Name: "main",
		Kind: "managed",
		Type: "google_compute_network",
		File: filepath.Join(
			nosqltestdata,
			"deploystack-nosql-client-server",
			"terraform",
			"main.tf",
		),
		Start: 34,
	},
	{
		Name: "default-allow-http",
		Kind: "managed",
		Type: "google_compute_firewall",
		File: filepath.Join(
			nosqltestdata,
			"deploystack-nosql-client-server",
			"terraform",
			"main.tf",
		),
		Start: 41,
	},

	{
		Name: "default-allow-internal",
		Kind: "managed",
		Type: "google_compute_firewall",
		File: filepath.Join(
			nosqltestdata,
			"deploystack-nosql-client-server",
			"terraform",
			"main.tf",
		),
		Start: 56,
	},

	{
		Name: "default-allow-ssh",
		Kind: "managed",
		Type: "google_compute_firewall",
		File: filepath.Join(
			nosqltestdata,
			"deploystack-nosql-client-server",
			"terraform",
			"main.tf",
		),
		Start: 79,
	},

	{
		Name: "server",
		Kind: "managed",
		Type: "google_compute_instance",
		File: filepath.Join(
			nosqltestdata,
			"deploystack-nosql-client-server",
			"terraform",
			"main.tf",
		),
		Start: 95,
	},

	{
		Name: "client",
		Kind: "managed",
		Type: "google_compute_instance",
		File: filepath.Join(
			nosqltestdata,
			"deploystack-nosql-client-server",
			"terraform",
			"main.tf",
		),
		Start: 136,
	},

	{
		Name:  "project_id",
		Kind:  "variable",
		Type:  "string",
		File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
		Start: 17,
	},

	{
		Name:  "project_number",
		Kind:  "variable",
		Type:  "string",
		File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
		Start: 21,
	},

	{
		Name:  "zone",
		Kind:  "variable",
		Type:  "string",
		File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
		Start: 25,
	},

	{
		Name:  "region",
		Kind:  "variable",
		Type:  "string",
		File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
		Start: 29,
	},

	{
		Name:  "basename",
		Kind:  "variable",
		Type:  "string",
		File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
		Start: 33,
	},

	{
		Name:  "gcp_service_list",
		Kind:  "variable",
		Type:  "list(string)",
		File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
		Start: 37,
	},
}
