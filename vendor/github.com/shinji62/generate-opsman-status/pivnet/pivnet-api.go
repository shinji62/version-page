package pivnet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/blang/semver"
	gopivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logshim"
	"github.com/pivotal-cf/pivnet-resource/gp"

	"github.com/shinji62/generate-opsman-status/opsmanager"
)

type PivnetApi struct {
	Client  *gp.Client
	bo      *opsmanager.BoshDiagnostic
	InfoPCF *InfoPCF
}

func NewPivnetClient(config gopivnet.ClientConfig, boshOpsMan *opsmanager.BoshDiagnostic) Pivnet {

	stdoutLogger := log.New(os.Stdout, "", log.LstdFlags)
	stderrLogger := log.New(os.Stderr, "", log.LstdFlags)
	logger := logshim.NewLogShim(stdoutLogger, stderrLogger, false)

	return &PivnetApi{
		Client:  gp.NewClient(config, logger),
		bo:      boshOpsMan,
		InfoPCF: &InfoPCF{},
	}
}

func (pa *PivnetApi) GenerateJson(outputFile string) {
	//
	//f, err := os.Create(pa.file)
	b, err := json.Marshal(pa.InfoPCF)
	if err == nil {
		err = ioutil.WriteFile(outputFile, b, 0644)
	}
}

func (pa *PivnetApi) getVersion(v string) string {
	version, err := semver.Parse(v)
	if err != nil {
		return v
	}
	return fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)
}

func (pa *PivnetApi) CreateInfoPCF() {
	pa.InfoPCF.OpsManagerVersion = pa.bo.Versions.ReleaseVersion
	pa.InfoPCF.TileResources = pa.formatTiles()
}

func (pa *PivnetApi) formatTiles() []TileResource {
	tileResources := []TileResource{}
	for _, v := range pa.bo.AddedProducts.Deployed {
		tileResource := TileResource{}
		tileResource.CleanVersion = pa.getVersion(v.Version)
		tileResource.OriginalVersion = v.Version
		tileResource.Name = pa.mapProductTile(v.Name)
		tileResource.OpsManProductName = v.Name
		tileResource.Release = pa.GetPivnetInfo(tileResource.Name, tileResource.CleanVersion)
		tileResources = append(tileResources, tileResource)

	}
	return tileResources
}

func (pa *PivnetApi) mapProductTile(opsmanName string) string {
	switch opsmanName {
	case "cf":
		return "elastic-runtime"
	default:
		return opsmanName
	}
}

func (pa *PivnetApi) GetPivnetInfo(productSlug string, release string) gopivnet.Release {
	var foundRelease gopivnet.Release
	foundRelease, err := pa.Client.GetRelease(productSlug, release)
	if err != nil {
		return gopivnet.Release{}
	}
	return foundRelease
}
