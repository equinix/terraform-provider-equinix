---
page_title: "Equinix: equinix_metal_gateway"
subcategory: ""
description: |-
  Create Equinix Metal Gateways
---

# metal\_gateway

Use this resource to create Metal Gateway resources in Equinix Metal.

## Example Usage

```hcl
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

```hcl
# Create Metal Gateway for a VLAN and reserved IP address block

resource "equinix_metal_vlan" "test" {
  description = "test VLAN in SV"
  metro       = "sv"
  project_id  = local.project_id
}

resource "equinix_metal_reserved_ip_block" "test" {
  project_id = local.project_id
  metro      = "sv"
  quantity   = 2
}

resource "equinix_metal_gateway" "test" {
  project_id        = local.project_id
  vlan_id           = equinix_metal_vlan.test.id
  ip_reservation_id = equinix_metal_reserved_ip_block.test.id
}
```

## Argument Reference

* `project_id` - UUID of the project where the gateway is scoped to
* `vlan_id` - UUID of the VLAN where the gateway is scoped to
* `ip_reservation_id` - (Optional) UUID of IP reservation block to bind to the gateway, the reservation must be in the same metro as the VLAN, conflicts with `private_ipv4_subnet_size`
* `private_ipv4_subnet_size` - (Optional) Size of the private IPv4 subnet to create for this metal gateway, must be one of (8, 16, 32, 64, 128), conflicts with `ip_reservation_id`

## Attributes Reference

* `state` - Status of the gateway resource
