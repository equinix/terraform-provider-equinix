---
page_title: "Equinix Metal: metal_device_bgp_neighbors"
subcategory: ""
description: |-
  Provides a datasource for listing BGP neighbors of an Equinix Metal device
---

# metal_device_bgp_neighbors

Use this datasource to retrieve list of BGP neighbors of a device in the Equinix Metal host.

To have any BGP neighbors listed, the device must be in [BGP-enabled project](../r/project.html) and have a [BGP session](../r/bgp_session.html) assigned.

To learn more about using BGP in Equinix Metal, see the [metal_bgp_session](../r/bgp_session.html) resource documentation.

## Example Usage

```hcl
# Get Project by name and print UUIDs of its users

data "metal_device_bgp_neighbors" "test" {
  device_id = "4c641195-25e5-4c3c-b2b7-4cd7a42c7b40"
}

output "bgp_neighbors_listing" {
  value = data.metal_device_bgp_neighbors.test.bgp_neighbors
}
```

## Argument Reference

The following arguments are supported:

* `device_id` - UUID of BGP-enabled device whose neighbors to list

## Attributes Reference

The following attributes are exported:

* `bgp_neighbors` - array of BGP neighbor records with attributes:
  * `address_family` - IP address version, 4 or 6
  * `customer_as` - Local autonomous system number
  * `customer_ip` - Local used peer IP address
  * `md5_enabled` - Whether BGP session is password enabled
  * `md5_password` - BGP session password in plaintext (not a checksum)
  * `multihop` - Whether the neighbor is in EBGP multihop session
  * `peer_as` - Peer AS number (different than customer_as for EBGP)
  * `peer_ips` - Array of IP addresses of this neighbor's peers
  * `routes_in` - Array of incoming routes. Each route has attributes:
    * `route` - CIDR expression of route (IP/mask)
    * `exact` - (bool) Whether the route is exact
  * `routes_out` - Array of outgoing routes in the same format
  