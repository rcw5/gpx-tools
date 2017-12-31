package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rcw5/gpx-tools/gpx"
)

func Simplify(inputFile, outputFile string, numPoints int) error {
	if _, err := os.Stat(inputFile); err != nil {
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

	if outputFile == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		outputFile = fmt.Sprintf("%s/%s_simplify%s", pwd,
			strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile)), filepath.Ext(inputFile))
	}
	track.SimplifyTo(numPoints)
	trackXML := track.ToXML()

	err = ioutil.WriteFile(outputFile, []byte(trackXML), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
