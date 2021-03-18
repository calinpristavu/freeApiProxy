package tester

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

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

	requestSet.testAll()
}

func loadRequestsFromJson(fileName string) (*requestSet, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	var request requestSet

	if err := json.NewDecoder(file).Decode(&request); err != nil {
		return nil, err
	}

	return &request, nil
}
