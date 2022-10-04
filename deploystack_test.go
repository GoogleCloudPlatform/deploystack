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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/kylelemons/godebug/diff"
	"google.golang.org/api/option"
)

var (
	projectID = ""
	creds     map[string]string
)

func TestMain(m *testing.M) {
	var err error
	opts = option.WithCredentialsFile("creds.json")

	dat, err := os.ReadFile("creds.json")
	if err != nil {
		log.Fatalf("unable to handle the json config file: %v", err)
	}

	json.Unmarshal(dat, &creds)

	projectID = creds["project_id"]
	if err != nil {
		log.Fatalf("could not get environment project id: %s", err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestReadConfig(t *testing.T) {
	tests := map[string]struct {
		file string
		desc string
		want Stack
		err  error
	}{
		"error": {
			file: "z.json",
			desc: "z.txt",
			want: Stack{},
			err:  fmt.Errorf("unable to read config file: open z.json: no such file or directory"),
		},
		"no_custom": {
			file: "test_files/no_customs/deploystack.json",
			desc: "test_files/no_customs/deploystack.txt",
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
			file: "test_files/customs/deploystack.json",
			desc: "test_files/customs/deploystack.txt",
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
			file: "test_files/customs_options/deploystack.json",
			desc: "test_files/customs_options/deploystack.txt",
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
			err := s.ReadConfig(tc.file, tc.desc)

			if err != tc.err {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

			compareValues(tc.want.Config.Title, s.Config.Title, t)
			compareValues(tc.want.Config.Description, s.Config.Description, t)
			compareValues(tc.want.Config.Duration, s.Config.Duration, t)
			compareValues(tc.want.Config.Project, s.Config.Project, t)
			compareValues(tc.want.Config.Region, s.Config.Region, t)
			compareValues(tc.want.Config.RegionType, s.Config.RegionType, t)
			compareValues(tc.want.Config.RegionDefault, s.Config.RegionDefault, t)
			for i, v := range s.Config.CustomSettings {
				compareValues(tc.want.Config.CustomSettings[i], v, t)
			}
		})
	}
}

func compareValues(want interface{}, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected: \n|%v|\ngot: \n|%v|", want, got)
	}
}

func TestProcessCustoms(t *testing.T) {
	tests := map[string]struct {
		file string
		desc string
		want string
		err  error
	}{
		"custom_options": {
			file: "test_files/customs_options/deploystack.json",
			desc: "test_files/customs_options/deploystack.txt",
			want: `********************************************************************************[1;36mDeploystack [0m
Deploystack will walk you through setting some options for the  
stack this solutions installs. 
Most questions have a default that you can choose by hitting the Enter key  
********************************************************************************[1;36mPress the Enter Key to continue [0m
********************************************************************************
[1;36mTESTCONFIG[0m
A test string for usage with this stuff.
It's going to take around [0;36m5 minutes[0m
********************************************************************************
[1;36mNodes: [0m
 1) 1 
 2) 2 
[1;36m 3) 3 [0m
Choose number from list, or just [enter] for [1;36m3[0m
> 
[46mProject Details [0m 
Nodes: [1;36m3[0m
`,
			err: nil,
		},
		"custom": {
			file: "test_files/customs/deploystack.json",
			desc: "test_files/customs/deploystack.txt",
			want: `********************************************************************************[1;36mDeploystack [0m
Deploystack will walk you through setting some options for the  
stack this solutions installs. 
Most questions have a default that you can choose by hitting the Enter key  
********************************************************************************[1;36mPress the Enter Key to continue [0m
********************************************************************************
[1;36mTESTCONFIG[0m
A test string for usage with this stuff.
It's going to take around [0;36m5 minutes[0m
********************************************************************************
[1;36mNodes: [0m
Enter value, or just [enter] for [1;36m3[0m
> 
[46mProject Details [0m 
Nodes: [1;36m3[0m
`,
			err: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewStack()
			err := s.ReadConfig(tc.file, tc.desc)

			if err != tc.err {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

			got := captureOutput(func() {
				if err := s.Process("terraform.tfvars"); err != nil {
					log.Fatalf("problemn collecting the configurations: %s", err)
				}
			})

			if !reflect.DeepEqual(tc.want, got) {
				fmt.Println(diff.Diff(got, tc.want))
				t.Fatalf("expected: \n|%v|\ngot: \n|%v|", tc.want, got)
			}
		})
	}
}

func TestStackTFvars(t *testing.T) {
	s := NewStack()
	s.AddSetting("project", "testproject")
	s.AddSetting("boolean", "true")
	s.AddSetting("set", "[item1,item2]")
	got := s.Terraform()

	want := `boolean="true"
project="testproject"
set=["item1","item2"]
`

	if got != want {
		fmt.Println(diff.Diff(want, got))
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}

func TestStackTFvarsWithProjectNAme(t *testing.T) {
	s := NewStack()
	s.AddSetting("project", "testproject")
	s.AddSetting("boolean", "true")
	s.AddSetting("project_name", "dontshow")
	s.AddSetting("set", "[item1,item2]")
	got := s.Terraform()

	want := `boolean="true"
project="testproject"
set=["item1","item2"]
`

	if got != want {
		fmt.Println(diff.Diff(want, got))
		t.Fatalf("expected: %v, got: %v", want, got)
	}
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

func blockOutput() (string, *os.File) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	w.Close()
	out, _ := ioutil.ReadAll(r)
	return string(out), rescueStdout
}

func randSeq(n int) string {
	rand.Seed(time.Now().Unix())

	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func remove(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func TestMassgePhoneNumber(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
		err   error
	}{
		"Good":  {"800 555 1234", "+1.8005551234", nil},
		"Weird": {"d746fd83843", "+1.74683843", nil},
		"BAD":   {"dghdhdfuejfhfhfhrghfhfhdhgreh", "", ErrorCustomNotValidPhoneNumber},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := massagePhoneNumber(tc.input)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v", tc.err, err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestCustomCollect(t *testing.T) {
	_, rescueStdout := blockOutput()
	defer func() { os.Stdout = rescueStdout }()
	tests := map[string]struct {
		input  string
		custom Custom
		want   string
	}{
		"UserEntry":      {"test_input", Custom{Name: "test", Default: "broken_test"}, "test_input"},
		"Default":        {"", Custom{Name: "test", Default: "working_test"}, "working_test"},
		"Phone":          {"215-555-5321", Custom{Name: "test", Default: "215-555-5321", Validation: "phonenumber"}, "+1.2155555321"},
		"PhoneDefault":   {"", Custom{Name: "test", Default: "215-555-5321", Validation: "phonenumber"}, "+1.2155555321"},
		"Integer":        {"30", Custom{Name: "test", Default: "50", Validation: "integer"}, "30"},
		"IntegerDefault": {"", Custom{Name: "test", Default: "50", Validation: "integer"}, "50"},
		"YNYes":          {"yes", Custom{Name: "test", Default: "yes", Validation: "yesorno"}, "yes"},
		"YNYesDefault":   {"", Custom{Name: "test", Default: "yes", Validation: "yesorno"}, "yes"},
		"YNY":            {"y", Custom{Name: "test", Default: "yes", Validation: "yesorno"}, "yes"},
		"YNNo":           {"no", Custom{Name: "test", Default: "yes", Validation: "yesorno"}, "no"},
		"YNn":            {"n", Custom{Name: "test", Default: "yes", Validation: "yesorno"}, "no"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			content := []byte(tc.input)

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("error setting up environment for testing %v", err)
			}

			_, err = w.Write(content)
			if err != nil {
				t.Error(err)
			}
			w.Close()

			stdin := os.Stdin
			// Restore stdin right after the test.
			defer func() { os.Stdin = stdin }()
			os.Stdin = r

			if err := tc.custom.Collect(); err != nil {
				t.Errorf("custom.Collect failed: %v", err)
			}

			if !reflect.DeepEqual(tc.want, tc.custom.Value) {
				t.Fatalf("expected: %v, got: %v", tc.want, tc.custom.Value)
			}
		})
	}
}

func TestFindAndReadRequired(t *testing.T) {
	testdata := "test_files/configs"

	tests := map[string]struct {
		pwd       string
		terraform string
		scripts   string
		messages  string
	}{
		"Original":  {pwd: "original", terraform: ".", scripts: "scripts", messages: "messages"},
		"Perferred": {pwd: "preferred", terraform: "terraform", scripts: ".deploystack/scripts", messages: ".deploystack/messages"},
		"Configed":  {pwd: "configed", terraform: "tf", scripts: "ds/scripts", messages: "ds/messages"},
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

			if !reflect.DeepEqual(tc.terraform, s.Config.PathTerraform) {
				t.Fatalf("expected: %s, got: %s", tc.terraform, s.Config.PathTerraform)
			}

			if !reflect.DeepEqual(tc.scripts, s.Config.PathScripts) {
				t.Fatalf("expected: %s, got: %s", tc.scripts, s.Config.PathScripts)
			}

			if !reflect.DeepEqual(tc.messages, s.Config.PathMessages) {
				t.Fatalf("expected: %s, got: %s", tc.messages, s.Config.PathMessages)
			}
		})
		if err := os.Chdir(wd); err != nil {
			t.Errorf("failed to reset the wd: %v", err)
		}
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
