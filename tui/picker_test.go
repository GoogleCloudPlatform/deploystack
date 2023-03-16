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
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestPicker(t *testing.T) {

	tests := map[string]struct {
		listLabel      string
		spinnerLabel   string
		key            string
		defaultValue   string
		preProcessor   tea.Cmd
		postProcessor  func(string, *Queue) tea.Cmd
		state          string
		outputFile     string
		msg            tea.Msg
		exlistLabel    string
		exspinnerLabel string
		exkey          string
		exstate        string
		content        string
		slowQueryText  string
	}{
		"basic": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: nil,
			exstate:      "idle",
			msg:          tea.MouseEvent{},
			outputFile:   "picker_basic.txt",
		},
		"basic_with_content": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: nil,
			state:        "idle",
			msg:          tea.MouseEvent{},
			outputFile:   "picker_basic_with_content.txt",
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
			outputFile: "picker_spinner.txt",
		},
		"items": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: nil,
			state:        "displaying",
			msg:          tea.Msg([]list.Item{item{label: "Choice", value: "choice"}}),
			outputFile:   "picker_items.txt",
		},
		"items_with_default": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: nil,
			state:        "displaying",
			msg: tea.Msg([]list.Item{
				item{label: "Choice", value: "choice"},
				item{label: "Choice1", value: "choice1"},
				item{label: "Choice2", value: "choice2"},
				item{label: "Choice3", value: "choice3"},
			}),
			defaultValue: "choice3",
			outputFile:   "picker_items_with_default.txt",
		},
		"error": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: nil,
			state:        "idle",
			msg:          errMsg{err: fmt.Errorf("error")},
			outputFile:   "picker_error.txt",
		},
		"slowQueryText": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: func() tea.Cmd {
				return func() tea.Msg {
					items := []list.Item{}
					return items
				}
			}(),
			state:         "querying",
			msg:           tea.MouseEvent{},
			outputFile:    "picker_slowquerytext.txt",
			slowQueryText: "A slow query came through here",
		},

		"success": {
			listLabel:      "test",
			spinnerLabel:   "test",
			key:            "test",
			preProcessor:   nil,
			state:          "idle",
			msg:            successMsg{},
			outputFile:     "picker_success.txt",
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
			outputFile:     "picker_send_enter.txt",
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
			outputFile:   "picker_send_ctrl_c.txt",
			exkey:        "",
		},

		"post_processor": {
			listLabel:    "test",
			spinnerLabel: "test",
			key:          "test",
			preProcessor: func() tea.Cmd {
				return func() tea.Msg {
					items := []list.Item{}
					return items
				}
			}(),
			state:      "displaying",
			exstate:    "querying",
			msg:        tea.KeyMsg{Type: tea.KeyEnter},
			outputFile: "picker_post_processor.txt",
			postProcessor: func(projectID string, q *Queue) tea.Cmd {
				return func() tea.Msg {

					return successMsg{}
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			dummyPicker := newPicker("dummy", "dummy", "dummy", "", nil)

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

			ptmp := newPicker(tc.listLabel, tc.spinnerLabel, tc.key, tc.defaultValue, tc.preProcessor)

			if tc.content != "" {
				ptmp.addContent(tc.content)
			}

			ptmp.querySlowText = tc.slowQueryText

			if tc.postProcessor != nil {
				ptmp.addPostProcessor(tc.postProcessor)
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
				testdata := filepath.Join(testFilesDir, "tui/testdata", tc.outputFile)
				tcOutput := readTestFile(testdata)
				if content != tcOutput {
					writeDebugFile(content, testdata)
					t.Fatalf("text wasn't the same. Look in testdata for expected and debug/testdata for got")
				}

			}
		})
	}
}

func TestPositionDefault(t *testing.T) {
	tests := map[string]struct {
		items        []list.Item
		defaultValue string
		wantItems    []list.Item
		wantIndex    int
	}{
		"empty": {
			items:        []list.Item{},
			defaultValue: "",
			wantItems:    []list.Item{},
			wantIndex:    0,
		},
		"basic": {
			items: []list.Item{
				item{label: "label1", value: "value1"},
				item{label: "label2", value: "value2"},
				item{label: "label3", value: "value3"},
				item{label: "label4", value: "value4"},
				item{label: "label5", value: "value5"},
				item{label: "label6", value: "value6"},
			},
			defaultValue: "label3",
			wantItems: []list.Item{
				item{label: "label1", value: "value1"},
				item{label: "label2", value: "value2"},
				item{label: "label3 (Default Value)", value: "value3"},
				item{label: "label4", value: "value4"},
				item{label: "label5", value: "value5"},
				item{label: "label6", value: "value6"},
			},
			wantIndex: 2,
		},
		"MoreThan12": {
			items: []list.Item{
				item{label: "label1", value: "value1"},
				item{label: "label2", value: "value2"},
				item{label: "label3", value: "value3"},
				item{label: "label4", value: "value4"},
				item{label: "label5", value: "value5"},
				item{label: "label6", value: "value6"},
				item{label: "label7", value: "value7"},
				item{label: "label8", value: "value8"},
				item{label: "label9", value: "value9"},
				item{label: "label10", value: "value10"},
				item{label: "label11", value: "value11"},
				item{label: "label12", value: "value12"},
			},
			defaultValue: "label3",
			wantItems: []list.Item{
				item{label: "label3 (Default Value)", value: "value3"},
				item{label: "label1", value: "value1"},
				item{label: "label2", value: "value2"},
				item{label: "label4", value: "value4"},
				item{label: "label5", value: "value5"},
				item{label: "label6", value: "value6"},
				item{label: "label7", value: "value7"},
				item{label: "label8", value: "value8"},
				item{label: "label9", value: "value9"},
				item{label: "label10", value: "value10"},
				item{label: "label11", value: "value11"},
				item{label: "label12", value: "value12"},
			},
			wantIndex: 0,
		},
		"MoreThan12Default1st": {
			items: []list.Item{
				item{label: "label1", value: "value1"},
				item{label: "label2", value: "value2"},
				item{label: "label3", value: "value3"},
				item{label: "label4", value: "value4"},
				item{label: "label5", value: "value5"},
				item{label: "label6", value: "value6"},
				item{label: "label7", value: "value7"},
				item{label: "label8", value: "value8"},
				item{label: "label9", value: "value9"},
				item{label: "label10", value: "value10"},
				item{label: "label11", value: "value11"},
				item{label: "label12", value: "value12"},
			},
			defaultValue: "label1",
			wantItems: []list.Item{
				item{label: "label1 (Default Value)", value: "value1"},
				item{label: "label2", value: "value2"},
				item{label: "label3", value: "value3"},
				item{label: "label4", value: "value4"},
				item{label: "label5", value: "value5"},
				item{label: "label6", value: "value6"},
				item{label: "label7", value: "value7"},
				item{label: "label8", value: "value8"},
				item{label: "label9", value: "value9"},
				item{label: "label10", value: "value10"},
				item{label: "label11", value: "value11"},
				item{label: "label12", value: "value12"},
			},
			wantIndex: 0,
		},
		"MoreThan12WithCreate": {
			items: []list.Item{
				item{label: "label1", value: "value1"},
				item{label: "label2", value: "value2"},
				item{label: "label3", value: "value3"},
				item{label: "label4", value: "value4"},
				item{label: "label5", value: "value5"},
				item{label: "label6", value: "value6"},
				item{label: "label7", value: "value7"},
				item{label: "label8", value: "value8"},
				item{label: "label9", value: "value9"},
				item{label: "label10", value: "value10"},
				item{label: "label11", value: "value11"},
				item{label: "label12", value: "value12"},
				item{label: "Create New Project", value: ""},
			},
			defaultValue: "label3",
			wantItems: []list.Item{
				item{label: "Create New Project", value: ""},
				item{label: "label3 (Default Value)", value: "value3"},
				item{label: "label1", value: "value1"},
				item{label: "label2", value: "value2"},
				item{label: "label4", value: "value4"},
				item{label: "label5", value: "value5"},
				item{label: "label6", value: "value6"},
				item{label: "label7", value: "value7"},
				item{label: "label8", value: "value8"},
				item{label: "label9", value: "value9"},
				item{label: "label10", value: "value10"},
				item{label: "label11", value: "value11"},
				item{label: "label12", value: "value12"},
			},
			wantIndex: 1,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotItems, gotIndex := positionDefault(tc.items, tc.defaultValue)
			assert.Equal(t, tc.wantItems, gotItems)
			assert.Equal(t, tc.wantIndex, gotIndex)

		})
	}
}
