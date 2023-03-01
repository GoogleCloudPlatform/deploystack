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

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

func TestFindAndReadConfig(t *testing.T) {
	wd, err := filepath.Abs("../")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}
	testdata := fmt.Sprintf("%s/test_files/configs", wd)

	tests := map[string]struct {
		pwd string
		err error
	}{
		"Original": {
			pwd: "original",
		},
		"Perferred": {
			pwd: "preferred",
		},
		"PerferredYAML": {
			pwd: "preferredyaml",
		},
		"Configed": {
			pwd: "configed",
		},
		"Error": {
			pwd: "error",
			err: ErrConfigNotExist,
		},
		"ErrorNoPAth": {
			pwd: "errorNotexists",
			err: ErrConfigNotExist,
		},
		"ErrorBadFile": {
			pwd: "errorbadfile",
			err: errors.New("unable to parse config file: unable to convert content to Config: yaml: unmarshal errors:\n  line 15: cannot unmarshal !!str `Look at...` into config.Config"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if err := os.Chdir(fmt.Sprintf("%s/%s", testdata, tc.pwd)); err != nil {
				if tc.err == nil {
					t.Fatalf("failed to set the wd: %v", err)
				}
				t.SkipNow()
			}

			s := NewStack()

			if _, err := s.findAndReadConfig(); err != nil {
				if tc.err == nil {
					t.Fatalf("could not read config file: %s", err)
				}
				if err.Error() != tc.err.Error() {
					t.Fatalf("expected: \n'%s'\n, got: \n'%s'\n", tc.err, err)
				}
			}

		})
		if err := os.Chdir(wd); err != nil {
			t.Errorf("failed to reset the wd: %v", err)
		}
	}
}

func TestFindTFFolder(t *testing.T) {
	wd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}
	testdata := fmt.Sprintf("%s/test_files/terraform", wd)
	tests := map[string]struct {
		in   string
		want string
		err  error
	}{
		"toplevel": {
			in:   "toplevel",
			want: ".",
		},
		"secondlevel": {
			in:   "secondlevel",
			want: "terraform",
		},
		"notterraform": {
			in:   "notterraform",
			want: "other",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			oldWD, err := os.Getwd()
			if err != nil {
				t.Fatalf("expected no error, got: %+v", err)
			}

			if err := os.Chdir(fmt.Sprintf("%s/%s", testdata, tc.in)); err != nil {
				t.Fatalf("expected no error, got: %+v", err)
			}
			defer os.Chdir(oldWD)

			stack := NewStack()

			got, err := stack.findTFFolder(stack.Config)

			if tc.err == nil && err != nil {
				t.Fatalf("expected no error, got: %+v", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
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
		"Original": {
			pwd:       "original",
			terraform: ".",
			scripts:   "scripts",
			messages:  "messages"},

		"Perferred": {
			pwd:       "preferred",
			terraform: "terraform",
			scripts:   ".deploystack/scripts",
			messages:  ".deploystack/messages"},
		"PerferredYAML": {
			pwd:       "preferredyaml",
			terraform: "terraform",
			scripts:   ".deploystack/scripts",
			messages:  ".deploystack/messages"},

		"Configed": {
			pwd:       "configed",
			terraform: "tf",
			scripts:   "ds/scripts",
			messages:  "ds/messages"},
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

func TestStackTFvars(t *testing.T) {
	tests := map[string]struct {
		in   Settings
		want string
	}{
		"basic": {
			in: Settings{
				Setting{Name: "project", Value: "testproject", Type: "string"},
				Setting{Name: "boolean", Value: "true", Type: "string"},
				Setting{Name: "set", Value: "[item1,item2]", Type: "string"},
			},
			want: `boolean="true"
project="testproject"
set=["item1","item2"]
`,
		},
		"with basic types": {
			in: Settings{
				Setting{Name: "project", Value: "testproject", Type: "string"},
				Setting{Name: "boolean", Value: "true", Type: "boolean"},
				Setting{Name: "number", Value: "3", Type: "number"},
				Setting{Name: "set", Value: "[item1,item2]", Type: "string"},
			},
			want: `boolean=true
number=3
project="testproject"
set=["item1","item2"]
`,
		},
		"with complext types": {
			in: Settings{
				Setting{Name: "project", Value: "testproject", Type: "string"},
				Setting{Name: "boolean", Value: "true", Type: "boolean"},
				Setting{Name: "number", Value: "3", Type: "number"},
				Setting{Name: "set", List: []string{"item1", "item2"}, Type: "list"},
				Setting{Name: "object", Map: map[string]string{"nickname": "item2", "email": "item2@example.com"}, Type: "map"},
			},
			want: `boolean=true
number=3
object={email="item2@example.com",nickname="item2"}
project="testproject"
set=["item1","item2"]
`,
		},
		"ingnore fields": {
			in: Settings{
				Setting{Name: "project", Value: "testproject", Type: "string"},
				Setting{Name: "boolean", Value: "true", Type: "boolean"},
				Setting{Name: "project_name", Value: "dontshow", Type: "string"},
				Setting{Name: "stack_name", Value: "dontshow", Type: "string"},
				Setting{Name: "", Value: "empty", Type: "string"},
				Setting{Name: "empty", Value: "", Type: "string"},
				Setting{Name: "set", List: []string{"item1", "item2"}, Type: "list"},
				Setting{Name: "object", Map: map[string]string{"nickname": "item2", "email": "item2@example.com"}, Type: "map"},
			},
			want: `boolean=true
object={email="item2@example.com",nickname="item2"}
project="testproject"
set=["item1","item2"]
`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			s := NewStack()
			s.Settings = tc.in
			got := s.Terraform()
			if !reflect.DeepEqual(got, tc.want) {
				fmt.Printf("Case :%s\n", name)
				fmt.Println(diff.Diff(got, tc.want))
				t.Fatalf("Output Text different than expected")
			}
		})
	}
}

func TestTerraformFile(t *testing.T) {
	tests := map[string]struct {
		filename string
		want     error
	}{
		"Ok": {
			filename: "test_files/file/shouldwork.txt",
			want:     nil,
		},
		"fail": {
			filename: "test_files/file/shouldwork/dir.txt",
			want:     errors.New("open test_files/file/shouldwork/dir.txt: no such file or directory"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewStack()

			got := s.TerraformFile(tc.filename)
			os.Remove(tc.filename)
			if tc.want == nil {
				if got != nil {
					t.Fatalf("expected: no error got: %+v", got)
				}
				t.SkipNow()
			}

			if !strings.Contains(got.Error(), tc.want.Error()) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestStackAddSettings(t *testing.T) {
	tests := map[string]struct {
		in []struct {
			key   string
			value string
		}
		want Settings
	}{
		"basic": {
			in: []struct {
				key   string
				value string
			}{
				{key: "test1", value: "value1"},
				{key: "test_project", value: "project_name"},
			},
			want: Settings{
				Setting{Name: "test1", Value: "value1", Type: "string"},
				Setting{Name: "test_project", Value: "project_name", Type: "string"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			got := NewStack()
			for _, v := range tc.in {
				got.AddSetting(v.key, v.value)
			}

			if !reflect.DeepEqual(tc.want, got.Settings) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got.Settings)
			}
		})
	}
}

func TestStackDeleteSettings(t *testing.T) {
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
			deletekeys: []string{"another", "once"},
			want: Settings{
				Setting{Name: "test1", Value: "value1"},
				Setting{Name: "test_project", Value: "project_name"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			got := NewStack()
			got.Settings = tc.in

			for _, v := range tc.deletekeys {
				got.DeleteSetting(v)
			}

			if !reflect.DeepEqual(tc.want, got.Settings) {
				t.Fatalf("expected: \n%+v, \ngot: \n%+v", tc.want, got.Settings)
			}
		})
	}
}

func TestStackGetSettings(t *testing.T) {
	tests := map[string]struct {
		in   Settings
		key  string
		want string
	}{
		"basic": {
			in: Settings{
				Setting{Name: "test1", Value: "value1"},
				Setting{Name: "test_project", Value: "project_name"},
				Setting{Name: "another", Value: "thing"},
				Setting{Name: "once", Value: "more"},
			},
			key:  "test1",
			want: "value1",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			s := NewStack()
			s.Settings = tc.in
			got := s.GetSetting(tc.key)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: \n%+v, \ngot: \n%+v", tc.want, got)
			}
		})
	}
}
