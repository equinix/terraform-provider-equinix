<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/public/img/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Provider for Equinix Platform

The Terraform Equinix provider is a plugin for Terraform that allows for lifecycle
management of Equinix Platform resources.

[![Go Tests](https://github.com/equinix/terraform-provider-equinix/actions/workflows/test.yml/badge.svg)](https://github.com/equinix/terraform-provider-equinix/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/equinix/terraform-provider-equinix)](https://goreportcard.com/report/github.com/equinix/terraform-provider-equinix)
[![GoDoc](https://godoc.org/github.com/go-resty/resty?status.svg)](https://godoc.org/github.com/equinix/terraform-provider-equinix)
![GitHub](https://img.shields.io/github/license/equinix/terraform-provider-equinix)
[![Equinix Community](https://img.shields.io/badge/Equinix%20Community-E8272C.svg)](https://[https://community.equinix.com/](https://community.equinix.com/))

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
- collection of modules can be found in
[Terraform Equinix modules directory](modules/)

## Developing the provider

Check [Equinix provider development](DEVELOPMENT.md) for guides on building
and developing the provider.
