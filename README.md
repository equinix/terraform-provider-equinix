<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Provider for Equinix Platform

The Terraform Equinix provider is a plugin for Terraform that allows for lifecycle management of Equinix Platform resources.

[![Build Status](https://travis-ci.com/equinix/terraform-provider-equinix.svg?branch=master)](https://travis-ci.com/github/equinix/terraform-provider-equinix)
[![Go Report Card](https://goreportcard.com/badge/github.com/equinix/terraform-provider-equinix)](https://goreportcard.com/report/github.com/equinix/terraform-provider-equinix)
[![GoDoc](https://godoc.org/github.com/go-resty/resty?status.svg)](https://godoc.org/github.com/equinix/terraform-provider-equinix)
![GitHub](https://img.shields.io/github/license/equinix/terraform-provider-equinix)

---

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.14+ (to build provider plugin)

## Using the provider

See the [Equinix Provider documentation](docs/index.md) to get started using the
Equinix provider.

- [Data source: ECXF port](docs/data-sources/ecx_port.md)
- [Data source: ECXF layer 2 seller
  profile](docs/data-sources/ecx_l2_sellerprofile.md)
- [Resource: ECXF layer 2 connection](docs/resources/ecx_l2_connection.md)
- [Resource: ECXF layer 2 connection
  accepter](docs/resources/ecx_l2_connection_accepter.md)
- [Resource: ECXF layer 2 service
  profile](docs/resources/ecx_l2_serviceprofile.md)

## Building the provider

1. Clone Equinix Terraform Provider repository

   ```sh
   git clone https://github.com/equinix/terraform-provider-equinix.git
   ```

2. Build the provider

   Enter the provider directory and build the provider:

   ```sh
   cd terraform-provider-equinix
   make build
   ```

3. Install the provider

   Provider binary can be installed in terraform plugins directory `~/.terraform.d/plugins` by running make with _install_ target:

   ```sh
   make install
   ```

## Developing the provider

- use Go programming best practices, _gofmt, go_vet, golint, ineffassign_, etc.
- enter the provider directory

  ```sh
  cd terraform-provider-equinix
  ```

- to build, use make `build` target

  ```sh
  make build
  ```

- to run unit tests, use make `test` target

  ```sh
  make test
  ```

- to run acceptance tests, use make `testacc` target

  ```sh
  make testacc
  ```

  Check "Running acceptance tests" section for more details.

### Running acceptance tests

**NOTE**: acceptance tests create resources on real infrastructure, thus may be subject for costs. In order to run acceptance tests, you must set necessary provider configuration attributes.

```sh
export EQUINIX_API_ENDPOINT=https://api.equinix.com
export EQUINIX_API_CLIENTID=someID
export EQUINIX_API_CLIENTSECRET=someSecret
make testacc
```

#### ECX Port acceptance tests

ECX Port data source acceptance tests use below parameters, that can be set to match with desired testing environment. If not set, defaults values, **from Sandbox environment** are used.

- **TF_ACC_ECX_PORT_NAME** - sets name of the port used in data source

#### ECX L2 connection acceptance tests

ECX Layer 2 connection acceptance tests use below parameters, that can be set to match with desired testing environment. If not set, defaults values, **from Sandbox environment** are used.

- **TF_ACC_ECX_L2_AWS_SP_ID** - sets UUID of Layer2 service profile for AWS
- **TF_ACC_ECX_L2_AZURE_SP_ID** - sets UUID of Layer2 service profile for Azure
- **TF_ACC_ECX_PRI_DOT1Q_PORT_ID** - sets UUID of Dot1Q encapsulated port on primary device
- **TF_ACC_ECX_SEC_DOT1Q_PORT_ID** - sets UUID of Dot1Q encapsulated port on secondary device

Example - running tests on Sandbox environment but with defined ports:

```sh
export EQUINIX_API_ENDPOINT=https://sandboxapi.equinix.com
export EQUINIX_API_CLIENTID=someID
export EQUINIX_API_CLIENTSECRET=someSecret
export TF_ACC_ECX_PRI_DOT1Q_PORT_ID="6ca3704b-c660-4c6f-9e66-3282f8de787b"
export TF_ACC_ECX_SEC_DOT1Q_PORT_ID="7a80ab13-4e04-455c-82e3-79d962d0c0c3"
make testacc
```
