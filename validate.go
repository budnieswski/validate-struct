package validate

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	Rule "github.com/budnieswski/validate-struct/rule"
)

const TAG_NAME = "validate"

var TIME_TYPE = reflect.TypeOf(time.Time{})

func Validate(data []byte, schema interface{}) map[string]interface{} {
	var dataDecoded interface{}

	schemaReflectedType := reflect.TypeOf(schema)

	if schemaReflectedType.Kind() != reflect.Struct {
		panic("wrong schema, only accept STRUCT as schema type")
	}

	if err := json.Unmarshal(data, &dataDecoded); err != nil {
		return map[string]interface{}{
			"_": err.Error(),
		}
	}

	dataReflected := reflect.ValueOf(dataDecoded)

	if rs := realValidateType(dataReflected, schemaReflectedType); rs != nil {
		return rs.(map[string]interface{})
	}

	return map[string]interface{}{}
	// return realValidateType(dataReflected, schemaReflectedType).(map[string]interface{})
}

func realValidateType(data reflect.Value, schema reflect.Type) interface{} {
	switch schema.Kind() {
	case reflect.Struct:
		errs := make(map[string]interface{}, 0)

		for i := 0; i < schema.NumField(); i++ {
			field := schema.Field(i)
			fieldType := field.Type
			fieldTag := field.Tag.Get(TAG_NAME)
			fieldName := field.Tag.Get("json")
			dataValue := findOnMap(data, fieldName)

			if fieldTag != "" {
				failed := validateTag(dataValue, fieldTag)

				if reflect.ValueOf(failed).Len() > 0 {
					errs[fieldName] = failed
					continue
				}
			}

			if dataValue.Kind() == reflect.Invalid || fieldType == TIME_TYPE {
				continue
			}

			if !isCompatibleType(dataValue, fieldType) {
				errs[fieldName] = "invalid field type"
				continue
			}

			rs := realValidateType(dataValue, fieldType)

			if rs != nil && reflect.ValueOf(rs).Len() > 0 {
				errs[fieldName] = rs
			}
		}

		return errs
	case reflect.Array, reflect.Slice:
		errs := make(map[int]interface{}, 0)
		schemaElem := schema.Elem()

		for j := 0; j < data.Len(); j++ {
			itemValue := data.Index(j).Elem()

			if !isCompatibleType(itemValue, schemaElem) {
				errs[j] = "invalid field type"
				continue
			}

			rs := realValidateType(itemValue, schemaElem)
			if rs != nil && reflect.ValueOf(rs).Len() > 0 {
				errs[j] = rs
			}
		}

		return errs
	case reflect.Map:
		errs := make(map[string]interface{}, 0)
		schemaElem := schema.Elem()

		for _, v := range data.MapKeys() {
			itemValue := data.MapIndex(v).Elem()

			if !isCompatibleType(itemValue, schemaElem) {
				errs[v.String()] = "invalid field type"
				continue
			}

			rs := realValidateType(itemValue, schemaElem)
			if rs != nil && reflect.ValueOf(rs).Len() > 0 {
				errs[v.String()] = rs
			}
		}

		return errs
	default:
		// fmt.Printf("%*s >> Its a default\n", 4, "")
		// fmt.Printf("%*s DataValue: (%s) - %v\n", 8, "", data.Kind(), data)

		return nil
	}
}

func findOnMap(data reflect.Value, name string) reflect.Value {
	if data.Kind() != reflect.Map {
		return reflect.ValueOf(nil)
	}

	for _, e := range data.MapKeys() {
		if e.String() == name {
			return data.MapIndex(e).Elem()
		}
	}

	return reflect.ValueOf(nil)
}

func isCompatibleType(data reflect.Value, expected reflect.Type) bool {
	if expected.Kind() == data.Kind() {
		return true
	}

	if expected.Kind() == reflect.Struct && data.Kind() == reflect.Map {
		return true
	}

	return data.CanConvert(expected)
}

func validateTag(data reflect.Value, tag string) string {
	rules := strings.Split(tag, ",")

	for _, rule := range rules {
		valid, err := validateRule(data, rule)

		if !valid {
			return err.Error()
		}
	}

	return ""
}

func validateRule(data reflect.Value, rule string) (bool, error) {
	validator := getRuleValidator(rule)

	return validator.Validate(data)
}

func getRuleValidator(rule string) Rule.Interface {
	args := strings.Split(rule, "=")

	switch args[0] {
	case "required":
		return Rule.RequiredRule{}
	case "min":
		validator := Rule.MinRule{}
		min, _ := strconv.Atoi(args[1])
		validator.Min = min
		return validator
	case "max":
		validator := Rule.MaxRule{}
		max, _ := strconv.Atoi(args[1])
		validator.Max = max
		return validator
	case "date":
		var format string

		if len(args) > 1 {
			format = args[1]
		}

		return Rule.DateRule{
			Format: format,
		}
	case "enum":
		var validValues []string

		if len(args) > 1 {
			validValues = strings.Split(args[1], "|")
		}

		return Rule.ENumRule{
			ValidValues: validValues,
		}
	}

	return Rule.DefaultRule{}
}
