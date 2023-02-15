package tui

import (
	"log"
	"os"
)

func readTestFile(file string) string {
	dat, err := os.ReadFile(file)
	if err != nil {
		return "Couldn't read test file"
	}

	return string(dat)
}

func writeDebugFile(content string, target string) {
	d1 := []byte(content)
	err := os.WriteFile("debug/"+target, d1, 0o644)
	if err != nil {
		log.Printf("err: %s", err)
	}
}
