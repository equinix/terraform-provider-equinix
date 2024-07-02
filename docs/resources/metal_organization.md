---
subcategory: "Metal"
---

# equinix_metal_organization (Resource)

Provides a resource to manage organization resource in Equinix Metal.

## Example Usage

```terraform
# Create a new Organization
resource "equinix_metal_organization" "tf_organization_1" {
  name        = "foobar"
  description = "quux"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Organization.
* `address` - (Required) An object that has the address information. See [Address](#address) below for more details.
* `description` - (Optional) Description string.
* `website` - (Optional) Website link.
* `twitter` - (Optional) Twitter handle.
* `logo` - (Optional) Logo URL.

### Address

The `address` block contains:

* `address` - (Required) Postal address.
* `city` - (Required) City name.
* `country` - (Required) Two letter country code (ISO 3166-1 alpha-2), e.g. US.
* `zip_code` - (Required) Zip Code.
* `state` - (Optional) State name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the organization.
* `created` - The timestamp for when the organization was created.
* `updated` - The timestamp for the last time the organization was updated.

## Import

This resource can be imported using an existing organization ID:

```sh
terraform import equinix_metal_organization {existing_organization_id}
```
