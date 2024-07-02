---
subcategory: "Metal"
---

# equinix_metal_project_api_key (Resource)

Use this resource to create Metal Project API Key resources in Equinix Metal. Project API keys can be used to create and read resources in a single project. Each API key contains a token which can be used for authentication in Equinix Metal HTTP API (in HTTP request header `X-Auth-Token`).

Read-only keys only allow to list and view existing resources, read-write keys can also be used to create resources.

## Example Usage

```terraform
# Create a new read-only API key in existing project
resource "equinix_metal_project_api_key" "test" {
  project_id  = local.existing_project_id
  description = "Read-only key scoped to a projct"
  read_only   = true
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) UUID of the project where the API key is scoped to.
* `description` - (Required) Description string for the Project API Key resource.
* `read-only` - (Optional) Flag indicating whether the API key shoud be read-only.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `token` - API token which can be used in Equinix Metal API clients
