---
layout: "packet"
page_title: "Packet: packet_volume_attachment"
sidebar_current: "docs-packet-resource-volume-attachment"
description: |-
  Provides attachment of volumes to devices in the Packet Host.
---

# packet\_volume\_attachment

Provides attachment of Packet Block Storage Volume to Devices.

Device and volume must be in the same location (facility).

Once attached by Terraform, they must then be mounted using the `packet_block_attach` and `packet_block_detach` scripts.

## Example Usage

```hcl
  resource "packet_project" "test_project" {
      name = "test-project"
  }

  resource "packet_device" "test_device_va" {
      hostname         = "terraform-test-device-va"
      plan             = "t1.small.x86"
      facility         = "ewr1"
      operating_system = "ubuntu_16_04"
      billing_cycle    = "hourly"
      project_id       = "${packet_project.test_project.id}"
  }

  resource "packet_volume" "test_volume_va" {
      plan = "storage_1"
      billing_cycle = "hourly"
      size = 100
      project_id = "${packet_project.test_project.id}"
      facility = "ewr1"
      snapshot_policies = { snapshot_frequency = "1day", snapshot_count = 7 }
  }

  resource "packet_volume_attachment" "test_volume_attachment" {
      device_id = "${packet_device.test_device_va.id}"
      volume_id = "${packet_volume.test_volume_va.id}"
  }
```

## Argument Reference

The following arguments are supported:

* `volume_id` - (Required) The ID of the volume to attach
* `device_id` - (Required) The ID of the device to which the volume should be attached

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the volume attachment