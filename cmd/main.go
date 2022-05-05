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

package main

import (
	"log"

	"github.com/GoogleCloudPlatform/deploystack"
)

func main() {
	deploystack.ClearScreen()
	f := deploystack.HandleFlags()
	s := deploystack.NewStack()
	s.ProcessFlags(f)

	if err := s.ReadConfig("deploystack.json", "deploystack.txt"); err != nil {
		log.Fatalf("could not read config file: %s", err)
	}

	if err := s.Process("terraform.tfvars"); err != nil {
		log.Fatalf("problemn collecting the configurations: %s", err)
	}
}
