// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploystack

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/deploystack/config"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	"github.com/GoogleCloudPlatform/deploystack/github"
	"github.com/GoogleCloudPlatform/deploystack/terraform"
	"github.com/go-test/deep"
	"github.com/kylelemons/godebug/diff"
	cp "github.com/otiai10/copy"
)

var testFilesDir = filepath.Join(os.Getenv("DEPLOYSTACK_PATH"), "test_files")

var nosqltestdata = filepath.Join(testFilesDir, "reposformeta")

// Getting this right was such a huge pain in the ass, I'm using it a few times
// AND NEVER TOUCHING IT AGAIN

var nosqlMeta = Meta{
	DeployStack: config.Config{
		Title:         "NOSQL CLIENT SERVER",
		Name:          "nosql-client-server",
		Duration:      5,
		Project:       true,
		ProjectNumber: true,
		Region:        true,
		RegionType:    "compute",
		RegionDefault: "us-central1",
		Zone:          true,
		PathTerraform: "terraform",
		PathMessages:  ".deploystack/messages",
		PathScripts:   ".deploystack/scripts",

		DocumentationLink: "https://cloud.google.com/shell/docs/cloud-shell-tutorials/deploystack/nosql-client-server",
		AuthorSettings: config.Settings{
			{
				Name:  "basename",
				Value: "nosql-client-server",
				Type:  "string",
			},
		},
		Description: `This process will configure and create two Compute Engine instances:
* Server - which will run mongodb
* Client - which will run a custom go application that talks to mongo and then 
           exposes an API where you consumne data from mongo`,
	},
	Terraform: terraform.Blocks{
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
			Text: `
resource "google_project_service" "all" {
  for_each                   = toset(var.gcp_service_list)
  project                    = var.project_number
  service                    = each.key
  disable_dependent_services = false
  disable_on_destroy         = false
}`,
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
			Text: `
data "google_compute_network" "default" {
  project = var.project_id
  name = "default"
}`,
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
			Text: `
resource "google_compute_network" "main" {
  provider                = google-beta
  name                    = "${var.basename}-network"
  auto_create_subnetworks = true
  project                 = var.project_id
}`,
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
			Text: `
resource "google_compute_firewall" "default-allow-http" {
  name    = "deploystack-allow-http"
  project = var.project_number
  network = google_compute_network.main.name

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }

  source_ranges = ["0.0.0.0/0"]

  target_tags = ["http-server"]
}`,
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
			Text: `
resource "google_compute_firewall" "default-allow-internal" {
  name    = "deploystack-allow-internal"
  project = var.project_number
  network = google_compute_network.main.name

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow{
    protocol = "udp"
    ports    = ["0-65535"]
  }

  allow{
    protocol = "icmp"
  }

  source_ranges = ["10.128.0.0/20"]

}`,
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
			Text: `
resource "google_compute_firewall" "default-allow-ssh" {
  name    = "deploystack-allow-ssh"
  project = var.project_number
  network = google_compute_network.main.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]

  target_tags = ["ssh-server"]
}`,
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
			Text: `# Create Instances
resource "google_compute_instance" "server" {
  name         = "server"
  zone         = var.zone
  project      = var.project_id
  machine_type = "e2-standard-2"
  tags         = ["ssh-server", "http-server"]
  allow_stopping_for_update = true


  boot_disk {
    auto_delete = true
    device_name = "server"
    initialize_params {
      image = "family/ubuntu-1804-lts"
      size  = 10
      type  = "pd-standard"
    }
  }

  network_interface {
    network = google_compute_network.main.name
    access_config {
      // Ephemeral public IP
    }
  }

  service_account {
    scopes = ["https://www.googleapis.com/auth/logging.write"]
  }

  metadata_startup_script = <<SCRIPT
    apt-get update
    apt-get install -y mongodb
    service mongodb stop
    sed -i 's/bind_ip = 127.0.0.1/bind_ip = 0.0.0.0/' /etc/mongodb.conf
    iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 27017
    service mongodb start
  SCRIPT
  depends_on              = [google_project_service.all]
}`,
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
			Text: `
resource "google_compute_instance" "client" {
  name         = "client"
  zone         = var.zone
  project      = var.project_id
  machine_type = "e2-standard-2"
  tags         = ["http-server", "https-server", "ssh-server"]
  allow_stopping_for_update = true

  boot_disk {
    auto_delete = true
    device_name = "client"
    initialize_params {
      image = "family/ubuntu-1804-lts"
      size  = 10
      type  = "pd-standard"
    }
  }
  service_account {
    scopes = ["https://www.googleapis.com/auth/logging.write"]
  }

  network_interface {
    network = google_compute_network.main.name

    access_config {
      // Ephemeral public IP
    }
  }

  metadata_startup_script = <<SCRIPT
    add-apt-repository ppa:longsleep/golang-backports -y && \
    apt update -y && \
    apt install golang-go -y
    mkdir /modcache
    mkdir /go
    mkdir /app && cd /app
    curl https://raw.githubusercontent.com/GoogleCloudPlatform/golang-samples/main/compute/quickstart/compute_quickstart_sample.go --output main.go
    go mod init exec
    GOPATH=/go GOMODCACHE=/modcache GOCACHE=/modcache go mod tidy
    GOPATH=/go GOMODCACHE=/modcache GOCACHE=/modcache go get -u 
    sed -i 's/mongoport = "80"/mongoport = "27017"/' /app/main.go
    echo "GOPATH=/go GOMODCACHE=/modcache GOCACHE=/modcache HOST=${google_compute_instance.server.network_interface.0.network_ip} go run main.go"
    GOPATH=/go GOMODCACHE=/modcache GOCACHE=/modcache HOST=${google_compute_instance.server.network_interface.0.network_ip} go run main.go & 
  SCRIPT

  depends_on = [google_project_service.all]
}`,
		},

		{
			Name:  "project_id",
			Kind:  "variable",
			Type:  "string",
			File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
			Start: 17,
			Text: `
variable "project_id" {
  type = string
}`,
		},

		{
			Name:  "project_number",
			Kind:  "variable",
			Type:  "string",
			File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
			Start: 21,
			Text: `
variable "project_number" {
  type = string
}`,
		},

		{
			Name:  "zone",
			Kind:  "variable",
			Type:  "string",
			File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
			Start: 25,
			Text: `
variable "zone" {
  type = string
}`,
		},

		{
			Name:  "region",
			Kind:  "variable",
			Type:  "string",
			File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
			Start: 29,
			Text: `
variable "region" {
  type = string
}`,
		},

		{
			Name:  "basename",
			Kind:  "variable",
			Type:  "string",
			File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
			Start: 33,
			Text: `
variable "basename" {
  type = string
}`,
		},

		{
			Name:  "gcp_service_list",
			Kind:  "variable",
			Type:  "list(string)",
			File:  filepath.Join(nosqltestdata, "deploystack-nosql-client-server", "terraform", "variables.tf"),
			Start: 37,
			Text: `
variable "gcp_service_list" {
  description = "The list of apis necessary for the project"
  type        = list(string)
  default = [
    "compute.googleapis.com",
  ]
}`,
		},
	},
}

func compareValues(label string, want interface{}, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("%s: expected: \n|%v|\ngot: \n|%v|", label, want, got)
	}
}

func TestPrecheck(t *testing.T) {
	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}

	testdata := fmt.Sprintf("%s/test_files/configs", wd)
	tests := map[string]struct {
		wd   string
		want string
	}{
		"single": {
			wd:   fmt.Sprintf("%s/preferred", testdata),
			want: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			oldWD, _ := os.Getwd()
			err := os.Chdir(tc.wd)
			if err != nil {
				t.Fatalf("error changing wd: %s", err)
			}

			out := captureOutput(func() {
				Precheck()
			})

			if !strings.Contains(tc.want, string(out)) {
				t.Fatalf("expected to contain: %+v, got: %+v", tc.want, string(out))
			}

			os.Chdir(oldWD)
		})
	}
}
func TestPrecheckMulti(t *testing.T) {
	// Precheck exits if it is called in testing with mutliple stacks
	// so make throwing the exit the test

	if os.Getenv("BE_CRASHER") == "1" {
		wd, err := filepath.Abs(".")
		if err != nil {
			t.Fatalf("error setting up environment for testing %v", err)
		}

		testdata := fmt.Sprintf("%s/test_files/configs", wd)
		path := fmt.Sprintf("%s/multi", testdata)
		oldWD, _ := os.Getwd()
		if err := os.Chdir(path); err != nil {
			t.Fatalf("error changing wd: %s", err)
		}

		Precheck()
		os.Chdir(oldWD)
		return
	}
	// So this is what I got addvice to do, but it caused panics
	// Setting it to explicilty "go test" fixed it
	// cmd := exec.Command(os.Args[0], "-test.run=TestPrecheckMulti")
	cmd := exec.Command("go", "test", "-timeout", "20s", "-test.run=TestPrecheckMulti")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)

}

func captureOutput(f func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}

func TestCacheContact(t *testing.T) {
	tests := map[string]struct {
		in  gcloud.ContactData
		err error
	}{
		"basic": {
			in: gcloud.ContactData{
				AllContacts: gcloud.DomainRegistrarContact{
					Email: "test@example.com",
					Phone: "+155555551212",
					PostalAddress: gcloud.PostalAddress{
						RegionCode:         "US",
						PostalCode:         "94502",
						AdministrativeArea: "CA",
						Locality:           "San Francisco",
						AddressLines:       []string{"345 Spear Street"},
						Recipients:         []string{"Googler"},
					},
				},
			},
			err: nil,
		},
		"err": {
			in:  gcloud.ContactData{},
			err: fmt.Errorf("stat contact.yaml: no such file or directory"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ContactSave(tc.in)

			if tc.err == nil {
				if _, err := os.Stat(contactfile); errors.Is(err, os.ErrNotExist) {
					t.Fatalf("expected no error,  got: %+v", err)
				}
			} else {
				if _, err := os.Stat(contactfile); err.Error() != tc.err.Error() {
					t.Fatalf("expected %+v, got: %+v", tc.err, err)
				}

			}

			os.Remove(contactfile)

		})
	}
}

func TestCheckForContact(t *testing.T) {
	tests := map[string]struct {
		in   string
		want gcloud.ContactData
	}{
		"basic": {
			in: "test_files/contact/contact.yaml",
			want: gcloud.ContactData{
				AllContacts: gcloud.DomainRegistrarContact{
					Email: "test@example.com",
					Phone: "+155555551212",
					PostalAddress: gcloud.PostalAddress{
						RegionCode:         "US",
						PostalCode:         "94502",
						AdministrativeArea: "CA",
						Locality:           "San Francisco",
						AddressLines:       []string{"345 Spear Street"},
						Recipients:         []string{"Googler"},
					},
				},
			},
		},

		"empty": {
			in:   contactfile,
			want: gcloud.ContactData{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			oldContactFile := contactfile
			contactfile = tc.in

			got := ContactCheck()
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}

			contactfile = oldContactFile
		})
	}
}

func TestInit(t *testing.T) {
	errUnableToRead := errors.New("unable to read config file: ")
	tests := map[string]struct {
		path string
		want config.Stack
		err  error
	}{
		"error": {
			path: "sadasd",
			want: config.Stack{},
			err:  errUnableToRead,
		},
		"no_custom": {
			path: "test_files/dsfolders/no_customs",
			want: config.Stack{
				Config: config.Config{
					Title:         "TESTCONFIG",
					Name:          "test",
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
		"no_name": {
			path: "test_files/dsfolders/no_name",
			want: config.Stack{
				Config: config.Config{
					Title:         "NONAME",
					Name:          "",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       true,
					Region:        true,
					RegionType:    "functions",
					RegionDefault: "us-central1",
				},
			},
			err: fmt.Errorf("could retrieve name of stack: could not open local git directory: repository does not exist \nDeployStack author: fix this by adding a 'name' key and value to the deploystack config"),
		},
		"custom": {
			path: "test_files/dsfolders/customs",
			want: config.Stack{
				Config: config.Config{
					Title:         "TESTCONFIG",
					Name:          "test",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       false,
					Region:        false,
					RegionType:    "",
					RegionDefault: "",
					CustomSettings: []config.Custom{
						{Name: "nodes", Description: "Nodes", Default: "3"},
					},
				},
			},
			err: nil,
		},
		"custom_options": {
			path: "test_files/dsfolders/customs_options",
			want: config.Stack{
				Config: config.Config{
					Title:         "TESTCONFIG",
					Name:          "test",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       false,
					Region:        false,
					RegionType:    "",
					RegionDefault: "",

					CustomSettings: []config.Custom{
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

			s, err := Init(tc.path)

			if tc.err == nil {
				if err != nil {
					t.Fatalf("expected: no error got: %+v", err)
				}
			}

			if errors.Is(err, tc.err) {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

			compareValues("Name", tc.want.Config.Name, s.Config.Name, t)
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

func TestShortName(t *testing.T) {
	tests := map[string]struct {
		in   string
		want string
	}{
		"deploystack-repo":     {in: "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry", want: "cost-sentry"},
		"non-deploystack-repo": {in: "https://github.com/tpryan/microservices-demo", want: "microservices-demo"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := Meta{}
			m.Github.Name = tc.in

			got := m.ShortName()
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestShortNameUnderscore(t *testing.T) {
	tests := map[string]struct {
		in   string
		want string
	}{
		"deploystack-repo":     {in: "https://github.com/GoogleCloudPlatform/deploystack-cost-sentry", want: "cost_sentry"},
		"non-deploystack-repo": {in: "https://github.com/tpryan/microservices-demo", want: "microservices_demo"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := Meta{}
			m.Github.Name = tc.in

			got := m.ShortNameUnderscore()
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestGetRepo(t *testing.T) {
	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}

	testdata := fmt.Sprintf("%s/test_files/repoforgithub", wd)
	tests := map[string]struct {
		repo github.Repo
		path string
		want string
		err  error
	}{
		"deploystack-nosql-client-server": {
			repo: github.Repo{
				Name:   "deploystack-nosql-client-server",
				Owner:  "GoogleCloudPlatform",
				Branch: "main",
			},
			path: testdata,
			want: fmt.Sprintf("%s/deploystack-nosql-client-server", testdata),
		},

		"deploystack-cost-sentry": {
			repo: github.Repo{
				Name:   "deploystack-cost-sentry",
				Owner:  "GoogleCloudPlatform",
				Branch: "main",
			},
			path: testdata,
			want: fmt.Sprintf("%s/deploystack-cost-sentry_1", testdata),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			// Defered to make sure it runs even if things fatal
			defer func() {
				err = os.RemoveAll(tc.want)
				if err != nil {
					t.Logf(err.Error())
				}
			}()

			got, err := DownloadRepo(tc.repo, tc.path)

			if tc.err == nil && err != nil {
				t.Fatalf("expected: no error got: %+v", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}

			if _, err := os.Stat(tc.want); os.IsNotExist(err) {
				t.Fatalf("expected: %s to exist it does not", err)
			}

		})
	}
}

func TestGetAcceptableDir(t *testing.T) {
	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}
	testdata := fmt.Sprintf("%s/test_files/repoforgithub", wd)

	tests := map[string]struct {
		in   string
		want string
	}{
		"doesnotexist": {
			in:   fmt.Sprintf("%s/testfolder", testdata),
			want: fmt.Sprintf("%s/testfolder", testdata),
		},
		"exists": {
			in:   fmt.Sprintf("%s/alreadyexists", testdata),
			want: fmt.Sprintf("%s/alreadyexists_2", testdata),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := UniquePath(tc.in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestNewMeta(t *testing.T) {
	testdata := filepath.Join(testFilesDir, "reposformeta")

	tests := map[string]struct {
		path string
		want Meta
		err  error
	}{
		"deploystack-nosql-client-server": {
			path: filepath.Join(testdata, "deploystack-nosql-client-server"),
			want: nosqlMeta,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewMeta(tc.path)

			got.Terraform.Sort()
			tc.want.Terraform.Sort()

			if tc.err == nil && err != nil {
				t.Fatalf("expected no error, got %s", err)
			}

			if !reflect.DeepEqual(tc.want, got) {

				for i := range tc.want.Terraform {
					fmt.Println(diff.Diff(
						tc.want.Terraform[i].Text,
						got.Terraform[i].Text,
					))
				}

				fmt.Println(diff.Diff(
					tc.want.DeployStack.Description,
					got.DeployStack.Description,
				))

				diff := deep.Equal(tc.want, got)
				t.Errorf("compare failed: %v", diff)
			}
		})
	}
}

func TestSuggest(t *testing.T) {
	tests := map[string]struct {
		in   Meta
		want config.Config
		err  error
	}{
		"nosql-client-server": {
			in: nosqlMeta,
			want: config.Config{
				Title:         "NOSQL CLIENT SERVER",
				Name:          "nosql-client-server",
				Duration:      5,
				Project:       true,
				ProjectNumber: true,
				Region:        true,
				RegionType:    "compute",
				RegionDefault: "us-central1",
				Zone:          true,
				PathTerraform: "terraform",
				PathMessages:  ".deploystack/messages",
				PathScripts:   ".deploystack/scripts",

				DocumentationLink: "https://cloud.google.com/shell/docs/cloud-shell-tutorials/deploystack/nosql-client-server",
				AuthorSettings: config.Settings{
					{
						Name:  "basename",
						Value: "nosql-client-server",
						Type:  "string",
					},
				},
				Description: `This process will configure and create two Compute Engine instances:
* Server - which will run mongodb
* Client - which will run a custom go application that talks to mongo and then 
           exposes an API where you consumne data from mongo`,
				Products: []config.Product{
					{Product: "Compute Engine"},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.in.Suggest()

			if tc.err == nil && err != nil {
				t.Fatalf("expected no error, got %s", err)
			}

			if !reflect.DeepEqual(tc.want, got) {

				fmt.Println(diff.Diff(
					tc.want.Description,
					got.Description,
				))

				diff := deep.Equal(tc.want, got)
				t.Errorf("compare failed: %v", diff)
			}
		})
	}
}

func TestAttemptRepo(t *testing.T) {
	tempName, err := os.MkdirTemp("", "testrepos")
	defer os.RemoveAll(tempName)
	if err != nil {
		t.Fatalf("could not get a temp directory for test: %s", err)
	}

	err = os.Mkdir(filepath.Join(tempName, "deploystack-single-vm"), 0666)
	if err != nil {
		t.Fatalf("could not get a make a directory for test: %s", err)
	}

	tests := map[string]struct {
		name     string
		wd       string
		wantdir  string
		wantrepo github.Repo
		err      error
	}{
		"deploystack-nosql-client-server": {
			name:    "deploystack-nosql-client-server",
			wd:      tempName,
			wantdir: filepath.Join(tempName, "deploystack-nosql-client-server"),
			wantrepo: github.Repo{
				Name:   "deploystack-nosql-client-server",
				Owner:  "GoogleCloudPlatform",
				Branch: "main",
			},
		},
		"deploystack-single-vm": {
			name:    "deploystack-single-vm",
			wd:      tempName,
			wantdir: filepath.Join(tempName, "deploystack-single-vm_1"),
			wantrepo: github.Repo{
				Name:   "deploystack-single-vm",
				Owner:  "GoogleCloudPlatform",
				Branch: "main",
			},
		},
		"cost-sentry": {
			name:    "cost-sentry",
			wd:      tempName,
			wantdir: filepath.Join(tempName, "deploystack-cost-sentry"),
			wantrepo: github.Repo{
				Name:   "deploystack-cost-sentry",
				Owner:  "GoogleCloudPlatform",
				Branch: "main",
			},
		},

		"notvalid": {
			name:    "badreponame",
			wd:      tempName,
			wantdir: filepath.Join(tempName, "deploystack-badreponame"),
			wantrepo: github.Repo{
				Name:   "deploystack-badreponame",
				Owner:  "GoogleCloudPlatform",
				Branch: "main",
			},
			err: fmt.Errorf("cannot clone repo"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotdir, gotrepo, err := AttemptRepo(tc.name, tc.wd)

			if tc.err == nil && err != nil {
				t.Fatalf("expected no error, got: %s", err)
			}

			if tc.err != nil && err != nil {
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("expected error %s, got: %s", tc.err, err)
				}
			}

			if !reflect.DeepEqual(tc.wantdir, gotdir) {
				diff := deep.Equal(tc.wantdir, gotdir)
				t.Errorf("compare failed: %v", diff)
			}

			if !reflect.DeepEqual(tc.wantrepo, gotrepo) {
				diff := deep.Equal(tc.wantrepo, gotrepo)
				t.Errorf("compare failed: %v", diff)
			}
		})
	}
}

func TestWriteConfig(t *testing.T) {
	tempName, err := os.MkdirTemp("", "testreposwc")
	defer os.RemoveAll(tempName)
	if err != nil {
		t.Fatalf("could not get a temp directory for test: %s", err)
	}

	src := filepath.Join(testFilesDir, "reposformeta", "terraform-google-load-balanced-vms")
	dest := filepath.Join(tempName, "terraform-google-load-balanced-vms")

	if err := cp.Copy(src, dest); err != nil {
		t.Fatalf("could create a test directory for test: %s", err)
	}

	gh := github.Repo{
		Name:   "terraform-google-load-balanced-vms",
		Owner:  "GoogleCloudPlatform",
		Branch: "main",
	}

	if err := WriteConfig(dest, gh); err != nil {
		t.Fatalf("writeconfig: failed %s", err)
	}

	target := filepath.Join(dest, ".deploystack", "deploystack.yaml")
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("writeconfig: failed to write file %s", err)
	}
}
