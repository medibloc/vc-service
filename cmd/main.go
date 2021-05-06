package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/gorilla/mux"
	"github.com/medibloc/vc-service/pkg"
	log "github.com/sirupsen/logrus"
)

func main() {
	var config pkg.Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("failed to process env vars: %v", err)
	}

	pkg.InitLog(&config)

	r := mux.NewRouter()
	pkg.RegisterHandlers(r)

	svr := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", config.Port),
		WriteTimeout: config.ReadTimeout,
		ReadTimeout:  config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
		Handler:      r,
	}

	go func() {
		log.Infof("listening HTTP: %v", svr.Addr)
		if err := svr.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	log.Infof("shutting down gracefully...")
	svr.Shutdown(ctx)
}
