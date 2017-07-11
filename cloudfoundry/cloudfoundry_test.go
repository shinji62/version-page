package cloudfoundry_test

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/shinji62/version-page/cloudfoundry"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cloudfoundry", func() {

	var cf Cloudfoundry
	BeforeEach(func() {
		cf = &CloudfoundryApi{}

	})

	Context("with Buildpacks", func() {
		It("Should Parse Correct Buildpacks", func() {
			bp := cfclient.Buildpack{
				Guid:     "1035316f-08a2-4f76-af17-43dc132c6ffa",
				Name:     "binary_buildpack",
				Enabled:  true,
				Locked:   true,
				Filename: "binary_buildpack-cached-v1.0.9.zip",
			}
			bpInfo := BuildpackInfo{
				Name:            "binary_buildpack",
				Version:         "v1.0.9",
				Filename:        "binary_buildpack-cached-v1.0.9.zip",
				ReleaseNotesUrl: "https://github.com/cloudfoundry/binary-buildpack/releases/tag/v1.0.9",
			}
			//bps := []cfclient.Buildpack{bp}
			infobp := cf.EnhancedBuildPacks(bp)
			Expect(infobp).To(Equal(bpInfo))
		})
		It("Should Parse Correct PHP Buildpacks", func() {
			bp := cfclient.Buildpack{
				Guid:     "1035316f-08a2-4f76-af17-43dc132c6ffa",
				Name:     "php_buildpack",
				Enabled:  true,
				Locked:   true,
				Filename: "php_buildpack-cached-v4.3.36+1498674037.zip",
			}
			bpInfo := BuildpackInfo{
				Name:            "php_buildpack",
				Version:         "v4.3.36",
				Filename:        "php_buildpack-cached-v4.3.36+1498674037.zip",
				ReleaseNotesUrl: "https://github.com/cloudfoundry/php-buildpack/releases/tag/v4.3.36",
			}
			//bps := []cfclient.Buildpack{bp}
			infobp := cf.EnhancedBuildPacks(bp)
			Expect(infobp).To(Equal(bpInfo))
		})
	})

})
