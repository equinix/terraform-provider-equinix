---
subcategory: "Metal"
---

~> **Deprecation Notice** This data source is deprecated and scheduled for removal. Equinix Metal's operational end-of-life is June 30, 2026. This data source will be removed in the next major provider version (5.0.0). For ongoing Metal service usage until the sunset, please continue with version 4.x of the Equinix Terraform provider. Additional information regarding the Metal platform conclusion is available at: https://docs.equinix.com/metal/


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
* `logo` - (Deprecated) Logo URL.
* `address` - Address information
  * `address` - Postal address.
  * `city` - City name.
  * `country` - Two letter country code (ISO 3166-1 alpha-2), e.g. US.
  * `zip_code` - Zip Code.
  * `state` - State name.
