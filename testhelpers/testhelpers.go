package testhelpers

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rcw5/gpx-tools/gpx"
)

func GenerateTrack(numPoints int) gpx.Track {
	var track = gpx.Track{Creator: "Sample", Title: "Sample"}
	var trackPoints = []*gpx.TrackPoint{}
	for i := 1; i <= numPoints; i++ {
		trackPoints = append(trackPoints, &gpx.TrackPoint{Lat: float64(i), Lon: float64(i)})
	}
	track.TrackPoints = trackPoints
	return track
}

func GenerateGpx(numPoints int) string {
	str := `<?xml version="1.0"?>
	<gpx version="1.0" creator="gpx-tools"
	  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	  xmlns="http://www.topografix.com/GPX/1/0"
	  xsi:schemaLocation="http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd">
	  <trk>
		<name>Sample Track</name>
		<trkseg>`
	for i := 1; i <= numPoints; i++ {
		str += fmt.Sprintf("<trkpt lat=\"%f\" lon=\"%f\" />\n", float64(i), float64(i))
	}
	str += `</trkseg>
	</trk>
	</gpx>`
	return str
}

func GetAsset(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func WriteStringToTempFile(folder, contents string) (string, error) {
	return WriteBytesToTempFile(folder, []byte(contents))
}
func WriteBytesToTempFile(folder string, contents []byte) (string, error) {
	tempFile, err := ioutil.TempFile(folder, "")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	err = ioutil.WriteFile(tempFile.Name(), contents, os.ModePerm)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), err
}
