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
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTextInput(t *testing.T) {
	tests := map[string]struct {
		outputFile     string
		label          string
		defaultValue   string
		key            string
		state          string
		spinnerLabel   string
		content        []component
		msg            tea.Msg
		exlabel        string
		exdefaultValue string
		exkey          string
		exstate        string
		exspinnerLabel string
	}{
		"basic": {
			outputFile:   "testdata/page_custom_basic.txt",
			label:        "test",
			spinnerLabel: "test",
			key:          "test",
			state:        "idle ",
			msg:          tea.MouseEvent{},
		},

		"send_enter": {
			label:          "test",
			spinnerLabel:   "test",
			defaultValue:   "test",
			key:            "test",
			msg:            tea.KeyMsg{Type: tea.KeyEnter},
			outputFile:     "testdata/page_custom_send_enter.txt",
			exlabel:        "dummy",
			exspinnerLabel: "loading dummy",
			exkey:          "dummy",
			exstate:        "idle",
		},

		"success": {
			outputFile:     "testdata/page_custom_success.txt",
			label:          "test",
			spinnerLabel:   "test",
			key:            "test",
			state:          "idle",
			msg:            successMsg{},
			exlabel:        "dummy",
			exspinnerLabel: "loading dummy",
			exkey:          "dummy",
		},

		"spinner": {
			outputFile:   "testdata/page_custom_spinner.txt",
			label:        "test",
			spinnerLabel: "test",
			key:          "test",
			state:        "querying",
			msg:          tea.MouseEvent{},
		},

		"error": {
			outputFile:   "testdata/page_custom_error.txt",
			label:        "test",
			spinnerLabel: "test",
			key:          "test",
			state:        "idle",
			msg:          errMsg{err: fmt.Errorf("error")},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.exkey == "" {
				tc.exkey = tc.key
			}

			if tc.exlabel == "" {
				tc.exlabel = tc.label
			}

			if tc.exstate == "" {
				tc.exstate = tc.state
			}

			if tc.exspinnerLabel == "" {
				tc.exspinnerLabel = tc.spinnerLabel
			}

			q := getTestQueue(appTitle, "test")
			dummyTi := newTextInput("dummy", "dummy", "dummy", "loading dummy")
			testTi := newTextInput(tc.label, tc.defaultValue, tc.key, tc.spinnerLabel)

			q.add(&testTi)
			q.add(&dummyTi)

			ti := q.models[0].(*textInput)

			ti.Init()
			if tc.state != "" {
				ti.state = tc.state
			}

			rawTi, cmd := ti.Update(tc.msg)
			tea.Batch(cmd)

			var nextTi textInput

			switch v := rawTi.(type) {
			case textInput:
				nextTi = v
			case *textInput:
				nextTi = *v
			}

			nextTi.Init()

			if tc.exkey != nextTi.key {
				t.Fatalf("key - want '%s' got '%s'", tc.exkey, nextTi.key)
			}

			if tc.exlabel != nextTi.label {
				t.Fatalf("listLabel - want '%s' got '%s'", tc.exlabel, nextTi.label)
			}

			if tc.exspinnerLabel != nextTi.spinnerLabel {
				t.Fatalf("spinnerLabel - want '%s' got '%s'", tc.exspinnerLabel, nextTi.spinnerLabel)
			}

			if tc.exstate != nextTi.state {
				t.Fatalf("state - want '%s' got '%s'", tc.exstate, nextTi.state)
			}

			content := nextTi.View()
			tcOutput := readTestFile(tc.outputFile)
			if content != tcOutput {
				writeDebugFile(content, tc.outputFile)
				t.Fatalf("text wasn't the same")
			}
		})
	}
}
