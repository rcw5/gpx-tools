package simplify_test

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rcw5/gpx-simplifier-cli/simplify"
	. "github.com/rcw5/gpx-simplifier-cli/simplify"
)

func generateTrack(numPoints int) Track {
	var track = Track{Creator: "Sample", Title: "Sample"}
	var trackPoints = []*TrackPoint{}
	for i := 1; i <= numPoints; i++ {
		trackPoints = append(trackPoints, &TrackPoint{Lat: float64(i), Lon: float64(i)})
	}
	track.TrackPoints = trackPoints
	Expect(len(track.TrackPoints)).To(Equal(numPoints))
	return track
}

var _ = Describe("Simplify", func() {

	It("Should load XML file", func() {
		theTrack := generateTrack(20)
		tmpfile, err := ioutil.TempFile("", "example")
		Expect(err).ToNot(HaveOccurred())
		defer os.Remove(tmpfile.Name())
		content, err := xml.MarshalIndent(theTrack, "  ", "    ")
		Expect(err).ToNot(HaveOccurred())
		_, err = tmpfile.Write(content)
		Expect(err).ToNot(HaveOccurred())

		track, err := simplify.Load(tmpfile.Name())
		Expect(err).ToNot(HaveOccurred())
		Expect(track).To(Equal(theTrack))

	})

	It("Should split track into equal parts", func() {
		theTrack := generateTrack(20)

		splitTrack, err := theTrack.Split(4)
		Expect(err).ToNot(HaveOccurred())

		Expect(len(splitTrack)).To(Equal(4))

		for i := 0; i < 4; i++ {
			Expect(len(splitTrack[i].TrackPoints)).To(Equal(5))
		}
	})

})
