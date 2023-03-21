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

package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/deploystack/config"
	"github.com/kylelemons/godebug/diff"
)

func TestDrawProgress(t *testing.T) {
	tests := map[string]struct {
		in   int
		want string
		len  int
	}{
		"50%": {
			in:   50,
			want: "[0;37m   Progress [1;36m████████████████████████████████████████████[0;37m░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░",
			len:  hardWidthLimit,
		},
		"0%": {
			in:   0,
			want: "[0;37m   Progress [1;36m[0;37m░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░",
			len:  hardWidthLimit,
		},
		"100%": {
			in:   100,
			want: "[0;37m   Progress [1;36m████████████████████████████████████████████████████████████████████████████████████████[0;37m",
			len:  hardWidthLimit,
		},
		"75%": {
			in:   75,
			want: "[0;37m   Progress [1;36m██████████████████████████████████████████████████████████████████[0;37m░░░░░░░░░░░░░░░░░░░░░░",
			len:  hardWidthLimit,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := drawProgress(tc.in)

			got = strings.ReplaceAll(got, "\x1b[1;m", "")
			got = strings.ReplaceAll(got, "\x1b[0m", "")

			if tc.want != got {
				t.Fatalf("want \n%s\n got\n%s\n", tc.want, got)
			}

		})
	}

}

func TestProductListLongest(t *testing.T) {
	tests := map[string]struct {
		configPath  string
		wantItem    int
		wantProduct int
	}{
		"simple": {
			configPath:  "config_basic.yaml",
			wantItem:    14,
			wantProduct: 22,
		},

		"long_description": {
			configPath:  "config_long_description.yaml",
			wantItem:    20,
			wantProduct: 38,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testdata := filepath.Join(testFilesDir, "tui/testdata", tc.configPath)
			s := readTestFile(testdata)

			stack := config.NewStack()
			config, err := config.NewConfigYAML([]byte(s))
			if err != nil {
				t.Fatalf("could not read in config %s:", err)
			}
			stack.Config = config

			d := newDescription(&stack)

			prods, _ := d.parse()

			gotItem := prods.longest("item")
			gotProduct := prods.longest("product")

			if tc.wantItem != gotItem {
				t.Fatalf("item - want '%d' got '%d'", tc.wantItem, gotItem)
			}

			if tc.wantProduct != gotProduct {
				t.Fatalf("roduct - want '%d' got '%d'", tc.wantProduct, gotProduct)
			}
		})
	}
}

func TestDescriptionRender(t *testing.T) {
	tests := map[string]struct {
		configPath string
		outputFile string
	}{
		"simple": {
			configPath: "config_basic.yaml",
			outputFile: "description_basic.txt",
		},

		"one_min": {
			configPath: "config_one_min.yaml",
			outputFile: "description_one_min.txt",
		},

		"long_description": {
			configPath: "config_long_description.yaml",
			outputFile: "description_long_description.txt",
		},
		"products_description": {
			configPath: "config_product_description.yaml",
			outputFile: "description_product_description.txt",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			testdata := filepath.Join(testFilesDir, "tui/testdata", tc.configPath)

			s := readTestFile(testdata)

			stack := config.NewStack()
			config, err := config.NewConfigYAML([]byte(s))
			if err != nil {
				t.Fatalf("could not read in config %s:", err)
			}
			stack.Config = config

			d := newDescription(&stack)

			outputdata := filepath.Join(testFilesDir, "tui/testdata", tc.outputFile)
			want := readTestFile(outputdata)
			got := d.render()

			if want != got {
				fmt.Println(diff.Diff(want, got))
				writeDebugFile(got, outputdata)
				t.Fatalf("text wasn't the same. Look in testdata for expected and debug/testdata for got")
			}
		})
	}
}

func TestErrorAlertRender(t *testing.T) {
	tests := map[string]struct {
		errMsg     errMsg
		outputFile string
	}{
		"NoUserMessage": {
			errMsg:     errMsg{err: fmt.Errorf("Everything broke")},
			outputFile: "error_alert_no_user_message.txt",
		},

		"UserMessage": {
			errMsg: errMsg{
				err:     fmt.Errorf("Everything broke"),
				usermsg: "It was probably something you said",
			},
			outputFile: "error_alert_user_message.txt",
		},

		"TargetQuit": {
			errMsg: errMsg{
				err:    fmt.Errorf("Everything broke"),
				target: "quit",
			},
			outputFile: "error_alert_target_quit.txt",
		},
		"TargetOther": {
			errMsg: errMsg{
				err:    fmt.Errorf("Everything broke"),
				target: "other",
			},
			outputFile: "error_alert_target_other.txt",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			e := errorAlert{tc.errMsg}

			testdata := filepath.Join(testFilesDir, "tui/testdata", tc.outputFile)

			want := readTestFile(testdata)
			got := e.Render()

			if want != got {
				fmt.Println(diff.Diff(want, got))
				writeDebugFile(got, testdata)
				t.Fatalf("text wasn't the same. Look in testdata for expected and debug/testdata for got")
			}
		})
	}
}

func TestSettingsTableRender(t *testing.T) {
	tests := map[string]struct {
		settings   map[string]string
		outputFile string
	}{
		"simple": {
			settings: map[string]string{
				"testkey": "testvalue",
			},
			outputFile: "settingstable_basic.txt",
		},
		"average": {
			settings: map[string]string{
				"project_id":     "test-id",
				"project_number": "123344567",
				"project_name":   "test-project",
				"stack_name":     "test-stack-value",
				"testkey":        "testvalue",
			},
			outputFile: "settingstable_average .txt",
		},
		"outliers": {
			settings: map[string]string{
				"project_id":     "test-id",
				"project_number": "123344567",
				"project_name":   "test-project",
				"stack_name":     "test-stack-value",
				"testkey":        "testvalue",
				"testkey2":       "12345678901234567890123456789012345678901234567890",
				"empty":          "",
			},
			outputFile: "settingstable_outliers .txt",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			stack := config.NewStack()

			for key, value := range tc.settings {
				stack.AddSetting(key, value)
			}

			table := newSettingsTable(&stack)

			testdata := filepath.Join(testFilesDir, "tui/testdata", tc.outputFile)
			want := readTestFile(testdata)
			got := table.render()

			if want != got {
				fmt.Println(diff.Diff(want, got))
				writeDebugFile(got, tc.outputFile)
				t.Fatalf("text wasn't the same. Look in testdata for expected and debug/testdata for got")
			}
		})
	}
}
