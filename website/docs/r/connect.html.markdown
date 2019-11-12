---
layout: "packet"
page_title: "Packet: packet_connect [DEPRECATED]"
sidebar_current: "docs-packet-resource-connect"
description: |-
  Provides a resource for Packet Connect, which is now deprecated and will be fully removed in a later release.
---

# packet_connect [DEPRECATED]

Provides a resource for [Packet Connect](https://www.packet.com/cloud/all-features/packet-connect/), a link between Packet VLANs and VLANs in other cloud providers, which is now deprecated. Packet Connect will be fully removed in a later release.

## Example Usage

```hcl
# Create a new VLAN in ewr1 and connect it to Azure ExpressRoute 

resource "packet_vlan" "vlan1" {
  description = "VLAN in New Jersey"
  facility    = "ewr1"
  project_id  = "${local.project_id}"
}

resource "packet_connect" "my_expressroute" {
  name        = "test"
  facility    = "ewr1"
  project_id  = "${local.project_id}"
  # provider ID for Azure ExpressRoute is ed5de8e0-77a9-4d3b-9de0-65281d3aa831
  provider_id = "ed5de8e0-77a9-4d3b-9de0-65281d3aa831"
  # provider_payload for Azure ExpressRoute provider is your ExpressRoute
  # authorization key (in UUID format)
  provider_payload = "58b4ec12-af34-4435-5435-db3bde4a4b3a"
  port_speed  = 100
  vxlan       = "${packet_vlan.vlan1.vxlan}"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name for the Connect resource
* `facility` - (Required) Facility where to create the VLAN
* `project_id` - (Required) ID of parent project
* `provider_id` - (Required) ID of Connect Provider. Provider IDs are
  * Azure ExpressRoute - "ed5de8e0-77a9-4d3b-9de0-65281d3aa831"
* `provider_payload` - (Required) Authorization key for the Connect provider
* `port_speed` - (Required) Port speed in Mbps
* `vxlan` - (Required) VXLAN Network identifier of the linked Packet VLAN

## Attributes Reference

The following attributes are exported:

* `status` - Status of the Connect resource, one of PROVISIONING, PROVISIONED, DEPROVISIONING, DEPROVISIONED
