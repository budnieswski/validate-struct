# Validate Struct
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/budnieswski/validate-struct?style=flat-square)](https://github.com/budnieswski/validate-struct/tags)
[![Coveralls](https://img.shields.io/coveralls/github/budnieswski/validate-struct?style=flat-square)](https://coveralls.io/github/budnieswski/validate-struct?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/budnieswski/validate-struct)](https://goreportcard.com/report/github.com/budnieswski/validate-struct)
[![GoDoc](https://godoc.org/github.com/budnieswski/validate-struct?status.svg)](https://pkg.go.dev/github.com/budnieswski/validate-struct)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/budnieswski/validate-struct?style=flat-square)])(#)
[![License](https://img.shields.io/dub/l/vibe-d.svg?style=flat-square)](#)




## Contents
- [Install](#install)
- [Example](#example)
- [Goals](#goals)

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

## Goals
- [ ] Create validate tests
- [ ] Create rule tests
- [ ] Standardize validate error return
- [ ] Create Github templates
	- [ ] Contributing
	- [ ] Issue
	- [ ] Pull request
- [ ] Create Github workflows
	- [X] Tests
	- [X] Coverage
	- [ ] CI
- [ ] Create doc.go
- [ ] Make validate benchmarks
- [ ] Improve performance
- [ ] Create a way for the user add custom rules