package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"aigc.dev/appstore/apple"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var (
	flagRegion      string
	flagSearchLimit int
	flagCredPath    string
	flagMac         string
	flagTrackID     string
	flagOutput      string

	AppleCommand = &cli.Command{
		Name:    "apple",
		Usage:   "Subcommand for interacting with Apple App Store",
		Aliases: []string{"ios", "i"},
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "cred",
				Aliases:     []string{"c"},
				EnvVars:     []string{"APPLE_CRED"},
				Usage:       "Path to the Apple App Store credential file",
				Value:       "apple.cred.json",
				Destination: &flagCredPath,
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:    "login",
				Aliases: []string{"l"},
				Usage:   "Login to the Apple App Store",
				Action:  appleLoginAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "mac",
						Aliases:     []string{"m"},
						EnvVars:     []string{"APPLE_MAC"},
						Usage:       "Specify the MAC address of the device, if empty, it will use the device's actual MAC address",
						Destination: &flagMac,
					},
					&cli.StringFlag{
						Name:        "region",
						Aliases:     []string{"c"},
						EnvVars:     []string{"APPLE_REGION"},
						Usage:       "Region of the Apple App Store, ISO 3166-1 alpha-2 country code",
						Destination: &flagRegion,
					},
				},
			},
			{
				Name:      "search",
				Aliases:   []string{"s"},
				Usage:     "Search for an app in the Apple App Store",
				Action:    appleSearchAction,
				Args:      true,
				ArgsUsage: "<query>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "region",
						Aliases:     []string{"c"},
						EnvVars:     []string{"APPLE_REGION"},
						Usage:       "Region of the Apple App Store, ISO 3166-1 alpha-2 country code",
						Value:       "US",
						Destination: &flagRegion,
					},
					&cli.IntFlag{
						Name:        "limit",
						Aliases:     []string{"l"},
						EnvVars:     []string{"APPLE_SEARCH_LIMIT"},
						Usage:       "Limit the number of search results",
						Value:       50,
						Destination: &flagSearchLimit,
					},
				},
			},
			{
				Name:    "download",
				Aliases: []string{"d"},
				Usage:   "Download an app from the Apple App Store",
				Action:  appleDownloadAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "bundle-id",
						Aliases:     []string{"b"},
						EnvVars:     []string{"APPLE_BUNDLE_ID"},
						Usage:       "Bundle ID of the app to download",
						Required:    true,
						Destination: &flagTrackID,
					},
					&cli.PathFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						EnvVars:     []string{"APPLE_OUTPUT"},
						Usage:       "Output path to save the downloaded app, if empty, it will save to the current directory",
						Destination: &flagOutput,
					},
				},
			},
		},
	}
)

func appleLoginAction(cCtx *cli.Context) error {
	if flagCredPath == "" {
		log.Fatal("cred is required")
	}

	// check if the credential file exists
	fs := afero.NewOsFs()
	if info, err := fs.Stat(flagCredPath); err == nil && !info.IsDir() {
		log.Fatalf("credential file %s already exists", flagCredPath)
	}

	// user input
	reader := bufio.NewReader(cCtx.App.Reader)
	fmt.Print("Apple ID: ")
	appleId, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	appleId = strings.TrimSpace(appleId)

	fmt.Print("Password: ")
	pwdBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	pwd := strings.TrimSpace(string(pwdBytes))
	fmt.Println()

	if flagMac == "" {
		fmt.Print("MAC Address (optional): ")
		flagMac, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		flagMac = strings.TrimSpace(flagMac)
	} else {
		fmt.Printf("Specified MAC Address: %s\n", flagMac)
	}

	if flagRegion == "" {
		fmt.Print("Account Region (optional): ")
		flagRegion, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		flagRegion = strings.TrimSpace(flagRegion)
	} else {
		fmt.Printf("Specified Region: %s\n", flagRegion)
	}

	// login
	client := apple.NewAppleClient()
	if err := client.Login(apple.LoginOptions{AppleID: appleId, Password: pwd, MacAddress: flagMac, Region: flagRegion}); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Login successful, welcome %s\n", client.Cred.AppleID)

	// save credential
	if err := client.SaveCredentials(flagCredPath); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Credential saved to %s\n", flagCredPath)

	return nil
}

func appleSearchAction(cCtx *cli.Context) error {
	client := apple.NewAppleClient()

	query := cCtx.Args().First()
	if query == "" {
		log.Fatal("query is required")
	}

	opt := apple.SearchOptions{
		Region: flagRegion,
		Query:  query,
		Limit:  flagSearchLimit,
	}

	items, err := client.Search(opt)
	if err != nil {
		log.Fatal(err)
	}

	tw := table.NewWriter()
	tw.SetOutputMirror(cCtx.App.Writer)

	tw.AppendHeader(table.Row{"#", "ID", "Name", "Bundle ID", "Genre"})
	for i, item := range items {
		tw.AppendRow([]any{i + 1, item.TrackID, item.TrackName, item.BundleID, item.PrimaryGenre})
	}

	tw.Render()
	return nil
}

func appleDownloadAction(cCtx *cli.Context) error {
	client := apple.NewAppleClient()
	if err := client.LoadCredentials(flagCredPath); err != nil {
		log.Fatal(err)
	}

	if flagTrackID == "" {
		log.Fatal("track-id is required")
	}

	if flagOutput == "" {
		flagOutput = fmt.Sprintf("./%s.ipa", flagTrackID)
	}

	return client.Download(flagTrackID, flagOutput)
}
