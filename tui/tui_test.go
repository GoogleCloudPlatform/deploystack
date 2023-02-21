// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tui

import (
	"log"
	"os"
)

func readTestFile(file string) string {
	dat, err := os.ReadFile(file)
	if err != nil {
		return "Couldn't read test file"
	}

	return string(dat)
}

func writeDebugFile(content string, target string) {
	d1 := []byte(content)
	err := os.WriteFile("debug/"+target, d1, 0o644)
	if err != nil {
		log.Printf("err: %s", err)
	}
}
