package main

import (
	"flag"
	"log"

	"github.com/kazdevl/prelviz"
)

var (
	projectDirectoryPath string
	outputFilePath       string
	dotLayout            string
)

func main() {
	flag.StringVar(&projectDirectoryPath, "i", "", `requreid: "true", description: "input project directory path"`)
	flag.StringVar(&outputFilePath, "o", "", `requreid: "false", description: "output file path(default is stdout)"`)
	flag.StringVar(&dotLayout, "l", "dot", `requreid: "false", description: "dot layout. ex) dot, neato, fdp, sfdp, twopi, circo"`)
	flag.Parse()

	if projectDirectoryPath == "" {
		log.Fatal("project directory path is required")
	}

	prelviz, err := prelviz.NewPrelviz(projectDirectoryPath, outputFilePath, dotLayout)
	if err != nil {
		log.Fatal(err)
	}
	if err = prelviz.Run(); err != nil {
		log.Fatal(err)
	}
}
