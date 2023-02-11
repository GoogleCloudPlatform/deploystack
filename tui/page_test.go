package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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
			outputFile: "testdata/page_basic.txt",
			content:    []component{newTextBlock(explainText)},
			msg:        successMsg{},
		},
		"send_enter": {
			key:        "test",
			outputFile: "testdata/page_send_enter.txt",
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
			tcOutput := readTestFile(tc.outputFile)
			if content != tcOutput {
				writeDebugFile(content, tc.outputFile)
				t.Fatalf("text wasn't the same")
			}
		})
	}
}
