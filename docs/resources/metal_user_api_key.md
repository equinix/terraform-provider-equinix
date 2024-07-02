---
subcategory: "Metal"
---

# equinix_metal_user_api_key (Resource)

Use this resource to create Metal User API Key resources in Equinix Metal. Each API key contains a token which can be used for authentication in Equinix Metal HTTP API (in HTTP request header `X-Auth-Token`).

Read-only keys only allow to list and view existing resources, read-write keys can also be used to create resources.

## Example Usage

```terraform
# Create a new read-only user API key

resource "equinix_metal_user_api_key" "test" {
  description = "Read-only user key"
  read_only   = true
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required) Description string for the User API Key resource.
* `read-only` - (Required) Flag indicating whether the API key shoud be read-only.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `user_id` - UUID of the owner of the API key.
* `token` - API token which can be used in Equinix Metal API clients.
