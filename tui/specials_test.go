package tui

import (
	"testing"

	"github.com/GoogleCloudPlatform/deploystack"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/api/cloudbilling/v1"
)

func TestNewProjectCreator(t *testing.T) {
	tests := map[string]struct {
		key        string
		outputFile string
	}{
		"basic": {
			key:        "project_id",
			outputFile: "testdata/project_creator_basic.txt",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			out := newProjectCreator(tc.key)
			q.add(&out)

			got := out.View()
			want := readTestFile(tc.outputFile)

			if want != got {
				writeDebugFile(got, tc.outputFile)
				t.Fatalf("text wasn't the same")
			}
		})
	}
}

func TestNewProjectSelector(t *testing.T) {
	tests := map[string]struct {
		key          string
		listLabel    string
		preProcessor tea.Cmd
		outputFile   string
		update       bool
	}{
		"waiting": {
			key:        "project_id",
			listLabel:  "Selecte a project to use",
			outputFile: "testdata/project_selector_waiting.txt",
		},
		"updated": {
			key:        "project_id",
			listLabel:  "Selecte a project to use",
			outputFile: "testdata/project_selector_updated.txt",
			update:     true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")

			out := newProjectSelector(tc.key, tc.listLabel, getProjects(&q))
			q.add(&out)

			if tc.update {
				cmd := out.Init()
				for i := 0; i < 2; i++ {

					msg := cmd()

					switch v := msg.(type) {
					case tea.BatchMsg:
						msgs := msg.(tea.BatchMsg)

						for _, v2 := range msgs {
							var tmp tea.Model
							tmp, cmd = out.Update(v2())
							out = tmp.(picker)
						}
					default:
						var tmp tea.Model
						tmp, cmd = out.Update(v)
						out = tmp.(picker)
					}

				}

			}

			got := out.View()
			want := readTestFile(tc.outputFile)

			if want != got {
				writeDebugFile(got, tc.outputFile)
				t.Fatalf("text wasn't the same")
			}
		})
	}
}

func TestNewBillingSelector(t *testing.T) {
	tests := map[string]struct {
		key        string
		outputFile string
		state      string
		single     bool
	}{
		"basic": {
			key:        "billing_account",
			outputFile: "testdata/billing_selector_basic.txt",
			state:      "idle",
		},
		"displaying": {
			key:        "project_id",
			outputFile: "testdata/project_selector_displaying.txt",
			state:      "displaying",
		},
		"displaying_single": {
			key:        "project_id",
			outputFile: "testdata/project_selector_displaying_single.txt",
			state:      "displaying",
			single:     true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")

			out := newBillingSelector(tc.key, getBillingAccounts(&q), nil)
			q.add(&out)
			p := newBillingSelector("dummy", getBillingAccounts(&q), nil)
			p.spinnerLabel = "dummy"
			q.add(&p)

			if tc.single {
				m := mock{}
				baList := []*cloudbilling.BillingAccount{
					{
						DisplayName: "Very Limted Funds",
						Name:        "billingAccounts/000000-000000-00000Y",
					},
				}
				m.save("BillingAccountList", baList)

				q.client = m
			}

			if tc.state == "displaying" {
				cmd := out.Init()
				for i := 0; i < 2; i++ {

					msg := cmd()

					switch v := msg.(type) {
					case tea.BatchMsg:
						msgs := msg.(tea.BatchMsg)

						for _, v2 := range msgs {
							var tmp tea.Model
							tmp, cmd = out.Update(v2())

							switch p := tmp.(type) {
							case picker:
								out = p
							case *picker:
								out = *p
							}

						}
					default:
						var tmp tea.Model
						tmp, cmd = out.Update(v)
						switch p := tmp.(type) {
						case picker:
							out = p
						case *picker:
							out = *p
						}
					}

				}

			}

			got := out.View()
			want := readTestFile(tc.outputFile)

			if want != got {
				writeDebugFile(got, tc.outputFile)
				t.Fatalf("text wasn't the same")
			}
		})
	}
}

func TestNewProjectFlow(t *testing.T) {
	tests := map[string]struct {
		want             string
		createNewProject bool
	}{
		"createProject": {
			want:             "project_id" + projNewSuffix,
			createNewProject: true,
		},
		"doNotCreateProject": {
			want:             "dummy",
			createNewProject: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			key := "project_id"

			q := getTestQueue(appTitle, "test")
			p1 := newProjectSelector(key, "", getProjects(&q))
			p2 := newProjectCreator(key + projNewSuffix)
			p3 := newPage("dummy", []component{})
			q.add(&p1, &p2, &p3)

			p := q.Start()
			if !tc.createNewProject {
				q.stack.AddSetting(key, "nonnilvalue")
				p2.value = "nonnilvalue"
				p2.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
			}

			tmp, _ := q.next()
			p = tmp.(QueueModel)

			got := p.getKey()

			if tc.want != got {
				t.Fatalf("want '%s' got '%s'", tc.want, got)
			}
		})
	}
}

func TestNewCustom(t *testing.T) {
	tests := map[string]struct {
		c          deploystack.Custom
		outputFile string
	}{
		"basic": {
			c: deploystack.Custom{
				Name:        "test",
				Description: "A test option",
				Default:     "Test",
			},
			outputFile: "testdata/custom_basic.txt",
		},
		"phone": {
			c: deploystack.Custom{
				Name:        "test",
				Description: "A test phone",
				Default:     "1-555-555-4040",
				Validation:  validationPhoneNumber,
			},
			outputFile: "testdata/custom_phone.txt",
		},
		"yesorno": {
			c: deploystack.Custom{
				Name:        "test",
				Description: "Yay or Nay",
				Default:     "Yes",
				Validation:  validationYesOrNo,
			},
			outputFile: "testdata/custom_yesorno.txt",
		},
		"integer": {
			c: deploystack.Custom{
				Name:        "test",
				Description: "a number",
				Default:     "5",
				Validation:  validationInteger,
			},
			outputFile: "testdata/custom_integer.txt",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			out := newCustom(tc.c)
			q.add(out)

			got := out.View()
			want := readTestFile(tc.outputFile)

			if want != got {
				writeDebugFile(got, tc.outputFile)
				t.Fatalf("text wasn't the same")
			}
		})
	}
}

func TestQueueBatch(t *testing.T) {
	tests := map[string]struct {
		f     func(*Queue)
		count int
		keys  []string
	}{
		"region": {
			f:     newRegion,
			count: 1,
			keys:  []string{"region"},
		},
		"zone": {
			f:     newZone,
			count: 1,
			keys:  []string{"zone"},
		},

		"domain": {
			f:     newDomain,
			count: 10,
			keys: []string{
				"domain",
				"domain_email",
				"domain_phone",
				"domain_country",
				"domain_postalcode",
				"domain_state",
				"domain_city",
				"domain_address",
				"domain_name",
				"domain_consent",
			},
		},

		"GCEInstance": {
			f:     newGCEInstance,
			count: 12,
			keys: []string{
				"gce-use-defaults",
				"instance-name",
				"region",
				"zone",
				"instance-machine-type-family",
				"instance-machine-type",
				"instance-image-project",
				"instance-image-family",
				"instance-image",
				"instance-disktype",
				"instance-disksize",
				"instance-webserver",
			},
		},
		"MachineTypeManager": {
			f:     newMachineTypeManager,
			count: 2,
			keys: []string{
				"instance-machine-type-family",
				"instance-machine-type",
			},
		},

		"DiskImageManager": {
			f:     newDiskImageManager,
			count: 3,
			keys: []string{
				"instance-image-project",
				"instance-image-family",
				"instance-image",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")
			tc.f(&q)

			if tc.count != len(q.models) {
				t.Fatalf("count - want '%d' got '%d'", tc.count, len(q.models))
			}

			for _, v := range tc.keys {
				q.removeModel(v)
			}

			if 0 != len(q.models) {
				t.Logf("Models remain")
				for _, v := range q.models {
					t.Logf("%s", v.getKey())
				}

				t.Fatalf("key check - want '%d' got '%d'", 0, len(q.models))

			}
		})
	}
}

func TestCustomPages(t *testing.T) {
	tests := map[string]struct {
		config string
		count  int
		keys   []string
	}{
		"region": {
			config: "testdata/config_multicustom.yaml",
			count:  5,
			keys: []string{
				"nodes",
				"label",
				"location",
				"budgetamount",
				"yesorno",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			q := getTestQueue(appTitle, "test")

			s := readTestFile(tc.config)

			config, err := deploystack.NewConfigYAML([]byte(s))
			if err != nil {
				t.Fatalf("could not read in config %s:", err)
			}
			q.stack.Config = config

			newCustomPages(&q)

			if tc.count != len(q.models) {
				t.Logf("Models ")
				for _, v := range q.models {
					t.Logf("%s", v.getKey())
				}
				t.Fatalf("count - want '%d' got '%d'", tc.count, len(q.models))
			}

			for _, v := range tc.keys {
				q.removeModel(v)
			}

			if len(q.models) != 0 {
				t.Logf("Models remain")
				for _, v := range q.models {
					t.Logf("%s", v.getKey())
				}

				t.Fatalf("key check - want '%d' got '%d'", 0, len(q.models))

			}
		})
	}
}
