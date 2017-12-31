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

var _ = Describe("Split Command", func() {
	var tempDir string
	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "gpx-simplifier")
		Expect(err).ToNot(HaveOccurred())
	})
	AfterEach(func() {
		os.RemoveAll(tempDir)
	})
	Context("Split", func() {
		It("Splits a GPX file into multiple parts", func() {
			commands.Split("fixtures/split_file.gpx", tempDir, 2)
			for i := 1; i <= 2; i++ {
				expectedFile, err := testhelpers.GetAsset(fmt.Sprintf("fixtures/expected_split_file_%d.gpx", i))
				Expect(err).ToNot(HaveOccurred())
				actualFile, err := ioutil.ReadFile(fmt.Sprintf("%s/split_file_%d.gpx", tempDir, i))
				Expect(err).ToNot(HaveOccurred())
				Expect(string(actualFile)).To(MatchXML(expectedFile))
			}
		})
		PIt("Provides a suitable default for the output directory")
	})
})
