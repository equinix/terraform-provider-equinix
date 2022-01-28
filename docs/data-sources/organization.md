---
page_title: "Equinix Metal: metal_organization"
subcategory: ""
description: |-
  Provides an Equinix Metal Organization datasource. This can be used to read existing Organizations.
---

# metal_organization

Provides an Equinix Metal organization datasource.

## Example Usage

```hcl
# Fetch a organization data and show projects which belong to it
data "metal_organization" "test" {
  organization_id = local.org_id
}

output "projects_in_the_org" {
  value = data.metal_organization.test.project_ids
}
```

## Argument Reference

The following arguments are supported:

* `name` - The organization name
* `organization_id` - The UUID of the organization resource

Exactly one of `name` or `organization_id` must be given.

## Attributes Reference

The following attributes are exported:

* `project_ids` - UUIDs of project resources which belong to this organization
* `description` - Description string
* `website` - Website link
* `twitter` - Twitter handle
* `logo` - Logo URL
