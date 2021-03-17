package main

import (
	"github.com/calinpristavu/freeApiProxy/pkg/tester"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)

	tester.Run()
}
