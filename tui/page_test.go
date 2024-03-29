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
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestPage(t *testing.T) {
	tests := map[string]struct {
		key        string
		outputFile string
		content    []component
		msg        tea.Msg
	}{
		"basic": {
			key:        "test",
			outputFile: "page_basic.txt",
			content:    []component{newTextBlock(explainText)},
			msg:        successMsg{},
		},
		"send_enter": {
			key:        "test",
			outputFile: "page_send_enter.txt",
			content:    []component{newTextBlock(explainText)},
			msg:        tea.KeyMsg{Type: tea.KeyEnter},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			dummyPage := newPage("dummy", []component{newTextBlock("dummy")})
			p := newPage(tc.key, tc.content)

			q.add(&p)
			q.add(&dummyPage)

			rawP, _ := p.Update(tc.msg)
			var newP page

			switch v := rawP.(type) {
			case page:
				newP = v
			case *page:
				newP = *v
			}

			newP.Init()

			content := newP.View()

			testdata := filepath.Join(testFilesDir, "tui/testdata", tc.outputFile)

			tcOutput := readTestFile(testdata)
			if content != tcOutput {
				writeDebugFile(content, testdata)
				t.Fatalf("text wasn't the same. Look in testdata for expected and debug/testdata for got")
			}
		})
	}
}

func TestDynamicPageValue(t *testing.T) {
	tests := map[string]struct {
		in   page
		want string
	}{
		"basic": {
			in:   newPage("testkey", nil),
			want: "testvalue",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Empty(t, tc.in.getValue())
			tc.in.setValue(tc.want)
			assert.Equal(t, tc.want, tc.in.getValue())
			tc.in.clear()
			assert.Equal(t, "", tc.in.getValue())
		})
	}
}

func TestDynamicPageClearContent(t *testing.T) {
	tests := map[string]struct {
		in   page
		want []component
	}{
		"basic": {
			in: newPage("testkey", []component{
				newTextBlock("TestContent"),
			}),
			want: []component{
				newTextBlock("TestContent"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.in.content)
			tc.in.clearContent()
			assert.Equal(t, []component{}, tc.in.content)
		})
	}
}

func TestPagePreviewFunc(t *testing.T) {
	q := getTestQueue(appTitle, "test")
	page := newPage("test", nil)
	q.add(&page)

	f := func(q *Queue) {
		p := q.models[q.current]
		p.setValue("test")

	}
	page.addPreView(f)

	assert.Empty(t, page.getValue())

	page.View()
	assert.Equal(t, "test", page.getValue())

}
