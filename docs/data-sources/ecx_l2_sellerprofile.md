---
layout: "equinix"
page_title: "Equinix: equinix_ecx_l2_sellerprofile"
subcategory: ""
description: |-
 Get information on ECX Fabric Layer2 Seller Profile
---

# Data Source: equinix_ecx_l2_sellerprofile

Use this data source to get details of ECX Fabric seller profile with a given name.

## Example usage

```hcl
data "equinix_ecx_l2_sellerprofile" "aws" {
  name = "AWS Direct Connect"
}

output "id" {
  value = data.equinix_ecx_l2_sellerprofile.aws.id
}
```

## Argument Reference

- `name` - (Optional) Name of seller profile
- `organization_name` - (Optional) Name of seller's organization
- `organization_global_name` - (Optional) Name of seller's global organization

## Attributes Reference

- `uuid` - Unique identifier of seller profile
- `speed_from_api` - information if seller is deriving connection speed
from an API call
- `speed_customization_allowed` - information if seller allows customer to enter
a custom connection speed
- `redundancy_required` - information if seller requires connections to be redundant
- `encapsulation` - seller profile's encapsulation (Dot1q or QinQ)
- `speed_band` - one or more specifications of speed/bandwidth supported by
seller profile
  - `speed` - speed/bandwidth supported by this profile
  - `speed_unit` - unit of the speed/bandwidth supported by this profile
- `metro` - one or more specifications of metro locations supported by seller profile
  - `code` - metro code
  - `name` - metro name
  - `ibxes` - list of IBXes supported within given metro
  - `regions` - list of regions supported within given metro
- `additional_info` - one or more specifications of additional buyer information
attributes that can be provided in connection definition that uses given seller profile
  - `name` - name of an attribute
  - `description` - textual description of an attribute
  - `data_type` - data type of an attribute _(BOOLEAN / INTEGER / STRING)_
  - `mandatory` - specifies if attribute is mandatory to create connection
