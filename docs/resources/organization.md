---
page_title: "Equinix Metal: metal_organization"
subcategory: ""
description: |-
  Provides an Equinix Metal Organization resource.
---

# metal\_organization

Provides a resource to manage organization resource in Equinix Metal.

## Example Usage

```hcl
# Create a new Project
resource "metal_organization" "tf_organization_1" {
  name        = "foobar"
  description = "quux"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Organization
* `description` - Description string
* `website` - Website link
* `twitter` - Twitter handle
* `logo` - Logo URL

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the organization
* `name` - The name of the Organization
* `description` - Description string
* `website` - Website link
* `twitter` - Twitter handle
* `logo` - Logo URL
