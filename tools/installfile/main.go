package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/deploystack"
)

var lineLength = 80

func main() {
	err := printFile()
	if err != nil {
		log.Fatalf("err %s", err)
	}
}

func printFile() error {
	f, err := os.Create("output.txt")
	if err != nil {
		return err
	}

	deploystack.Divider, err = deploystack.BuildDivider(lineLength)
	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(f)

	fmt.Fprintf(w, "%s%s%s\n", deploystack.TERMCYANB, deploystack.Divider, deploystack.TERMCLEAR)
	fmt.Fprintf(w, "%s%s%s\n", deploystack.TERMCYANB, padString("Deploystack INSTALL", lineLength), deploystack.TERMCLEAR)
	fmt.Fprintf(w, "%s%s%s\n", deploystack.TERMCYANB, deploystack.Divider, deploystack.TERMCLEAR)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "%s\n", "To continue, run this command at the command line below.")
	fmt.Fprintf(w, "%s\n", "This process should walk you through setting up and deploying")
	fmt.Fprintf(w, "%s\n", "this deploystack application.")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "%s%s%s\n", deploystack.TERMCYANREV, deploystack.Divider, deploystack.TERMCLEAR)
	fmt.Fprintf(w, "%s%s%s\n", deploystack.TERMCYANREV, padString("./deploystack install", lineLength), deploystack.TERMCLEAR)
	fmt.Fprintf(w, "%s%s%s\n", deploystack.TERMCYANREV, deploystack.Divider, deploystack.TERMCLEAR)

	w.Flush()

	return nil
}

func padString(s string, linelength int) string {
	width := linelength - len(s)

	var sb strings.Builder
	sb.WriteString(s)

	for i := 0; i < width; i++ {
		sb.WriteString(" ")
	}

	return sb.String()
}
