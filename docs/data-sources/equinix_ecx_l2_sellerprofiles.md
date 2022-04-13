---
subcategory: "Fabric"
---

# equinix_ecx_l2_sellerprofiles (Data Source)

Use this data source to get details of available Equinix Fabric layer 2 seller profiles. It is
possible to apply filtering criteria for returned list of profiles.

## Example usage

```hcl
data "equinix_ecx_l2_sellerprofiles" "aws" {
  organization_global_name = "AWS"
  metro_codes              = ["SV", "DC"]
  speed_bands              = ["1GB", "500MB"]
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply on returned seller profile names and filter
search results.
* `metro_codes` - (Optional) List of metro codes of locations that should be served by resulting
profiles.
* `speed_bands` - (Optional) List of speed bands that should be supported by resulting profiles.
* `organization_name` - (Optional) Name of seller's organization.
* `organization_global_name` - (Optional) Name of seller's global organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `profiles` - List of resulting profiles. Each element in the `profiles` list exports all
[Service Profile Attributes](./equinix_ecx_l2_sellerprofile.md#attributes-reference).
