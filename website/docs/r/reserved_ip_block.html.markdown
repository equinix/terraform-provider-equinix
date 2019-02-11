---
layout: "packet"
page_title: "Packet: packet_reserved_ip_block"
sidebar_current: "docs-packet-resource-resevered-ip-block"
description: |-
  Provides a Resource for reserving IP addresses in the Packet Host
---

# packet\_reserved\_ip\_block

Provides a resource to create and manage blocks of reserved IP addresses in a project.

When a user provisions first device in a facility, Packet API automatically allocates IPv6/56 and private IPv4/25 blocks.
The new device then gets IPv6 and private IPv4 addresses from those block. It also gets a public IPv4/31 address.
Every new device in the project and facility will automatically get IPv6 and private IPv4 addresses from these pre-allocated blocks.
The IPv6 and private IPv4 blocks can't be created, only imported. With this resource, it's possible to create either public IPv4 blocks or global IPv4 blocks.

Public blocks are allocated in a facility. Addresses from public blocks can only be assigned to devices in the facility. Public blocks can have mask from /24 (256 addresses) to /32 (1 address). If you create public block with this resource, you must fill the facility argmument.

Addresses from global blocks can be assigned in any facility. Global blocks can have mask from /30 (4 addresses), to /32 (1 address). If you create global block with this resource, you must specify type = "global_ipv4" and you must omit the facility argument.

Once IP block is allocated or imported, an address from it can be assigned to device with the `packet_ip_attachment` resource.

## Example Usage

```hcl
# Allocate /30 block of max 2 public IPv4 addresses in Parsippany, NJ (ewr1) for myproject

resource "packet_reserved_ip_block" "two_elastic_addresses" {
    project_id = "${packet_project.myproject.id}"
    facility = "ewr1"
    quantity = 2
}

# Allocate 1 global floating IP, which can be assigned to device in any facility

resource "packet_reserved_ip_block" "test" {
    project_id = "${packet_project.myproject.id}"
    type     = "global_ipv4"
	quantity = 1
}`
```


## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The packet project ID where to allocate the address block
* `quantity` - (Required) The number of allocated /32 addresses, a power of 2
* `facility` - (Optional) Facility where to allocate the public IP address block, makes sense only for type==public_ipv4
* `type` - (Optional) Either "global_ipv4" or "public_ipv4", defaults to "public_ipv4" for backward compatibility


## Attributes Reference

The following attributes are exported:

* `facility` - The facility where the block was allocated, empty for global blocks
* `project_id` - To which project the addresses beling
* `quantity` - Number of /32 addresses in the block
* `id` - The unique ID of the block
* `cidr_notation` - Address and mask in CIDR notation, e.g. "147.229.15.30/31"
* `network` - Network IP address portion of the block specification
* `netmask` - Mask in decimal notation, e.g. "255.255.255.0"
* `cidr` - length of CIDR prefix of the block as integer
* `address_family` - Address family as integer (4 or 6)
* `public` - boolean flag whether addresses from a block are public
* `global` - boolean flag whether addresses from a block are global (i.e. can be assigned in any facility)

Idempotent reference to a first /32 address from a reserved block might look like 
`"${cidrhost(packet_reserved_ip_block.test.cidr_notation,0)}/32"`.
