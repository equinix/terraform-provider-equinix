---
subcategory: "Metal"
---

# Resource: equinix_metal_port

Use this resource to set up network ports on an Equnix Metal device. This resource can control both physical and bond ports.

This Terraform resource doesn't create an API resource in Equinix Metal, but rather provides finer control for (Layer 2 networking)[https://metal.equinix.com/developers/docs/layer2-networking/].

The port resource referred is created together with device and accessible either via the device resource or over `/port/<uuid>` API path.

## Example Usage

See the [Network Types Guide](../guides/network_types.md) for examples of this resource.

## Argument Reference

* `port_id` - (Required) ID of the port to read
* `bonded` - (Required) Whether the port should be bonded
* `layer2` - (Optional) Whether to put the port to Layer 2 mode, valid only for bond ports
* `vlan_ids` - (Optional) List of VLAN UUIDs to attach to the port, valid only for L2 and Hybrid ports
* `vxlan_ids` - (Optional) List of VXLAN IDs to attach to the port, valid only for L2 and Hybrid ports
* `native_vlan_id` - (Optional) UUID of a VLAN to assign as a native VLAN. It must be one of attached VLANs (from `vlan_ids` parameter), valid only for physical (non-bond) ports
* `reset_on_delete` - (Optional) Behavioral setting to reset the port to default settings. For a bond port it means layer3 without vlans attached, eth ports will be bonded without native vlan and vlans attached

## Attributes Reference

* `name` - Name of the port, e.g. `bond0` or `eth0`
* `network_type` - One of layer2-bonded, layer2-individual, layer3, hybrid and hybrid-bonded. This attribute is only set on bond ports.
* `type` - Type is either "NetworkBondPort" for bond ports or "NetworkPort" for bondable ethernet ports
* `mac` - MAC address of the port
* `bond_id` - UUID of the bond port
* `bond_name` - Name of the bond port
* `bonded` - Flag indicating whether the port is bonded
* `disbond_supported` - Flag indicating whether the port can be removed from a bond
* `vlan_ids` - List of VLAN UUIDs to attach to the port
* `vxlan_ids` - List of VXLAN IDs to attach to the port

