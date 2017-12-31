package commands_test

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rcw5/gpx-tools/commands"
	"github.com/rcw5/gpx-tools/testhelpers"
)

var _ = Describe("Simplify Command", func() {
	var tempDir string
	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "gpx-simplifier")
		Expect(err).ToNot(HaveOccurred())
	})
	AfterEach(func() {
		os.RemoveAll(tempDir)
	})
	Context("Simplify", func() {
		It("Simplifies a GPX file", func() {
			gpxFilePath := fmt.Sprintf("%s/simplify_file.gpx", tempDir)
			commands.Simplify("fixtures/simplify_file.gpx", gpxFilePath, 10)

			expectedSimplify, err := testhelpers.GetAsset("fixtures/expected_simplify_file.gpx")
			Expect(err).ToNot(HaveOccurred())
			actualSimplify, err := ioutil.ReadFile(gpxFilePath)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(actualSimplify)).To(MatchXML(expectedSimplify))
		})
		PIt("Provides a suitable default for the output file")
	})
})
