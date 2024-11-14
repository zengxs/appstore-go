package main

import (
	"log"

	"github.com/urfave/cli/v2"
)

var (
	AppleCommand = &cli.Command{
		Name:    "apple",
		Usage:   "Subcommand for interacting with Apple App Store",
		Aliases: []string{"ios", "i"},
		Action:  appleAction,
	}
)

func appleAction(cCtx *cli.Context) error {
	log.Println("Google Play Store command is not implemented yet")
	return nil
}
