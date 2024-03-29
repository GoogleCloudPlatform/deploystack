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
	"flag"
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/tui"
)

const (
	contactfile = "contact.yaml"
)

func main() {
	verify := flag.Bool("verify", false, "Whether or not to be in verify mode")
	name := flag.Bool("name", false, "Whether or not to be in drop the name of the stack")

	flag.Parse()

	s, err := deploystack.Init(".")
	if err != nil {
		log.Fatalf("could not read initialize deploystack: %s", err)
	}

	if *verify {
		fmt.Printf("%s|%s|%s\n", s.Config.PathTerraform, s.Config.PathMessages, s.Config.PathScripts)
		return
	}

	if *name {
		fmt.Printf("%s\n", s.Config.Name)
		return
	}

	tui.Run(s, false)
}
