# Validate Struct

## Install
```
go get github.com/budnieswski/validate-struct
```

## Example
```golang
package main

import (
	"encoding/json"
	"fmt"

	ValidateStruct "github.com/budnieswski/validate-struct"
)

type User struct {
	ID     uint64 `json:"id" validate:"required,min=2,max=4"`
	Name   string `json:"name" validate:"required,min=3"`
	Status bool   `json:"status"`
}

func main() {
	input := []byte(`
		{
			"id": 12,
			"name": "Mr Dummy",
			"status": true
		}
	`)

	valid := ValidateStruct.Validate(input, User{})

	if len(valid) > 0 {
		fmt.Println("Validation failed")
		resp, _ := json.Marshal(valid)
		fmt.Println(string(resp))
	}
}
```