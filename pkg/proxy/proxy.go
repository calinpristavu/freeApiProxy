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
	srv.Router.Path("/weather").
		Methods("GET").
		Queries(
			"q", "",
			"lat", "{lat:[0-9]+}",
			"lon", "{lon:[0-9]+}",
			"lang", "",
			"units", "",
			"mode", "",
		).
		HandlerFunc(handleWeatherRequest)

	srv.Router.Path("/find").
		Methods("GET").
		Queries(
			"q", "",
			"cnt", "{cnt:[0-9]+}",
			"mode", "",
			"lat", "{lat:[0-9]+}",
			"lon", "{lon:[0-9]+}",
			"type", "",
			"units", "",
		).
		HandlerFunc(handleFindRequest)
}

const baseUrl = "https://community-open-weather-map.p.rapidapi.com"

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

func handleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/weather?%s", baseUrl, r.URL.RawQuery),
		nil,
	)

	req.URL.Query().Encode()

	logrus.WithFields(logrus.Fields{
		"req": req,
	}).Info("Performing request to external service")

	req.Header.Add("x-rapidapi-key", credentials.key)
	req.Header.Add("x-rapidapi-host", credentials.host)

	performRequest(w, req)
}

func handleFindRequest(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/find?%s", baseUrl, r.URL.RawQuery),
		nil,
	)

	req.Header.Add("x-rapidapi-key", credentials.key)
	req.Header.Add("x-rapidapi-host", credentials.host)

	performRequest(w, req)
}

func performRequest(w http.ResponseWriter, req *http.Request) {
	res, _ := http.DefaultClient.Do(req)

	if res.StatusCode != 200 {
		logrus.WithFields(logrus.Fields{
			"request": req,
			"respons": res,
		}).Error("received status code != 200")

		w.WriteHeader(res.StatusCode)

		return
	}

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
