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
	"log"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/gcloud"
	"github.com/GoogleCloudPlatform/deploystack/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	s := deploystack.NewStack()

	if err := s.FindAndReadRequired(); err != nil {
		log.Fatalf("could not read config file: %s", err)
	}

	client := gcloud.NewClient(context.Background(), "")

	q := tui.NewQueue(&s, &client)
	q.Save("contact", deploystack.CheckForContact())
	q.InitializeUI()

	p := tea.NewProgram(q.Start(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalf(err.Error())
	}

	s.TerraformFile("terraform.tfvars")
}
