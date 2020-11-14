---
page_title: "Equinix Metal: precreated_ip_block"
subcategory: ""
description: |-
  Load automatically created IP blocks from your Equinix Metal project
---

# metal\_precreated\_ip\_block

Use this data source to get CIDR expression for precreated IPv6 and IPv4 blocks in Equinix Metal.
You can then use the cidrsubnet TF builtin function to derive subnets.

## Example Usage

```hcl
# Create device in your project and then assign /64 subnet from precreated block
# to the new device

# Declare your project ID
locals {
  project_id = "<UUID_of_your_project>"
}

resource "metal_device" "web1" {
  hostname         = "web1"
  plan             = "t1.small.x86"
  facilities       = ["ewr1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id

}

data "metal_precreated_ip_block" "test" {
  facility       = "ewr1"
  project_id     = local.project_id
  address_family = 6
  public         = true
}

# The precreated IPv6 blocks are /56, so to get /64, we specify 8 more bits for network.
# The cirdsubnet interpolation will pick second /64 subnet from the precreated block.

resource "metal_ip_attachment" "from_ipv6_block" {
  device_id     = metal_device.web1.id
  cidr_notation = cidrsubnet(data.metal_precreated_ip_block.test.cidr_notation, 8, 2)
}
```

## Argument Reference

* `project_id` - (Required) ID of the project where the searched block should be.
* `address_family` - (Required) 4 or 6, depending on which block you are looking for.
* `public` - (Required) Whether to look for public or private block.
* `global` - (Optional) Whether to look for global block. Default is false for backward compatibility.
* `facility` - (Optional) Facility of the searched block. (Optional) Only allowed for non-global blocks.

## Attributes Reference

* `cidr_notation` - CIDR notation of the looked up block.
