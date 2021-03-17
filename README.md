# Equinix Metal Terraform Provider

[![GitHub release](https://img.shields.io/github/release/equinix/terraform-provider-metal/all.svg?style=flat-square)](https://github.com/equinix/terraform-provider-metal/releases)
![](https://img.shields.io/badge/Stability-Maintained-green.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/equinix/terraform-provider-metal)](https://goreportcard.com/report/github.com/equinix/terraform-provider-metal)

[![Slack](https://slack.equinixmetal.com/badge.svg)](https://slack.equinixmetal.com)
[![Twitter Follow](https://img.shields.io/twitter/follow/equinixmetal.svg?style=social&label=Follow)](https://twitter.com/intent/follow?screen_name=equinixmetal)

<img src="https://metal.equinix.com/metal/images/logo/equinix-metal-full.svg" width="600px">

This repository is [Maintained](https://github.com/packethost/standards/blob/master/maintained-statement.md) meaning that this software is supported by Equinix Metal and its community - available to use in production environments.

## Using the provider

The Equinix Metal provider will be installed on `terraform init` of a template using any of the `metal_*` resources.

See <https://registry.terraform.io/providers/equinix/metal/latest/docs> for documentation on the resources included in this provider.

### Migrating from Packet

[Packet is now Equinix Metal!](https://blog.equinix.com/blog/2020/10/06/equinix-metal-metal-and-more/) See [Issue #1](https://github.com/equinix/terraform-provider-metal/issues/1) for more details on migrating existing projects.
## Requirements

- [Terraform 0.12+](https://www.terraform.io/downloads.html) (for v1.0.0 of this provider and newer)
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Building the provider

Clone the repository, enter the provider directory, and build the provider.

```sh
git clone git@github.com:equinix/terraform-provider-metal
cd terraform-provider-metal
make build
```

## Developing the provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-metal
...
```

## Testing provider code

We have mostly acceptance tests in the provider. There's no point for you to run them all, but you should run the one covering the functionality which you change. The acceptance test run will cost you some money, so feel free to abstain. The acceptance test suite will be run for your PR during the review process.

To run an acceptance test, find the relevant test function in `*_test.go` (for example TestAccMetalDevice_Basic), and run it as

```sh
TF_ACC=1 go test -v -timeout=20m ./... -run=TestAccMetalDevice_Basic
```

If you want to see HTTP traffic, set `TF_LOG=DEBUG`, i.e.

```sh
TF_LOG=DEBUG TF_ACC=1 go test -v -timeout=20m ./... -run=TestAccMetalDevice_Basic
```

## Testing the provider with Terraform

Once you've built the plugin binary (see [Developing the provider](#developing-the-provider) above), it can be incorporated within your Terraform environment using the `-plugin-dir` option. Subsequent runs of Terraform will then use the plugin from your development environment.

```sh
terraform init -plugin-dir $GOPATH/bin
```
