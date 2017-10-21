package simplify

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"

	"github.com/golangplus/fmt"
)

type Track struct {
	Creator     string        `xml:"creator,attr"`
	Title       string        `xml:"trk>name"`
	TrackPoints []*TrackPoint `xml:"trk>trkseg>trkpt"`
}

type TrackPoint struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
}

type xteRec struct {
	xte        float64
	trackpoint *TrackPoint
	prevXte    *xteRec
	nextXte    *xteRec
	ordinal    int
}

//Load a GPX file into a Track
func Load(fileName string) (Track, error) {
	if _, err := os.Stat(fileName); err != nil {
		return Track{}, fmt.Errorf("%s does not exist", fileName)
	}

	xmlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return Track{}, fmt.Errorf("Error reading %s: %s", fileName, err)
	}

	var track Track
	xml.Unmarshal(xmlFile, &track)
	return track, nil
}

//SplitInto splits a single track into the specified number of smaller tracks
func (track *Track) SplitInto(numFiles int) []Track {
	trackpointsPerFile := len(track.TrackPoints) / numFiles

	var tracks = []Track{}
	for i := 1; i <= numFiles; i++ {
		if track.Title == "" {
			track.Title = "file"
		}
		fileName := fmt.Sprintf("%s-part%d", track.Title, i)
		newTrack := Track{Creator: "gpx-simplifier", Title: fileName}
		start := (i - 1) * trackpointsPerFile
		end := start + trackpointsPerFile
		if i < numFiles {
			newTrack.TrackPoints = track.TrackPoints[start:end]
		} else {
			newTrack.TrackPoints = track.TrackPoints[start:]
		}
		tracks = append(tracks, newTrack)
	}
	return tracks
}

//Save saves the track definition to a file
func (track *Track) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	var buffer bytes.Buffer
	header := `<?xml version="1.0"?>
	<gpx version="1.0" creator="gpx-simplifier-cli"
	  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	  xmlns="http://www.topografix.com/GPX/1/0"
	  xsi:schemaLocation="http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd">
	  <trk>
		<name>%s</name>
		<trkseg>`
	_, err = buffer.WriteString(fmt.Sprintf(header, track.Title))
	if err != nil {
		return err
	}
	for _, val := range track.TrackPoints {
		trackPoint := fmt.Sprintf("<trkpt lat=\"%f\" lon=\"%f\" />\n", val.Lat, val.Lon)
		buffer.WriteString(trackPoint)
	}
	footer := `    </trkseg>
	</trk>
  </gpx>`
	buffer.WriteString(footer)
	if err != nil {
		return err
	}

	_, err = w.WriteString(buffer.String())
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}

//SimplifyTo simplifies a track into a certain number of points
func (track *Track) SimplifyTo(numPoints int) {
	if len(track.TrackPoints) < numPoints {
		return
	}

	var list []*xteRec
	var prevXte *xteRec

	for idx, t := range track.TrackPoints {
		var xteForTrackPoint float64
		if idx == 0 || idx == len(track.TrackPoints)-1 {
			xteForTrackPoint = math.MaxFloat64
		} else {
			prev := track.TrackPoints[idx-1]
			curr := track.TrackPoints[idx]
			next := track.TrackPoints[idx+1]
			xteForTrackPoint = calculateXte(*prev, *curr, *next)
		}
		tmp := t //not sure why this is needed yet but otherwise everything points to last trackpoint
		xte := xteRec{trackpoint: tmp, ordinal: idx, xte: xteForTrackPoint, prevXte: prevXte}
		//Update xtePointers but skip the first entry
		if prevXte != nil {
			prevXte.nextXte = &xte
		}
		prevXte = &xte
		list = insertXte(list, &xte)
	}

	//Remove smallest from thetrack
	pointsToRemove := int(math.Abs(float64(numPoints - len(track.TrackPoints))))
	for i := 0; i < pointsToRemove; i++ {
		lowestElem := list[0]

		prevElem := lowestElem.prevXte
		prevPrevElem := prevXte.prevXte
		nextElem := lowestElem.nextXte
		nextNextElem := nextElem.nextXte

		if prevPrevElem != nil {
			prevElem.xte = calculateXte(*prevPrevElem.trackpoint, *prevElem.trackpoint, *nextElem.trackpoint)
		}
		if nextNextElem != nil {
			nextElem.xte = calculateXte(*prevElem.trackpoint, *nextElem.trackpoint, *nextNextElem.trackpoint)
		}
		//Remove this from the linked list
		prevElem.nextXte = nextElem
		nextElem.prevXte = prevElem

		list = removeXte(list, prevElem)
		list = removeXte(list, nextElem)
		list = removeXte(list, lowestElem)
		list = insertXte(list, prevElem)
		list = insertXte(list, nextElem)

		for i, item := range track.TrackPoints {
			if item == lowestElem.trackpoint {
				track.TrackPoints = append(track.TrackPoints[:i], track.TrackPoints[i+1:]...)
			}
		}
	}
}

//Length returns the number of trackpoints in this track
func (track *Track) Length() int {
	return len(track.TrackPoints)
}

func insertXte(list []*xteRec, xte *xteRec) []*xteRec {
	if len(list) == 0 {
		return append(list, xte)
	}
	for i, item := range list {
		if item.xte >= xte.xte {
			return append(list[:i], append([]*xteRec{xte}, list[i:]...)...)
		}
	}
	return list
}

func removeXte(list []*xteRec, xte *xteRec) []*xteRec {
	for i, item := range list {
		if item == xte {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func printList(list []*xteRec) {
	for i, item := range list {
		fmtp.Printfln("%d %g", i, item.xte)
	}
}

//Calculate the Cross Track Error of three trackpoints
//That is, the height of a triangle formed of lines AC, AB, BC
//The lower the XTE the lower the impact on the shape of the final route
func calculateXte(a, b, c TrackPoint) float64 {
	aTob := distance(a.Lat, a.Lon, b.Lat, b.Lon)
	aToc := distance(a.Lat, a.Lon, c.Lat, c.Lon)
	bToc := distance(b.Lat, b.Lon, c.Lat, c.Lon)

	area := area(aTob, aToc, bToc)
	xte := height(aToc, area)
	return xte
}

func height(base, area float64) float64 {
	return 2 * (area / base)
}

//heron's formula: http://www.mathopenref.com/heronsformula.html
func area(aTob, aToc, bToc float64) float64 {
	p := (aTob + aToc + bToc) / 2
	return math.Sqrt(p * (p - aTob) * (p - aToc) * (p - bToc))
}

//Taken from https://gist.github.com/cdipaolo/d3f8db3848278b49db68
// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return math.Abs(2 * r * math.Asin(math.Sqrt(h)))
}
