<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Provider for Equinix Platform

The Terraform Equinix provider is a plugin for Terraform that allows for lifecycle
management of Equinix Platform resources.

[![Build Status](https://travis-ci.com/equinix/terraform-provider-equinix.svg?branch=master)](https://travis-ci.com/github/equinix/terraform-provider-equinix)
[![Go Report Card](https://goreportcard.com/badge/github.com/equinix/terraform-provider-equinix)](https://goreportcard.com/report/github.com/equinix/terraform-provider-equinix)
[![GoDoc](https://godoc.org/github.com/go-resty/resty?status.svg)](https://godoc.org/github.com/equinix/terraform-provider-equinix)
![GitHub](https://img.shields.io/github/license/equinix/terraform-provider-equinix)

[![Equinix Metal Slack](https://slack.equinixmetal.com/badge.svg)](https://slack.equinixmetal.com)
[![Equinix Metal Twitter](https://img.shields.io/twitter/follow/equinixmetal.svg?style=social&label=Follow)](https://twitter.com/intent/follow?screen_name=equinixmetal)

---

This repository is [Maintained](https://github.com/equinix-labs/equinix-labs/blob/main/maintained-statement.md) meaning that this software is supported by Equinix Metal and its community - available to use in production environments.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+

## Using the provider

The Equinix provider will be installed automatically from official Terraform
registry on `terraform init`.

See the [Equinix Provider documentation](https://registry.terraform.io/providers/equinix/equinix/latest/docs)
to get started using the Equinix provider.

## Documentation

- full documentation is available on [Terraform Registry website](https://registry.terraform.io/providers/equinix/equinix/latest/docs)
- use case driven guides can be found in [Equinix Provider examples directory](examples/)
- collection of modules can be found in
[Terraform Equinix modules directory](modules/)

## Developing the provider

Check [Equinix provider development](DEVELOPMENT.md) for guides on building
and developing the provider.
