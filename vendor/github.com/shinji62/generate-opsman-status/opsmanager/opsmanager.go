package opsmanager

import "time"

const DeployedProductEndpoint = "/api/v0/deployed/products"
const DiagnosticReport = "/api/v0/diagnostic_report"
const Sessions = "/api/v0/sessions"
const UserAgent = "opsman-GenerateStatus"

type OpsManager interface {
	GetDeployedProduct()
	GetDiagnosticReport() (*BoshDiagnostic, error)
	CloseSessions() bool
}

type VersionsResource struct {
	InstallationSchemaVersion   string `json:"installation_schema_version"`
	MetadataVersion             string `json:"metadata_version"`
	ReleaseVersion              string `json:"release_version"`
	JavascriptMigrationsVersion string `json:"javascript_migrations_version"`
}

type DirectorConfigurationResource struct {
	BoshRecreateOnNextDeploy bool        `json:"bosh_recreate_on_next_deploy"`
	ResurrectorEnabled       bool        `json:"resurrector_enabled"`
	BlobstoreType            string      `json:"blobstore_type"`
	MaxThreads               interface{} `json:"max_threads"`
	DatabaseType             string      `json:"database_type"`
	NtpServers               []string    `json:"ntp_servers"`
	HmPagerDutyEnabled       bool        `json:"hm_pager_duty_enabled"`
	HmEmailerEnabled         bool        `json:"hm_emailer_enabled"`
	VMPasswordType           string      `json:"vm_password_type"`
}

type AddedProductsResource struct {
	Deployed []struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		Stemcell string `json:"stemcell"`
	} `json:"deployed"`
	Staged []struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		Stemcell string `json:"stemcell"`
	} `json:"staged"`
}

type BoshDiagnostic struct {
	Versions              VersionsResource              `json:"versions"`
	GenerationTime        time.Time                     `json:"generation_time"`
	InfrastructureType    string                        `json:"infrastructure_type"`
	DirectorConfiguration DirectorConfigurationResource `json:"director_configuration"`
	Releases              []string                      `json:"releases"`
	Stemcells             []string                      `json:"stemcells"`
	ProductTemplates      []string                      `json:"product_templates"`
	AddedProducts         AddedProductsResource         `json:"added_products"`
}
