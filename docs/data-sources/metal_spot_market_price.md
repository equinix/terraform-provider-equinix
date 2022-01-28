---
page_title: "Equinix: equinix_metal_spot_market_price"
subcategory: ""
description: |-
  Get an Equinix Metal Spot Market Price
---

# metal\_operating\_system

Use this data source to get Equinix Metal Spot Market Price for a plan.

## Example Usage

Lookup by facility:

```hcl
data "equinix_metal_spot_market_price" "example" {
  facility = "ny5"
  plan     = "c3.small.x86"
}
```

Lookup by metro:

```hcl
data "equinix_metal_spot_market_price" "example" {
  metro    = "sv"
  plan     = "c3.small.x86"
}
```

## Argument Reference

* `plan` - (Required) Name of the plan.
* `facility` - (Optional) Name of the facility.
* `metro` - (Optional) Name of the metro.

## Attributes Reference

* `price` - Current spot market price for given plan in given facility.
