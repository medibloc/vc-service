package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/medibloc/vc-service/pkg/aries"
	"github.com/medibloc/vc-service/pkg/config"
	"github.com/medibloc/vc-service/pkg/rest/kms"
	"github.com/medibloc/vc-service/pkg/rest/verifiable"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := envconfig.Process("", config.Cfg); err != nil {
		log.Fatal(err)
	}

	ariesProvider, err := aries.NewProvider()
	if err != nil {
		log.Fatal(err)
	}

	restRouter := echo.New()
	verifiable.RegisterHandlers(restRouter, ariesProvider)
	kms.RegisterHandlers(restRouter, ariesProvider)

	//TODO: graceful shutdown
	log.Fatal(restRouter.Start(":9991"))
}
