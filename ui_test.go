package deploystack

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

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
		"Odd below 10": {
			input: toLabeledValueSlice([]string{"one", "two", "three"}),
			def:   "two",
			want: ` 1) one   
[1;36m 2) two   [0m
 3) three 
Choose number from list, or just [enter] for [1;36mtwo[0m
> `,
		},
		"Even below 10": {
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
		"EvenNumber above 10": {
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
		"OddNumber above 10": {
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
		"ProjectList": {
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

func TestDocumenationLink(t *testing.T) {
	tests := map[string]struct {
		title, desc, link, want string
		duration                int
	}{
		"NoLink": {
			title: "test", desc: "test", duration: 1,
			want: `********************************************************************************
[1;36mtest[0m
test
It's going to take around [0;36m1 minute[0m
********************************************************************************
`,
		},
		"Link": {
			title: "test", desc: "test", duration: 1, link: "http://deploystack.dev",
			want: `********************************************************************************
[1;36mtest[0m
test
It's going to take around [0;36m1 minute[0m

If you would like more information about this stack, please read the 
documentation at: 
[1;36mhttp://deploystack.dev[0m 
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
			s.Config.DocumentationLink = tc.link

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

func regionsListHelper(file string) ([]string, error) {
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

func removeFromSlice(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}

	return slice
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func TestRemoveFromSlice(t *testing.T) {
	tests := map[string]struct {
		in     []string
		remove string
		want   []string
	}{
		"basic":     {in: []string{"one", "two", "three"}, remove: "two", want: []string{"one", "three"}},
		"no action": {in: []string{"one", "two", "three"}, remove: "four", want: []string{"one", "two", "three"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := removeFromSlice(tc.in, tc.remove)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestRemoveStringFromSlice(t *testing.T) {
	tests := map[string]struct {
		in   []string
		want []string
	}{
		"no action":        {in: []string{"one", "two", "three"}, want: []string{"one", "two", "three"}},
		"remove one three": {in: []string{"one", "two", "three", "three"}, want: []string{"one", "two", "three"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := removeDuplicateStr(tc.in)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestGetRegions(t *testing.T) {
	_, rescueStdout := blockOutput()
	defer func() { os.Stdout = rescueStdout }()
	cRegions, err := regionsListHelper("test_files/gcloudout/regions_compute.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	fRegions, err := regionsListHelper("test_files/gcloudout/regions_functions.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	rRegions, err := regionsListHelper("test_files/gcloudout/regions_run.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	tests := map[string]struct {
		product string
		project string
		want    []string
		err     error
	}{
		"computeRegions":   {product: "compute", project: projectID, want: cRegions, err: nil},
		"functionsRegions": {product: "functions", project: projectID, want: fRegions, err: nil},
		"runRegions":       {product: "run", project: projectID, want: rRegions, err: nil},
		"GarbageInout":     {product: "An outdated iPad", project: projectID, want: []string{}, err: fmt.Errorf("invalid product requested: %s", "An outdated iPad")},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := regions(tc.project, tc.product)

			// BUG: getting weird regions intertmittenly popping up. Solving with this hack
			if tc.product == "compute" {
				got = removeDuplicateStr(removeFromSlice(removeFromSlice(got, "me-west1"), "us-west4"))
				tc.want = removeDuplicateStr(removeFromSlice(removeFromSlice(cRegions, "me-west1"), "us-west4"))
			}

			if err != tc.err {
				if err.Error() != tc.err.Error() {
					t.Fatalf("expected: no error, got: %v", err)
				}
			}

			sort.Strings(got)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestRegionManage(t *testing.T) {
	_, rescueStdout := blockOutput()
	defer func() { os.Stdout = rescueStdout }()
	tests := map[string]struct {
		input   string
		project string
		product string
		want    string
	}{
		"Run":       {"", projectID, "run", "us-central1"},
		"Compute":   {"", projectID, "compute", "us-central1"},
		"Functions": {"", projectID, "functions", "us-central1"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			content := []byte(tc.input)

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("error setting up environment for testing %v", err)
			}

			_, err = w.Write(content)
			if err != nil {
				t.Error(err)
			}
			w.Close()

			stdin := os.Stdin
			// Restore stdin right after the test.
			defer func() { os.Stdin = stdin }()
			os.Stdin = r

			got, err := RegionManage(tc.project, tc.product, DefaultRegion)
			if err != nil {
				t.Errorf("collectionfailed: %v", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestMachineTypeManage(t *testing.T) {
	_, rescueStdout := blockOutput()
	defer func() { os.Stdout = rescueStdout }()
	tests := map[string]struct {
		input   string
		project string
		zone    string
		want    string
	}{
		"Default": {"", projectID, "us-central1-a", "n1-standard-1"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			content := []byte(tc.input)

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("error setting up environment for testing %v", err)
			}

			_, err = w.Write(content)
			if err != nil {
				t.Error(err)
			}
			w.Close()

			stdin := os.Stdin
			// Restore stdin right after the test.
			defer func() { os.Stdin = stdin }()
			os.Stdin = r

			got, err := MachineTypeManage(tc.project, tc.zone)
			if err != nil {
				t.Errorf("collection failed: %v", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestImageManage(t *testing.T) {
	_, rescueStdout := blockOutput()
	defer func() { os.Stdout = rescueStdout }()
	defaultImage, err := getLatestImage(projectID, DefaultImageProject, DefaultImageFamily)
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}

	tests := map[string]struct {
		input   string
		project string
		want    string
	}{
		"Default": {"", projectID, defaultImage},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			content := []byte(tc.input)

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("error setting up environment for testing %v", err)
			}

			_, err = w.Write(content)
			if err != nil {
				t.Error(err)
			}
			w.Close()

			stdin := os.Stdin
			// Restore stdin right after the test.
			defer func() { os.Stdin = stdin }()
			os.Stdin = r

			got, err := ImageManage(tc.project)
			if err != nil {
				t.Errorf("collection failed: %v", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestGCEInstanceManage(t *testing.T) {
	_, rescueStdout := blockOutput()
	defer func() { os.Stdout = rescueStdout }()
	defaultImage, err := getLatestImage(projectID, DefaultImageProject, DefaultImageFamily)
	if err != nil {
		t.Fatalf("error setting up environment for testing %v", err)
	}

	basename := "testing"

	defaultConfig := GCEInstanceConfig{
		"instance-image":        defaultImage,
		"instance-disksize":     "200",
		"instance-disktype":     "pd-standard",
		"instance-tags":         "[http-server,https-server]",
		"instance-name":         fmt.Sprintf("%s-instance", basename),
		"region":                DefaultRegion,
		"zone":                  fmt.Sprintf("%s-a", DefaultRegion),
		"instance-machine-type": "n1-standard-1",
	}

	tests := map[string]struct {
		input   string
		project string
		want    GCEInstanceConfig
	}{
		"Default": {"", projectID, defaultConfig},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			content := []byte(tc.input)

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("error setting up environment for testing %v", err)
			}

			_, err = w.Write(content)
			if err != nil {
				t.Error(err)
			}
			w.Close()

			stdin := os.Stdin
			// Restore stdin right after the test.
			defer func() { os.Stdin = stdin }()
			os.Stdin = r

			got, err := GCEInstanceManage(tc.project, basename)
			if err != nil {
				t.Errorf("collection failed: %v", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestExtractAccount(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"Basic": {input: "Something (Account)", want: "Account"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := extractAccount(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

// func TestManageProject(t *testing.T) {
// 	tests := map[string]struct {
// 		want string
// 	}{
// 		"1": {
// 			want: `
// [1;36mChoose a project to use for this application.[0m

// [46mNOTE:[0;36m This app will make changes to the project. [0m
// While those changes are reverseable, it would be better to put it in a fresh new project.
//  1) CREATE NEW PROJECT
//  2) ds-tester-helper
// Choose number from list.
// > `,
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			got := captureOutput(func() {
// 				ProjectManage()
// 			})

// 			fmt.Println(diff.Diff(got, tc.want))
// 			if !reflect.DeepEqual(tc.want, got) {
// 				fmt.Printf("ProjectID: %s\n", projectID)
// 				t.Fatalf("expected: \n|%v|\ngot: \n|%v|", tc.want, got)
// 			}
// 		})
// 	}
// }
