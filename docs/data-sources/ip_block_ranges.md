---
page_title: "Equinix Metal: ip_block_ranges"
subcategory: ""
description: |-
  List IP address ranges allocated to a project
---

# metal\_ip\_block\_ranges

Use this datasource to get CIDR expressions for allocated IP blocks of all the types in a project, optionally filtered by facility.

There are four types of IP blocks in Equinix Metal: global IPv4, public IPv4, private IPv4 and IPv6. Both global and public IPv4 are routable from the Internet. Public IPv4 block is allocated in a facility, and addresses from it can only be assigned to devices in that facility. Addresses from Global IPv4 block can be assigned to a device in any facility.

The datasource has 4 list attributes: `global_ipv4`, `public_ipv4`, `private_ipv4` and `ipv6`, each listing CIDR notation (`<network>/<mask>`) of respective blocks from the project.

## Example Usage

```hcl
# List CIDR expressions of all the allocated IP block in you project.

# Declare your project ID
locals {
  project_id = "<UUID_of_your_project>"
}

data "metal_ip_block_ranges" "test" {
  project_id = local.project_id
}

output "out" {
  value = data.metal_ip_block_ranges.test
}
```

## Argument Reference

* `project_id` - (Required) ID of the project from which to list the blocks.
* `facility` - (Optional) Facility code filtering the IP blocks. Global IPv4 blcoks will be listed anyway. If you omit this, all the block from the project will be listed.

## Attributes Reference

* `global_ipv4` - list of CIDR expressions for Global IPv4 blocks in the project
* `public_ipv4` - list of CIDR expressions for Public IPv4 blocks in the project
* `private_ipv4` - list of CIDR expressions for Private IPv4 blocks in the project
* `ipv6` - list of CIDR expressions for IPv6 blocks in the project
