---
page_title: "Equinix Metal: spot_market_price"
subcategory: ""
description: |-
  Get an Equinix Metal Spot Market Price
---

# packet\_operating\_system

Use this data source to get Equinix Metal Spot Market Price.

## Example Usage

```hcl
data "packet_spot_market_price" "example" {
  facility = "ewr1"
  plan     = "c1.small.x86"
}
```

## Argument Reference

* `facility` - (Required) Name of the facility.
* `plan` - (Required) Name of the plan.

## Attributes Reference

* `price` - Current spot market price for given plan in given facility.
