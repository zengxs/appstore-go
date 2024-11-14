package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "appstore",
		Usage: "A CLI tool to interact with Apple App Store and Google Play Store (WIP)",
		Commands: []*cli.Command{
			AppleCommand,
			GooglePlayCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
