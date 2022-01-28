---
page_title: "Equinix Metal: Metal Project API Key"
subcategory: ""
description: |-
  Create Equinix Metal Project API Keys
  ---

# metal_project_api_key

Use this resource to create Metal Project API Key resources in Equinix Metal. Project API keys can be used to create and read resources in a single project. Each API key contains a token which can be used for authentication in Equinix Metal HTTP API (in HTTP request header `X-Auth-Token`).


Read-only keys only allow to list and view existing resources, read-write keys can also be used to create resources.


## Example Usage

```hcl

# Create a new read-only API key in existing project

resource "equinix_metal_project_api_key" "test" {
  project_id  = local.existing_project_id
  description = "Read-only key scoped to a projct"
  read_only   = true
}
```

## Argument Reference

* `project_id` - UUID of the project where the API key is scoped to
* `description` - Description string for the Project API Key resource
* `read-only` - Flag indicating whether the API key shoud be read-only

## Attributes Reference

* `token` - API token which can be used in Equinix Metal API clients
