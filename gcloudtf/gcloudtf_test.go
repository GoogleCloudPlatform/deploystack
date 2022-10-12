package gcloudtf

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/kylelemons/godebug/diff"
)

func TestExtract(t *testing.T) {
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
			start: 1,
		},
		Block{
			Name: "project_id",
			Text: `variable "project_id" {
  type = string
}`,
			Kind:  "variable",
			Type:  "string",
			file:  "testdata/extracttest/variables.tf",
			start: 1,
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
			start: 10,
		},
	}

	for i := 0; i < len(*got); i++ {

		if want[i].Name != (*got)[i].Name {
			t.Fatalf("expected: %+v, got: %+v", want[i].Name, (*got)[i].Name)
		}

		if want[i].Text != strings.TrimSpace((*got)[i].Text) {
			fmt.Println(diff.Diff((*got)[i].Text, strings.TrimSpace((*got)[i].Text)))
			t.Fatalf("expected: \n%+v, got: \n%+v", want[i].Text, strings.TrimSpace((*got)[i].Text))
		}

		if want[i].Kind != (*got)[i].Kind {
			t.Fatalf("expected: %+v, got: %+v", want[i].Kind, (*got)[i].Kind)
		}

		if want[i].Type != (*got)[i].Type {
			t.Fatalf("expected: %+v, got: %+v", want[i].Type, (*got)[i].Type)
		}

		if want[i].file != (*got)[i].file {
			t.Fatalf("expected: %+v, got: %+v", want[i].file, (*got)[i].file)
		}

		if want[i].start != (*got)[i].start {
			t.Fatalf("expected: %+v, got: %+v", want[i].start, (*got)[i].start)
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
			Text: `variable "project_id" {
  type = string
}`,
			Kind:  "variable",
			Type:  "string",
			file:  "testdata/variables/variables.tf",
			start: 1,
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
			file:  "testdata/resources/main.tf",
			start: 1,
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
			file:  "testdata/modules/main.tf",
			start: 1,
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
