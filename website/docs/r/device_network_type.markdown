---
layout: "packet"
page_title: "Packet: packet_device_network_type"
sidebar_current: "docs-packet-resource-device-network-type"
description: |-
  Provides a resource to manage network type of Packet devices.
---

# packet_device_network_type

This resource controls network type of Packet devices.

To learn more about Layer 2 networking in Packet, refer to

* https://www.packet.com/resources/guides/layer-2-configurations/
* https://www.packet.com/developers/docs/network/advanced/layer-2/

## Example Usage

```
resource "packet_device" "test" {
  hostname         = "tfacc-device-port-vlan-attachment-test"
  plan             = "s1.large.x86"
  facilities       = ["nrt1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

resource "packet_device_network_type" "test" {
  device_id = packet_device.test.id
  type      = "hybrid"
}
```

If you are attaching VLAN to a device (i.e. using packet_port_vlan_attachment), link the device ID from this resource, in order to make the port attachment implicitly dependent on the state of the network type. If you link the device ID from the packet_device resource, Terraform will not wait for the network type change. See examples in [packet_port_vlan_attachment](port_vlan_attachment.html).


## Import

This resource can also be imported using existing device ID:

```
$ terraform import packet_device_network_type {existing device_id}
```


## Argument Reference

The following arguments are supported:

* `device_id` - (Required) The ID of the device on which the network type should be set.
* `type` - (Required) Network type to set. Must be one of `layer3`, `hybrid`, `layer2-individual` and `layer2-bonded`.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the controlled device. Use this in linked resources, if you need to wait for the network type change. It is the same as `device_id`.
