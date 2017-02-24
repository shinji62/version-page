package opsmanager

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type OpsManagerApi struct {
	client       *http.Client
	targetOpsMan string
}

func NewOpsManagerApi(client *http.Client, target string, clientId string, clientSecret string, skipSsLValidation bool) OpsManager {
	ctx := context.Background()
	if skipSsLValidation == false {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, client)
	} else {

		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, client)
	}

	authConfig := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     target + "/uaa/oauth/token",
	}

	return &OpsManagerApi{
		targetOpsMan: target,
		client:       authConfig.Client(ctx),
	}

}

func (o *OpsManagerApi) GetDeployedProduct() {
	req, err := http.NewRequest("GET", o.targetOpsMan+DeployedProductEndpoint, nil)
	resp, err := o.doRequest(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(o.decodeBody(resp, nil))
}

func (o *OpsManagerApi) GetDiagnosticReport() (*BoshDiagnostic, error) {
	var boshDiagnostic BoshDiagnostic
	req, err := http.NewRequest("GET", o.targetOpsMan+DiagnosticReport, nil)
	resp, err := o.doRequest(req)
	if err != nil {
		return &BoshDiagnostic{}, err
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &BoshDiagnostic{}, err
	}
	defer resp.Body.Close()
	err = json.Unmarshal(resBody, &boshDiagnostic)
	if err != nil {
		return &BoshDiagnostic{}, err
	}
	return &boshDiagnostic, nil
}

// DoRequest runs a request with our client
func (o *OpsManagerApi) doRequest(req *http.Request) (*http.Response, error) {

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-type", "application/json")
	//fmt.Println(req)
	resp, err := o.client.Do(req)
	return resp, err
}

// decodeBody is used to JSON decode a body
func (o *OpsManagerApi) decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}

func (o *OpsManagerApi) CloseSessions() bool {

	req, err := http.NewRequest("DELETE", o.targetOpsMan+Sessions, bytes.NewBufferString("{}"))
	resp, err := o.doRequest(req)

	if err != nil {
		fmt.Println(resp)
		return false
	}

	return resp.StatusCode == 200

}

func (o *OpsManagerApi) GetErt(BoshD *BoshDiagnostic) {

	return

}
