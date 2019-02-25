package main

import (
	"context"
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/plusspeed/payments-api/internal/api"
	"github.com/plusspeed/payments-api/internal/repository"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	appName = "payment-api"
	appDesc = "Allows create, read, update, get and list operations. Uses postgresql and got a health check."
)

func main() {
	app := cli.App(appName, appDesc)

	// HTTP
	port := app.Int(cli.IntOpt{
		Name:   "port",
		Desc:   "HTTP port for the app",
		EnvVar: "PORT",
		Value:  8081,
	})
	pathPrefix := app.String(cli.StringOpt{
		Name:   "path-prefix",
		Desc:   "Version of the API start with a /. The endpoints created will start by it.",
		EnvVar: "path-prefix",
		Value:  "/v1",
	})
	writeTimeSec := app.Int(cli.IntOpt{
		Name:   "write-timeout",
		Desc:   "number of seconds the http call waits writing until it times out.",
		EnvVar: "WRITE_TIMEOUT",
		Value:  10,
	})
	readTimeSec := app.Int(cli.IntOpt{
		Name:   "read-timeout",
		Desc:   "number of seconds the http call waits reading until it times out.",
		EnvVar: "READ_TIMEOUT",
		Value:  10,
	})
	idleTimeSec := app.Int(cli.IntOpt{
		Name:   "idle-timeout",
		Desc:   "number of seconds the http call waits idling until it times out.",
		EnvVar: "IDLE_TIMEOUT",
		Value:  10,
	})

	//Postgres
	pgAddress := app.String(cli.StringOpt{
		Name:   "db-address",
		Desc:   "the db address with the port number - eg.  127.0.0.1:5432",
		EnvVar: "DB_ADDRESS",
		Value:  "127.0.0.1:5432",
	})
	pgUsername := app.String(cli.StringOpt{
		Name:   "db-username",
		Desc:   "postgresql username",
		EnvVar: "DB_USERNAME",
		Value:  "test",
	})
	pgPassword := app.String(cli.StringOpt{
		Name:   "db-password",
		Desc:   "postgresql password",
		EnvVar: "DB_PASSWORD",
		Value:  "example",
	})
	dbName := app.String(cli.StringOpt{
		Name:   "db-name",
		Desc:   "the name of the database",
		EnvVar: "DB_NAME",
		Value:  "test",
	})

	//Service
	logLevel := app.String(cli.StringOpt{
		Name:   "log-level",
		Desc:   "Desired log level, - eg. info, warn, error",
		EnvVar: "LOG_LEVEL",
		Value:  "debug",
	})
	gracefulTimeSec := app.Int(cli.IntOpt{
		Name:   "graceful-timeout",
		Desc:   "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m",
		EnvVar: "GRACEFUL_TIMEOUT",
		Value:  10,
	})

	app.Before = func() {
		lvl, err := log.ParseLevel(*logLevel)
		if err != nil {
			log.WithError(err).Panic("error setting loglevel")
		}
		log.SetLevel(lvl)
	}
	app.Action = func() {

		//Created a new Repository.
		repo := repository.New(*pgAddress, *dbName, *pgUsername, *pgPassword)
		defer func() {
			err := repo.Database.Close()
			if err != nil {
				log.WithError(err).Panic("error closing db")
			}
		}()

		//Create a mux router
		router := api.NewRouter(*pathPrefix, repo)

		//Creates a http server with handler as the router
		addr := fmt.Sprintf("127.0.0.1:%d", *port)
		srv := http.Server{
			Addr:         addr,
			WriteTimeout: time.Duration(*writeTimeSec) * time.Second,
			ReadTimeout:  time.Duration(*readTimeSec) * time.Second,
			IdleTimeout:  time.Duration(*idleTimeSec) * time.Second,
			Handler:      router,
		}
		defer func() {
			err := srv.Close()
			if err != nil {
				log.WithError(err).Panic("error closing http server")
			}
		}()

		//Starts http server
		go func() {
			log.Infof("starting the http server in the address %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				log.WithError(err).Panicf("http server error")
			}
		}()

		waitForShutdown(gracefulTimeSec, &srv)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Panicf("app failed to run")
	}
}

func waitForShutdown(gracefulTimeSec *int, srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*gracefulTimeSec)*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	log.Info("shutting down")
	os.Exit(0)
}
