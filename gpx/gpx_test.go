package gpx_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rcw5/gpx-tools/gpx"
	"github.com/rcw5/gpx-tools/testhelpers"
)

var _ = Describe("gpx-tools", func() {
	var track gpx.Track
	var _ = BeforeEach(func() {
		track = testhelpers.GenerateTrack(20)
	})

	Context("Load", func() {
		BeforeEach(func() {
			sampleTrack := testhelpers.GenerateGpx(5)
			var err error
			track, err = gpx.Load(sampleTrack)
			Expect(err).ToNot(HaveOccurred())
		})
		It("Retrieves the creator from the GPX file", func() {
			Expect(track.Creator).To(Equal("gpx-tools"))
		})
		It("Retrieves the title from the GPX file", func() {
			Expect(track.Title).To(Equal("Sample Track"))
		})
		It("Retrieves the trackpoints from the GPX file", func() {
			Expect(track.Length()).To(Equal(5))
			Expect(track.TrackPoints[0]).To(Equal(&gpx.TrackPoint{Lat: 1, Lon: 1}))
			Expect(track.TrackPoints[1]).To(Equal(&gpx.TrackPoint{Lat: 2, Lon: 2}))
			Expect(track.TrackPoints[2]).To(Equal(&gpx.TrackPoint{Lat: 3, Lon: 3}))
			Expect(track.TrackPoints[3]).To(Equal(&gpx.TrackPoint{Lat: 4, Lon: 4}))
			Expect(track.TrackPoints[4]).To(Equal(&gpx.TrackPoint{Lat: 5, Lon: 5}))
		})
	})

	Context("ToXML", func() {
		It("Returns the track as XML", func() {
			sampleTrack := testhelpers.GenerateGpx(5)
			track, err := gpx.Load(sampleTrack)
			Expect(err).ToNot(HaveOccurred())
			savedTrack := track.ToXML()
			Expect(savedTrack).To(MatchXML(sampleTrack))
		})
	})

	Context("SplitInto", func() {
		It("splits a track with an even number of points into equal parts", func() {
			splitTracks := track.SplitInto(2)
			Expect(len(splitTracks)).To(Equal(2))
			for _, t := range splitTracks {
				Expect(t.Length()).To(Equal(10))
			}
		})

		It("includes the file number in the track's title", func() {
			splitTrack := track.SplitInto(2)
			Expect(splitTrack[0].Title).To(Equal(fmt.Sprintf("Sample-part1")))
			Expect(splitTrack[1].Title).To(Equal(fmt.Sprintf("Sample-part2")))
		})

		It("splits a track with an odd number of points without losing data", func() {
			oddTrack := testhelpers.GenerateTrack(21)
			splitTracks := oddTrack.SplitInto(2)
			Expect(len(splitTracks)).To(Equal(2))
			Expect(splitTracks[0].Length()).To(Equal(10))
			Expect(splitTracks[1].Length()).To(Equal(11))
		})
	})

	Context("SimplifyTo", func() {
		It("Simplifies a track to the correct number of points", func() {
			track.SimplifyTo(10)
			Expect(track.Length()).To(Equal(10))
		})

		It("Removes points with the smallest cross-track-error", func() {
			track = gpx.Track{
				TrackPoints: []*gpx.TrackPoint{
					&gpx.TrackPoint{Lat: 1, Lon: 1},
					&gpx.TrackPoint{Lat: 6, Lon: 6},
					&gpx.TrackPoint{Lat: 2, Lon: 2},
					&gpx.TrackPoint{Lat: 4, Lon: 4},
					&gpx.TrackPoint{Lat: 3, Lon: 3},
				},
			}
			track.SimplifyTo(4)
			Expect(*track.TrackPoints[0]).To(Equal(gpx.TrackPoint{Lat: 1, Lon: 1}))
			Expect(*track.TrackPoints[1]).To(Equal(gpx.TrackPoint{Lat: 6, Lon: 6}))
			Expect(*track.TrackPoints[2]).To(Equal(gpx.TrackPoint{Lat: 2, Lon: 2}))
			Expect(*track.TrackPoints[3]).To(Equal(gpx.TrackPoint{Lat: 3, Lon: 3}))
		})

		It("Doesn't simplify a track when track length is smaller than the requested number of points", func() {
			track.SimplifyTo(25)
			Expect(track.Length()).To(Equal(20))
		})
	})

})
