package tui

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/kylelemons/godebug/diff"
)

func TestProductListLongest(t *testing.T) {
	tests := map[string]struct {
		configPath  string
		wantItem    int
		wantProduct int
	}{
		"simple": {
			configPath:  "testdata/config_basic.yaml",
			wantItem:    14,
			wantProduct: 22,
		},

		"long_description": {
			configPath:  "testdata/config_long_description.yaml",
			wantItem:    20,
			wantProduct: 38,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := readTestFile(tc.configPath)

			stack := deploystack.NewStack()
			config, err := deploystack.NewConfigYAML([]byte(s))
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
			configPath: "testdata/config_basic.yaml",
			outputFile: "testdata/description_basic.txt",
		},

		"one_min": {
			configPath: "testdata/config_one_min.yaml",
			outputFile: "testdata/description_one_min.txt",
		},

		"long_description": {
			configPath: "testdata/config_long_description.yaml",
			outputFile: "testdata/description_long_description.txt",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := readTestFile(tc.configPath)

			stack := deploystack.NewStack()
			config, err := deploystack.NewConfigYAML([]byte(s))
			if err != nil {
				t.Fatalf("could not read in config %s:", err)
			}
			stack.Config = config

			d := newDescription(&stack)

			want := readTestFile(tc.outputFile)
			got := d.render()

			if want != got {
				fmt.Println(diff.Diff(want, got))
				writeDebugFile(got, tc.outputFile)
				t.Fatalf("text wasn't the same")
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
			outputFile: "testdata/error_alert_no_user_message.txt",
		},

		"UserMessage": {
			errMsg: errMsg{
				err:     fmt.Errorf("Everything broke"),
				usermsg: "It was probably something you said",
			},
			outputFile: "testdata/error_alert_user_message.txt",
		},

		"TargetQuit": {
			errMsg: errMsg{
				err:    fmt.Errorf("Everything broke"),
				target: "quit",
			},
			outputFile: "testdata/error_alert_target_quit.txt",
		},
		"TargetOther": {
			errMsg: errMsg{
				err:    fmt.Errorf("Everything broke"),
				target: "other",
			},
			outputFile: "testdata/error_alert_target_other.txt",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			e := errorAlert{tc.errMsg}

			want := readTestFile(tc.outputFile)
			got := e.Render()

			if want != got {
				fmt.Println(diff.Diff(want, got))
				writeDebugFile(got, tc.outputFile)
				t.Fatalf("text wasn't the same")
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
			outputFile: "testdata/settingstable_basic.txt",
		},
		"average": {
			settings: map[string]string{
				"project_id":     "test-id",
				"project_number": "123344567",
				"project_name":   "test-project",
				"stack_name":     "test-stack-value",
				"testkey":        "testvalue",
			},
			outputFile: "testdata/settingstable_average .txt",
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
			outputFile: "testdata/settingstable_outliers .txt",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			stack := deploystack.NewStack()

			for key, value := range tc.settings {
				stack.AddSetting(key, value)
			}

			table := newSettingsTable(&stack)

			want := readTestFile(tc.outputFile)
			got := table.render()

			if want != got {
				fmt.Println(diff.Diff(want, got))
				writeDebugFile(got, tc.outputFile)
				t.Fatalf("text wasn't the same")
			}
		})
	}
}
