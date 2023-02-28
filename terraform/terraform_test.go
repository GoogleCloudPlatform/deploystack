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
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/kylelemons/godebug/diff"
)

func TestExtract2(t *testing.T) {
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
					file:  fmt.Sprintf("%s/main.tf", testdata),
					start: 15,
				},
				Block{
					Name: "project_id",
					Text: `variable "project_id" {
  type = string
}`,
					Kind:  "variable",
					Type:  "string",
					file:  fmt.Sprintf("%s/variables.tf", testdata),
					start: 15,
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
					file:  fmt.Sprintf("%s/main.tf", testdata),
					start: 24,
				},
				Block{
					Name:  "project",
					Type:  "google_project",
					Kind:  "data",
					start: 37,
					file:  fmt.Sprintf("%s/main.tf", testdata),
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

				if tc.want[i].file != (*got)[i].file {
					t.Fatalf("file expected: %+v, got: %+v", tc.want[i].file, (*got)[i].file)
				}

				if tc.want[i].start != (*got)[i].start {
					t.Fatalf("start expected: %+v, got: %+v", tc.want[i].start, (*got)[i].start)
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
			file:  "testdata/extracttest/main.tf",
			start: 15,
		},
		Block{
			Name: "project_id",
			Text: `variable "project_id" {
  type = string
}`,
			Kind:  "variable",
			Type:  "string",
			file:  "testdata/extracttest/variables.tf",
			start: 15,
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
			file:  "testdata/extracttest/main.tf",
			start: 24,
		},
		Block{
			Name:  "project",
			Type:  "google_project",
			Kind:  "data",
			start: 37,
			file:  "testdata/extracttest/main.tf",
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

		if want[i].file != (*got)[i].file {
			t.Fatalf("file expected: %+v, got: %+v", want[i].file, (*got)[i].file)
		}

		if want[i].start != (*got)[i].start {
			t.Fatalf("start expected: %+v, got: %+v", want[i].start, (*got)[i].start)
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
			file:  "testdata/variables/variables.tf",
			start: 15,
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

	if want[0].file != (*got)[0].file {
		t.Fatalf("expected: %+v, got: %+v", want[0].file, (*got)[0].file)
	}

	if want[0].start != (*got)[0].start {
		t.Fatalf("expected: %+v, got: %+v", want[0].start, (*got)[0].start)
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
			file:  "testdata/resources/main.tf",
			start: 15,
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

	if want[0].file != (*got)[0].file {
		t.Fatalf("file expected: %+v, got: %+v", want[0].file, (*got)[0].file)
	}

	if want[0].start != (*got)[0].start {
		t.Fatalf("start expected: %+v, got: %+v", want[0].start, (*got)[0].start)
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
			file:  "testdata/modules/main.tf",
			start: 15,
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

	if want[0].file != (*got)[0].file {
		t.Fatalf("expected: %+v, got: %+v", want[0].file, (*got)[0].file)
	}

	if want[0].start != (*got)[0].start {
		t.Fatalf("expected: %+v, got: %+v", want[0].start, (*got)[0].start)
	}
}

func TestFindClosingBracket(t *testing.T) {
	tests := map[string]struct {
		start   int
		content string
		want    int
	}{
		"1": {start: 1, content: "", want: 0},
		"2": {start: 4, content: `

		# Enabling services in your GCP project
		variable "gcp_service_list" {
		  description = "The list of apis necessary for the project"
		  type        = list(string)
		  default = [
			"compute.googleapis.com",
		  ]
		}`, want: 9},
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

func TestNewGCPReources(t *testing.T) {
	wd, err := filepath.Abs("../")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}
	testdata := fmt.Sprintf("%s/terraform/testdata/yaml", wd)

	tests := map[string]struct {
		in   string
		want GCPResources
		err  error
	}{
		"basic": {
			in: fmt.Sprintf("%s/bigquery.yaml", testdata),
			want: GCPResources{
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
		},
		"nofile": {
			in:  fmt.Sprintf("%s/noexist.yaml", testdata),
			err: fmt.Errorf("unable to find or read config file"),
		},
		"badfile": {
			in:  fmt.Sprintf("%s/bad.yaml", testdata),
			err: fmt.Errorf("unable to convert content to GCPResources"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewGCPResources(tc.in)
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
