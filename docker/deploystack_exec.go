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
	"os"
	"strings"

	_ "embed"

	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/github"
	"github.com/GoogleCloudPlatform/deploystack/tui"
)

//go:embed versionDS
var versionDS string

//go:embed buildTime
var buildTime string

func main() {
	verify := flag.Bool("verify", false, "Whether or not to be in verify mode")
	name := flag.Bool("name", false, "Whether or not to be in drop the name of the stack")
	version := flag.Bool("version", false, "Shows version information")
	repo := flag.String("repo", "", "The name only of a Google Cloud Platform repo to download")
	suggest := flag.Bool("suggest", false, "Weather or not you want DeployStack to reccomend a config")

	flag.Parse()

	if *version {
		fmt.Printf("deploystack:        %s\n", strings.TrimSpace(versionDS))
		fmt.Printf("buildTime:          %s\n", strings.TrimSpace(buildTime))
		return
	}

	if *repo != "" {
		wd, err := os.Getwd()
		if err != nil {
			tui.Fatal(err)
		}

		dir, gh, err := deploystack.AttemptRepo(*repo, wd)

		if err != nil {
			tui.Fatal(err)
		}

		dir = strings.ReplaceAll(dir, "//", "/")

		// If init works, that means there is a valid deploystack config here
		// exit, and let the hosting script rerun itself
		if _, err = deploystack.Init(dir); err == nil {
			fmt.Printf("%s\n", dir)
			return
		}

		if err := deploystack.WriteConfig(dir, gh); err != nil {
			tui.Fatal(err)
		}

		fmt.Printf("%s\n", dir)
		return
	}

	if *suggest {
		wd, err := os.Getwd()
		if err != nil {
			tui.Fatal(err)
		}

		if err := deploystack.WriteConfig(wd, github.Repo{}); err != nil {
			tui.Fatal(err)
		}

		fmt.Printf("Made a reccomended config for you\n")
		return
	}

	err := deploystack.Precheck()
	if err != nil {
		tui.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		tui.Fatal(err)
	}

	s, err := deploystack.Init(wd)
	if err != nil {
		tui.Fatal(fmt.Errorf("could not initalize: %s", err))
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
