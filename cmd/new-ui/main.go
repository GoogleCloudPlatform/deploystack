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
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	"github.com/GoogleCloudPlatform/deploystack/tui"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	contactfile = "contact.yaml"
)

func main() {
	verify := flag.Bool("verify", false, "Whether or not to be in verify mode")
	name := flag.Bool("name", false, "Whether or not to be in drop the name of the stack")

	flag.Parse()

	s := deploystack.NewStack()

	if err := s.FindAndReadRequired(); err != nil {
		log.Fatalf("could not read config file: %s", err)
	}

	if *verify {
		fmt.Printf("%s|%s|%s\n", s.Config.PathTerraform, s.Config.PathMessages, s.Config.PathScripts)
		return
	}

	if *name {
		if err := s.Config.ComputeName(); err != nil {
			log.Fatalf("could retrieve name of stack: %s", err)
		}

		fmt.Printf("%s\n", s.Config.Name)
		return
	}

	s.AddSetting("stack_name", s.Config.Name)

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	defaultUserAgent := fmt.Sprintf("deploystack/%s", s.Config.Name)
	client := gcloud.NewClient(context.Background(), defaultUserAgent)

	q := tui.NewQueue(&s, &client)
	q.Save("contact", deploystack.CheckForContact())
	q.InitializeUI()

	p := tea.NewProgram(q.Start(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalf(err.Error())
	}

	s.TerraformFile("terraform.tfvars")

	deploystack.CacheContact(q.Get("contact"))
}
