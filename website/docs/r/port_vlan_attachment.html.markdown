---
layout: "packet"
page_title: "Packet: packet_port_vlan_attachment"
sidebar_current: "docs-packet-resource-port-vlan-attachment"
description: |-
  Provides a Resource for Attaching VLANs to Device Ports
---

# packet_port_vlan_attachment

Provides a resource to attach device ports to VLANs.

Device and VLAN must be in the same facility.

If you need this resource to add the port back to bond on removal, set `force_bond = true`.

To learn more about Layer 2 networking in Packet, refer to
* https://support.packet.com/kb/articles/layer-2-configurations
* https://support.packet.com/kb/articles/layer-2-overview

## Example Usage

```hcl
# Hybrid network type

resource "packet_vlan" "test" {
  description = "VLAN in New Jersey"
  facility    = "ewr1"
  project_id  = "${local.project_id}"
}

resource "packet_device" "test" {
  hostname         = "test"
  plan             = "m1.xlarge.x86"
  facilities       = ["ewr1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${local.project_id}"
  network_type     = "hybrid"
}

resource "packet_port_vlan_attachment" "test" {
  device_id = "${packet_device.test.id}"
  port_name = "eth1"
  vlan_id = "${packet_vlan.test.id}"
}


# Layer 2 network

resource "packet_device" "test" {
  hostname         = "test"
  plan             = "m1.xlarge.x86"
  facility         = "ewr1"
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${local.project_id}"
  network_type     = "layer2-individual"
}

resource "packet_vlan" "test1" {
  description = "VLAN in New Jersey"
  facility    = "ewr1"
  project_id  = "${packet_project.test.id}"
}

resource "packet_vlan" "test2" {
  description = "VLAN in New Jersey"
  facility    = "ewr1"
  project_id  = "${local.project_id}"
}

resource "packet_port_vlan_attachment" "test1" {
  device_id = "${packet_device.test.id}"
  vlan_id = "${packet_vlan.test1.id}"
  port_name = "eth1"
}

resource "packet_port_vlan_attachment" "test2" {
  device_id = "${packet_device.test.id}"
  vlan_id = "${packet_vlan.test2.id}"
  port_name = "eth1"
}
```

## Argument Reference

The following arguments are supported:

* `device_id` - (Required) ID of device to be assigned to the VLAN
* `port_name` - (Required) Name of network port to be assigned to the VLAN
* `vlan_id` - (Required) ID of target VLAN
* `force_bond` - Add port back to the bond when this resource is removed. Default is false.

## Attribute Referece

* `id` - UUID of device port used in the assignment
