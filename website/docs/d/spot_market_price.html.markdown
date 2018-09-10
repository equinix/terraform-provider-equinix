---
layout: "packet"
page_title: "Packet: spot_market_price"
sidebar_current: "docs-packet-datasource-spot-market-price"
description: |-
  Get a Packet Spot Market Prices
---

# packet\_operating\_system

Use this data source to get Packet Spot Market Prices.

## Example Usage

```hcl
data "packet_spot_market_price" "example" {
  facility = "ewr1"
  plan     = "baremetal_1"
}
```

## Argument Reference

 * `facility` - (Required) Name of the facility.
 * `plan` - (Required) Name of the plan.

