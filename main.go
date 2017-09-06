package main

import (
	"fmt"
	"os"

	"github.com/rcw5/gpx-simplifier-cli/simplify"
	"github.com/urfave/cli"
)

func main() {
	var numFiles int
	var pointsPerFile int
	var filename string
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
		cli.StringFlag{
			Name:        "filename, f",
			Usage:       "File to simplify",
			Destination: &filename,
		},
	}
	app.Name = "gpx-simplifier"
	app.Usage = "Simplify (and split) GPX files"
	app.Version = "0.0.1"
	app.Action = func(c *cli.Context) error {
		t, err := simplify.Load(filename)
		if err != nil {
			panic(err)
		}

		tracks, _ := t.Split(numFiles)

		for idx, track := range tracks {
			track.Simplify(pointsPerFile)
			err = track.Save(fmt.Sprintf("%s-part%d.gpx", filename, idx))
			if err != nil {
				panic(err)
			}
		}
		return nil
	}

	app.Run(os.Args)
}
