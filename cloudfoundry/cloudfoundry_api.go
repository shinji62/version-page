package cloudfoundry

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

type CloudfoundryApi struct {
	client        *cfclient.Client
	mutex         *sync.Mutex
	buildpackInfo []BuildpackInfo
}

func NewCloudFoundryApi(clientCF *cfclient.Client) Cloudfoundry {
	return &CloudfoundryApi{
		client:        clientCF,
		mutex:         &sync.Mutex{},
		buildpackInfo: []BuildpackInfo{},
	}
}

func (cf *CloudfoundryApi) FetchBuildpacksApi() {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()
	listBps, err := cf.client.ListBuildpacks()
	// TODO: Login Error
	if err != nil {
		fmt.Println(err)
		return

	}
	cf.buildpackInfo = nil
	for _, bp := range listBps {
		cf.buildpackInfo = append(cf.buildpackInfo, cf.EnhancedBuildPacks(bp))
	}
}

func (cf *CloudfoundryApi) GetBuildpacks() []BuildpackInfo {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()
	return cf.buildpackInfo
}

func (cf *CloudfoundryApi) PerformPoollingCaching(tickerTime time.Duration) {
	// Ticker Pooling the CC every X sec
	ccPooling := time.NewTicker(tickerTime)

	go func() {
		for range ccPooling.C {
			cf.FetchBuildpacksApi()
		}
	}()
}

func (cf *CloudfoundryApi) EnhancedBuildPacks(bpCf cfclient.Buildpack) BuildpackInfo {
	re := regexp.MustCompile(buildpackRegexp)
	var bp BuildpackInfo
	bp.Filename = bpCf.Filename
	bp.Name = bpCf.Name
	if founds := re.FindStringSubmatch(bp.Filename); founds != nil {
		bp.Version = founds[3]
		bp.ReleaseNotesUrl = buildPackReleaseBaseUrl + strings.Replace(founds[1], "_", "-", 2) + "/releases/tag/" + bp.Version
	}
	return bp
}
