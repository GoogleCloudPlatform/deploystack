package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type dynamicPage struct {
	queue         *Queue
	spinner       spinner.Model
	spinnerLabel  string
	key           string
	value         string
	err           error
	state         string
	content       []component
	preProcessor  tea.Cmd
	postProcessor func(string, *Queue) tea.Cmd
	preViewFunc   func(*Queue)
}

func (p *dynamicPage) getKey() string {
	return p.key
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

	for _, v := range p.content {
		doc.WriteString(bodyStyle.Render(v.render()))
		doc.WriteString("\n")
	}

	doc.WriteString("\n")
	doc.WriteString(bodyStyle.Render(promptStyle.Render(" Press the Enter Key to continue ")))

	return docStyle.Render(doc.String())
}

// TODO: a test for this is pretty straight forward
func (p page) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case successMsg:
		return p.queue.next()
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).String() {
		case "ctrl+c", "q":
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
