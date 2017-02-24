package pivnet

import (
	gopivnet "github.com/pivotal-cf/go-pivnet"
)

type Pivnet interface {
	GenerateJson(string)
	CreateInfoPCF()
}

type InfoPCF struct {
	OpsManagerVersion string         `json:"opsman_version"`
	TileResources     []TileResource `json:"tiles"`
}

type TileResource struct {
	Name              string           `json:"name"`
	OpsManProductName string           `json:"name_opsman"`
	CleanVersion      string           `json:"clean_version"`
	OriginalVersion   string           `json:"original_version"`
	Release           gopivnet.Release `json:"release"`
}
