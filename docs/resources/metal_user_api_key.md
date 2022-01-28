---
page_title: "Equinix: equinix_metal_user_api_key"
subcategory: ""
description: |-
  Create Equinix Metal User API Keys
  ---

# Resource: equinix_metal_user_api_key

Use this resource to create Metal User API Key resources in Equinix Metal. Each API key contains a token which can be used for authentication in Equinix Metal HTTP API (in HTTP request header `X-Auth-Token`).

Read-only keys only allow to list and view existing resources, read-write keys can also be used to create resources.

## Example Usage

```hcl

# Create a new read-only API key

resource "equinix_metal_user_api_key" "test" {
  description = "Read-only user key"
  read_only   = true
}
```

## Argument Reference

* `description` - Description string for the User API Key resource
* `read-only` - Flag indicating whether the API key shoud be read-only

## Attributes Reference

* `user_id` - UUID of the owner of the API key
* `token` - API token which can be used in Equinix Metal API clients
