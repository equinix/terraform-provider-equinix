---
page_title: "Equinix Metal: precreated_port"
subcategory: ""
description: |-
  Fetch device ports
---

# metal_port

Use this data source to read ports of existing devices. You can read port by either its UUID, or by a device UUID and port name.

## Example Usage

Create a device and read it's eth0 port to the datasource.

```hcl
locals {
  project_id = "<UUID_of_your_project>"
}

resource "metal_device" "test" {
  hostname         = "tfacc-test-device-port"
  plan             = "c3.medium.x86"
  facilities       = ["sv15"]
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

data "metal_port" "test" {
    device_id        = metal_device.test.id
    port_name        = "eth0"
}
```

## Argument Reference

* `id` - (Required) ID of the port to read, conflicts with device_id.
* `device_id` - (Required) 4 or 6, depending on which block you are looking for.
* `port_name` - (Required) Whether to look for public or private block.

## Attributes Reference

* `bond_port` - Whether port is a bond port or physical port
* `network_type` - One of layer2-bonded, layer2-individual, layer3, hybrid, hybrid-bonded
