package validate

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestValidate(test *testing.T) {
	var fixtureInputMap interface{}
	fixtureInputByte := []byte(`
		{
			"id": 12,
			"name": "Mr Dummy",
			"status": true
		}
	`)
	json.Unmarshal(fixtureInputByte, &fixtureInputMap)
	fixtureInputReflected := reflect.ValueOf(fixtureInputMap)

	test.Run("findOnMap", func(test *testing.T) {
		test.Run("Should return a correctly field value", func(test *testing.T) {
			field := "id"
			expectedValue := int(12)
			find := findOnMap(fixtureInputReflected, field)

			if find.Kind() == reflect.Invalid {
				test.Error("Couldn't get the value")
				return
			}

			if int(find.Float()) != expectedValue {
				test.Errorf("Expected: %d -- Found: %d", expectedValue, int(find.Float()))
			}
		})

		test.Run("Should return a invalid value", func(test *testing.T) {
			field := "_id"
			find := findOnMap(fixtureInputReflected, field)

			if find.Kind() != reflect.Invalid {
				test.Errorf("Found value: %d", int(find.Float()))
			}
		})
	})
}
