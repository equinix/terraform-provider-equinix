---
page_title: "Equinix Metal: metal_facility"
subcategory: ""
description: |-
  Provides an Equinix Metal facility datasource. This can be used to read facilities.
---

# metal_facility

Provides an Equinix Metal facility datasource.

## Example Usage

```hcl
# Fetch a facility by code and show its ID

data "equinix_metal_facility" "ny5" {
  code = "ny5"
}

output "id" {
  value = data.metal_facility.ny5.id
}
```

```hcl
# Verify that facility "dc13" has capacity for provisioning 2 c3.small.x86 
  devices and 1 c3.medium.x86 device

data "equinix_metal_facility" "test" {
  code = "dc13"
  capacity {
    plan = "c3.small.x86"
    quantity = 2
  }
  capacity {
    plan = "c3.medium.x86"
    quantity = 1
  }
}

```

## Argument Reference

The following arguments are supported:

* `code` - The facility code
* `features_required` - Set of feature strings that the facility must have

Facilities can be looked up by `code`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the facility
* `name` - The name of the facility
* `features` - The features of the facility
* `metro` - The metro code the facility is part of
* `capacity` - (Optional) Ensure that queried facility has capacity for specified number of given plans
  - `plan` - device plan to check
  - `quantity` - number of device to check

