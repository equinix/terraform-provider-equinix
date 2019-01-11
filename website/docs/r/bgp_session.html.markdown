---
layout: "packet"
page_title: "Packet: packet_bgp_session"
sidebar_current: "docs-packet-resource-bgp-session"
description: |-
  BGP session in Packet Host
---

# packet\_bgp\_session

Provides a resource to manage BGP session in Packet Host. Refer to [Packet BGP documentation](https://support.packet.com/kb/articles/bgp) for more details.

You need to have BGP config enabled in your project.

BGP session must be linked to a device running BIRD or other BGP routing client which will control route advertisements via the session to Packet's upstream routers. 

## Example Usage

```hcl

```

## Argument Reference

The following arguments are supported:

* `device_id` - (Required) ID of device to which to assign the subnet
* `cidr_notation` - (Required) CIDR notation of subnet from block reserved in the same
  project and facility as the device

## Attributes Reference

The following attributes are exported:
* `route_object`: Specifies AS-MACRO to use when building client route filters (as opposed to AS number and read only)
* `ranges`: The IP block ranges associated to your ASN (only populated in global BGP)
* `max_prefix`: The maximum number of route filters allowed per server
