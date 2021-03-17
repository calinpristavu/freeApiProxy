package main

import (
	"net/http"

	"github.com/calinpristavu/freeApiProxy/pkg/proxy"
	"github.com/calinpristavu/freeApiProxy/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	srv := server.New()

	proxy.Boot(srv)

	logrus.SetReportCaller(true)

	logrus.WithFields(logrus.Fields{
		"host": "localhost",
		"port": "8080",
	}).Info("starting webserver...")

	if err := http.ListenAndServe(":8080", srv.Router); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Fatal("could not start ws")
	}
}
