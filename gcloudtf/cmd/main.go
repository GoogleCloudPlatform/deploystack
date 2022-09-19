package main

import (
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/deploystack/gcloudtf"
)

func main() {
	b, err := gcloudtf.Extract("../testdata/extracttest")
	if err != nil {
		log.Fatalf("couldn't extract from TF file: %s", err)
	}

	for _, v := range *b {
		fmt.Printf("%s %s %s\n", v.Kind, v.Type, v.Name)
		fmt.Printf("%+v\n", v.Attr)

	}
}
