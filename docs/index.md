---
page_title: "Provider: Equinix"
---

# Equinix Provider

The Equinix provider is used to interact with the resources provided by Equinix Platform. The provider needs to be configured with the proper credentials before
it can be used.

For information about obtaining API key and secret required for Equinix Fabric and Network Edge refer to
[Generating Client ID and Client Secret key](https://developer.equinix.com/docs/ecx-getting-started#generating-client-id-and-client-secret-key)
from [Equinix Developer Platform portal](https://developer.equinix.com).

Interacting with Equinix Metal requires an API auth token that can be generated at [Project-level](https://metal.equinix.com/developers/docs/accounts/projects/#api-keys) or [User-level](https://metal.equinix.com/developers/docs/accounts/users/#api-keys) tokens can be used.

If you are only using Equinix Metal resources, you may omit the Client ID and Client Secret provider configuration parameters needed to access other Equinix resource types (Network Edge, Fabric, etc).

Use the navigation to the left to read about the available resources.

## Example Usage

Example HCL with [provider configuration](https://www.terraform.io/docs/configuration/providers.html)
and a [required providers definition](https://www.terraform.io/language/settings#specifying-a-required-terraform-version):

```hcl
terraform {
  required_providers {
    equinix = {
      source = "equinix/equinix"
    }
  }
}

# Credentials for all Equinix resources
provider "equinix" {
  client_id     = "someEquinixAPIClientID"
  client_secret = "someEquinixAPIClientSecret"
  auth_token    = "someEquinixMetalToken"
}
```

Client ID and Client Secret can be omitted when the only Equinix resources
consumed are Equinix Metal resources.

```hcl
# Credentials for only Equinix Metal resources
provider "equinix" {
  auth_token    = "someEquinixMetalToken"
}
```

Example provider configuration using `environment variables`:

```sh
export EQUINIX_API_CLIENTID=someEquinixAPIClientID
export EQUINIX_API_CLIENTSECRET=someEquinixAPIClientSecret
export METAL_AUTH_TOKEN=someEquinixMetalToken
```

### Token Authentication

Token's can be generated for the API Client using the OAuth2 Token features described in the
[OAuth2 API](https://developer.equinix.com/catalog/accesstokenv1#operation/GetOAuth2AccessToken) documentation.

API tokens can be provided using the `token` provider argument, or the `EQUINIX_API_TOKEN` evironment variable.
The `client_id` and `client_secret` arguments will be ignored in the presence of a `token` argument.

When testing against the [Equinix Sandbox API](https://developer.equinix.com/environment/sandbox), tokens must be used.

```hcl
provider "equinix" {
  token = "someToken"
}
```

## Argument Reference

The Equinix provider requires a few basic parameters. While the authentication arguments are
individually optionally, either `token` or `client_id` and `client_secret` must be defined
through arguments or environment settings to interact with Equinix Fabric and Network Edge
services, and `auth_token` to interact with Equinix Metal.

* `client_id` - (Optional) API Consumer Key available under "My Apps" in
  developer portal. This argument can also be specified with the
  `EQUINIX_API_CLIENTID` shell environment variable.

* `client_secret` (Optional) API Consumer secret available under "My Apps" in
  developer portal. This argument can also be specified with the
  `EQUINIX_API_CLIENTSECRET` shell environment variable.

* `token` (Optional) API tokens are generated from API Consumer clients using
  the [OAuth2
  API](https://developer.equinix.com/docs/ecx-getting-started#requesting-access-and-refresh-tokens).
  This argument can also be specified with the `EQUINIX_API_TOKEN` shell
  environment variable.

* `auth_token` - (Optional) This is your Equinix Metal API Auth token. This can
  also be specified with the `METAL_AUTH_TOKEN` environment variable.

* `endpoint` (Optional) The Equinix API base URL to point out desired environment.
   This argument can also be specified with the `EQUINIX_API_ENDPOINT`
   shell environment variable. (Defaults to `https://api.equinix.com`)

* `request_timeout` (Optional) The duration of time, in seconds, that the
  Equinix Platform API Client should wait before canceling an API request.
  Canceled requests may still result in provisioned resources. (Defaults to `30`)

* `response_max_page_size` (Optional) The maximum number of records in a single response
  for REST queries that produce paginated responses. (Default is client specific)

* `max_retries` (Optional) Maximum number of retries in case of network failure.

* `max_retry_wait_seconds` (Optional) Maximum time to wait in case of network failure.

These parameters can be provided in [Terraform variable
files](https://www.terraform.io/docs/configuration/variables.html#variable-definitions-tfvars-files)
or as environment variables. Nevertheless, please note that it is [not
recommended to keep sensitive data in plain text
files](https://www.terraform.io/docs/state/sensitive-data.html).
