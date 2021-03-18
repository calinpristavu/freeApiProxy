package tester

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/calinpristavu/freeApiProxy/pkg/jsonCompare"
	"github.com/sirupsen/logrus"
)

type responseBody map[string]interface{}

type response struct {
	Code int          `json:"code,omitempty"`
	Body responseBody `json:"body,omitempty"`
}

type path struct {
	Method     string                   `json:"method,omitempty"`
	Path       string                   `json:"path,omitempty"`
	Responses  []response               `json:"responses,omitempty"`
	Parameters []map[string]interface{} `json:"parameters,omitempty"`
}

type requestSet struct {
	BasePath string `json:"basePath,omitempty"`
	Paths    []path `json:"paths,omitempty"`
}

func (rb responseBody) matchesResponse(response io.ReadCloser) bool {
	var actualBody map[string]interface{}

	if err := json.NewDecoder(response).Decode(&actualBody); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("could not decode response body")

		return false
	}

	if !jsonCompare.HaveSameKeys(rb, actualBody) {
		logrus.WithFields(logrus.Fields{
			"expected": rb,
			"actual":   actualBody,
		}).Warn("actual response does not match expected response")

		return false
	}

	return true
}

func (p path) test(basePath string, wg *sync.WaitGroup) {
	defer wg.Done()
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

	defer res.Body.Close()

	if res.StatusCode != p.Responses[0].Code {
		logrus.WithFields(logrus.Fields{
			"response": res,
			"expectedCode": p.Responses[0].Code,
			"actualCode": res.StatusCode,
		}).Error("request failed")

		return
	}

	if p.Responses[0].Body.matchesResponse(res.Body) {
		logrus.WithFields(logrus.Fields{
			"path": p.Path,
		}).Info("Passed!")
	}
}

func (rs requestSet) testAll() {
	var wg sync.WaitGroup

	wg.Add(len(rs.Paths))

	for _, p := range rs.Paths {
		go p.test(rs.BasePath, &wg)
	}

	wg.Wait()
}
