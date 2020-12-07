---
page_title: "Equinix Metal: metal_ip_attachment"
subcategory: ""
description: |-
  Provides a Resource for Attaching IP Subnets from a Reserved Block to a Device
---

# metal\_ip\_attachment

Provides a resource to attach elastic IP subnets to devices.

To attach an IP subnet from a reserved block to a provisioned device, you must derive a subnet CIDR belonging to
one of your reserved blocks in the same project and facility as the target device.

For example, you have reserved IPv4 address block 147.229.10.152/30, you can choose to assign either the whole
block as one subnet to a device; or 2 subnets with CIDRs 147.229.10.152/31' and 147.229.10.154/31; or 4 subnets
with mask prefix length 32. More about the elastic IP subnets is [here](https://metal.equinix.com/developers/docs/networking/elastic-ips/).

Device and reserved block must be in the same facility.

## Example Usage

```hcl
# Reserve /30 block of max 2 public IPv4 addresses in Parsippany, NJ (ewr1) for myproject
resource "metal_reserved_ip_block" "myblock" {
  project_id = local.project_id
  facility   = "ewr1"
  quantity   = 2
}

# Assign /32 subnet (single address) from reserved block to a device
resource "metal_ip_attachment" "first_address_assignment" {
  device_id = metal_device.mydevice.id
  # following expression will result to sth like "147.229.10.152/32"
  cidr_notation = join("/", [cidrhost(metal_reserved_ip_block.myblock.cidr_notation, 0), "32"])
}
```

## Argument Reference

The following arguments are supported:

* `device_id` - (Required) ID of device to which to assign the subnet
* `cidr_notation` - (Required) CIDR notation of subnet from block reserved in the same
  project and facility as the device

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the assignment
* `device_id` - ID of device to which subnet is assigned
* `cidr_notation` - Assigned subnet in CIDR notation, e.g. "147.229.15.30/31"
* `gateway` - IP address of gateway for the subnet
* `network` - Subnet network address
* `netmask` - Subnet mask in decimal notation, e.g. "255.255.255.0"
* `cidr` - length of CIDR prefix of the subnet as integer
* `address_family` - Address family as integer (4 or 6)
* `public` - boolean flag whether subnet is reachable from the Internet
