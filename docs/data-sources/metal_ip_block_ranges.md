---
subcategory: "Metal"
---

# equinix_metal_ip_block_ranges (Data Source)

Use this datasource to get CIDR expressions for allocated IP blocks of all the types in a project, optionally filtered by facility or metro.

There are four types of IP blocks in Equinix: equinix_metal_global IPv4, public IPv4, private IPv4 and IPv6. Both global and public IPv4 are routable from the Internet. Public IPv4 blocks are allocated in a facility or metro, and addresses from it can only be assigned to devices in that location. Addresses from Global IPv4 block can be assigned to a device in any metro.

The datasource has 4 list attributes: `global_ipv4`, `public_ipv4`, `private_ipv4` and `ipv6`, each listing CIDR notation (`<network>/<mask>`) of respective blocks from the project.

## Example Usage

```terraform
# List CIDR expressions of all the allocated IP block in you project.

# Declare your project ID
locals {
  project_id = "<UUID_of_your_project>"
}

data "equinix_metal_ip_block_ranges" "test" {
  project_id = local.project_id
}

output "out" {
  value = data.equinix_metal_ip_block_ranges.test
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) ID of the project from which to list the blocks.
* `facility` - (**Deprecated**) Facility code filtering the IP blocks. Global IPv4 blocks will be listed anyway. If you omit this and metro, all the block from the project will be listed. Use metro instead; read the [facility to metro migration guide](https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices)
* `metro` - (Optional) Metro code filtering the IP blocks. Global IPv4 blocks will be listed anyway. If you omit this and facility, all the block from the project will be listed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `global_ipv4` - list of CIDR expressions for Global IPv4 blocks in the project.
* `public_ipv4` - list of CIDR expressions for Public IPv4 blocks in the project.
* `private_ipv4` - list of CIDR expressions for Private IPv4 blocks in the project.
* `ipv6` - list of CIDR expressions for IPv6 blocks in the project.
