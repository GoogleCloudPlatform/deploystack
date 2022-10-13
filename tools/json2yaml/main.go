package main

import (
	"flag"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/deploystack"
	"gopkg.in/yaml.v2"
)

func main() {
	filePtr := flag.String("file", "", "the file path of a local deploystack json file")
	flag.Parse()
	file := *filePtr

	dat, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("could not read config file: %s", err)
	}

	cfg, err := deploystack.NewConfigJSON(dat)
	if err != nil {
		log.Fatalf("could not parse as a DeployStack config option: %s", err)
	}

	out, err := yaml.Marshal(&cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	f, err := os.Create("deploystack.yaml")
	if err != nil {
		log.Fatalf("could not create a file: %v", err)
	}

	defer f.Close()

	if _, err := f.Write(out); err != nil {
		log.Fatalf("could not write file to yaml: %s", err)
	}
}
