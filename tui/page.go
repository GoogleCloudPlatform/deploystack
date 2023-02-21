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
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type dynamicPage struct {
	queue            *Queue
	spinner          spinner.Model
	spinnerLabel     string
	key              string
	value            string
	err              error
	state            string
	content          []component
	preProcessor     tea.Cmd
	postProcessor    func(string, *Queue) tea.Cmd
	preViewFunc      func(*Queue)
	showProgress     bool
	omitFromSettings bool
	querySlowText    string
}

func (p *dynamicPage) getKey() string {
	return p.key
}

func (p *dynamicPage) setValue(s string) {
	p.value = s
}

func (p *dynamicPage) getValue() string {
	return p.value
}

func (p *dynamicPage) clear() {
	p.value = ""
}

func (p *dynamicPage) clearContent() {
	p.content = []component{}
}

func (p *dynamicPage) addPostProcessor(f func(string, *Queue) tea.Cmd) {
	p.postProcessor = f
}

func (p *dynamicPage) addPreProcessor(f tea.Cmd) {
	p.preProcessor = f
}

func (p *dynamicPage) addQueue(q *Queue) {
	p.queue = q
}

func (p *dynamicPage) addContent(s ...string) {
	for _, v := range s {
		p.content = append(p.content, textBlock(v))
	}
}

func (p *dynamicPage) addPreView(f func(*Queue)) {
	p.preViewFunc = f
}

type page struct {
	dynamicPage
}

func newPage(key string, content []component) page {
	p := page{}
	p.key = key
	p.content = content
	p.showProgress = true
	return p
}

func (p page) Init() tea.Cmd {
	return p.preProcessor
}

func (p page) View() string {
	if p.preViewFunc != nil {
		p.preViewFunc(p.queue)
	}
	doc := strings.Builder{}
	doc.WriteString(p.queue.header.render())
	if p.showProgress {
		doc.WriteString(drawProgress(p.queue.calcPercent()))
		doc.WriteString("\n\n")
	}

	for _, v := range p.content {
		doc.WriteString(bodyStyle.Render(v.render()))
		doc.WriteString("\n")
	}

	doc.WriteString("\n")
	doc.WriteString(bodyStyle.Render(promptStyle.Render(" Press the Enter Key to continue ")))

	test := docStyle.Render(doc.String())

	return test
}

// TODO: a test for this is pretty straight forward
func (p page) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case successMsg:
		return p.queue.next()
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).String() {

		case "alt+b", "ctrl+b":
			return p.queue.prev()
		case "ctrl+c", "q":
			if p.queue.Get("halted") != nil {
				os.Exit(1)
			}
			return p.queue.exitPage()
		case "enter":
			if p.postProcessor != nil {
				if p.state != "querying" {
					p.state = "querying"
					p.err = nil
					return p, p.postProcessor(p.value, p.queue)
				}

				return p, nil
			}

			return p.queue.next()
		}

	}
	return p, nil
}
