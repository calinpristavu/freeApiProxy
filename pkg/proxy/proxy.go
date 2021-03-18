package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/calinpristavu/freeApiProxy/pkg/server"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func Boot(srv *server.Server) {
	srv.Router.Methods("GET").Path("/weather").HandlerFunc(handleWeatherRequest)
	srv.Router.Methods("GET").Path("/find").HandlerFunc(handleFindRequest)
}

const baseUrl = "https://community-open-weather-map.p.rapidapi.com/"

var credentials = struct {
	key  string
	host string
}{}

func init() {
	if err := godotenv.Load("../../.env.local"); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Info("could not load .env.local file")
	}

	if err := godotenv.Load("../../.env"); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Fatal("could not load .env file")
	}

	credentials.key = os.Getenv("X-RAPIDAPI-KEY")
	credentials.host = os.Getenv("X-RAPIDAPI-HOST")
}

func handleWeatherRequest(w http.ResponseWriter, _ *http.Request) {
	req, _ := http.NewRequest(
		"GET",
		baseUrl+"weather?q=London%2Cuk&lat=0&lon=0&id=2172797&lang=null&units=%22metric%22%20or%20%22imperial%22&mode=xml%2C%20html",
		nil,
	)

	logrus.WithFields(logrus.Fields{
		"req": req,
	}).Info("Performing request to external service")

	req.Header.Add("x-rapidapi-key", credentials.key)
	req.Header.Add("x-rapidapi-host", credentials.host)

	performRequest(w, req)
}

func handleFindRequest(w http.ResponseWriter, _ *http.Request) {
	req, _ := http.NewRequest(
		"GET",
		baseUrl+"/find?q=london&cnt=2&mode=null&lon=0&type=link%2C%20accurate&lat=0&units=imperial%2C%20metric",
		nil,
	)

	req.Header.Add("x-rapidapi-key", "dadb941facmsh19814c5df2350f9p1eacd3jsne11fcea3b143")
	req.Header.Add("x-rapidapi-host", "community-open-weather-map.p.rapidapi.com")

	performRequest(w, req)
}

func performRequest(w http.ResponseWriter, req *http.Request) {
	res, _ := http.DefaultClient.Do(req)

	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("could not close res.Body stream")
		}
	}()
	body, _ := ioutil.ReadAll(res.Body)

	w.Header().Add("Accept", "application/json")

	if _, err := fmt.Fprint(w, string(body)); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("could not write response to writer")
	}
}
