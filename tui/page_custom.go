package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type textInput struct {
	dynamicPage

	label string
	ti    textinput.Model
}

func newTextInput(label, defaultValue, key, spinnerLabel string) textInput {
	t := textInput{}
	t.key = key
	t.label = label

	t.state = "idle"
	t.spinnerLabel = spinnerLabel

	ti := textinput.New()
	ti.Placeholder = defaultValue
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	t.ti = ti

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle
	t.spinner = s

	return t
}

func (p textInput) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, p.spinner.Tick)
}

func (p textInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	keyTarget := strings.ReplaceAll(p.key, projNewSuffix, "")

	// if the intended key for this setting is already set, skip
	if p.queue.stack.GetSetting(p.key) != "" ||
		p.queue.stack.GetSetting(keyTarget) != "" {
		return p.queue.next()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return p, tea.Quit
		case "enter":
			val := p.ti.Value()
			if val == "" {
				val = p.ti.Placeholder
			}
			// TODO: see if you can figure out a test for these empty bits
			if val == "" {
				p.err = fmt.Errorf("You must enter a value")
				return p, nil
			}
			p.value = val

			// TODO: see if you can figure out a test for these untested bits
			if p.postProcessor != nil {
				if p.state != "querying" {
					p.state = "querying"
					p.err = nil
					return p, p.postProcessor(p.value, p.queue)
				}

				return p, nil
			}
			p.queue.stack.AddSetting(p.key, p.value)
			return p.queue.next()
		}

	// We handle errors just like any other message
	case errMsg:
		p.err = msg
		p.state = "idle"

		if msg.quit {
			return p, tea.Quit
		}

		var cmdSpin tea.Cmd
		p.spinner, cmdSpin = p.spinner.Update(msg)
		return p, cmdSpin
	case successMsg:
		// Filter project creation screens screeens
		newKey := strings.ReplaceAll(p.key, projNewSuffix, "")

		p.queue.stack.AddSetting(newKey, p.value)
		return p.queue.next()

	}
	var cmdSpin tea.Cmd
	p.spinner, cmdSpin = p.spinner.Update(msg)
	p.ti, cmd = p.ti.Update(msg)
	return p, tea.Batch(cmd, cmdSpin)
}

func (p textInput) View() string {
	if p.preViewFunc != nil {
		p.preViewFunc(p.queue)
	}

	doc := strings.Builder{}
	doc.WriteString(p.queue.header.render())

	doc.WriteString(bodyStyle.Render(titleStyle.Render(fmt.Sprintf("%s: ", p.label))))
	doc.WriteString("\n")

	inst := strings.Builder{}
	for _, v := range p.content {
		inst.WriteString(v.render())
	}
	doc.WriteString(instructionStyle.Width(width).Render(inst.String()))

	doc.WriteString("\n")
	doc.WriteString(inputText.Render(p.ti.View()))
	doc.WriteString("\n")

	if p.err != nil {
		height := len(p.err.Error()) / width
		doc.WriteString("\n")
		doc.WriteString(alertStyle.Width(width).Height(height).Render(fmt.Sprintf("Error: %s", p.err)))
		doc.WriteString("\n")
	}

	if p.state == "querying" && p.err == nil {
		doc.WriteString("\n")
		doc.WriteString(bodyStyle.Render(fmt.Sprintf("%s %s", p.spinnerLabel, p.spinner.View())))
		doc.WriteString("\n")
	}

	if p.state != "querying" {
		if p.ti.Placeholder != "" {
			styledPlaceHolder := textInputDefaultStyle.Render(p.ti.Placeholder)
			doc.WriteString(textInputPrompt.Render(fmt.Sprintf("Type a value or hit enter for '%s'", styledPlaceHolder)))
		} else {
			doc.WriteString(textInputPrompt.Render("Type a value and hit enter to continue"))
		}
	}

	return docStyle.Render(doc.String())
}
