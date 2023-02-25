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
	"github.com/GoogleCloudPlatform/deploystack/config"
	tea "github.com/charmbracelet/bubbletea"
)

// QueueModel is an extented version of tea.Model modified to work with the
// queue system
type QueueModel interface {
	tea.Model
	addQueue(*Queue)
	getKey() string
	setValue(string)
	addContent(...string)
	clearContent()
	clear()
}

// Queue represents the flow of the application from screen to screen, or
// a the developer level tea.Model to tea.Model. It allows for progression
// and even going back through the queue all to manage the population of
// a deploystack setting and tfvars file
type Queue struct {
	models  []QueueModel
	current int
	header  component
	stack   *config.Stack
	store   map[string]interface{}
	index   []string
	client  UIClient
}

// NewQueue creates a new queue. You should need only one per app
func NewQueue(s *config.Stack, client UIClient) Queue {
	q := Queue{stack: s, store: map[string]interface{}{}}
	q.client = client
	q.index = []string{}

	currentProject, _ := client.ProjectIDGet()

	q.Save("currentProject", currentProject)
	return q
}

// Model retrieves a give model by key from the queue
func (q *Queue) Model(key string) QueueModel {

	for _, v := range q.models {
		if v.getKey() == key {
			return v
		}
	}

	return nil
}

// Save stores a value in a simple cache for communicating between operations
// in the same process
func (q *Queue) Save(key string, val interface{}) {
	q.store[key] = val
}

// Get returns a previously stored value from the Queue cache
func (q *Queue) Get(key string) interface{} {
	val, ok := q.store[key]
	if !ok {
		return nil
	}
	return val
}

func (q *Queue) removeModel(key string) {
	for i, v := range q.index {
		if v == key {
			q.models = append(q.models[:i], q.models[i+1:]...)
			q.index = append(q.index[:i], q.index[i+1:]...)
		}
	}
}

func (q *Queue) goToModel(key string) (tea.Model, tea.Cmd) {
	if key == "quit" {
		return q.models[q.current], tea.Quit
	}

	for i, v := range q.models {
		if v.getKey() == key {
			q.current = i
			r := q.models[q.current]
			return r, r.Init()
		}
	}

	r := q.models[q.current]
	return r, nil
}

func (q *Queue) clear(key string) {
	for i, v := range q.models {
		if v.getKey() == key {
			q.current = i
			r := q.models[q.current]
			q.stack.DeleteSetting(key)
			r.clear()

		}
	}
}

func (q *Queue) next() (tea.Model, tea.Cmd) {
	q.current++
	if q.current >= len(q.models) {
		return q.models[len(q.models)-1], tea.Quit
	}

	r := q.models[q.current]
	return r, r.Init()
}

func (q *Queue) prev() (tea.Model, tea.Cmd) {
	q.current--
	if q.current <= 0 {
		return q.models[0], nil
	}

	r := q.models[q.current]
	r.setValue("")
	return r, r.Init()
}

func (q *Queue) currentKey() string {
	if len(q.models) == 0 {
		return ""
	}

	r := q.models[q.current].getKey()
	return r
}

// InitializeUI spins up everything we need to have a working queue in the
// hosting application
func (q *Queue) InitializeUI() {
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
	endpage.addPreProcessor(cleanUp(q))

	q.header = appHeader
	q.add(&firstPage)
	q.add(&descPage)
	q.ProcessConfig()
	q.add(&endpage)
}

func (q *Queue) getSettings() string {
	r := newSettingsTable(q.stack)

	return r.render()
}

func (q *Queue) exitPage() (tea.Model, tea.Cmd) {
	page := newPage("exit", []component{
		newTextBlock("You've chosen to stop moving forward through DeployStack. \n"),
		newTextBlock("If this was an error, you can try again by typing 'deploystack install' at the command prompt \n"),
	})
	page.showProgress = false
	q.add(&page)
	q.Save("halted", true)

	quit := func(string, *Queue) tea.Cmd {
		return tea.Quit
	}

	page.addPostProcessor(quit)

	return page, nil
}

func (q *Queue) countTotalSteps() int {
	total := len(q.models)

	for _, v := range q.models {
		if v.getKey() == "firstpage" {
			total--
		}

		if v.getKey() == "descpage" {
			total--
		}

		if v.getKey() == "endpage" {
			total--
		}
	}
	return total
}

func (q *Queue) calcPercent() int {

	if q.current == 2 {
		return 0
	}

	if q.current == len(q.models)-1 {
		return 100
	}
	total := q.countTotalSteps()
	current := q.current + 1 - 2
	percentage := int((float32(current) / float32(total)) * 100)

	if percentage >= 100 && q.current != len(q.models)-1 {
		return 90
	}

	return percentage
}

// ProcessConfig does the work of turning a DeployStack config file to a set
// of tui screens. It's separate from Initialize in case we want to be able
// to populate setting and variables with other information before running
// the genreation of those screens
func (q *Queue) ProcessConfig() error {
	var project, name string
	var err error

	s := q.stack

	sets := s.Config.GetAuthorSettings()

	for _, v := range sets {
		s.AddSettingComplete(v)

	}

	project = s.GetSetting("project_id")
	region := s.GetSetting("region")
	zone := s.GetSetting("zone")
	name = s.Config.Name

	if name == "" {
		err = s.Config.ComputeName()
		if err != nil {
			return err
		}
	}
	s.AddSetting("stack_name", s.Config.Name)

	if s.Config.Project && len(project) == 0 {
		p := config.Project{
			Name:       "project_id",
			UserPrompt: "Choose a project to use for this application.",
		}
		s.Config.Projects.Items = append(s.Config.Projects.Items, p)
	}

	if len(s.Config.Projects.Items) > 0 {

		currentProject := q.Get("currentProject").(string)

		for _, v := range s.Config.Projects.Items {
			s := newProjectSelector(v.Name, v.UserPrompt, currentProject, getProjects(q))
			c := newProjectCreator(v.Name + projNewSuffix)
			b := newBillingSelector(v.Name+billNewSuffix, getBillingAccounts(q), attachBilling)
			q.add(&s, &c, &b)
		}
	}

	if s.Config.BillingAccount {
		b := newBillingSelector("billing_account", getBillingAccounts(q), nil)
		b.list.Title = "Choose a billing account to use for with this application"
		q.add(&b)
	}

	if s.Config.ConfigureGCEInstance {
		newGCEInstance(q)
	}

	region = s.GetSetting("region")
	if s.Config.Region && len(region) == 0 {
		newRegion(q)
	}

	zone = s.GetSetting("zone")
	if s.Config.Zone && len(zone) == 0 {
		newZone(q)
	}

	if s.Config.Domain {
		newDomain(q)
	}

	newCustomPages(q)

	return err
}

func (q *Queue) add(m ...QueueModel) {
	uniques := map[string]bool{}

	for _, v := range q.models {
		uniques[v.getKey()] = true
	}

	for _, v := range m {
		// Basically if something dumb happens we don't rewrite queue
		// And since this is a author issue, not a user one, we should
		// fail silently
		if _, ok := uniques[v.getKey()]; ok {
			continue
		}

		v.addQueue(q)
		q.models = append(q.models, v)
		q.index = append(q.index, v.getKey())
	}
}

// method only used for testing
func (q *Queue) insert(m ...QueueModel) {
	tmp := q.models[:len(q.models)-1]
	tmp = append(tmp, m...)
	tmp = append(tmp, q.models[len(q.models)-1])

	q.models = tmp

}

// Start returns the first model to the hosting application so that it can
// be run through tea.NewProgram
func (q *Queue) Start() QueueModel {
	return q.models[0]
}
