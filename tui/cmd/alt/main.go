package main

import (
	"github.com/GoogleCloudPlatform/deploystack"
	"github.com/GoogleCloudPlatform/deploystack/tui"
)

func main() {

	s, err := deploystack.Init()
	if err != nil {
		tui.Fatal(err)
	}

	tui.AltRun(s, true)
}
