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

	fmt.Fprintln(w, deploystack.TERMCLEARSCREEN)
	fmt.Fprintf(w, "%s%s%s\n", deploystack.TERMCYANREV, deploystack.Divider, deploystack.TERMCLEAR)
	fmt.Fprintf(w, "%s%s%s\n", deploystack.TERMCYANREV, padString("Deploystack INSTALL", lineLength), deploystack.TERMCLEAR)
	fmt.Fprintf(w, "%s%s%s\n", deploystack.TERMCYANREV, deploystack.Divider, deploystack.TERMCLEAR)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "%s\n", "The install process should walk you through setting up and deploying")
	fmt.Fprintf(w, "%s\n", "this DeployStack project.")
	fmt.Fprintf(w, "%s\n", "To start, run the following command:")
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "%s%s%s\n", deploystack.TERMCYANB, padString("./deploystack install", lineLength), deploystack.TERMCLEAR)
	fmt.Fprintln(w, "")
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
