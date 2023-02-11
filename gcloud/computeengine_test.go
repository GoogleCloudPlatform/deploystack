package gcloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/deploystack"
	"google.golang.org/api/compute/v1"
)

func TestGetComputeRegions(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	cRegions, err := regionsListHelper("test_files/gcloudout/regions_compute.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	tests := map[string]struct {
		project string
		want    []string
	}{
		"computeRegions": {projectID, cRegions},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.ComputeRegionList(tc.project)
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

func TestZones(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
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
			got, err := c.ZoneList(tc.project, tc.region)
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
	c := NewClient(ctx, defaultUserAgent)
	uscaTypes, err := typefileHelper("test_files/gcloudout/types_uscentral1a.txt")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	tests := map[string]struct {
		zone    string
		project string
		want    *compute.MachineTypeList
	}{
		"computeRegions": {zone: "us-central1-a", project: projectID, want: uscaTypes},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.MachineTypeList(tc.project, tc.zone)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}

			sort.Slice(got.Items, func(i, j int) bool {
				return got.Items[i].Name < got.Items[j].Name
			})

			for i, v := range got.Items {
				want := tc.want.Items[i]

				if !reflect.DeepEqual(want.Name, v.Name) {
					t.Fatalf("%s: expected: %+v, got: %+v", v.Name, want.Name, v.Name)
				}

				if !reflect.DeepEqual(want.GuestCpus, v.GuestCpus) {
					t.Fatalf("%s: expected: %+v, got: %+v", v.Name, want.GuestCpus, v.GuestCpus)
				}

				if !closeEnough(want.MemoryMb, v.MemoryMb, 1) {
					t.Fatalf("%s: expected: %+v, got: %+v", v.Name, want.MemoryMb, v.MemoryMb)
				}
			}
		})
	}
}

func closeEnough(int1, int2, threshold int64) bool {
	return math.Abs(float64(int1)-float64(int2)) <= float64(threshold)
}

func typefileHelper(file string) (*compute.MachineTypeList, error) {
	result := &compute.MachineTypeList{}
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read region file (%s): %s", file, err)
	}

	temp := strings.Split(string(dat), "\n")

	for _, v := range temp {
		if v == "" {
			continue
		}

		items := strings.Split(v, "\t")
		name := items[0]
		procs, err := strconv.Atoi(items[2])
		if err != nil {
			return nil, err
		}

		mem, err := strconv.ParseFloat(items[1], 64)
		if err != nil {
			return nil, err
		}

		mt := compute.MachineType{
			Name:      name,
			GuestCpus: int64(procs),
			MemoryMb:  int64(mem * 1024),
		}
		result.Items = append(result.Items, &mt)
	}

	sort.Slice(result.Items, func(i, j int) bool {
		return result.Items[i].Name < result.Items[j].Name
	})

	return result, nil
}

func TestGetListOfDiskFamilies(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		input *compute.ImageList
		want  deploystack.LabeledValues
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
			want: deploystack.LabeledValues{
				deploystack.LabeledValue{
					Value:     "centos-cloud",
					Label:     "centos-cloud",
					IsDefault: false,
				},

				deploystack.LabeledValue{
					Value:     "debian-cloud",
					Label:     "debian-cloud",
					IsDefault: false,
				},

				deploystack.LabeledValue{
					Value:     "windows-cloud",
					Label:     "windows-cloud",
					IsDefault: false,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := c.ImageFamilyList(tc.input)

			got.Sort()

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestGetListOfImageTypesByFamily(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		input           *compute.ImageList
		family, project string
		want            deploystack.LabeledValues
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
			want: deploystack.LabeledValues{
				deploystack.LabeledValue{
					Value:     "centos-cloud/centos-server-1",
					Label:     "centos-server-1",
					IsDefault: false,
				},
				deploystack.LabeledValue{
					Value:     "centos-cloud/centos-server-2",
					Label:     "centos-server-2",
					IsDefault: false,
				},
				deploystack.LabeledValue{
					Value:     "centos-cloud/centos-server-3",
					Label:     "centos-server-3",
					IsDefault: false,
				},
				deploystack.LabeledValue{
					Value:     "centos-cloud/centos-server-4",
					Label:     "centos-server-4 (Latest)",
					IsDefault: true,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := c.ImageTypeListByFamily(tc.input, tc.project, tc.family)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestGetListOfMachineeTypesByFamily(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		input  *compute.MachineTypeList
		family string
		want   deploystack.LabeledValues
	}{
		"DiskFamilies": {
			input: &compute.MachineTypeList{
				Items: []*compute.MachineType{
					{Name: "n1-standard-1", Description: "1 Proc"},
					{Name: "n1-standard-2", Description: "2 Proc"},
					{Name: "n1-standard-4", Description: "4 Proc"},
					{Name: "n1-standard-8", Description: "8 Proc"},
					{Name: "n1-standard-16", Description: "16 Proc"},
					{Name: "n1-standard-32", Description: "32 Proc"},
					{Name: "n1-highmem-32", Description: "32 Proc"},
					{Name: "a1-highmem-32", Description: "32 Proc"},
				},
			},
			family: "n1-standard",
			want: deploystack.LabeledValues{
				deploystack.LabeledValue{
					Value:     "n1-standard-1",
					Label:     "n1-standard-1 1 Proc",
					IsDefault: true,
				},
				deploystack.LabeledValue{
					Value:     "n1-standard-2",
					Label:     "n1-standard-2 2 Proc",
					IsDefault: false,
				},
				deploystack.LabeledValue{
					Value:     "n1-standard-4",
					Label:     "n1-standard-4 4 Proc",
					IsDefault: false,
				},
				deploystack.LabeledValue{
					Value:     "n1-standard-8",
					Label:     "n1-standard-8 8 Proc",
					IsDefault: false,
				},
				deploystack.LabeledValue{
					Value:     "n1-standard-16",
					Label:     "n1-standard-16 16 Proc",
					IsDefault: false,
				},
				deploystack.LabeledValue{
					Value:     "n1-standard-32",
					Label:     "n1-standard-32 32 Proc",
					IsDefault: false,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := c.MachineTypeListByFamily(tc.input, tc.family)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestGetListOfMachineTypeFamily(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	tests := map[string]struct {
		input *compute.MachineTypeList
		want  deploystack.LabeledValues
	}{
		"DiskFamilies": {
			input: &compute.MachineTypeList{
				Items: []*compute.MachineType{
					{Name: "n1-standard-1", Description: "1 Proc"},
					{Name: "n1-standard-2", Description: "2 Proc"},
					{Name: "n1-standard-4", Description: "4 Proc"},
					{Name: "n1-standard-8", Description: "8 Proc"},
					{Name: "n1-standard-16", Description: "16 Proc"},
					{Name: "n1-standard-32", Description: "32 Proc"},
					{Name: "n1-highmem-32", Description: "32 Proc"},
					{Name: "a1-highmem-32", Description: "32 Proc"},
				},
			},
			want: deploystack.LabeledValues{
				deploystack.LabeledValue{
					Value:     "n1-standard",
					Label:     "n1 standard",
					IsDefault: false,
				},

				deploystack.LabeledValue{
					Value:     "n1-highmem",
					Label:     "n1 highmem",
					IsDefault: false,
				},

				deploystack.LabeledValue{
					Value:     "a1-highmem",
					Label:     "a1 highmem",
					IsDefault: false,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := c.MachineTypeFamilyList(tc.input)

			tc.want.Sort()

			for i, v := range got {
				if !reflect.DeepEqual(tc.want[i].Value, v.Value) {
					t.Fatalf("Value expected: %+v, got: %+v", tc.want[i].Value, v.Value)
				}

				if !reflect.DeepEqual(tc.want[i].Label, v.Label) {
					t.Fatalf("Label expected: %+v, got: %+v", tc.want[i].Label, v.Label)
				}
			}
		})
	}
}

func getImageByProjectFromFile(imgs []*compute.Image, imageproject string) []*compute.Image {
	result := []*compute.Image{}
	for _, v := range imgs {
		if strings.Contains(v.SelfLink, fmt.Sprintf("/%s/", imageproject)) {
			result = append(result, v)
		}
	}

	return result
}

func getLatestImageByProjectFromFile(imgs []*compute.Image, imageproject, imagefamily string) string {
	result := []*compute.Image{}
	for _, v := range imgs {
		if strings.Contains(v.SelfLink, fmt.Sprintf("/%s/", imageproject)) {
			result = append(result, v)
		}
	}

	result2 := []*compute.Image{}
	for _, v := range imgs {
		if v.Family == imagefamily {
			result2 = append(result2, v)
		}
	}

	answer := fmt.Sprintf("%s/%s", imageproject, result2[len(result2)-1].Name)

	return answer
}

func TestImages(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	dat, err := os.ReadFile("test_files/gcloudout/images.json")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	imgs := []*compute.Image{}
	err = json.Unmarshal(dat, &imgs)
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	tests := map[string]struct {
		project      string
		imageProject string
		want         []*compute.Image
	}{
		"debian": {projectID, "debian-cloud", getImageByProjectFromFile(imgs, "debian-cloud")},
		"rhel":   {projectID, "rhel-cloud", getImageByProjectFromFile(imgs, "rhel-cloud")},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.ImageList(tc.project, tc.imageProject)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}

			// sort.Strings(got)

			if !reflect.DeepEqual(tc.want, got.Items) {
				fmt.Printf("\n\nWant\n")
				for _, v := range tc.want {
					fmt.Printf("%+v\n", v.Name)
				}
				fmt.Printf("\n\nGot\n")
				for _, v := range got.Items {
					fmt.Printf("%+v\n", v.Name)
				}

				t.Fatalf("expected: %+v, got: %+v", tc.want, got.Items)
			}
		})
	}
}

func TestGetLatestImage(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
	dat, err := os.ReadFile("test_files/gcloudout/images.json")
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	imgs := []*compute.Image{}
	err = json.Unmarshal(dat, &imgs)
	if err != nil {
		t.Fatalf("got error during preloading: %s", err)
	}

	tests := map[string]struct {
		project      string
		imageProject string
		imageFamily  string
		want         string
	}{
		"debian": {projectID, "debian-cloud", "debian-11", getLatestImageByProjectFromFile(imgs, "debian-cloud", "debian-11")},
		"rhel":   {projectID, "rhel-cloud", "rhel-9", getLatestImageByProjectFromFile(imgs, "rhel-cloud", "rhel-9")},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.ImageLatestGet(tc.project, tc.imageProject, tc.imageFamily)
			if err != nil {
				t.Fatalf("expected: no error, got: %v", err)
			}

			// sort.Strings(got)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}
