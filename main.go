package main

import (
	"os"

	"github.com/abemedia/appcast/pkg/cmd"
)

var version = "dev"

func main() {
	if err := cmd.Execute(version, nil); err != nil {
		os.Exit(1)
	}
}
