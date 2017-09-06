package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rcw5/gpx-simplifier-cli/simplify"
	"github.com/urfave/cli"
)

func main() {
	var numFiles int
	var pointsPerFile int
	var filename string
	var debug bool
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "number-of-files, n",
			Value:       1,
			Usage:       "Number of smaller files to split the larger file into",
			Destination: &numFiles,
		},
		cli.IntFlag{
			Name:        "points-per-file, p",
			Value:       1000,
			Usage:       "Number of GPX trackpoints per file",
			Destination: &pointsPerFile,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debugging",
			Destination: &debug,
		},
		cli.StringFlag{
			Name:        "filename, f",
			Usage:       "File to simplify",
			Destination: &filename,
		},
	}
	app.Name = "gpx-simplifier"
	app.Usage = "Simplify (and split) GPX files"
	app.Version = "0.0.2"
	app.Action = func(c *cli.Context) error {
		if filename == "" {
			return cli.NewExitError("ERROR: --filename must be specified", 1)
		}
		t, err := simplify.Load(filename)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error simplifying file: %s", err), 1)
		}

		tracks, _ := t.Split(numFiles)
		for idx, track := range tracks {
			track.Simplify(pointsPerFile)
			filenameNoExtension := strings.TrimSuffix(filename, filepath.Ext(filename))
			err = track.Save(fmt.Sprintf("%s-part%d.gpx", filenameNoExtension, idx+1))
			if err != nil {
				panic(err)
			}
		}
		return nil
	}

	app.Run(os.Args)
}
