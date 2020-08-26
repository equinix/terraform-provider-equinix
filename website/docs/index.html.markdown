---
layout: "equinix"
page_title: "Provider: Equinix"
sidebar_current: "docs-equinix-index"
description: |-
  The Terraform Equinix provider allows for lifecycle management of Equinix Platform resources.
---

# Equinix Provider

The Equinix provider is used to interact with the resources supported by
Equinix. The provider needs to be configured with the proper credentials before
it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

Example [provider configuration](https://www.terraform.io/docs/configuration/providers.html) in `main.tf` file:

```hcl
provider equinix {
  endpoint = "https://api.equinix.com"
  client_id = "someID"
  client_secret = "someSecret"
}
```

Example provider configuration using `environment variables`:

```sh
export EQUINIX_API_ENDPOINT="https://api.equinix.com"
export EQUINIX_API_CLIENTID=someID
export EQUINIX_API_CLIENTSECRET=someSecret
```

## Argument Reference

The Equinix provider requires a few basic parameters:

- `client_id` - (Required) API Consumer Key available under "My Apps" in
  developer portal. This can be specified with the `EQUINIX_API_CLIENTID` shell
  environment variable.

- `client_secret` (Required) API Consumer secret available under "My Apps" in
  developer portal. This can be specified with the `EQUINIX_API_CLIENTSECRET`
  shell environment variable.

- `endpoint` (Optional) The Equinix API base URL, for example
  `https://playgroundapi.equinix.com` or `https://sandboxapi.equinix.com`. This
  can be specified with the `EQUINIX_API_ENDPOINT` shell environment variable.
  (Defaults to `https://api.equinix.com`)

- `request_timeout` (Optional) The duration of time, in seconds, that the
  Equinix Platform API Client should wait before canceling an API request.
  Canceled requests may still result in provisioned resources. (Defaults to `5`)

These parameters can be provided in [Terraform variable
files](https://www.terraform.io/docs/configuration/variables.html#variable-definitions-tfvars-files)
or as environment variables. Nevertheless, please note that it is [not
recommended to keep sensitive data in plain text
files](https://www.terraform.io/docs/state/sensitive-data.html).
