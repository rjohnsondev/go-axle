ApiAxle client for Go (golang)
==============================

[![Build Status](https://travis-ci.org/rjohnsondev/go-axle.png)](https://travis-ci.org/rjohnsondev/go-axle) 
[![Coverage Status](https://coveralls.io/repos/rjohnsondev/go-axle/badge.png?branch=HEAD)](https://coveralls.io/r/rjohnsondev/go-axle?branch=HEAD)

## Features

Should provide complete programmatic access to all API methods as documented here: http://apiaxle.com/api.html

## Installation

This should get you started:

    go get github.com/rjohnsondev/go-axle

## Docs

Generated documentation can be viewed by running:

    godoc -http :3000

then navigating to http://localhost:3000/pkg/go-axle/

## Usage

This example program will connect to the ApiAxle server running on port 28902, create a new test API, a new key and link the key to the API.

```go
package main

import (
	"fmt"
	"github.com/rjohnsondev/go-axle"
	"os"
)

const (
	TEST_API_AXLE_SERVER = "http://localhost:28902/"
	TEST_API_NAME        = "goaxletestapi"
	TEST_KEY_NAME        = "goaxletestkey"
	TEST_API_ENDPOINT    = "localhost:80"
)

func main() {

	// Ping the server
	err := goaxle.Ping(TEST_API_AXLE_SERVER)
	if err != nil {
		fmt.Printf("Unable to connect to apiaxle at: %v\n", TEST_API_AXLE_SERVER)
		os.Exit(1)
	}

	// Provision a new API
	api := goaxle.NewApi(TEST_API_AXLE_SERVER, TEST_API_NAME, TEST_API_ENDPOINT)
	api.GlobalCache = 30 // cache requests for 30 seconds
	err = api.Save()
	if err != nil {
		panic(err)
	}

	// Create a new key
	key := goaxle.NewKey(TEST_API_AXLE_SERVER, TEST_KEY_NAME)
	key.Qpd = 200000 // Queries per day
	err = key.Save()
	if err != nil {
		panic(err)
	}

	// Link it up
	_, err = api.LinkKey(TEST_KEY_NAME)
	if err != nil {
		panic(err)
	}
}
```
