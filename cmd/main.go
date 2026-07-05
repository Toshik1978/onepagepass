// Command onepagepass converts a PDF scan of a document into a
// two-pages-per-A4 layout that is handier to print.
package main

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/Toshik1978/onepagepass/pkg/converter"
)

const (
	version = "Commit %s, built on %s"
)

// Buildstamp and Commit are populated at build time via -ldflags.
var (
	Buildstamp = "undefined" //nolint:gochecknoglobals
	Commit     = "undefined" //nolint:gochecknoglobals
)

func main() {
	if err := newRootCommand().Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// newRootCommand builds the onepagepass CLI command tree.
func newRootCommand() *cli.Command {
	return &cli.Command{
		Name:    "onepagepass",
		Version: fmt.Sprintf(version, Commit, Buildstamp),
		Authors: []any{
			&mail.Address{Name: "Anton Krivenko", Address: "anton@krivenko.dev"},
		},
		Copyright: "(c) 2021-2026 Anton Krivenko",
		Usage:     "Run to convert scan of your documents to more useful PDF",
		Action: func(_ context.Context, cmd *cli.Command) error {
			return cli.ShowAppHelp(cmd)
		},
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
						Name:  "dpi",
						Value: 300,
					},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					convert := converter.New(cmd.String("pdf"), cmd.Float64("dpi"))
					return convert.Run()
				},
			},
		},
	}
}
