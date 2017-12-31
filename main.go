package main

import (
	"os"

	"github.com/rcw5/gpx-tools/commands"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gpx-tools"
	app.Usage = "Simple suite of tools to manipulate GPX files"
	app.Version = "0.0.3"

	app.Commands = []cli.Command{
		{
			Name:      "simplify",
			Aliases:   []string{"sim"},
			Usage:     "Simplify a GPX File",
			ArgsUsage: "FILE",
			Action: func(c *cli.Context) error {
				if c.Args().First() == "" {
					cli.ShowCommandHelp(c, "simplify")
					os.Exit(1)
				}
				err := commands.Simplify(c.Args().First(), c.String("output-file"), c.Int("number-of-points"))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "output-file", Usage: "Output file"},
				cli.IntFlag{Name: "number-of-points", Value: 500},
			},
		},
		{
			Name:    "split",
			Aliases: []string{"spl"},
			Usage:   "Split a GPX file into a number of smaller files",
			Action: func(c *cli.Context) error {
				if c.Args().First() == "" {
					cli.ShowCommandHelp(c, "split")
					os.Exit(1)
				}
				err := commands.Split(c.Args().First(), c.String("output-path"), c.Int("number-of-files"))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
			ArgsUsage: "FILE",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "output-path", Usage: "Output files will be saved to this directory (default: current directory)"},
				cli.IntFlag{Name: "number-of-files", Usage: "Split the file into this number of smaller files", Value: 2},
			},
		},
	}

	app.Run(os.Args)
}
