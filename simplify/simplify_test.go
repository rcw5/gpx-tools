package simplify_test

import (
	"encoding/xml"
	"fmt"
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
	Expect(track.Length()).To(Equal(numPoints))
	return track
}

var _ = Describe("gpx-simplifier-cli", func() {
	var track Track
	var _ = BeforeEach(func() {
		track = generateTrack(20)
	})

	Context("Track Splitting", func() {
		It("splits a track with an even number of points into equal parts", func() {
			splitTracks := track.SplitInto(2)
			Expect(len(splitTracks)).To(Equal(2))
			for _, t := range splitTracks {
				Expect(t.Length()).To(Equal(10))
			}
		})

		It("sets the track title", func() {
			splitTrack := track.SplitInto(2)
			Expect(splitTrack[0].Title).To(Equal(fmt.Sprintf("Sample-part1")))
			Expect(splitTrack[1].Title).To(Equal(fmt.Sprintf("Sample-part2")))
		})

		It("splits a track with an odd number of points without losing data", func() {
			oddTrack := generateTrack(21)
			splitTracks := oddTrack.SplitInto(2)
			Expect(len(splitTracks)).To(Equal(2))
			Expect(splitTracks[0].Length()).To(Equal(10))
			Expect(splitTracks[1].Length()).To(Equal(11))
		})
	})

	Context("Simplification", func() {
		It("Simplifies a track to the correct number of points", func() {
			track.SimplifyTo(10)
			Expect(track.Length()).To(Equal(10))
		})

		It("Removes points with the smallest cross-track-error", func() {
			track = Track{
				TrackPoints: []*TrackPoint{
					&TrackPoint{Lat: 1, Lon: 1},
					&TrackPoint{Lat: 6, Lon: 6},
					&TrackPoint{Lat: 2, Lon: 2},
					&TrackPoint{Lat: 4, Lon: 4},
					&TrackPoint{Lat: 3, Lon: 3},
				},
			}
			track.SimplifyTo(4)
			Expect(*track.TrackPoints[0]).To(Equal(TrackPoint{Lat: 1, Lon: 1}))
			Expect(*track.TrackPoints[1]).To(Equal(TrackPoint{Lat: 6, Lon: 6}))
			Expect(*track.TrackPoints[2]).To(Equal(TrackPoint{Lat: 2, Lon: 2}))
			Expect(*track.TrackPoints[3]).To(Equal(TrackPoint{Lat: 3, Lon: 3}))
		})

		It("Doesn't simplify a track when track length is too small", func() {
			track.SimplifyTo(25)
			Expect(track.Length()).To(Equal(20))
		})

	})

	It("Should load XML file", func() {
		track := generateTrack(20)
		tmpfile, err := ioutil.TempFile("", "example")
		Expect(err).ToNot(HaveOccurred())
		defer os.Remove(tmpfile.Name())
		content, err := xml.MarshalIndent(track, "  ", "    ")
		Expect(err).ToNot(HaveOccurred())
		_, err = tmpfile.Write(content)
		Expect(err).ToNot(HaveOccurred())

		simplfiedTrack, err := simplify.Load(tmpfile.Name())
		Expect(err).ToNot(HaveOccurred())
		Expect(track).To(Equal(simplfiedTrack))

	})
})
