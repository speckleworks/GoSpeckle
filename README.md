# GoSpeckle
> YASC - Yet Another Speckle Client

## Overview
This is a client for the Speckle Server API written in Go. The main reason for building another client in another language is that Golang offers some neat multi platform binary build capabilities (if a UI or a CLI tool is built in Go using this client as a base then it can be compiled to run on Windows, OSX, Linux and even mobile).

For the moment the package interacts with:
* Accounts
* Clients
* Projects
* Streams
* Objects

And cannot interact with:
* Websockets API

## Examples
To create a new Stream:

```go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	gospeckle "github.com/speckleworks/GoSpeckle/pkg"
)

// PrintJson is a utility function to pretty print a struct as a
// JSON
func PrintJson(body interface{}) {
	var prettyJSON bytes.Buffer
	// var jsonData []byte
	jsonData, _ := json.Marshal(body)

	error := json.Indent(&prettyJSON, jsonData, "", "\t")
	if error != nil {
		log.Println("JSON parse error: ", error)
		return
	}

	log.Println(string(prettyJSON.Bytes()))
}

func main() {
	ctx := context.TODO()

  httpClient := new(http.Client)
  apiURL := "https://hestia.speckle.works/api/v1/"
	client := gospeckle.NewClient(httpClient, apiURL, nil, "", "")

	client.Login(ctx, "go@speckle.come", "some-secret-password")

	streamBase := gospeckle.RequestMetadata{
		Private:           true,
		AnonymousComments: false,
	}
	newStream := gospeckle.StreamRequest{
		// Metadata:    &streamBase,
		RequestMetadata: &streamBase,
		Description:     "A fancy new go stream",
		Tags:            []string{"test", "golang"},
	}

	stream, _, err := client.Stream.Create(ctx, newStream)
	// stream, _, err := client.Stream.Get(ctx, "ufF4mPRuZ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	PrintJson(stream)
}

```
