package main

import (
	"log"

	"github.com/tpryan/deploystack"
)

func main() {
	deploystack.ClearScreen()
	f := deploystack.HandleFlags()
	s := deploystack.NewStack()
	s.ProcessFlags(f)

	if err := s.ReadConfig("config.json", "description.txt"); err != nil {
		log.Fatalf("could not read config file: %s", err)
	}

	if err := s.Process("variables.tfvars"); err != nil {
		log.Fatalf("problemn collecting the configurations: %s", err)
	}
}
