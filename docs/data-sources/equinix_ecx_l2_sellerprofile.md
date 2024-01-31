---
subcategory: "Fabric"
---

# equinix_ecx_l2_sellerprofile (Data Source)

!> **DEPRECATED** End of Life will be June 30th, 2024. Use `equinix_fabric_service_profile` instead.

Use this data source to get details of Equinix Fabric layer 2 seller profile with a given name
and / or organization.

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

The following arguments are supported:

* `name` - (Optional) Name of the seller profile.
* `organization_name` - (Optional) Name of seller's organization.
* `organization_global_name` - (Optional) Name of seller's global organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - Unique identifier of the seller profile.
* `description` - Seller Profile text description.
* `speed_from_api` - Boolean that indicates if seller is deriving connection speed from an API call.
* `speed_customization_allowed` - Boolean that indicates if seller allows customer to enter a
custom connection speed.
* `redundancy_required` - Boolean that indicate if seller requires connections to be redundant
* `encapsulation` - Seller profile's encapsulation (either Dot1q or QinQ).
* `speed_band` - One or more specifications of speed/bandwidth supported by given seller profile.
See [Speed Band Attribute](#speed-band-attribute) below for more details.
* `metro` - One or more specifications of metro locations supported by seller profile.
See [Metro Attribute](#metro-attribute) below for more details.

* `additional_info` - One or more specifications of additional buyer information attributes that
can be provided in connection definition that uses given seller profile.
See [Additional Info Attribute](#additional-info-attribute) below for more details.

### Speed Band Attribute

Each element in the `speed_band` set exports:

* `speed` - Speed/bandwidth supported by given service profile.
* `speed_unit` - Unit of the speed/bandwidth supported by given service profile.

### Metro Attribute

Each element in the `metro` set exports:

* `code` - Location metro code.
* `name` - Location metro name.
* `ibxes` - List of IBXes supported within given metro.
* `regions` - List of regions supported within given.

### Additional Info Attribute

Each element in the `additional_info` set exports:

* `name` - Name of additional information attribute.
* `description` - Textual description of additional information attribute.
* `data_type` - Data type of additional information attribute. One of `BOOLEAN`, `INTEGER` or
`STRING`.
* `mandatory` - Specifies if additional information is mandatory to create
connection.
