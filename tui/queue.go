package tui

import (
	"github.com/GoogleCloudPlatform/deploystack"
	tea "github.com/charmbracelet/bubbletea"
)

// QueueModel is an extented version of tea.Model modified to work with the
// queue system
type QueueModel interface {
	tea.Model
	addQueue(*Queue)
	getKey() string
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
	stack   *deploystack.Stack
	store   map[string]interface{}
	index   []string
	client  UIClient
}

// NewQueue creates a new queue. You should need only one per app
func NewQueue(s *deploystack.Stack, client UIClient) Queue {
	q := Queue{stack: s, store: map[string]interface{}{}}
	q.client = client
	q.index = []string{}
	return q
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

func (q *Queue) currentKey() string {
	if len(q.models) == 0 {
		return ""
	}

	r := q.models[q.current].getKey()
	return r
}

func (q *Queue) nextKey() string {
	i := q.current + 1
	if i >= len(q.models) {
		return ""
	}
	r := q.models[i].getKey()
	return r
}

// InitializeUI spins up everything we need to have a working queue in the
// hosting application
func (q *Queue) InitializeUI() {
	desc := newDescription(q.stack)
	appHeader := newHeader(appTitle, q.stack.Config.Title)

	firstPage := newPage("firstpage", []component{newTextBlock(explainText)})
	descPage := newPage("descpage", []component{desc})

	endpage := newPage("endpage", []component{
		newTextBlock(titleStyle.Render("Project Settings")),
		newSettingsTable(q.stack),
	})
	endpage.addPostProcessor(cleanUp)

	q.header = appHeader
	q.add(&firstPage)
	q.add(&descPage)
	q.ProcessConfig()
	q.add(&endpage)
}

// ProcessConfig does the work of turning a DeployStack config file to a set
// of tui screens. It's separate from Initialize in case we want to be able
// to populate setting and variables with other information before running
// the genreation of those screens
func (q *Queue) ProcessConfig() error {
	var project, region, zone, projectnumber, name string
	var err error

	s := q.stack

	for i, v := range s.Config.HardSet {
		s.AddSetting(i, v)
	}

	project = s.GetSetting("project_id")
	region = s.GetSetting("region")
	zone = s.GetSetting("zone")
	name = s.Config.Name

	if name == "" {
		err = s.Config.ComputeName()
		if err != nil {
			return err
		}
	}
	s.AddSetting("stack_name", s.Config.Name)

	if s.Config.Project && len(project) == 0 {
		p := deploystack.Project{
			Name:       "project_id",
			UserPrompt: "Choose a project to use for this application.",
		}
		s.Config.Projects.Items = append(s.Config.Projects.Items, p)
	}

	if len(s.Config.Projects.Items) > 0 {
		for _, v := range s.Config.Projects.Items {

			projectsPage := newProjectSelector(
				v.Name,
				v.UserPrompt,
				getProjects(q),
			)
			q.add(&projectsPage)

			projectCreator := newProjectCreator(v.Name + "_new")
			q.add(projectCreator)

		}
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

	if s.Config.ProjectNumber {

		proj, err := deploystack.ProjectIDGet()
		if err != nil {
			return err
		}

		projectnumber, err = deploystack.ProjectNumberGet(proj)
		if err != nil {
			return err
		}
		s.AddSetting("project_number", projectnumber)
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
		// fails silently
		if _, ok := uniques[v.getKey()]; ok {
			continue
		}

		v.addQueue(q)
		q.models = append(q.models, v)
		q.index = append(q.index, v.getKey())
	}
}

// Start returns the first model to the hosting application so that it can
// be run through tea.NewProgram
func (q *Queue) Start() QueueModel {
	return q.models[0]
}
