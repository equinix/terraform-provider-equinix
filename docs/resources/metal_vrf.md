---
subcategory: "Metal"
---

# equinix_metal_vrf (Resource)

Use this resource to manage a VRF.

See the [Virtual Routing and Forwarding documentation](https://deploy.equinix.com/developers/docs/metal/layer2-networking/vrf/) for product details and API reference material.

## Example Usage

Create a VRF in your desired metro and project with any IP ranges that you want the VRF to route and forward.

```terraform
resource "equinix_metal_project" "example" {
    name = "example"
}

resource "equinix_metal_vrf" "example" {
    description = "VRF with ASN 65000 and a pool of address space that includes 192.168.100.0/25"
    name        = "example-vrf"
    metro       = "da"
    local_asn   = "65000"
    ip_ranges   = ["192.168.100.0/25", "192.168.200.0/25"]
    project_id  = equinix_metal_project.example.id
}
```

Create IP reservations and assign them to a Metal Gateway resources. The Gateway will be assigned the first address in the block.

```terraform
resource "equinix_metal_reserved_ip_block" "example" {
    description = "Reserved IP block (192.168.100.0/29) taken from on of the ranges in the VRF's pool of address space."
    project_id  = equinix_metal_project.example.id
    metro       = equinix_metal_vrf.example.metro
    type        = "vrf"
    vrf_id      = equinix_metal_vrf.example.id
    cidr        = 29
    network     = "192.168.100.0"
}

resource "equinix_metal_vlan" "example" {
    description = "A VLAN for Layer2 and Hybrid Metal devices"
    metro       = equinix_metal_vrf.example.metro
    project_id  = equinix_metal_project.example.id
}

resource "equinix_metal_gateway" "example" {
    project_id        = equinix_metal_project.example.id
    vlan_id           = equinix_metal_vlan.example.id
    ip_reservation_id = equinix_metal_reserved_ip_block.example.id
}
```

Attach a Virtual Circuit from a Dedicated Metal Connection to the Metal Gateway.

```terraform
data "equinix_metal_connection" "example" {
    connection_id = var.metal_dedicated_connection_id
}

resource "equinix_metal_virtual_circuit" "example" {
    name          = "example-vc"
    description   = "Virtual Circuit"
    connection_id = data.equinix_metal_connection.example.id
    project_id    = equinix_metal_project.example.id
    port_id       = data.equinix_metal_connection.example.ports[0].id
    nni_vlan      = 1024
    vrf_id        = equinix_metal_vrf.example.id
    peer_asn      = 65530
    subnet        = "192.168.100.16/31"
    metal_ip      = "192.168.100.16"
    customer_ip   = "192.168.100.17"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) User-supplied name of the VRF, unique to the project
* `metro` - (Required) Metro ID or Code where the VRF will be deployed.
* `project_id` - (Required) Project ID where the VRF will be deployed.
* `description` - (Optional) Description of the VRF.
* `local_asn` - (Optional) The 4-byte ASN set on the VRF.
* `ip_ranges` - (Optional) All IPv4 and IPv6 Ranges that will be available to BGP Peers. IPv4 addresses must be /8 or smaller with a minimum size of /29. IPv6 must be /56 or smaller with a minimum size of /64. Ranges must not overlap other ranges within the VRF.

## Attributes Reference

No additional attributes are exported.

## Import

This resource can be imported using an existing VRF ID:

```sh
terraform import equinix_metal_vrf {existing_id}
```
