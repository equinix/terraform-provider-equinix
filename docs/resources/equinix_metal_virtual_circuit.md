---
subcategory: "Metal"
---

# equinix_metal_virtual_circuit (Resource)

Use this resource to associate VLAN with a Dedicated Port from
[Equinix Fabric - software-defined interconnections](https://metal.equinix.com/developers/docs/networking/fabric/#associating-a-vlan-with-a-dedicated-port).

See the [Virtual Routing and Forwarding documentation](https://deploy.equinix.com/developers/docs/metal/layer2-networking/vrf/) for product details and API reference material.

## Example Usage

Pick an existing Project and dedicated Connection, create a VLAN and use `equinix_metal_virtual_circuit`
to associate it with a Primary Port of the Connection.

```hcl
locals {
  project_id = "52000fb2-ee46-4673-93a8-de2c2bdba33c"
  conn_id = "73f12f29-3e19-43a0-8e90-ae81580db1e0"
}

data "equinix_metal_connection" test {
  connection_id = local.conn_id
}

resource "equinix_metal_vlan" "test" {
  project_id = local.project_id
  metro      = data.equinix_metal_connection.test.metro
}

resource "equinix_metal_virtual_circuit" "test" {
  connection_id = local.conn_id
  project_id = local.project_id
  port_id = data.equinix_metal_connection.test.ports[0].id
  vlan_id = equinix_metal_vlan.test.id
  nni_vlan = 1056
}
```

## Argument Reference

The following arguments are supported:

* `connection_id` - (Required) UUID of Connection where the VC is scoped to.
* `project_id` - (Required) UUID of the Project where the VC is scoped to.
* `port_id` - (Required) UUID of the Connection Port where the VC is scoped to.
* `nni_vlan` - (Required) Equinix Metal network-to-network VLAN ID.
* `vlan_id` - (Required) UUID of the VLAN to associate.
* `name` - (Optional) Name of the Virtual Circuit resource.
* `description` - (Optional) Description for the Virtual Circuit resource.
* `tags` - (Optional) Tags for the Virtual Circuit resource.
* `speed` - (Optional) Speed of the Virtual Circuit resource.
* `vrf_id` - (Optional) UUID of the VRF to associate.
* `peer_asn` - (Optional, required with `vrf_id`) The BGP ASN of the peer. The same ASN may be the used across several VCs, but it cannot be the same as the local_asn of the VRF.
* `subnet` - (Optional, required with `vrf_id`) A subnet from one of the IP
  blocks associated with the VRF that we will help create an IP reservation for. Can only be either a /30 or /31.
  * For a /31 block, it will only have two IP addresses, which will be used for
  the metal_ip and customer_ip.
  * For a /30 block, it will have four IP addresses, but the first and last IP addresses are not usable. We will default to the first usable IP address for the metal_ip.
* `metal_ip` - (Optional, required with `vrf_id`) The Metal IP address for the SVI (Switch Virtual Interface) of the VirtualCircuit. Will default to the first usable IP in the subnet.
* `customer_ip` - (Optional, required with `vrf_id`) The Customer IP address which the CSR switch will peer with. Will default to the other usable IP in the subnet.
* `md5` - (Optional, only valid with `vrf_id`) The password that can be set for the VRF BGP peer

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Status of the virtal circuit.
* `vnid` - VNID VLAN parameter, see the [documentation for Equinix Fabric](https://metal.equinix.com/developers/docs/networking/fabric/).
* `nni_vnid` - NNI VLAN parameters, see the [documentation for Equinix Fabric](https://metal.equinix.com/developers/docs/networking/fabric/).

## Import

This resource can be imported using an existing Virtual Circuit ID:

```sh
terraform import equinix_metal_virtual_circuit {existing_id}
```
