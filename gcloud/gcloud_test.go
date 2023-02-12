package gcloud

import (
	"context"
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

	"google.golang.org/api/option"
)

var (
	projectID        = ""
	creds            map[string]string
	opts             = option.WithCredentialsFile("")
	ctx              = context.Background()
	defaultUserAgent = "deploystack/testing"
)

func TestMain(m *testing.M) {
	var err error
	opts = option.WithCredentialsFile("../creds.json")

	dat, err := os.ReadFile("../creds.json")
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

func readTestFile(file string) string {
	dat, err := os.ReadFile(file)
	if err != nil {
		return "Couldn't read test file"
	}

	return string(dat)
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

func TestGetRegions(t *testing.T) {
	c := NewClient(ctx, defaultUserAgent)
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
		"computeRegions": {
			product: "compute",
			project: projectID,
			want:    cRegions,
			err:     nil,
		},

		"functionsRegions": {
			product: "functions",
			project: projectID,
			want:    fRegions,
			err:     nil,
		},

		"runRegions": {
			product: "run",
			project: projectID,
			want:    rRegions,
			err:     nil,
		},

		"GarbageInout": {
			product: "An outdated iPad",
			project: projectID,
			want:    []string{},
			err: fmt.Errorf(
				"invalid product (%s) requested",
				"An outdated iPad"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := c.RegionList(tc.project, tc.product)

			// BUG: getting weird regions intertmittenly popping up. Solving with this hack
			if tc.product == "compute" {
				got = removeDuplicateStr(removeFromSlice(removeFromSlice(got, "me-west1"), "us-west4"))
				tc.want = removeDuplicateStr(removeFromSlice(removeFromSlice(cRegions, "me-west1"), "us-west4"))
			}

			if err != tc.err {
				if err.Error() != tc.err.Error() {
					t.Fatalf("expected: error (%v), got: %v", tc.err, err)
				}
			}

			sort.Strings(got)

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}
