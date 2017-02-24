package cloudfoundry

import (
	"time"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

const buildPackReleaseBaseUrl = "https://github.com/cloudfoundry/"
const buildpackRegexp = "(.*)-(cached|offline)-(v[0-9-.]+).zip"

type BuildpackInfo struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Filename        string `json:"filename"`
	ReleaseNotesUrl string `json:"release_notes_url"`
}

type Cloudfoundry interface {
	GetBuildpacks() []BuildpackInfo
	EnhancedBuildPacks(cfclient.Buildpack) BuildpackInfo
	PerformPoollingCaching(time.Duration)
	FetchBuildpacksApi()
}
