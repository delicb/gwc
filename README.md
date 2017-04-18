# GWC (Go Web Client)
[![Go Report Card](https://goreportcard.com/badge/github.com/delicb/gwc)](https://goreportcard.com/report/github.com/delicb/gwc)
[![Build Status](https://travis-ci.org/delicb/gwc.svg?branch=master)](https://travis-ci.org/delicb/gwc)
[![codecov](https://codecov.io/gh/delicb/gwc/branch/master/graph/badge.svg)](https://codecov.io/gh/delicb/gwc)
![status](https://img.shields.io/badge/status-beta-red.svg)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/delicb/gwc)



GWC is GoLang HTTP client based on [Cliware](https://github.com/delicb/cliware)
client middleware management library and its [middlewares](https://github.com/delicb/cliware-middlewares).

Basic idea is to use same middleware mechanism on client side as many projects
use it for server development. 

Because it is based on middlewares client is very pluggable. However, even out 
of box it should support most common use cases.

# Install
Run `go get github.com/delicb/gwc` in terminal.

# Example
There is testable example in this repository. It basically does following:
```go
package gwc

import (
	"fmt"
	"net/http"

    "go.delic.rs/gwc"
	"go.delic.rs/cliware-middlewares/headers"
	"go.delic.rs/cliware-middlewares/errors"
)

type HTTPBinResponse struct {
	Headers struct {
		UserAgent string `json:"User-Agent"`
	} `json:"headers"`
}

func main() {
	respBody := new(HTTPBinResponse)

	client := gwc.New(
		http.DefaultClient,
		// Add user agent header, it will be applied to all requests
		headers.Add("User-Agent", "example-client"),
		errors.Errors(),
	)
	// send request
	resp, err := client.Get().URL("https://httpbin.org/get").Send()
	
	// check errors
	// because of errors middleware included in client, ever status codes 
	// 400+ will be turned into errors.
	if err != nil {
		panic(err)
	}

	resp.JSON(respBody)
	fmt.Println(respBody.Headers.UserAgent)
	// output: example-client
}

```

In this example we create new GWC instance with some middlewares (in this case
we set `User-Agent` header to each request sent using this middleware).

We use client to send request and (after some error checking) we deserialize
JSON body to structure.

More complex example (small part of client for big API) can be found 
[here](https://github.com/delicb/sevenbridges-go).

# State
This is early development, not stable, backward compatibility not guarantied.

# Contribution
Any contribution is welcome. If you find this code useful, please let me know.
If you find bugs - feel free to open an issue, or, even better, new pull request.

# Credits
Idea and bunch of implementation details were taken from cool GoLang HTTP client
[Gentleman](https://github.com/h2non/gentleman).

Difference is that GWC is based on Cliware, which supports `context` for client
requests and has more simple idea of middleware. Also, GWC is lacking some
features that Gentleman has (like mux). For now I do not plan on adding them
to GWC, but I might write middleware that support similar functionality in the
future. 
