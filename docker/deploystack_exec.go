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

// The dsexec binary that runs on Cloud Shell to present all of the options to
// users
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "embed"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/tui"
)

//go:embed versionTUI
var versionDS string

//go:embed versionTUI
var versionTUI string

//go:embed versionGcloud
var versionGcloud string

//go:embed versionConfig
var versionConfig string

//go:embed buildTime
var buildTime string

func main() {
	verify := flag.Bool("verify", false, "Whether or not to be in verify mode")
	name := flag.Bool("name", false, "Whether or not to be in drop the name of the stack")
	version := flag.Bool("version", false, "Shows version information")
	repo := flag.String("repo", "", "The name only of a Google Cloud Platform repo to download")

	flag.Parse()

	if *version {
		fmt.Printf("deploystack:        %s\n", strings.TrimSpace(versionDS))
		fmt.Printf("deploystack/tui:    %s\n", strings.TrimSpace(versionTUI))
		fmt.Printf("deploystack/gcloud: %s\n", strings.TrimSpace(versionGcloud))
		fmt.Printf("deploystack/config: %s\n", strings.TrimSpace(versionConfig))
		fmt.Printf("buildTime:          %s\n", strings.TrimSpace(buildTime))
		return
	}

	if *repo != "" {
		wd, err := os.Getwd()
		if err != nil {
			tui.Fatal(err)
		}

		dir, err := deploystack.DownloadRepo(*repo, wd)
		if err != nil {
			tui.Fatal(err)
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			tui.Fatal(err)
		}

		fmt.Printf("%s\n", strings.ReplaceAll(dir, "//", "/"))
		return
	}

	err := deploystack.Precheck()
	if err != nil {
		tui.Fatal(err)
	}

	s, err := deploystack.Init()
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
