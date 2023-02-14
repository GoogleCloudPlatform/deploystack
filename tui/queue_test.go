package tui

import (
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/deploystack"
	tea "github.com/charmbracelet/bubbletea"
)

func getTestQueue(title, subtitle string) Queue {
	appHeader := newHeader(title, subtitle)
	stack := deploystack.NewStack()
	mock := mock{}
	q := NewQueue(&stack, mock)
	q.header = appHeader

	return q
}

func TestQueueKeyValue(t *testing.T) {
	tests := map[string]struct {
		key   string
		value interface{}
	}{
		"string": {
			key:   "test",
			value: "alsotest",
		},
		"struct": {
			key:   "test",
			value: struct{ item string }{item: "test"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			q.Save(tc.key, tc.value)

			got := q.Get(tc.key)

			if tc.value != got {
				t.Fatalf("key - want '%s' got '%s'", tc.value, got)
			}
		})
	}
}

func TestQueueStart(t *testing.T) {
	firstPage := newPage("firstpage", []component{newTextBlock(explainText)})
	secondPage := newPage("secondpage", []component{newTextBlock(explainText)})

	tests := map[string]struct {
		models []QueueModel
		exkey  string
	}{
		"single": {
			models: []QueueModel{&firstPage},
			exkey:  "firstpage",
		},
		"multiple": {
			models: []QueueModel{&firstPage, &secondPage},
			exkey:  "firstpage",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			q.add(tc.models...)

			got := q.Start()

			if tc.exkey != got.getKey() {
				t.Fatalf("key - want '%s' got '%s'", tc.exkey, got.getKey())
			}
		})
	}
}

func TestQueueRemoveModel(t *testing.T) {
	firstPage := newPage("firstpage", []component{newTextBlock(explainText)})
	secondPage := newPage("secondpage", []component{newTextBlock(explainText)})
	thirdPage := newPage("thirdpage", []component{newTextBlock(explainText)})
	fourthPage := newPage("fourthpage", []component{newTextBlock(explainText)})

	tests := map[string]struct {
		models []QueueModel
		target string
		want   int
	}{
		"one": {
			models: []QueueModel{&firstPage},
			target: "firstpage",
			want:   0,
		},
		"two": {
			models: []QueueModel{&firstPage, &secondPage},
			target: "firstpage",
			want:   1,
		},
		"four": {
			models: []QueueModel{&firstPage, &secondPage, &thirdPage, &fourthPage},
			target: "thirdpage",
			want:   3,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			q.add(tc.models...)

			q.removeModel(tc.target)

			got := len(q.models)
			if tc.want != got {
				t.Fatalf("want '%d' got '%d'", tc.want, got)
			}
		})
	}
}

func TestQueueProcess(t *testing.T) {
	tests := map[string]struct {
		config string
		keys   []string
	}{
		"basic": {
			config: "testdata/config_basic.yaml",
			keys:   []string{},
		},
		"complex": {
			config: "testdata/config_complex.yaml",
			keys: []string{
				"project_id",
				"project_id_2",
				"project_id" + projNewSuffix,
				"project_id_2" + projNewSuffix,
				"project_id" + billNewSuffix,
				"project_id_2" + billNewSuffix,
				"billing_account",
				"gce-use-defaults",
				"instance-name",
				"region",
				"zone",
				"instance-machine-type-family",
				"instance-machine-type",
				"instance-image-project",
				"instance-image-family",
				"instance-image",
				"instance-disksize",
				"instance-disktype",
				"instance-webserver",
				"domain",
				"domain_email",
				"domain_phone",
				"domain_country",
				"domain_postalcode",
				"domain_state",
				"domain_city",
				"domain_address",
				"domain_name",
				"domain_consent",
				"nodes",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			s := readTestFile(tc.config)

			config, err := deploystack.NewConfigYAML([]byte(s))
			if err != nil {
				t.Fatalf("could not read in config %s:", err)
			}
			q.stack.Config = config

			if err := q.ProcessConfig(); err != nil {
				t.Fatalf("expected no error, got %s", err)
			}

			if len(tc.keys) != len(q.models) {
				t.Logf("Models")
				for i, v := range q.models {
					t.Logf("%d:%s", i, v.getKey())
				}

				t.Fatalf("count - want '%d' got '%d'", len(tc.keys), len(q.models))
			}

			for _, v := range tc.keys {
				q.removeModel(v)
			}

			if 0 != len(q.models) {
				t.Logf("Models remain")
				for _, v := range q.models {
					t.Logf("%s", v.getKey())
				}

				t.Fatalf("key check - want '%d' got '%d'", 0, len(q.models))

			}
		})
	}
}

func TestQueueInitialize(t *testing.T) {
	tests := map[string]struct {
		keys []string
	}{
		"basic": {
			keys: []string{
				"firstpage",
				"descpage",
				"endpage",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")

			q.InitializeUI()

			if len(tc.keys) != len(q.models) {
				t.Logf("Models")
				for i, v := range q.models {
					t.Logf("%d:%s", i, v.getKey())
				}

				t.Fatalf("count - want '%d' got '%d'", len(tc.keys), len(q.models))
			}

			for _, v := range tc.keys {
				q.removeModel(v)
			}

			if 0 != len(q.models) {
				t.Logf("Models remain")
				for _, v := range q.models {
					t.Logf("%s", v.getKey())
				}

				t.Fatalf("key check - want '%d' got '%d'", 0, len(q.models))

			}
		})
	}
}

func TestQueueGoToModel(t *testing.T) {
	firstPage := newPage("firstpage", []component{newTextBlock("A 1st page")})
	secondPage := newPage("secondpage", []component{newTextBlock("A 2nd page")})
	thirdPage := newPage("thirdpage", []component{newTextBlock("A 3rd page")})
	fourthPage := newPage("fourthpage", []component{newTextBlock("A last page")})

	tests := map[string]struct {
		models   []QueueModel
		target   string
		want     string
		wanttype string
	}{
		"one": {
			models:   []QueueModel{&firstPage},
			target:   "firstpage",
			want:     "A 1st page",
			wanttype: "nil",
		},
		"two": {
			models:   []QueueModel{&firstPage, &secondPage},
			target:   "firstpage",
			want:     "A 1st page",
			wanttype: "nil",
		},
		"four": {
			models:   []QueueModel{&firstPage, &secondPage, &thirdPage, &fourthPage},
			target:   "thirdpage",
			want:     "A 3rd page",
			wanttype: "nil",
		},

		"quit": {
			models:   []QueueModel{&firstPage, &secondPage, &thirdPage, &fourthPage},
			target:   "quit",
			want:     "A 3rd page",
			wanttype: "quitMsg",
		},
		"invalidkey": {
			models:   []QueueModel{&firstPage, &secondPage, &thirdPage, &fourthPage},
			target:   "aninvalidkey",
			want:     "A 1st page",
			wanttype: "nil",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			q.add(tc.models...)

			got, cmd := q.goToModel(tc.target)

			if tc.wanttype == "nil" && cmd != nil {
				t.Fatalf("wanted '%s' to be nil got '%+v'", tc.want, cmd)

				if !strings.Contains(got.View(), tc.want) {
					t.Fatalf("wanted '%s' to be contained in got '%s'", tc.want, got.View())
				}
			}

			if tc.wanttype != "nil" {
				gotmsg := cmd()
				wantmsg := tea.Quit()

				if gotmsg != wantmsg {
					t.Fatalf("wanted '%+v' got '%+v'", wantmsg, gotmsg)
				}

			}
		})
	}
}

func TestQueueClear(t *testing.T) {
	firstPage := newPage("firstpage", []component{newTextBlock("A 1st page")})

	tests := map[string]struct {
		model page
		key   string
		value string
		want  string
	}{
		"one": {
			model: firstPage,
			key:   "firstpage",
			value: "A value",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			q.stack.AddSetting(tc.key, tc.value)
			tc.model.value = tc.value
			q.add(&tc.model)

			if q.stack.GetSetting(tc.key) != tc.value {
				t.Fatalf("stack setting did not happen properly")
			}

			q.clear(tc.key)

			if q.stack.GetSetting(tc.key) != "" {
				t.Fatalf("stack clear did not happen properly")
			}

			if tc.model.value != "" {
				t.Fatalf("model clear did not happen properly")
			}
		})
	}
}