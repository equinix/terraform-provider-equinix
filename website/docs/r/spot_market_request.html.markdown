---
layout: "packet"
page_title: "Packet: packet_spot_market_request"
sidebar_current: "docs-packet-spot-market-request"
description: |-
  Provides a Packet Spot Market Request Resource.
---

# packet\_volume

Provides a Packet Spot Market Request resource to allow you to
manage spot market requests on your account.

## Example Usage

```hcl
# Create a spot market request
resource "packet_spot_market_request" "req" {
  project_id      = "${packet_project.cool_project.id}"
  "max_bid_price" = 0.03
  "facilities"    = ["ewr1"]
  "devices_min"   = 1
  "devices_max"   = 1

  "instance_parameters" {
    "hostname"         = "testspot"
    "billing_cycle"    = "hourly"
    "operating_system" = "coreos_stable"
    "plan"             = "baremetal_0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `devices_max` - (Required) Maximum number devices to be created
* `devices_min` - (Required) Miniumum number devices to be created
* `facilities` - (Required) Facility IDs where devices should be created
* `instance_parameters` - (Required) Device parameters. See device resource for details
* `project_id` - (Required) Project ID
   