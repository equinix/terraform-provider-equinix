# Equinix Fabric Go client

Equinix Fabric client library written in Go.

[![Build Status](https://travis-ci.com/equinix/ecx-go.svg?branch=master)](https://travis-ci.com/github/equinix/ecx-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/equinix/ecx-go)](https://goreportcard.com/report/github.com/equinix/ecx-go)
[![GoDoc](https://godoc.org/github.com/equinix/ecx-go?status.svg)](https://godoc.org/github.com/equinix/ecx-go)
![GitHub](https://img.shields.io/github/license/equinix/ecx-go)

---

## Purpose

Equinix Fabric client library was written in Go for purpose of managing Fabric
resources from Terraform provider plugin.

Library gives possibility to manage layer two connections and service profiles
on Equinix Fabric and connect to any Cloud Service Provider, other Enterprise
or between own ports.

## Features

Client library consumes Equinix Fabric's REST API and allows to:

- manage Fabric L2 connections
  - retrieve L2 connection details
  - create non redundant L2 connection
  - create redundant L2 connection
  - delete L2 connection
  - update L2 connection (name and speed)
- manage Fabric L2 service profiles
- retrieve list of Fabric user ports
- retrieve list of Fabric L2 seller profiles

**NOTE**: scope of this library is limited to needs of Terraform provider plugin
and it is not providing full capabilities of Equinix Fabric API

## Usage

### Code

1. Add ecx-go module to import statement.
   In below example, Equinix `oauth2-go` module is imported as well

   ```go
   import (
    "github.com/equinix/oauth2-go"
    "github.com/equinix/ecx-go"
   )
   ```

2. Define baseURL that will be used in all REST API requests

    ```go
    baseURL := "https://sandboxapi.equinix.com"
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

4. Create Equinix Fabric REST client with a given `baseURL` and oauth's `http.Client`

    ```go
    var ecxClient ecx.Client = ecx.NewClient(ctx, baseURL, authClient)
    ```

5. Use Equinix Fabric client to perform some operation `i.e. fetch`

    ```go
    l2conn, err := ecxClient.GetL2Connection("myUUID")
    if err != nil {
      log.Printf("Error while fetching connection - %v", err)
    } else {
      log.Printf("Retrieved connection - %+v", l2conn)
    }
    ```
