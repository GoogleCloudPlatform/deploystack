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
	"log"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

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

func compareValues(label string, want interface{}, got interface{}, t *testing.T) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("%s: expected: \n|%v|\ngot: \n|%v|", label, want, got)
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
