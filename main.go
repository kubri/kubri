package main

import (
	"os"

	"github.com/abemedia/appcast/pkg/cmd"
	_ "github.com/joho/godotenv/autoload"
)

var version = "dev"

func main() {
	if err := cmd.Execute(version); err != nil {
		os.Exit(1)
	}
}
