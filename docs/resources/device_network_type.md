---
layout: "metal"
page_title: "Equinix Metal: metal_device_network_type"
sidebar_current: "docs-metal-resource-device-network-type"
description: |-
  Provides a resource to manage network type of Equinix Metal devices.
---

# metal_device_network_type

This resource controls network type of Equinix Metal devices.

To learn more about Layer 2 networking in Equinix Metal, refer to

* <https://metal.equinix.com/developers/docs/networking/layer2/>
* <https://metal.equinix.com/developers/docs/networking/layer2-configs/>

If you are attaching VLAN to a device (i.e. using metal_port_vlan_attachment), link the device ID from this resource, in order to make the port attachment implicitly dependent on the state of the network type. If you link the device ID from the metal_device resource, Terraform will not wait for the network type change. See examples in [metal_port_vlan_attachment](port_vlan_attachment).

## Example Usage

### Create one s1.large device and put it to hybrid network mode

```hcl
resource "metal_device" "test" {
  hostname         = "tfacc-device-port-vlan-attachment-test"
  plan             = "s1.large.x86"
  facilities       = ["nrt1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

resource "metal_device_network_type" "test" {
  device_id = metal_device.test.id
  type      = "hybrid"
}
```

### Create two devices in hybrid mode and add a VLAN to their eth1 ports

```hcl
locals {
    project_id = "<uuid>"
    device_count = 2
}

resource "metal_vlan" "test" {
  facility    = "nrt1"
  project_id  = local.project_id
}


resource "metal_device" "test" {
  count            = local.device_count
  hostname         = "test${count.index}"
  plan             = "s1.large.x86"
  facilities       = ["nrt1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

resource "metal_device_network_type" "test" {
  count     = local.device_count
  device_id = metal_device.test[count.index].id
  type      = "hybrid"
}


resource "metal_port_vlan_attachment" "test" {
  count     = local.device_count
  device_id = metal_device_network_type.test[count.index].id
  port_name = "eth1"
  vlan_vnid = metal_vlan.test.vxlan
}

```


## Import

This resource can also be imported using existing device ID:

```sh
terraform import metal_device_network_type {existing device_id}
```

## Argument Reference

The following arguments are supported:

* `device_id` - (Required) The ID of the device on which the network type should be set.
* `type` - (Required) Network type to set. Must be one of `layer3`, `hybrid`, `layer2-individual` and `layer2-bonded`.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the controlled device. Use this in linked resources, if you need to wait for the network type change. It is the same as `device_id`.
