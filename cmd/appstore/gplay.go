package main

import (
	"log"

	"github.com/urfave/cli/v2"
)

var (
	GooglePlayCommand = &cli.Command{
		Name:    "gplay",
		Usage:   "Subcommand for interacting with Google Play Store",
		Aliases: []string{"googleplay", "android", "g"},
		Action:  gplayAction,
	}
)

func gplayAction(cCtx *cli.Context) error {
	log.Println("Google Play Store command is not implemented yet")
	return nil
}
