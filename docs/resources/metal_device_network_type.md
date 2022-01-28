---
layout: "metal"
page_title: "Equinix: equinix_metal_device_network_type"
sidebar_current: "docs-metal-resource-device-network-type"
description: |-
  Provides a resource to manage network type of Equinix Metal devices.
---

# Resource: equinix_metal_device_network_type

This resource controls network type of Equinix Metal devices.

To learn more about Layer 2 networking in Equinix Metal, refer to

* <https://metal.equinix.com/developers/docs/networking/layer2/>
* <https://metal.equinix.com/developers/docs/networking/layer2-configs/>

If you are attaching VLAN to a device (i.e. using equinix_metal_port_vlan_attachment), link the device ID from this resource, in order to make the port attachment implicitly dependent on the state of the network type. If you link the device ID from the equinix_metal_device resource, Terraform will not wait for the network type change. See examples in [metal_port_vlan_attachment](port_vlan_attachment).

## Example Usage

See the [Network Types Guide](../guides/network_types.md) for examples of this resource and to learn about the recommended `equinix_metal_port` alternative.

## Import

This resource can also be imported using existing device ID:

```sh
terraform import equinix_metal_device_network_type {existing device_id}
```

## Argument Reference

The following arguments are supported:

* `device_id` - (Required) The ID of the device on which the network type should be set.
* `type` - (Required) Network type to set. Must be one of `layer3`, `hybrid`, `layer2-individual` and `layer2-bonded`.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the controlled device. Use this in linked resources, if you need to wait for the network type change. It is the same as `device_id`.
