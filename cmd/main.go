package main

import (
	"flag"
	"log"

	"github.com/tpryan/deploystack"
)

func main() {
	projectPtr := flag.String("project", "", "A Google Cloud Project ID")
	regionPtr := flag.String("region", "", "A Google Cloud Region")
	zonePtr := flag.String("zone", "", "A Google Cloud Zone")

	flag.Parse()

	// fmt.Println("\033[2J")
	s := deploystack.NewStack()

	s.AddSetting("project", *projectPtr)
	s.AddSetting("region", *regionPtr)
	s.AddSetting("zone", *zonePtr)

	if err := s.ReadConfig("config.json"); err != nil {
		log.Fatalf("could not read config file: %s", err)
	}

	if err := s.Process("terraform.tfvars"); err != nil {
		log.Fatalf("problemn collecting the configurations: %s", err)
	}
}
