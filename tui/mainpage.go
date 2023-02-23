package tui

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/config"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	tea "github.com/charmbracelet/bubbletea"
)

type model interface {
	tea.Model
	getKey() string
	setValue(s string)
	getValue() string
	clear()
	clearContent()
	addPostProcessor(f func(string, *Queue) tea.Cmd)
	addPreProcessor(f tea.Cmd)
	addQueue(q *Queue)
	addContent(s ...string)
	addPreView(f func(*Queue))
}

type mainpage struct {
	header      component
	queue       *Queue
	currentPage *model
}

func newMainPage(h header, q *Queue) mainpage {
	return mainpage{header: h, queue: q}
}

func (m mainpage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	mod := *m.currentPage
	return mod.Update(msg)
}

func (m mainpage) View() string {
	mod := *m.currentPage
	doc := strings.Builder{}

	doc.WriteString(m.header.render())
	doc.WriteString(mod.View())

	return docStyle.Render(doc.String())
}

func (m mainpage) Init() tea.Cmd {
	if m.currentPage == nil {
		curr := m.queue.Start().(model)

		m.currentPage = &curr
	}

	mod := *m.currentPage
	return mod.Init()
}

func AltRun(s *config.Stack, useMock bool) {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	defaultUserAgent := fmt.Sprintf("deploystack/%s", s.Config.Name)

	client := gcloud.NewClient(context.Background(), defaultUserAgent)
	q := NewQueue(s, &client)

	if useMock {
		q = NewQueue(s, GetMock(1))
	}

	// q.Save("contact", deploystack.CheckForContact())
	// q.InitializeUI()
	desc := newDescription(q.stack)
	appHeader := newHeader(appTitle, q.stack.Config.Title)
	firstPage := newPage("firstpage", []component{newTextBlock(explainText)})
	descPage := newPage("descpage", []component{desc})
	firstPage.showProgress = false
	descPage.showProgress = false
	endpage := newPage("endpage", []component{
		newTextBlock(titleStyle.Render("Project Settings")),
		newSettingsTable(q.stack),
	})

	q.header = appHeader
	q.add(&firstPage)
	q.add(&descPage)
	q.ProcessConfig()
	q.add(&endpage)

	m := newMainPage(appHeader, &q)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		Fatal(err)
	}

	if q.Get("halted") != nil {
		Fatal(nil)
	}

	s.TerraformFile("terraform.tfvars")
	deploystack.CacheContact(q.Get("contact"))

	fmt.Print("\n\n")
	fmt.Print(titleStyle.Render("Deploystack"))
	fmt.Print("\n")
	fmt.Print(subTitleStyle.Render(s.Config.Title))
	fmt.Print("\n")
	fmt.Print(strong.Render("Installation will proceed with these settings"))
	fmt.Print(q.getSettings())
}
