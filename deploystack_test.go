// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploystack

import (
	"encoding/json"
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
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/option"
)

var (
	projectID = ""
	creds     map[string]string
)

func TestMain(m *testing.M) {
	var err error
	opts = option.WithCredentialsFile("creds.json")

	dat, err := os.ReadFile("creds.json")
	if err != nil {
		log.Fatalf("unable to handle the json config file: %v", err)
	}

	json.Unmarshal(dat, &creds)

	projectID = creds["project_id"]
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
		err  error
	}{
		"error": {
			file: "z.json",
			desc: "z.txt",
			want: Stack{},
			err:  fmt.Errorf("unable to read config file: open z.json: no such file or directory"),
		},
		"no_custom": {
			file: "test_files/no_customs/deploystack.json",
			desc: "test_files/no_customs/deploystack.txt",
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
			err: nil,
		},
		"custom": {
			file: "test_files/customs/deploystack.json",
			desc: "test_files/customs/deploystack.txt",
			want: Stack{
				Config: Config{
					Title:         "TESTCONFIG",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       false,
					Region:        false,
					RegionType:    "",
					RegionDefault: "",
					CustomSettings: []Custom{
						{Name: "nodes", Description: "Nodes", Default: "3"},
					},
				},
			},
			err: nil,
		},
		"custom_options": {
			file: "test_files/customs_options/deploystack.json",
			desc: "test_files/customs_options/deploystack.txt",
			want: Stack{
				Config: Config{
					Title:         "TESTCONFIG",
					Description:   "A test string for usage with this stuff.",
					Duration:      5,
					Project:       false,
					Region:        false,
					RegionType:    "",
					RegionDefault: "",

					CustomSettings: []Custom{
						{
							Name:        "nodes",
							Description: "Nodes",
							Default:     "3",
							Options:     []string{"1", "2", "3"},
						},
					},
				},
			},
			err: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewStack()
			err := s.ReadConfig(tc.file, tc.desc)

			if err != tc.err {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

			compareValues(tc.want.Config.Title, s.Config.Title, t)
			compareValues(tc.want.Config.Description, s.Config.Description, t)
			compareValues(tc.want.Config.Duration, s.Config.Duration, t)
			compareValues(tc.want.Config.Project, s.Config.Project, t)
			compareValues(tc.want.Config.Region, s.Config.Region, t)
			compareValues(tc.want.Config.RegionType, s.Config.RegionType, t)
			compareValues(tc.want.Config.RegionDefault, s.Config.RegionDefault, t)
			for i, v := range s.Config.CustomSettings {
				compareValues(tc.want.Config.CustomSettings[i], v, t)
			}
		})
	}
}

func compareValues(want interface{}, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected: \n|%v|\ngot: \n|%v|", want, got)
	}
}

func TestProcessCustoms(t *testing.T) {
	tests := map[string]struct {
		file string
		desc string
		want string
		err  error
	}{
		"custom_options": {
			file: "test_files/customs_options/deploystack.json",
			desc: "test_files/customs_options/deploystack.txt",
			want: `********************************************************************************
[1;36mTESTCONFIG[0m
A test string for usage with this stuff.
It's going to take around [0;36m5 minutes[0m
********************************************************************************
[1;36mNodes: [0m
 1) 1 
 2) 2 
[1;36m 3) 3 [0m
Choose number from list, or just [enter] for [1;36m3[0m
> 
[46mProject Details [0m 
Nodes: [1;36m3[0m
`,
			err: nil,
		},
		"custom": {
			file: "test_files/customs/deploystack.json",
			desc: "test_files/customs/deploystack.txt",
			want: `********************************************************************************
[1;36mTESTCONFIG[0m
A test string for usage with this stuff.
It's going to take around [0;36m5 minutes[0m
********************************************************************************
[1;36mNodes: [0m
Enter value, or just [enter] for [1;36m3[0m
> 
[46mProject Details [0m 
Nodes: [1;36m3[0m
`,
			err: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewStack()
			err := s.ReadConfig(tc.file, tc.desc)

			if err != tc.err {
				if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
					t.Fatalf("expected: error(%s) got: error(%s)", tc.err, err)
				}
			}

			got := captureOutput(func() {
				if err := s.Process("terraform.tfvars"); err != nil {
					log.Fatalf("problemn collecting the configurations: %s", err)
				}
			})

			if !reflect.DeepEqual(tc.want, got) {
				fmt.Println(diff.Diff(got, tc.want))
				t.Fatalf("expected: \n|%v|\ngot: \n|%v|", tc.want, got)
			}
		})
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
			want: `Polling for zones...
[1;36mChoose a valid zone to use for this application. [0m
[1;36m 1) us-central1-a [0m
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
				fmt.Printf("ProjectID: %s\n", projectID)
				t.Fatalf("expected: \n|%v|\ngot: \n|%v|", tc.want, got)
			}
		})
	}
}

func TestSelectFromListRender(t *testing.T) {
	tests := map[string]struct {
		input LabeledValues
		def   string
		want  string
	}{
		"1": {
			input: toLabeledValueSlice([]string{"one", "two", "three"}),
			def:   "two",
			want: ` 1) one   
[1;36m 2) two   [0m
 3) three 
Choose number from list, or just [enter] for [1;36mtwo[0m
> `,
		},
		"2": {
			input: toLabeledValueSlice([]string{"one", "two", "three", "four", "five", "six"}),
			def:   "six",
			want: ` 1) one   
 2) two   
 3) three 
 4) four  
 5) five  
[1;36m 6) six   [0m
Choose number from list, or just [enter] for [1;36msix[0m
> `,
		},
		"3": {
			input: toLabeledValueSlice([]string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve"}),
			def:   "six",
			want: ` 1) one     7) seven  
 2) two     8) eight  
 3) three   9) nine   
 4) four   10) ten    
 5) five   11) eleven 
[1;36m 6) six    [0m12) twelve 
Choose number from list, or just [enter] for [1;36msix[0m
> `,
		},
		"4": {
			input: toLabeledValueSlice([]string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven"}),
			def:   "six",
			want: ` 1) one     7) seven  
 2) two     8) eight  
 3) three   9) nine   
 4) four   10) ten    
 5) five   11) eleven 
[1;36m 6) six    [0m
Choose number from list, or just [enter] for [1;36msix[0m
> `,
		},
		"5": {
			input: toLabeledValueSlice([]string{
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
			}),
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
11) basiclb-tester-project        [1;36m30) stackinabox                   [0m
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

			if !reflect.DeepEqual(tc.want, got) {
				fmt.Println(diff.Diff(got, tc.want))
				t.Fatalf("expected: \n|%v|\ngot: \n|%v|", tc.want, got)
			}
		})
	}
}

func TestStackTFvars(t *testing.T) {
	s := NewStack()
	s.AddSetting("project", "testproject")
	s.AddSetting("boolean", "true")
	s.AddSetting("set", "[item1,item2]")
	got := s.Terraform()

	want := `boolean="true"
project="testproject"
set=["item1","item2"]
`

	if got != want {
		fmt.Println(diff.Diff(want, got))
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

	want := `
[46mProject Details [0m 
Region: [1;36mtest-a[0m
Zone:   [1;36mtest[0m
`

	if got != want {
		t.Fatalf("expected: \n|%v|\n, got: \n|%v|\n", want, got)
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

func listHelper(file string) ([]string, error) {
	result := []string{}
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return result, fmt.Errorf("unable to read region file (%s): %s", file, err)
	}

	temp := strings.Split(string(dat), "\n")

	for _, v := range temp {
		if v == "" {
			continue
		}
		full := strings.Split(v, "/")
		result = append(result, strings.TrimSpace(full[len(full)-1]))
	}

	sort.Strings(result)

	return result, nil
}

func TestGetRegions(t *testing.T) {
	cRegions, err := listHelper("test_files/regions_compute.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	fRegions, err := listHelper("test_files/regions_functions.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	rRegions, err := listHelper("test_files/regions_run.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	tests := map[string]struct {
		product string
		project string
		want    []string
	}{
		"computeRegions":   {product: "compute", project: projectID, want: cRegions},
		"functionsRegions": {product: "functions", project: projectID, want: fRegions},
		"runRegions":       {product: "run", project: projectID, want: rRegions},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := regions(tc.project, tc.product)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}

			sort.Strings(got)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestGetProjectNumbers(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"1": {input: creds["project_id"], want: creds["project_number"]},
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

func TestGetProjects(t *testing.T) {
	tests := map[string]struct {
		want []string
	}{
		"1": {want: []string{
			creds["project_id"],
		}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := Projects()

			gotfiltered := []string{}

			for _, v := range got {
				if !strings.Contains(v, "zprojectnamedelete") {
					gotfiltered = append(gotfiltered, v)
				}
			}

			sort.Strings(tc.want)
			sort.Strings(gotfiltered)

			if len(gotfiltered) != len(tc.want) {
				fmt.Printf("Expected:\n%s", tc.want)
				fmt.Printf("Got:\n%s", gotfiltered)
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

func TestGetBillingAccounts(t *testing.T) {
	tests := map[string]struct {
		want []*cloudbilling.BillingAccount
	}{
		"NoErrorNoAccounts": {want: []*cloudbilling.BillingAccount{}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := billingAccounts()

			sort.Slice(got[:], func(i, j int) bool {
				return got[i].DisplayName < got[j].DisplayName
			})

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
		"Too long":  {input: "zprojectnamedeletethisprojectnamehastoomanycharacters", err: ErrorProjectCreateTooLong},
		"Bad Chars": {input: "ALLUPERCASEDONESTWORK", err: ErrorProjectInvalidCharacters},
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

func TestLinkProjectToBillingAccount(t *testing.T) {
	tests := map[string]struct {
		project string
		account string
		err     error
	}{
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

func TestMassgePhoneNumber(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
		err   error
	}{
		"Good":  {"800 555 1234", "+1.8005551234", nil},
		"Weird": {"d746fd83843", "+1.74683843", nil},
		"BAD":   {"dghdhdfuejfhfhfhrghfhfhdhgreh", "", ErrorCustomNotValidPhoneNumber},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := massagePhoneNumber(tc.input)
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v", tc.err, err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
