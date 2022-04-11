package deploystack

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/kylelemons/godebug/diff"
)

var projectID = ""

func TestMain(m *testing.M) {
	var err error
	projectID, err = ProjectID()
	if err != nil {
		log.Fatalf("could not get environment project id: %s", err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestReadConfig(t *testing.T) {
	tests := map[string]struct {
		file string
		desc string
		want Stack
	}{
		"1": {
			file: "test_files/config.json",
			desc: "test_files/description.txt",
			want: Stack{
				Config: Config{
					Title:         "TESTCONFIG",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       true,
					Region:        true,
					RegionType:    "functions",
					RegionDefault: "us-central1",
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewStack()
			s.ReadConfig(tc.file, tc.desc)

			compareValues(tc.want.Config.Title, s.Config.Title, t)
			compareValues(tc.want.Config.Description, s.Config.Description, t)
			compareValues(tc.want.Config.Duration, s.Config.Duration, t)
			compareValues(tc.want.Config.Project, s.Config.Project, t)
			compareValues(tc.want.Config.Region, s.Config.Region, t)
			compareValues(tc.want.Config.RegionType, s.Config.RegionType, t)
			compareValues(tc.want.Config.RegionDefault, s.Config.RegionDefault, t)
		})
	}
}

func compareValues(want interface{}, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected: \n|%v|\ngot: \n|%v|", want, got)
	}
}

func TestManageZone(t *testing.T) {
	tests := map[string]struct {
		project string
		region  string
		want    string
	}{
		"1": {
			project: projectID,
			region:  "us-central1",
			want: `Enabling service to poll...
Polling for zones...
[1;36mChoose a valid zone to use for this application. [0m
[1;36m 1) us-central1-a  [0m
 2) us-central1-b 
 3) us-central1-c 
 4) us-central1-f 
Choose number from list, or just [enter] for [1;36mus-central1-a[0m
> `,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := captureOutput(func() {
				ZoneManage(tc.project, tc.region)
			})

			fmt.Println(diff.Diff(got, tc.want))
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: \n|%v|\ngot: \n|%v|", tc.want, got)
			}
		})
	}
}

func TestManageRegion(t *testing.T) {
	tests := map[string]struct {
		project string
		product string
		def     string
		want    string
		err     error
	}{
		"compute": {
			project: projectID,
			product: "compute",
			want: `Enabling service to poll...
Polling for regions...
[1;36mChoose a valid region to use for this application. [0m
 1) asia-east1              16) europe-west3            
 2) asia-east2              17) europe-west4            
 3) asia-northeast1         18) europe-west6            
 4) asia-northeast2         19) northamerica-northeast1 
 5) asia-northeast3         20) northamerica-northeast2 
 6) asia-south1             21) southamerica-east1      
 7) asia-south2             22) southamerica-west1      
 8) asia-southeast1         [1;36m23) us-central1              [0m
 9) asia-southeast2         24) us-east1                
10) australia-southeast1    25) us-east4                
11) australia-southeast2    26) us-west1                
12) europe-central2         27) us-west2                
13) europe-north1           28) us-west3                
14) europe-west1            29) us-west4                
15) europe-west2            
Choose number from list, or just [enter] for [1;36mus-central1[0m
> `,
			err: nil,
			def: "us-central1",
		},
		"functions": {
			project: projectID,
			product: "functions",
			want: `Enabling service to poll...
Polling for regions...
[1;36mChoose a valid region to use for this application. [0m
 1) asia-east1              13) europe-west3            
 2) asia-east2              14) europe-west6            
 3) asia-northeast1         15) northamerica-northeast1 
 4) asia-northeast2         16) southamerica-east1      
 5) asia-northeast3         [1;36m17) us-central1              [0m
 6) asia-south1             18) us-east1                
 7) asia-southeast1         19) us-east4                
 8) asia-southeast2         20) us-west1                
 9) australia-southeast1    21) us-west2                
10) europe-central2         22) us-west3                
11) europe-west1            23) us-west4                
12) europe-west2            
Choose number from list, or just [enter] for [1;36mus-central1[0m
> `,
			err: nil,
			def: "us-central1",
		},
		"run": {
			project: projectID,
			product: "run",
			want: `Enabling service to poll...
Polling for regions...
[1;36mChoose a valid region to use for this application. [0m
 1) asia-east1              16) europe-west3            
 2) asia-east2              17) europe-west4            
 3) asia-northeast1         18) europe-west6            
 4) asia-northeast2         19) northamerica-northeast1 
 5) asia-northeast3         20) northamerica-northeast2 
 6) asia-south1             21) southamerica-east1      
 7) asia-south2             22) southamerica-west1      
 8) asia-southeast1         [1;36m23) us-central1              [0m
 9) asia-southeast2         24) us-east1                
10) australia-southeast1    25) us-east4                
11) australia-southeast2    26) us-west1                
12) europe-central2         27) us-west2                
13) europe-north1           28) us-west3                
14) europe-west1            29) us-west4                
15) europe-west2            
Choose number from list, or just [enter] for [1;36mus-central1[0m
> `,
			err: nil,
			def: "us-central1",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := captureOutput(func() {
				RegionManage(tc.project, tc.product, tc.def)
			})

			fmt.Println(diff.Diff(got, tc.want))
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: \n|%v|\ngot: \n|%v|", tc.want, got)
			}
		})
	}
}

func TestSelectFromListRender(t *testing.T) {
	tests := map[string]struct {
		input []string
		def   string
		want  string
	}{
		"1": {
			input: []string{"one", "two", "three"},
			def:   "two",
			want: ` 1) one   
[1;36m 2) two    [0m
 3) three 
Choose number from list, or just [enter] for [1;36mtwo[0m
> `,
		},
		"2": {
			input: []string{"one", "two", "three", "four", "five", "six"},
			def:   "six",
			want: ` 1) one   
 2) two   
 3) three 
 4) four  
 5) five  
[1;36m 6) six    [0m
Choose number from list, or just [enter] for [1;36msix[0m
> `,
		},
		"3": {
			input: []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve"},
			def:   "six",
			want: ` 1) one     7) seven  
 2) two     8) eight  
 3) three   9) nine   
 4) four   10) ten    
 5) five   11) eleven 
[1;36m 6) six     [0m12) twelve 
Choose number from list, or just [enter] for [1;36msix[0m
> `,
		},
		"4": {
			input: []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven"},
			def:   "six",
			want: ` 1) one     7) seven  
 2) two     8) eight  
 3) three   9) nine   
 4) four   10) ten    
 5) five   11) eleven 
[1;36m 6) six     [0m
Choose number from list, or just [enter] for [1;36msix[0m
> `,
		},
		"5": {
			input: []string{
				"CREATE NEW PROJECT",
				"aiab-test-project",
				"appinabox-baslclb-tester",
				"appinabox-scaler-tester",
				"appinabox-todo-tester",
				"appinabox-yesornosite-demo",
				"appinabox-yesornosite-tester",
				"basiclb-test-project-delete",
				"basiclb-test-project2",
				"basiclb-tester-delete-me",
				"basiclb-tester-project",
				"bucketsite-test",
				"cloud-logging-generator",
				"cloudicons",
				"cost-sentry-experiments",
				"deploy-terraform",
				"deploystack-costsentry-tester",
				"deploystack-terraform",
				"deploystack-terraform-2",
				"microsites-appinabox",
				"microsites-deploystack",
				"microsites-stackables",
				"neos-log-test",
				"neos-test-304321",
				"neosregional",
				"nicholascain-starter-project",
				"scaler-microsite",
				"scaler-test-ui-delete",
				"stack-terraform",
				"stackinabox",
				"stackinaboxtester",
				"sustained-racer-323200",
				"todo-microsite",
				"vertexaitester",
				"zprojectnamedeletecbamp",
				"zprojectnamedeletefrzcl",
				"zprojectnamedeletehgzcu",
				"zprojectnamedeleteveday",
			},
			def: "stackinabox",
			want: ` 1) CREATE NEW PROJECT            20) microsites-appinabox          
 2) aiab-test-project             21) microsites-deploystack        
 3) appinabox-baslclb-tester      22) microsites-stackables         
 4) appinabox-scaler-tester       23) neos-log-test                 
 5) appinabox-todo-tester         24) neos-test-304321              
 6) appinabox-yesornosite-demo    25) neosregional                  
 7) appinabox-yesornosite-tester  26) nicholascain-starter-project  
 8) basiclb-test-project-delete   27) scaler-microsite              
 9) basiclb-test-project2         28) scaler-test-ui-delete         
10) basiclb-tester-delete-me      29) stack-terraform               
11) basiclb-tester-project        [1;36m30) stackinabox                    [0m
12) bucketsite-test               31) stackinaboxtester             
13) cloud-logging-generator       32) sustained-racer-323200        
14) cloudicons                    33) todo-microsite                
15) cost-sentry-experiments       34) vertexaitester                
16) deploy-terraform              35) zprojectnamedeletecbamp       
17) deploystack-costsentry-tester 36) zprojectnamedeletefrzcl       
18) deploystack-terraform         37) zprojectnamedeletehgzcu       
19) deploystack-terraform-2       38) zprojectnamedeleteveday       
Choose number from list, or just [enter] for [1;36mstackinabox[0m
> `,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := captureOutput(func() {
				listSelect(tc.input, tc.def)
			})

			fmt.Println(diff.Diff(got, tc.want))
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: \n|%v|\ngot: \n|%v|", tc.want, got)
			}
		})
	}
}

func TestStackTFvars(t *testing.T) {
	s := NewStack()
	s.AddSetting("project", "testproject")
	s.AddSetting("boolean", "true")
	got := s.Terraform()

	want := `boolean="true"
project="testproject"
`

	if got != want {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}

func TestTitle(t *testing.T) {
	tests := map[string]struct {
		title, desc, want string
		duration          int
	}{
		"Just1": {
			title: "test", desc: "test", duration: 1,
			want: `********************************************************************************
[1;36mtest[0m
test
It's going to take around [0;36m1 minute[0m
********************************************************************************
`,
		},
		"MoreThan1": {
			title: "test", desc: "test", duration: 2,
			want: `********************************************************************************
[1;36mtest[0m
test
It's going to take around [0;36m2 minutes[0m
********************************************************************************
`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewStack()
			s.Config.Title = tc.title
			s.Config.Description = tc.desc
			s.Config.Duration = tc.duration

			got := captureOutput(func() {
				s.Config.PrintHeader()
			})
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestSection(t *testing.T) {
	input := NewSection("test")

	got1 := captureOutput(func() {
		input.Open()
	})

	want1 := `********************************************************************************
[0;36mtest[0m
********************************************************************************
`

	if got1 != want1 {
		t.Fatalf("expected: %v, got: %v", want1, got1)
	}

	got2 := captureOutput(func() {
		input.Close()
	})

	want2 := `********************************************************************************
[0;36mtest - [1;36mdone[0m
********************************************************************************
`

	if got2 != want2 {
		t.Fatalf("expected: %v, got: %v", want2, got2)
	}
}

func TestStackPrintSettings(t *testing.T) {
	s := Stack{Settings: map[string]string{"zone": "test", "region": "test-a"}}

	got := captureOutput(func() {
		s.PrintSettings()
	})

	want := `[36mProject Details [0m 
Zone:   [1;36mtest[0m
Region: [1;36mtest-a[0m
`

	if got != want {
		t.Fatalf("expected: \n%v\n, got: \n%v\n", want, got)
	}
}

func captureOutput(f func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}

func TestGetRegions(t *testing.T) {
	tests := map[string]struct {
		product string
		project string
		want    []string
	}{
		"1": {product: "compute", project: projectID, want: []string{
			"asia-east1",
			"asia-east2",
			"asia-northeast1",
			"asia-northeast2",
			"asia-northeast3",
			"asia-south1",
			"asia-south2",
			"asia-southeast1",
			"asia-southeast2",
			"australia-southeast1",
			"australia-southeast2",
			"europe-central2",
			"europe-north1",
			"europe-west1",
			"europe-west2",
			"europe-west3",
			"europe-west4",
			"europe-west6",
			"northamerica-northeast1",
			"northamerica-northeast2",
			"southamerica-east1",
			"southamerica-west1",
			"us-central1",
			"us-east1",
			"us-east4",
			"us-west1",
			"us-west2",
			"us-west3",
			"us-west4",
		}},
		"2": {product: "functions", project: projectID, want: []string{
			"asia-east1",
			"asia-east2",
			"asia-northeast1",
			"asia-northeast2",
			"asia-northeast3",
			"asia-south1",
			"asia-southeast1",
			"asia-southeast2",
			"australia-southeast1",
			"europe-central2",
			"europe-west1",
			"europe-west2",
			"europe-west3",
			"europe-west6",
			"northamerica-northeast1",
			"southamerica-east1",
			"us-central1",
			"us-east1",
			"us-east4",
			"us-west1",
			"us-west2",
			"us-west3",
			"us-west4",
		}},
		"3": {product: "run", project: projectID, want: []string{
			"asia-east1",
			"asia-east2",
			"asia-northeast1",
			"asia-northeast2",
			"asia-northeast3",
			"asia-south1",
			"asia-south2",
			"asia-southeast1",
			"asia-southeast2",
			"australia-southeast1",
			"australia-southeast2",
			"europe-central2",
			"europe-north1",
			"europe-west1",
			"europe-west2",
			"europe-west3",
			"europe-west4",
			"europe-west6",
			"northamerica-northeast1",
			"northamerica-northeast2",
			"southamerica-east1",
			"southamerica-west1",
			"us-central1",
			"us-east1",
			"us-east4",
			"us-west1",
			"us-west2",
			"us-west3",
			"us-west4",
		}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := regions(tc.project, tc.product)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestGetZones(t *testing.T) {
	tests := map[string]struct {
		project string
		region  string
		want    []string
	}{
		"1": {project: projectID, region: "us-central1", want: []string{
			"us-central1-a",
			"us-central1-b",
			"us-central1-c",
			"us-central1-f",
		}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := zones(tc.project, tc.region)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestGetProjectID(t *testing.T) {
	tests := map[string]struct {
		want string
	}{
		"1": {want: projectID},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ProjectID()
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

// TODO: get rid of tests that are project specific
func TestGetProjectNumbers(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"1": {input: "stackinabox", want: "911905110276"},
		"2": {input: "stackinaboxtester", want: "130107067175"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ProjectNumber(tc.input)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

// TODO: get rid of tests that are project specific
func TestGetProjects(t *testing.T) {
	tests := map[string]struct {
		want []string
	}{
		"1": {want: []string{
			"deploystack-terraform-2",
			"deploystack-terraform",
			"deploy-terraform",
			"stack-terraform",
			"deploystack-costsentry-tester",
			"microsites-deploystack",
			"microsites-stackables",
			"vertexaitester",
			"stackinaboxtester",
			"stackinabox",
			"scaler-test-ui-delete",
			"basiclb-test-project-delete",
			"basiclb-tester-delete-me",
			"basiclb-tester-project",
			"basiclb-test-project2",
			"aiab-test-project",
			"appinabox-baslclb-tester",
			"cost-sentry-experiments",
			"appinabox-todo-tester",
			"appinabox-scaler-tester",
			"appinabox-yesornosite-demo",
			"appinabox-yesornosite-tester",
			"bucketsite-test",
			"microsites-appinabox",
			"scaler-microsite",
			"todo-microsite",
			"neosregional",
			"cloudicons",
			"sustained-racer-323200",
			"nicholascain-starter-project",
			"cloud-logging-generator",
			"neos-log-test",
			"neos-test-304321",
		}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := projects()

			gotfiltered := []string{}

			for _, v := range got {
				if !strings.Contains(v, "zprojectnamedelete") {
					gotfiltered = append(gotfiltered, v)
				}
			}

			sort.Strings(tc.want)
			sort.Strings(gotfiltered)

			if len(gotfiltered) != len(tc.want) {
				t.Fatalf("expected: %v, got: %v", len(tc.want), len(gotfiltered))
			}

			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, gotfiltered) {
				t.Fatalf("expected: %v, got: %v", tc.want, gotfiltered)
			}
		})
	}
}

// TODO: get rid of tests that are project specific
func TestGetBillingAccounts(t *testing.T) {
	tests := map[string]struct {
		want []string
	}{
		"1": {want: []string{"0145C0-557C58-C970F3"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := billingAccounts()

			sort.Strings(tc.want)

			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	tests := map[string]struct {
		input string
		err   error
	}{
		"1": {input: "zprojectnamedelete", err: nil},
		"2": {input: "zprojectnamedeletethisprojectnamehastoomanycharacters", err: ErrorProjectCreateTooLong},
		"3": {input: "ALLUPERCASEDONESTWORK", err: ErrorProjectInvalidCharacters},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			name := tc.input + randSeq(5)
			err := projectCreate(name)
			projectDelete(name)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v project: %s", tc.err, err, name)
			}
		})
	}
}

// TODO: get rid of tests that are project specific
func TestLinkProjectToBillingAccount(t *testing.T) {
	tests := map[string]struct {
		project string
		account string
		err     error
	}{
		"ShouldWork":  {project: projectID, account: "0145C0-557C58-C970F3", err: nil},
		"BadProject":  {project: "stackinaboxstackinabox", account: "0145C0-557C58-C970F3", err: ErrorBillingNoPermission},
		"BaddAccount": {project: projectID, account: "AAAAAA-BBBBBB-CCCCCC", err: ErrorBillingInvalidAccount},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := BillingAccountProjectAttach(tc.project, tc.account)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v", tc.err, err)
			}
		})
	}
}

func randSeq(n int) string {
	rand.Seed(time.Now().Unix())

	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func remove(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

// func TestDeleteProject(t *testing.T) {
// 	tests := map[string]struct {
// 		input string
// 	}{
// 		"1": {input: "aprojectthatisthirtycharacter"},
// 		"2": {input: "aprojecttodeleteplease"},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			err := DeleteProjectCall(tc.input)
// 			if err != nil {
// 				t.Fatalf("expected: no error, got: %v", err)
// 			}
// 		})
// 	}
// }
