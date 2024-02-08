package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/kubri/kubri/pkg/cmd"
)

var version = "dev"

func main() {
	if err := cmd.Execute(version); err != nil {
		os.Exit(1)
	}
}
