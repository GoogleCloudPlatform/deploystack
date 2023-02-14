package deploystack

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

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