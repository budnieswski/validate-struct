package validate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	Rule "github.com/budnieswski/validate-struct/rule"
)

const TAG_NAME = "validate"

var TIME_TYPE = reflect.TypeOf(time.Time{})

type ValidationResult struct {
	Errors   map[string]interface{} `json:"errors"`
	HasError bool                   `json:"hasError"`
}

func Validate(data []byte, schema interface{}) ValidationResult {
	var dataDecoded interface{}

	schemaReflectedValue := reflect.ValueOf(schema)
	schemaReflectedType := schemaReflectedValue.Type()

	if schemaReflectedType.Kind() == reflect.Ptr {
		schemaReflectedType = reflect.Indirect(schemaReflectedValue).Type()
	}

	if schemaReflectedType.Kind() != reflect.Struct {
		panic("wrong schema, only accept STRUCT as schema type")
	}

	if err := json.Unmarshal(data, &dataDecoded); err != nil {
		return ValidationResult{
			Errors: map[string]interface{}{
				"_": err.Error(),
			},
			HasError: true,
		}
	}

	dataReflected := reflect.ValueOf(dataDecoded)

	if rs := realValidateType(dataReflected, schemaReflectedType); rs != nil {
		return ValidationResult{
			Errors:   rs.(map[string]interface{}),
			HasError: true,
		}
	}

	return ValidationResult{
		Errors:   nil,
		HasError: false,
	}
}

func realValidateType(data reflect.Value, schema reflect.Type) interface{} {
	switch schema.Kind() {
	case reflect.Struct:
		errs := make(map[string]interface{}, 0)

		for i := 0; i < schema.NumField(); i++ {
			field := schema.Field(i)
			fieldType := field.Type
			fieldTag := field.Tag.Get(TAG_NAME)
			fieldName := strings.Split(field.Tag.Get("json"), ",")[0]
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

		if len(errs) < 1 {
			return nil
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
	case reflect.Ptr:
		schema = findRealTypeOfPointer(schema)

		return realValidateType(data, schema)
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

func findRealTypeOfPointer(pointer reflect.Type) reflect.Type {
	for ok := pointer.Kind() == reflect.Ptr; ok; ok = pointer.Kind() == reflect.Ptr {
		pointer = pointer.Elem()
	}

	return pointer
}

func isCompatibleType(data reflect.Value, expected reflect.Type) bool {
	if expected.Kind() == data.Kind() {
		return true
	}

	expected = findRealTypeOfPointer(expected)

	// Expected Struct and received Map
	if expected.Kind() == reflect.Struct && data.Kind() == reflect.Map {
		return true
	}

	// Check Int variations
	// As it is an internal method, every numeric data value
	// is float64 provided by json unmarshal
	switch expected.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dataString := strings.Replace(fmt.Sprintf("%f", data.Interface()), ".000000", "", 1)
		dataInt, err := strconv.ParseInt(dataString, 10, 64)

		return err == nil && !reflect.Zero(expected).OverflowInt(dataInt)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		dataString := strings.Replace(fmt.Sprintf("%f", data.Interface()), ".000000", "", 1)
		dataInt, err := strconv.ParseUint(dataString, 10, 64)

		return err == nil && !reflect.Zero(expected).OverflowUint(dataInt)
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
	case "ip":
		return Rule.IPRule{}
	case "ipv4":
		return Rule.IPv4Rule{}
	case "ipv6":
		return Rule.IPv6Rule{}
	}

	return Rule.DefaultRule{}
}
