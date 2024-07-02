---
subcategory: "Metal"
---

# equinix_metal_organization (Data Source)

Provides an Equinix Metal organization datasource.

## Example Usage

```terraform
# Fetch a organization data and show projects which belong to it
data "equinix_metal_organization" "test" {
  organization_id = local.org_id
}

output "projects_in_the_org" {
  value = data.equinix_metal_organization.test.project_ids
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The organization name.
* `organization_id` - (Optional) The UUID of the organization resource.

Exactly one of `name` or `organization_id` must be given.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `project_ids` - UUIDs of project resources which belong to this organization.
* `description` - Description string.
* `website` - Website link.
* `twitter` - Twitter handle.
* `logo` - Logo URL.
* `address` - Address information
  * `address` - Postal address.
  * `city` - City name.
  * `country` - Two letter country code (ISO 3166-1 alpha-2), e.g. US.
  * `zip_code` - Zip Code.
  * `state` - State name.
