package validate_test

import (
	"encoding/json"
	"fmt"

	ValidateStruct "github.com/budnieswski/validate-struct"
)

func ExampleValidate() {
	type User struct {
		ID     *uint64 `json:"id" validate:"required,min=2,max=4"`
		Name   string  `json:"name" validate:"required,min=3"`
		Status bool    `json:"status"`
	}
	jsonInput := []byte(`{
			"id": 12,
			"name": "Mr Dummy",
			"status": "true"
	}`)
	validate := ValidateStruct.Validate(jsonInput, User{})
	if validate.HasError {
		json, _ := json.Marshal(validate)
		fmt.Printf("%s", string(json))
	}
	// Output:
	// {"errors":{"status":"invalid field type"},"hasError":true}
}
