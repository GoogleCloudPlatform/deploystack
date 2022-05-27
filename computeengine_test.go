package deploystack

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
	"testing"

	"google.golang.org/api/compute/v1"
)

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
				t.Fatalf("expected: no error, got: project-%s:%v", projectID, err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFormatMBToGB(t *testing.T) {
	tests := map[string]struct {
		input int64
		want  string
	}{
		"32 GB":  {input: 32768, want: "32 GB"},
		"240 GB": {input: 245760, want: "240 GB"},
		"192 GB": {input: 196608, want: "192 GB"},
		"16 GB":  {input: 16384, want: "16 GB"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := formatMBToGB(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

// func TestGetMachineTypes(t *testing.T) {
// 	uscaTypes, err := typefileHelper("test_files/types_uscentral1a.txt")
// 	if err != nil {
// 		t.Fatalf("got error during preloading: %s", err)
// 	}

// 	tests := map[string]struct {
// 		zone    string
// 		project string
// 		want    labeledValues
// 	}{
// 		"computeRegions": {zone: "us-central1-a", project: projectID, want: uscaTypes},
// 	}
// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			got, err := machineTypes(tc.project, tc.zone)
// 			if err != nil {
// 				t.Fatalf("expected: no error, got: %v", err)
// 			}

// 			got.sort()

// 			if !reflect.DeepEqual(tc.want, got) {
// 				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
// 			}
// 		})
// 	}
// }

func typefileHelper(file string) (LabeledValues, error) {
	result := LabeledValues{}
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return result, fmt.Errorf("unable to read region file (%s): %s", file, err)
	}

	temp := strings.Split(string(dat), "\n")

	for _, v := range temp {
		if v == "" {
			continue
		}

		items := strings.Split(v, "\t")
		nums := strings.Split(items[1], ".")
		lv := LabeledValue{items[0], fmt.Sprintf("%s CPUs: %s Mem: %s GB", items[0], items[2], nums[0])}
		result = append(result, lv)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Label < result[j].Label
	})

	return result, nil
}

func TestGetListOfDiskFamilies(t *testing.T) {
	tests := map[string]struct {
		input *compute.ImageList
		want  LabeledValues
	}{
		"DiskFamilies": {
			input: &compute.ImageList{
				Items: []*compute.Image{
					{Family: "windows-cloud"},
					{Family: "centos-cloud"},
					{Family: "centos-cloud"},
					{Family: "centos-cloud"},
					{Family: "centos-cloud"},
					{Family: "debian-cloud"},
				},
			},
			want: LabeledValues{
				LabeledValue{"centos-cloud", "centos-cloud"},
				LabeledValue{"debian-cloud", "debian-cloud"},
				LabeledValue{"windows-cloud", "windows-cloud"},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := getListOfImageFamilies(tc.input)

			got.sort()

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestGetListOfDiskTypes(t *testing.T) {
	tests := map[string]struct {
		input           *compute.ImageList
		family, project string
		want            LabeledValues
	}{
		"DiskFamilies": {
			input: &compute.ImageList{
				Items: []*compute.Image{
					{Family: "windows-cloud", Name: "windows-server"},
					{Family: "centos-server-pro", Name: "centos-server-1"},
					{Family: "centos-server-pro", Name: "centos-server-2"},
					{Family: "centos-server-pro", Name: "centos-server-3"},
					{Family: "centos-server-pro", Name: "centos-server-4"},
					{Family: "debian-cloud", Name: "debian-server"},
				},
			},
			family:  "centos-server-pro",
			project: "centos-cloud",
			want: LabeledValues{
				LabeledValue{"centos-cloud/centos-server-1", "centos-server-1"},
				LabeledValue{"centos-cloud/centos-server-2", "centos-server-2"},
				LabeledValue{"centos-cloud/centos-server-3", "centos-server-3"},
				LabeledValue{"centos-cloud/centos-server-4", "centos-server-4 (Latest)"},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := getListOfImageTypesByFamily(tc.input, tc.project, tc.family)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}
