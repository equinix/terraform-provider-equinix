# Equinix REST client

Equinix REST client written in Go.
Implementation is based on [Resty client](https://github.com/go-resty/resty).

![Build Status](https://github.com/equinix/rest-go/workflows/Build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/equinix/rest-go)](https://goreportcard.com/report/github.com/equinix/rest-go)
[![GoDoc](https://godoc.org/github.com/equinix/rest-go?status.svg)](https://godoc.org/github.com/equinix/rest-go)
![GitHub](https://img.shields.io/github/license/equinix/rest-go)

---

## Purpose

Purpose of this module is to wrap Resty REST client with Equinix specific error handling.
In addition, module adds support for paginated query requests.

Module is used by other Equinix client libraries, like [ECX Fabric Go client](https://github.com/equinix/ecx-go)
or [Network Edge Go client](https://github.com/equinix/ne-go).

## Features

* parses Equinix standardized error response body contents
* `GetPaginated` function queries for data on APIs with paginated responses. Pagination
 options can be configured by setting up attributes of `PagingConfig`

## Usage

1. Get recent equinix/rest-go module

   ```sh
   go get -d github.com/equinix/rest-go
   ```

2. Create new Equinix REST client with default HTTP client

   ```go
   import (
       "context"
       "net/http"
       "github.com/equinix/rest-go"
   )

   func main() {
     c := rest.NewClient(
          context.Background(),
          "https://api.equinix.com",
          &http.Client{})
   }
   ```

3. Use Equinix HTTP client with Equinix APIs

   ```go
    respBody := api.AccountResponse{}
    req := c.R().SetResult(&respBody)
    if err := c.Execute(req, "GET", "/ne/v1/device/account"); err != nil {
     //Equinix application error details will be included
     log.Printf("Got error: %s", err) 
    }
   ```

## Debugging

Debug logging comes from Resty client and logs request and response details to stderr.

Such debug logging can be enabled by setting up `EQUINIX_REST_LOG` environmental
variable to `DEBUG`
