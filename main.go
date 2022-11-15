package main

import (
	"log"
	"os"

	"github.com/ashrithr/tana-readwise-exporter/readwise"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "tana-readwise-exporter",
		Version:              "v0.1",
		Usage:                "Exports highlights from Readwise and formats using Tana paste.",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List the books",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "category",
						Aliases: []string{"c"},
						Usage:   "Specify category to list from 'books', 'articles', 'tweets', 'podcasts', etc",
						Value:   "books",
					},
					&cli.StringFlag{
						Name:     "token",
						Aliases:  []string{"t"},
						Usage:    "Readiwse API Token",
						EnvVars:  []string{"READWISE_TOKEN"},
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					return readwise.List(ctx.String("token"), ctx.String("category"))
				},
			},
			{
				Name:    "export",
				Aliases: []string{"e"},
				Usage:   "Export the highlights",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "token",
						Aliases:  []string{"t"},
						Usage:    "Readiwse API Token",
						EnvVars:  []string{"READWISE_TOKEN"},
						Required: true,
					},
					&cli.IntFlag{
						Name:    "updated-after",
						Aliases: []string{"u"},
						Usage:   "Updated after",
					},
					&cli.StringSliceFlag{
						Name:    "ids",
						Aliases: []string{"i"},
						Usage:   "Ids of books/articles to fetch the highlights for",
					},
				},
				Action: func(ctx *cli.Context) error {
					return readwise.ListHighlights(ctx.String("token"), ctx.Int("updated-after"), ctx.StringSlice("ids"))
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "print debugging logs",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
