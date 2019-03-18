---
layout: "packet"
page_title: "Packet: packet_vlan"
sidebar_current: "docs-packet-resource-vlan"
description: |-
  Provides a resource for Packet Virtual Network.
---

# packet_vlan

Provides a resource to allow users to manage Virtual Networks in their projects. VLANs are used in [Layer 2 networking setup](https://support.packet.com/kb/articles/layer-2-overview).

## Example Usage

```hcl
# Create a new VLAN in datacenter "ewr1"

resource "packet_vlan" "vlan1" {
  description = "VLAN in New Jersey"
  facility    = "ewr1"
  project_id  = "${local.project_id}"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) ID of parent project
* `facility` - (Required) Facility where to create the VLAN
* `description` - Description string

## Attributes Reference

The following attributes are exported:

* `vxlan` - VXLAN segment ID
