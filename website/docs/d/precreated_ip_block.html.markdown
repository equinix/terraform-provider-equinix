---
layout: "packet"
page_title: "Packet: precreated_ip_block"
sidebar_current: "docs-packet-datasource-precreated-ip-block"
description: |-
  Load automatically created IP blocks from your Packet project
---

# packet\_precreated\_ip\_block

Use this data source to get CIDR expression for precreated IPv6 and IPv4 blocks in Packet.
You can then use the cidrsubnet TF builtin function to derive subnets.

## Example Usage

```hcl

# Create project, device in it, and then assign /64 subnet from precreated block
# to the new device

resource "packet_project" "test" {
    name = "testpro"
}


resource "packet_device" "web1" {
  hostname         = "tftest"
  plan             = "baremetal_0"
  facility         = "ewr1"
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
}

# we have to make the datasource depend on the device. Here I do it implicitly
# with the project_id param, because an explicity "depends_on" attribute in
# a datasource taints the state:
# https://github.com/hashicorp/terraform/issues/11806
data "packet_precreated_ip_block" "test" {
    facility         = "ewr1"
    project_id       = "${packet_device.test.project_id}"
    address_family   = 6
    public           = true
}

# The precreated IPv6 blocks are /56, so to get /64, we specify 8 more bits for network.
# The cirdsubnet interpolation will pick second /64 subnet from the precreated block.

resource "packet_ip_attachment" "from_ipv6_block" {
    device_id = "${packet_device.web1.id}"
    cidr_notation = "${cidrsubnet(data.packet_precreated_ip_block.test.cidr_notation,8,2)}"
}

```

## Argument Reference

 * `project_id` - (Required) ID of the project where the searched block should be.
 * `address_family` - (Required) 4 or 6, depending on which block you are looking for.
 * `public` - (Required) Whether to look for public or private block. 
 * `facility` - (Required) Facility of the searched block.

## Attributes Reference

 * `cidr_notation` - CIDR notation of the looked up block.

