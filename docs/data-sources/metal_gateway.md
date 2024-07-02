---
subcategory: "Metal"
---

# equinix_metal_gateway (Data Source)

Use this datasource to retrieve Metal Gateway resources in Equinix Metal.

See the [Virtual Routing and Forwarding documentation](https://deploy.equinix.com/developers/docs/metal/layer2-networking/vrf/) for product details and API reference material.

## Example Usage

```terraform
# Create Metal Gateway for a VLAN with a private IPv4 block with 8 IP addresses

resource "equinix_metal_vlan" "test" {
  description = "test VLAN in SV"
  metro       = "sv"
  project_id  = local.project_id
}

data "equinix_metal_gateway" "test" {
  gateway_id = local.gateway_id
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required) UUID of the metal gateway resource to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `project_id` - UUID of the project where the gateway is scoped to.
* `vlan_id` - UUID of the VLAN where the gateway is scoped to.
* `vrf_id` - UUID of the VRF associated with the IP Reservation.
* `ip_reservation_id` - UUID of IP reservation block bound to the gateway.
* `private_ipv4_subnet_size` - Size of the private IPv4 subnet bound to this metal gateway. One of `8`, `16`, `32`, `64`, `128`.
* `state` - Status of the gateway resource.
