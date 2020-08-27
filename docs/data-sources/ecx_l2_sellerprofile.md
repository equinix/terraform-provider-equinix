---
layout: "equinix"
page_title: "Equinix: ecx_l2_sellerprofile"
sidebar_current: "docs-equinix-datasource-ecx-l2-sellerprofile"
description: |-
  Get the named Equinix ECX Layer2 Seller Profile
---

# equinix\_ecx\_l2\_sellerprofile

Data source `equinix_ecx_l2_sellerprofile` is used to fetch attributes of ECX Layer2 Seller Profile (like UUID) with a given profile name.

## Example usage

```hcl
data "equinix_ecx_l2_sellerprofile" "tf-aws" {
  name = "AWS Direct Connect"
}
```

## Argument Reference

- `name` - _(Optional)_ Name of seller profile
- `organization_name` - _(Optional)_ Name of seller's organization
- `organization_global_name` - _(Optional)_ Name of seller's global organization

## Attributes Reference

- `uuid` - Unique identifier of seller profile
- `speed_from_api` - information if seller is deriving connection speed from an API call
- `speed_customization_allowed` - information if seller allows customer to enter a custom connection speed
- `redundancy_required` - information if seller requires connections to be redundant
- `encapsulation` - seller profile's encapsulation (Dot1q or QinQ)
- `speed_band` - one or more specifications of speed/bandwidth supported by seller profile
  - `speed` - speed/bandwidth supported by this profile
  - `speed_unit` - unit of the speed/bandwidth supported by this profile
- `metro` - one or more specifications of metro locations supported by seller profile
  - `code` - metro code
  - `name` - metro name
  - `ibxes` - list of IBXes supported within given metro
  - `regions` - list of regions supported within given metro
- `additional_info` - one or more specifications of additional buyer information attrubutes that can be provided in connection definition that uses given seller profile
  - `name` - name of an attribute
  - `description` - textual description of an attribute
  - `data_type` - data type of an attribute _(BOOLEAN / INTEGER / STRING)_
  - `mandatory` - specifies if attribute is mandatory to create connection
