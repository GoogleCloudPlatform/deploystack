package tui

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func TestPicker(t *testing.T) {
	tests := map[string]struct {
		listLabel      string
		spinnerLabel   string
		key            string
		preProcessor   tea.Cmd
		state          string
		outputFile     string
		msg            tea.Msg
		exlistLabel    string
		exspinnerLabel string
		exkey          string
		exstate        string
		content        string
	}{
		"basic": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: nil,
			exstate:      "idle",
			msg:          tea.MouseEvent{},
			outputFile:   "testdata/picker_basic.txt",
		},
		"basic_with_content": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: nil,
			state:        "idle",
			msg:          tea.MouseEvent{},
			outputFile:   "testdata/picker_basic_with_content.txt",
			content:      "Adding some basic content to test",
		},
		"spinner": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: func() tea.Cmd {
				return func() tea.Msg {
					items := []list.Item{}
					return items
				}
			}(),
			state:      "querying",
			msg:        tea.MouseEvent{},
			outputFile: "testdata/picker_spinner.txt",
		},
		"items": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: nil,
			state:        "displaying",
			msg:          tea.Msg([]list.Item{item{label: "Choice", value: "choice"}}),
			outputFile:   "testdata/picker_items.txt",
		},
		"error": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: nil,
			state:        "idle",
			msg:          errMsg{err: fmt.Errorf("error")},
			outputFile:   "testdata/picker_error.txt",
		},

		"success": {
			listLabel:      "test",
			spinnerLabel:   "test",
			key:            "test",
			preProcessor:   nil,
			state:          "idle",
			msg:            successMsg{},
			outputFile:     "testdata/picker_success.txt",
			exlistLabel:    "dummy",
			exspinnerLabel: "dummy",
			exkey:          "dummy",
			exstate:        "idle",
		},

		"send_enter": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: func() tea.Cmd {
				return func() tea.Msg {
					items := []list.Item{}
					return items
				}
			}(),
			state:          "displaying",
			msg:            tea.KeyMsg{Type: tea.KeyEnter},
			outputFile:     "testdata/picker_send_enter.txt",
			exlistLabel:    "dummy",
			exspinnerLabel: "dummy",
			exkey:          "dummy",
			exstate:        "idle",
		},

		"send_ctrl_c": {
			listLabel:    "",
			spinnerLabel: "",
			key:          "",
			preProcessor: nil,
			state:        "",
			msg:          tea.KeyMsg{Type: tea.KeyCtrlC},
			outputFile:   "testdata/picker_send_ctrl_c.txt",
			exkey:        "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			dummyPicker := newPicker("dummy", "dummy", "dummy", nil)

			if tc.exkey == "" {
				tc.exkey = tc.key
			}

			if tc.exlistLabel == "" {
				tc.exlistLabel = tc.listLabel
			}

			if tc.exstate == "" {
				tc.exstate = tc.state
			}

			if tc.exspinnerLabel == "" {
				tc.exspinnerLabel = tc.spinnerLabel
			}

			ptmp := newPicker(tc.listLabel, tc.spinnerLabel, tc.key, tc.preProcessor)

			if tc.content != "" {
				ptmp.addContent(tc.content)
			}
			q.add(&ptmp)
			q.add(&dummyPicker)

			p := q.models[0].(*picker)

			p.Init()
			if tc.state != "" {
				p.state = tc.state
			}

			rawP, cmd := p.Update(tc.msg)
			tea.Batch(cmd)
			var newP picker

			switch v := rawP.(type) {
			case picker:
				newP = v
			case *picker:
				newP = *v
			}

			newP.Init()

			if name == "send_ctrl_c" {
				fmt.Printf("%+v", newP)
			}

			if tc.exkey != newP.key {
				t.Fatalf("key - want '%s' got '%s'", tc.exkey, newP.key)
			}

			if tc.exlistLabel != newP.list.Title {
				t.Fatalf("listLabel - want '%s' got '%s'", tc.exlistLabel, newP.list.Title)
			}

			if tc.exspinnerLabel != newP.spinnerLabel {
				t.Fatalf("spinnerLabel - want '%s' got '%s'", tc.exspinnerLabel, newP.spinnerLabel)
			}

			if tc.exstate != newP.state {
				t.Fatalf("state - want '%s' got '%s'", tc.exstate, newP.state)
			}

			if newP.key != "" {
				content := newP.View()
				tcOutput := readTestFile(tc.outputFile)
				if content != tcOutput {
					writeDebugFile(content, tc.outputFile)
					t.Fatalf("text wasn't the same")
				}

			}
		})
	}
}
