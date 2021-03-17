package jsonCompare

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

func HaveSameKeys(first map[string]interface{}, second map[string]interface{}) bool {
	success := true

	for key, value1 := range first {
		value2, isSLice := second[key]

		if !isSLice {
			logrus.WithFields(logrus.Fields{
				"key": key,
			}).Error("key not found")

			success = false

			continue
		}

		if value1 == nil || value2 == nil {
			if value1 != value2 {
				logrus.WithFields(logrus.Fields{
					"key":    key,
					"value1": value1,
					"value2": value2,
				}).Error("value types don't match")

				success = false
			}

			continue
		}

		if reflect.TypeOf(value1).Name() != reflect.TypeOf(value2).Name() {
			logrus.WithFields(logrus.Fields{
				"key":    key,
				"value1": value1,
				"value2": value2,
			}).Error("value types don't match")

			success = false

			continue
		}

		childMap1, isMap := value1.(map[string]interface{})
		childMap2, _ := value2.(map[string]interface{})

		if isMap {
			success = success && HaveSameKeys(childMap1, childMap2)

			continue
		}

		slice1, isSLice := value1.([]interface{})
		slice2, isSLice := value2.([]interface{})

		if isSLice {
			success = success && compareSlices(slice1, slice2)

			continue
		}
	}

	return success
}

func compareSlices(slice1 []interface{}, slice2 []interface{}) bool {
	success := true

	for key, value1 := range slice1 {
		if reflect.TypeOf(value1).Name() != reflect.TypeOf(slice2[key]).Name() {
			logrus.WithFields(logrus.Fields{
				"key":    key,
				"value1": value1,
				"value2": slice2[key],
			}).Error("value types don't match")

			success = false

			continue
		}

		childMap1, isMap := value1.(map[string]interface{})
		childMap2, _ := slice2[key].(map[string]interface{})

		if isMap {
			success = success && HaveSameKeys(childMap1, childMap2)

			continue
		}

		childSlice1, isSLice := value1.([]interface{})
		childSlice2, isSLice := value1.([]interface{})

		if isSLice {
			success = success && compareSlices(childSlice1, childSlice2)

			continue
		}
	}

	return success
}
