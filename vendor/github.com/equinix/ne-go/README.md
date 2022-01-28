# Equinix Network Edge Go client

Equinix Network Edge (NE) client library written in Go.

[![Build Status](https://travis-ci.com/equinix/ne-go.svg?branch=master)](https://travis-ci.com/github/equinix/ne-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/equinix/ne-go)](https://goreportcard.com/report/github.com/equinix/ne-go)
[![GoDoc](https://godoc.org/github.com/equinix/ne-go?status.svg)](https://godoc.org/github.com/equinix/ne-go)
![GitHub](https://img.shields.io/github/license/equinix/ne-go)

---

## Purpose

Equinix Network Edge client library was written in Go for purpose of managing NE
resources from Terraform provider plugin.

Library gives possibility to manage virtual network devices and associated network
services.

**NOTE**: scope of this library is limited to needs of Terraform provider plugin
and it is not providing full capabilities of Equinix Network Edge API

## Usage

### Code

1. Add ne-go module to import statement.
   In below example, Equinix `oauth2-go` module is imported as well

   ```go
   import (
   "github.com/equinix/oauth2-go"
   "github.com/equinix/ne-go"
   )
   ```

2. Define baseURL that will be used in all REST API requests

    ```go
    baseURL := "https://api.equinix.com"
    ```

3. Create oAuth configuration and oAuth enabled `http.Client`

    ```go
    authConfig := oauth2.Config{
      ClientID:     "someClientId",
      ClientSecret: "someSecret",
      BaseURL:      baseURL}
    ctx := context.Background()
    authClient := authConfig.New(ctx)
    ```

4. Create NE REST client with a given `baseURL` and oauth's `http.Client`

    ```go
    var neClient ne.Client = ne.NewClient(ctx, baseURL, authClient)
    ```

5. Use NE client to perform some operation i.e. **get device** details

    ```go
    device, err := neClient.GetDevice("existingDeviceUUID")
    if err != nil {
      log.Printf("Error while fetching device - %v", err)
    } else {
      log.Printf("Retrieved device - %+v", device)
    }
    ```
