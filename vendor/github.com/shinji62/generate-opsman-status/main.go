package main

import (
	"fmt"
	"log"
	"net/http"

	gopivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/shinji62/generate-opsman-status/opsmanager"
	"github.com/shinji62/generate-opsman-status/pivnet"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	opsmanagerEndpoint     = kingpin.Flag("opsman-endpoint", "OpsMan Endpoint").OverrideDefaultFromEnvar("OPSMAN_ENDPOINT").Required().String()
	clientIDOpsManagager   = kingpin.Flag("client-id-opsman", "Client ID.").OverrideDefaultFromEnvar("OPSMAN_CLIENT_ID").Required().String()
	clientSecretOpsManager = kingpin.Flag("client-secret-opsman", "Client secret.").OverrideDefaultFromEnvar("OPSMAN_CLIENT_SECRET").Required().String()
	pivnetToken            = kingpin.Flag("pivnet-api-token", "Pivnet API token").OverrideDefaultFromEnvar("PIVNET_TOKEN").Required().String()
	skipSSLValidation      = kingpin.Flag("skip-ssl-validation", "Please don't").Default("false").OverrideDefaultFromEnvar("SKIP_SSL_VALIDATION").Bool()
	pathProf               = kingpin.Flag("file-output", "Set the File output like ./path/opsman.json by default ./result.json").Default("result").OverrideDefaultFromEnvar("RESULT_FILE").String()
)

func main() {

	kingpin.Parse()

	//SetUp HTTPClient
	hClientOpsMan := http.DefaultClient

	//opsMan Client
	opsClient := opsmanager.NewOpsManagerApi(hClientOpsMan, *opsmanagerEndpoint, *clientIDOpsManagager, *clientSecretOpsManager, *skipSSLValidation)
	if !opsClient.CloseSessions() {
		log.Fatal("Could not logout user from OpsManager")
	}
	fmt.Println("All users have been disconnected from OpsManager")
	opsClient = opsmanager.NewOpsManagerApi(hClientOpsMan, *opsmanagerEndpoint, *clientIDOpsManagager, *clientSecretOpsManager, *skipSSLValidation)

	boshDiag, err := opsClient.GetDiagnosticReport()
	if err != nil {
		log.Fatalf("Could not get Diagnostic Report from OpsManager ", err)
	}

	if !opsClient.CloseSessions() {
		fmt.Println("Could not logout user from OpsManager, after operation")
	}
	fmt.Println("Users have been disconnected from OpsManager")

	//Pivnet Client
	config := gopivnet.ClientConfig{
		Host:      gopivnet.DefaultHost,
		Token:     *pivnetToken,
		UserAgent: opsmanager.UserAgent,
	}
	client := pivnet.NewPivnetClient(config, boshDiag)
	client.CreateInfoPCF()
	client.GenerateJson(*pathProf)

}
