package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/Toshik1978/onepagepass/pkg/converter"
)

const (
	version = "Commit %s, built on %s"
)

var (
	Buildstamp = "undefined" //nolint:revive
	Commit     = "undefined" //nolint:revive
)

func main() {
	app := &cli.App{
		Name:     "onepagepass",
		Version:  fmt.Sprintf(version, Commit, Buildstamp),
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Anton Krivenko",
				Email: "anton@thedatron.ru",
			},
		},
		Copyright: "(c) 2021 Anton Krivenko",
		Usage:     "Run to convert scan of your documents to more useful PDF",
		Action:    cli.ShowAppHelp,
		Commands: []*cli.Command{
			{
				Name:    "convert",
				Aliases: []string{"c"},
				Usage:   "convert given document",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "pdf",
						Required: true,
					},
					&cli.Float64Flag{
						Name:        "dpi",
						DefaultText: "300",
					},
				},
				Action: func(c *cli.Context) error {
					convert := converter.New(c.String("pdf"), c.Float64("dpi"))
					return convert.Run()
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
