package main

import (
	"log"

	"github.com/tpryan/deploystack"
)

func main() {
	// fmt.Println("\033[2J")
	f := deploystack.HandleFlags()
	s := deploystack.NewStack()
	s.ProcessFlags(f)

	if err := s.ReadConfig("config.json", "description.txt"); err != nil {
		log.Fatalf("could not read config file: %s", err)
	}

	if err := s.Process("terraform.tfvars"); err != nil {
		log.Fatalf("problemn collecting the configurations: %s", err)
	}
}
