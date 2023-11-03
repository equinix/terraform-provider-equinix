---
subcategory: "Fabric"
---

# DEPRECATED RESOURCE

End of Life will be June 30th, 2024. Use equinix_fabric_service_profiles instead.

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
* `speed_bands` - (Optional) List of speed bands that should be supported by resulting profiles. (** Please see the note below for important information regarding the use of this argument**)
* `organization_name` - (Optional) Name of seller's organization.
* `organization_global_name` - (Optional) Name of seller's global organization.

-> **NOTE:** Some seller profiles (such as OCI) do not show the available `speed_bands` until the connection is validated. Therefore, using the `speed_bands` argument as a filter can cause false empty results. You can use the [equinix_ecx_l2_sellerprofile](./equinix_ecx_l2_sellerprofile.md) data source to ensure that the seller profile you need returns its `speed_bands` information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `profiles` - List of resulting profiles. Each element in the `profiles` list exports all
[Service Profile Attributes](./equinix_ecx_l2_sellerprofile.md#attributes-reference).
