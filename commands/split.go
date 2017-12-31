package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rcw5/gpx-tools/gpx"
)

func Split(inputFile, outputPath string, numFiles int) error {
	var err error
	if _, err = os.Stat(inputFile); err != nil {
		return fmt.Errorf("%s does not exist", inputFile)
	}

	fileContents, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("Error reading %s: %s", inputFile, err)
	}
	track, err := gpx.Load(string(fileContents))
	if err != nil {
		return err
	}
	if outputPath == "" {
		outputPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("Cannot get output directory: %s", err)
		}
	}
	tracks := track.SplitInto(numFiles)

	outputFilePrefix := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	for index, splitTrack := range tracks {
		outputFileName := fmt.Sprintf("%s/%s_%d.gpx", outputPath, outputFilePrefix, index+1)
		err = ioutil.WriteFile(outputFileName, []byte(splitTrack.ToXML()), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
