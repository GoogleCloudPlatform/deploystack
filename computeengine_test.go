package deploystack

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
	"testing"
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

func TestGetMachineTypes(t *testing.T) {
	uscaTypes, err := typefileHelper("test_files/types_uscentral1a.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	tests := map[string]struct {
		zone    string
		project string
		want    labeledValues
	}{
		"computeRegions": {zone: "us-central1-a", project: projectID, want: uscaTypes},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := machineTypes(tc.project, tc.zone)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}

			got.sort()

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func typefileHelper(file string) (labeledValues, error) {
	result := labeledValues{}
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
		lv := labeledValue{items[0], fmt.Sprintf("%s CPUs: %s Mem: %s GB", items[0], items[2], nums[0])}
		result = append(result, lv)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].label < result[j].label
	})

	return result, nil
}
