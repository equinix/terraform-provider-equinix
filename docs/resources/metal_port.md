---
subcategory: "Metal"
---

# equinix_metal_port (Resource)

Use this resource to configure network ports on an Equinix Metal device. This resource can control both physical and bond ports.

This Terraform resource doesn't create an API resource in Equinix Metal, but rather provides finer control for [Layer 2 networking](https://metal.equinix.com/developers/docs/layer2-networking/).

The port resource referred is created together with device and accessible either via the device resource or over `/port/<uuid>` API path.

-> To achieve the network configurations available in the portal it may require the creation and combination of various `equinix_metal_port` resources. See the [Network Types Guide](../guides/network_types.md) for examples of this resource.

## Argument Reference

The following arguments are supported:

* `port_id` - (Required) ID of the port to read.
* `bonded` - (Required) Whether the port should be bonded.
* `layer2` - (Optional) Whether to put the port to Layer 2 mode, valid only for bond ports.
* `vlan_ids` - (Optional) List of VLAN UUIDs to attach to the port, valid only for L2 and Hybrid ports.
* `vxlan_ids` - (Optional) List of VXLAN IDs to attach to the port, valid only for L2 and Hybrid ports.
* `native_vlan_id` - (Optional) UUID of a VLAN to assign as a native VLAN. It must be one of attached VLANs (from `vlan_ids` parameter).
* `reset_on_delete` - (Optional) Behavioral setting to reset the port to default settings (layer3 bonded mode without any vlan attached) before delete/destroy.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/configuration/resources#operation-timeouts) for certain actions:

These timeout includes the time to disbond, convert to L2/L3, bond and update native vLAN.

* `create` - (Defaults to 30 mins) Used when creating the Port.
* `update` - (Defaults to 30 mins) Used when updating the Port.
* `delete` - (Defaults to 30 mins) Used when deleting the Port.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - Name of the port, e.g. `bond0` or `eth0`.
* `network_type` - One of layer2-bonded, layer2-individual, layer3, hybrid and hybrid-bonded. This attribute is only set on bond ports.
* `type` - Type is either "NetworkBondPort" for bond ports or "NetworkPort" for bondable ethernet ports.
* `mac` - MAC address of the port.
* `bond_id` - UUID of the bond port.
* `bond_name` - Name of the bond port.
* `disbond_supported` - Flag indicating whether the port can be removed from a bond.
