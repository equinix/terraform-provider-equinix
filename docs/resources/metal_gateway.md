---
subcategory: "Metal"
---

# equinix_metal_gateway (Resource)

Use this resource to create Metal Gateway resources in Equinix Metal.

See the [Virtual Routing and Forwarding documentation](https://deploy.equinix.com/developers/docs/metal/layer2-networking/vrf/) for product details and API reference material.

## Example Usage

```terraform
# Create Metal Gateway for a VLAN with a private IPv4 block with 8 IP addresses

resource "equinix_metal_vlan" "test" {
  description = "test VLAN in SV"
  metro       = "sv"
  project_id  = local.project_id
}

resource "equinix_metal_gateway" "test" {
  project_id               = local.project_id
  vlan_id                  = equinix_metal_vlan.test.id
  private_ipv4_subnet_size = 8
}
```

```terraform
# Create Metal Gateway for a VLAN and reserved IP address block

resource "equinix_metal_vlan" "test" {
  description = "test VLAN in SV"
  metro       = "sv"
  project_id  = local.project_id
}

resource "equinix_metal_reserved_ip_block" "test" {
  project_id = local.project_id
  metro      = "sv"
  quantity   = 8
}

resource "equinix_metal_gateway" "test" {
  project_id        = local.project_id
  vlan_id           = equinix_metal_vlan.test.id
  ip_reservation_id = equinix_metal_reserved_ip_block.test.id
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) UUID of the project where the gateway is scoped to.
* `vlan_id` - (Required) UUID of the VLAN where the gateway is scoped to.
* `ip_reservation_id` - (Optional) UUID of Public or VRF IP Reservation to associate with the gateway, the reservation must be in the same metro as the VLAN, conflicts with `private_ipv4_subnet_size`.
* `private_ipv4_subnet_size` - (Optional) Size of the private IPv4 subnet to create for this metal gateway, must be one of `8`, `16`, `32`, `64`, `128`. Conflicts with `ip_reservation_id`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `state` - Status of the gateway resource.
* `vrf_id` - UUID of the VRF associated with the IP Reservation

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

* `delete` - (Default `20m`)
