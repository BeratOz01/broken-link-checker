package main

import (
	"log"
	"os"

	internals "github.com/BeratOz01/broken-link-checker/internals"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "broken",
		Usage: "This is a CLI tool to check broken links in a website with history and more",
		Commands: []*cli.Command{
			{
				Name:    "history",
				Aliases: []string{"h"},
				Action: func(c *cli.Context) error {
					page := c.Int("page")
					internals.PrintHistory(page)
					return nil
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "page",
						Usage: "specify the page to display",
						Value: 1,
					},
				},
			},
			{
				Name:    "start",
				Aliases: []string{"s"},
				Action: func(c *cli.Context) error {
					website := c.String("website")
					internals.Start(website)
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "website",
						Usage: "Specify the website to check",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
