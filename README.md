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

---

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

## Developing the provider

Check [Equinix provider development](DEVELOPMENT.md) for guides on building
and developing the provider.
