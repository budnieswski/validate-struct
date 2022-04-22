package validate

import (
	"encoding/json"
	"reflect"
	"testing"

	Rule "github.com/budnieswski/validate-struct/rule"
)

func TestValidate(test *testing.T) {
	var fixtureInputMap interface{}
	var fixtureInputReflected reflect.Value
	fixtureInputByte := []byte(`
		{
			"id": 12,
			"name": "Mr Dummy",
			"status": true
		}
	`)

	if err := json.Unmarshal(fixtureInputByte, &fixtureInputMap); err != nil {
		test.Errorf("Failed to decode fixture: %s", err.Error())
	}

	fixtureInputReflected = reflect.ValueOf(fixtureInputMap)

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

		test.Run("Should return a invalid value when input data is correctly", func(test *testing.T) {
			field := "_id"
			find := findOnMap(fixtureInputReflected, field)

			if find.Kind() != reflect.Invalid {
				test.Errorf("Found value: %d", int(find.Float()))
			}
		})

		test.Run("Should return a invalid value when input data is not correctly", func(test *testing.T) {
			field := "_id"
			input := reflect.ValueOf("")
			find := findOnMap(input, field)

			if find.Kind() != reflect.Invalid {
				test.Errorf("Found value: %d", int(find.Float()))
			}
		})
	})

	test.Run("validateRule", func(test *testing.T) {
		test.Run("Should return a true to Required validation", func(test *testing.T) {
			value := reflect.ValueOf("foo")
			valid, err := validateRule(value, "required")

			if err != nil {
				test.Errorf("Validation failed: %s", err.Error())
				return
			}

			if valid == false {
				test.Errorf("Failed on validation")
			}
		})

		test.Run("Should return a false to Required validation", func(test *testing.T) {
			value := reflect.ValueOf("")
			valid, err := validateRule(value, "required")

			if err == nil {
				test.Error("Validation must be failed but returned nil error")
				return
			}

			if valid {
				test.Errorf("Failed on validation")
			}
		})
	})

	test.Run("getRuleValidator", func(test *testing.T) {
		test.Run("Should return a DefaultRule", func(test *testing.T) {
			expectedType := reflect.TypeOf(Rule.DefaultRule{})
			givenType := reflect.TypeOf(getRuleValidator("_foo_"))

			if givenType != expectedType {
				test.Errorf("Expected: %s -- Found: %s", expectedType, givenType)
			}
		})

		test.Run("Should return a RequiredRule", func(test *testing.T) {
			expectedType := reflect.TypeOf(Rule.RequiredRule{})
			givenType := reflect.TypeOf(getRuleValidator("required"))

			if givenType != expectedType {
				test.Errorf("Expected: %s -- Found: %s", expectedType, givenType)
			}
		})

		test.Run("Should return a MinRule", func(test *testing.T) {
			expectedType := reflect.TypeOf(Rule.MinRule{})
			givenType := reflect.TypeOf(getRuleValidator("min=2"))

			if givenType != expectedType {
				test.Errorf("Expected: %s -- Found: %s", expectedType, givenType)
			}
		})

		test.Run("Should return a MaxRule", func(test *testing.T) {
			expectedType := reflect.TypeOf(Rule.MaxRule{})
			givenType := reflect.TypeOf(getRuleValidator("max=2"))

			if givenType != expectedType {
				test.Errorf("Expected: %s -- Found: %s", expectedType, givenType)
			}
		})

		test.Run("Should return a DateRule", func(test *testing.T) {
			expectedType := reflect.TypeOf(Rule.DateRule{})
			givenType := reflect.TypeOf(getRuleValidator("date=timestamp"))

			if givenType != expectedType {
				test.Errorf("Expected: %s -- Found: %s", expectedType, givenType)
			}
		})

		test.Run("Should return a ENumRule", func(test *testing.T) {
			expectedType := reflect.TypeOf(Rule.ENumRule{})
			givenType := reflect.TypeOf(getRuleValidator("enum=US|BR|CA"))

			if givenType != expectedType {
				test.Errorf("Expected: %s -- Found: %s", expectedType, givenType)
			}
		})
	})

	test.Run("isCompatibleType", func(test *testing.T) {
		test.Run("Should return true when given are compatible", func(test *testing.T) {
			value := reflect.ValueOf(123.0)
			expected := reflect.TypeOf(123.45)
			compatible := isCompatibleType(value, expected)

			if !compatible {
				test.Errorf("Values are not compatible, expected: %s - given: %s", expected.Kind(), value.Kind())
			}
		})

		test.Run("Should return false when given are not compatible", func(test *testing.T) {
			value := reflect.ValueOf(123.45)
			expected := reflect.TypeOf(123)
			compatible := isCompatibleType(value, expected)

			if compatible != false {
				test.Errorf("Values are compatible, expected: %s - given: %s", expected.Kind(), value.Kind())
			}
		})

		test.Run("Should return false when given negative int", func(test *testing.T) {
			value := reflect.ValueOf(-123.0)
			expected := reflect.TypeOf(uint(1))
			compatible := isCompatibleType(value, expected)

			if compatible != false {
				test.Errorf("Values are compatible, expected: %s - given: %s", expected.Kind(), value.Kind())
			}
		})

		test.Run("Should return false when given value overflow bit", func(test *testing.T) {
			value := reflect.ValueOf(256.0)
			expected := reflect.TypeOf(uint8(1))
			compatible := isCompatibleType(value, expected)

			if compatible != false {
				test.Errorf("Values are compatible, expected: %s - given: %s %v", expected.Kind(), value.Kind(), compatible)
			}
		})

		test.Run("Should return true when values types are equal", func(test *testing.T) {
			value := reflect.ValueOf(123.0)
			expected := reflect.TypeOf(123)
			compatible := isCompatibleType(value, expected)

			if !compatible {
				test.Errorf("Values are not compatible, expected: %s - given: %s", expected.Kind(), value.Kind())
			}
		})

		test.Run("Should return true when given Map as input and Struct as expected", func(test *testing.T) {
			value := reflect.ValueOf(map[string]string{})
			expected := reflect.TypeOf(struct{}{})
			compatible := isCompatibleType(value, expected)

			if !compatible {
				test.Errorf("Values are not compatible, expected: %s - given: %s", expected.Kind(), value.Kind())
			}
		})

		test.Run("Should return false when values are not compatible", func(test *testing.T) {
			value := reflect.ValueOf("123")
			expected := reflect.TypeOf(123.45)
			compatible := isCompatibleType(value, expected)

			if compatible {
				test.Errorf("Values are not compatible, expected: %s - given: %s", expected.Kind(), value.Kind())
			}
		})
	})
}
