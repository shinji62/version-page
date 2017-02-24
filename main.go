package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	//log "github.com/Sirupsen/logrus"
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/shinji62/generate-opsman-status/pivnet"
	"github.com/shinji62/version-page/cloudfoundry"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	iris "gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
	"gopkg.in/kataras/iris.v6/middleware/logger"
)

var (
	tickerTime        = kingpin.Flag("cc-pull-time", "CloudController Polling time in sec").Default("60s").OverrideDefaultFromEnvar("CF_PULL_TIME").Duration()
	debug             = kingpin.Flag("debug", "Enable debug mode. This disables forwarding to syslog").Default("false").OverrideDefaultFromEnvar("DEBUG").Bool()
	apiEndpoint       = kingpin.Flag("api-endpoint", "Api endpoint address. For bosh-lite installation of CF: https://api.10.244.0.34.xip.io").OverrideDefaultFromEnvar("API_ENDPOINT").Required().String()
	port              = kingpin.Flag("port", "Port to listen").Envar("PORT").Short('p').Required().Int()
	clientID          = kingpin.Flag("client-id", "Client ID.").OverrideDefaultFromEnvar("CF_CLIENT_ID").Required().String()
	clientSecret      = kingpin.Flag("client-secret", "Client secret.").OverrideDefaultFromEnvar("CF_CLIENT_SECRET").Required().String()
	skipSSLValidation = kingpin.Flag("skip-ssl-validation", "Please don't").Default("false").OverrideDefaultFromEnvar("SKIP_SSL_VALIDATION").Bool()
)

var (
	version = "0.0.0"
)

func main() {

	kingpin.Version(version)
	kingpin.Parse()
	//Cf.Client

	c := cfclient.Config{
		ApiAddress:        *apiEndpoint,
		ClientID:          *clientID,
		ClientSecret:      *clientSecret,
		SkipSslValidation: *skipSSLValidation,
		UserAgent:         "version-pages/",
	}
	cfClient, err := cfclient.NewClient(&c)
	if err != nil {
		log.Fatal("Failing setuping CfClient ", err)
		os.Exit(1)

	}

	cfApi := cloudfoundry.NewCloudFoundryApi(cfClient)
	cfApi.FetchBuildpacksApi()
	cfApi.PerformPoollingCaching(*tickerTime)
	//Load OpsManFile
	opsManFile, err := ioutil.ReadFile("./opsman/result.json")
	if err != nil {
		log.Fatal("Failing reading OpsMan file ", err)
		os.Exit(1)

	}

	var opsManJson pivnet.InfoPCF
	err = json.Unmarshal(opsManFile, &opsManJson)
	if err != nil {
		log.Fatal("Failing to Decode Json ", err)
		os.Exit(1)
	}

	app := iris.New(
		iris.OptionCharset("UTF-8"),
	)

	customLogger := logger.New()
	app.Use(customLogger)
	// output startup banner and error logs on os.Stdout
	if *debug {
		app.Adapt(iris.DevLogger())

	}

	// set the router, you can choose gorillamux too
	app.Adapt(httprouter.New())

	//Set templates
	app.Adapt(view.Django("./templates/", ".tpl"))

	app.Get("/buildpacks.json", func(ctx *iris.Context) {
		ctx.JSON(iris.StatusOK, cfApi.GetBuildpacks())

	})
	app.Get("/pcf.json", func(ctx *iris.Context) {
		ctx.JSON(iris.StatusOK, opsManJson)

	})
	app.Get("/", func(ctx *iris.Context) {
		ctx.MustRender("opsman.tpl",
			map[string]interface{}{
				"infoPcf":        opsManJson,
				"infoBuildpacks": cfApi.GetBuildpacks()})
	})
	app.Adapt(iris.EventPolicy{
		// Interrupt Event means when control+C pressed on terminal.
		Interrupted: func(*iris.Framework) {
			// shut down gracefully, but wait 5 seconds the maximum before closed
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			app.Shutdown(ctx)
		},
	})

	app.Listen(fmt.Sprintf(":%d", *port))
}
