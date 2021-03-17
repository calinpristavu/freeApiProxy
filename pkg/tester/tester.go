package tester

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/calinpristavu/freeApiProxy/pkg/jsonCompare"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Code int                    `json:"code,omitempty"`
	Body map[string]interface{} `json:"body,omitempty"`
}

type Path struct {
	Method     string                   `json:"method,omitempty"`
	Path       string                   `json:"path,omitempty"`
	Responses  []Response               `json:"responses,omitempty"`
	Parameters []map[string]interface{} `json:"parameters,omitempty"`
}

func (p Path) Test(basePath string) {
	request, err := http.NewRequest(p.Method, basePath+p.Path, nil)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path": p,
		}).Error("could not create request instance")

		return
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"response": res,
			"err":      err,
		}).Error("request failed")

		return
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("could not close res.Body stream")
		}
	}()

	var actualBody map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&actualBody); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("could not decode response body")
		return
	}

	if !jsonCompare.HaveSameKeys(p.Responses[0].Body, actualBody) {
		logrus.WithFields(logrus.Fields{
			"expected": p.Responses[0].Body,
			"actual":   actualBody,
		}).Warn("actual response does not match expected response")

		return
	}

	logrus.WithFields(logrus.Fields{
		"path": p.Path,
	}).Info("Passed!")
}

type RequestSet struct {
	BasePath string `json:"basePath,omitempty"`
	Paths    []Path `json:"paths,omitempty"`
}

func (rs RequestSet) TestAll() {
	for _, p := range rs.Paths {
		p.Test(rs.BasePath)
	}
}

func Run() {
	absPath, err := filepath.Abs("../../pkg/tester/requests.json")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("could not find requests.json file")

		return
	}

	requestSet, err := loadRequestsFromJson(absPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"fileName": "requests.json",
			"err":      err,
		}).Error("could not load request set")

		return
	}

	requestSet.TestAll()
}

func loadRequestsFromJson(fileName string) (*RequestSet, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	var request RequestSet

	if err := json.NewDecoder(file).Decode(&request); err != nil {
		return nil, err
	}

	return &request, nil
}
