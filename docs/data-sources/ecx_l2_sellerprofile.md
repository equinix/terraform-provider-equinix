---
layout: "equinix"
page_title: "Equinix: equinix_ecx_l2_sellerprofile"
subcategory: ""
description: |-
  Provides Equinix Fabric Layer2 Seller Profile resource
---

# Data Source: equinix_ecx_l2_sellerprofile

Use this data source to get details of Equinix Fabric layer 2
seller profile with a given name and / or organization.

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

- `name` - (Optional) Name of the seller profile
- `organization_name` - (Optional) Name of seller's organization
- `organization_global_name` - (Optional) Name of seller's global organization

## Attributes Reference

- `uuid` - Unique identifier of the seller profile
- `description` - Seller Profile text description
- `speed_from_api` - Boolean that indicates if seller is deriving connection speed
from an API call
- `speed_customization_allowed` - Boolean that indicates if seller allows customer
to enter a custom connection speed
- `redundancy_required` - Boolean that indicate if seller requires connections
to be redundant
- `encapsulation` - Seller profile's encapsulation (either Dot1q or QinQ)
- `speed_band` - One or more specifications of speed/bandwidth supported by given
seller profile
  - `speed_band.#.speed` - Speed/bandwidth supported by given service profile
  - `speed_band.#.speed_unit` - Unit of the speed/bandwidth supported by given
  service profile
- `metro` - One or more specifications of metro locations supported by seller profile
  - `metro.#.code` - Location metro code
  - `metro.#.name` - Location metro name
  - `metro.#.ibxes` - List of IBXes supported within given metro
  - `metro.#.regions` - List of regions supported within given metro
- `additional_info` - One or more specifications of additional buyer information
attributes that can be provided in connection definition that uses given seller profile
  - `additional_info.#.name` - Name of additional information attribute
  - `additional_info.#.description` - Textual description of additional information
 attribute
  - `additional_info.#.data_type` - Data type of additional information attribute.
  Either *BOOLEAN*, *INTEGER* or *STRING*
  - `additional_info.#.mandatory` - Specifies if additional information
  is mandatory to create connection
