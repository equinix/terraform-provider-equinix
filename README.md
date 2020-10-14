# Equinix Metal Terraform Provider

[![GitHub release](https://img.shields.io/github/release/packethost/terraform-provider-packet/all.svg?style=flat-square)](https://github.com/packethost/terraform-provider-packet/releases)
![](https://img.shields.io/badge/Stability-Maintained-green.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/packethost/terraform-provider-packet)](https://goreportcard.com/report/github.com/packethost/terraform-provider-packet)

[![Slack](https://slack.equinixmetal.com/badge.svg)](https://slack.equinixmetal.com)
[![Twitter Follow](https://img.shields.io/twitter/follow/equinixmetal.svg?style=social&label=Follow)](https://twitter.com/intent/follow?screen_name=equinixmetal)

<img src="https://metal.equinix.com/metal/images/logo/equinix-metal-full.svg" width="600px">

[Packet is now Equinix Metal!](https://blog.equinix.com/blog/2020/10/06/equinix-metal-metal-and-more/)

This repository is [Maintained](https://github.com/packethost/standards/blob/master/maintained-statement.md) meaning that this software is supported by Equinix Metal and its community - available to use in production environments.

## Using the provider

The Equinix Metal provider will be installed on `terraform init` of a template using any of the `packet_*` resources.

See <https://registry.terraform.io/providers/packethost/packet/latest/docs> for documentation on the resources included in this provider.

## Requirements

- [Terraform 0.12+](https://www.terraform.io/downloads.html) (for v3.0.0 of this provider and newer)
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Building the provider

Clone repository to: `$GOPATH/src/github.com/packethost/terraform-provider-packet`

```sh
mkdir -p $GOPATH/src/github.com/packethost; cd $GOPATH/src/github.com/packethost
git clone git@github.com:packethost/terraform-provider-packet
```

Enter the provider directory and build the provider

```sh
cd $GOPATH/src/github.com/packethost/terraform-provider-packet
make build
```

## Developing the provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-packet
...
```

## Testing provider code

We have mostly acceptance tests in the provider. There's no point for you to run them all, but you should run the one covering the functionality which you change. The acceptance test run will cost you some money, so feel free to abstain. The acceptance test suite will be run for your PR during the review process.

To run an acceptance test, find the relevant test function in `*_test.go` (for example TestAccPacketDevice_Basic), and run it as

```sh
TF_ACC=1 go test -v -timeout=20m -run=TestAccPacketDevice_Basic
```

If you want to see HTTP traffic, set `TF_LOG=DEBUG`, i.e.

```sh
TF_LOG=DEBUG TF_ACC=1 go test -v -timeout=20m -run=TestAccPacketDevice_Basic
```

## Testing the provider with Terraform

Once you've built the plugin binary (see [Developing the provider](#developing-the-provider) above), it can be incorporated within your Terraform environment using the `-plugin-dir` option. Subsequent runs of Terraform will then use the plugin from your development environment.

```sh
terraform init -plugin-dir $GOPATH/bin
```
