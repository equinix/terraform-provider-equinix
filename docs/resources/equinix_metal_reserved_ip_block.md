---
subcategory: "Metal"
---

# equinix_metal_reserved_ip_block (Resource)

Provides a resource to create and manage blocks of reserved IP addresses in a project.

When a user provisions first device in a facility, Equinix Metal API automatically allocates IPv6/56 and private IPv4/25 blocks.
The new device then gets IPv6 and private IPv4 addresses from those block. It also gets a public IPv4/31 address.
Every new device in the project and facility will automatically get IPv6 and private IPv4 addresses from these pre-allocated blocks.
The IPv6 and private IPv4 blocks can't be created, only imported. With this resource, it's possible to create either public IPv4 blocks or global IPv4 blocks.

Public blocks are allocated in a facility. Addresses from public blocks can only be assigned to devices in the facility. Public blocks can have mask from /24 (256 addresses) to /32 (1 address). If you create public block with this resource, you must fill the facility argmument.

Addresses from global blocks can be assigned in any facility. Global blocks can have mask from /30 (4 addresses), to /32 (1 address). If you create global block with this resource, you must specify type = "global_ipv4" and you must omit the facility argument.

Once IP block is allocated or imported, an address from it can be assigned to device with the `equinix_metal_ip_attachment` resource.

## Example Usage

Allocate reserved IP blocks:

```hcl
# Allocate /31 block of max 2 public IPv4 addresses in Silicon Valley (sv15) facility for myproject

resource "equinix_metal_reserved_ip_block" "two_elastic_addresses" {
  project_id = local.project_id
  facility   = "sv15"
  quantity   = 2
}

# Allocate 1 floating IP in Sillicon Valley (sv) metro

resource "equinix_metal_reserved_ip_block" "test" {
  project_id = local.project_id
  type       = "public_ipv4"
  metro      = "sv"
  quantity   = 1
}

# Allocate 1 global floating IP, which can be assigned to device in any facility

resource "equinix_metal_reserved_ip_block" "test" {
  project_id = local.project_id
  type       = "global_ipv4"
  quantity   = 1
}
```

Allocate a block and run a device with public IPv4 from the block

```hcl
# Allocate /31 block of max 2 public IPv4 addresses in Silicon Valley (sv15) facility
resource "equinix_metal_reserved_ip_block" "example" {
  project_id = local.project_id
  facility   = "sv15"
  quantity   = 2
}

# Run a device with both public IPv4 from the block assigned

resource "equinix_metal_device" "nodes" {
  project_id       = local.project_id
  facilities       = ["sv15"]
  plan             = "c3.small.x86"
  operating_system = "ubuntu_20_04"
  hostname         = "test"
  billing_cycle    = "hourly"

  ip_address {
    type            = "public_ipv4"
    cidr            = 31
    reservation_ids = [equinix_metal_reserved_ip_block.example.id]
  }

  ip_address {
    type = "private_ipv4"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The metal project ID where to allocate the address block.
* `quantity` - (Required) The number of allocated `/32` addresses, a power of 2.
* `type` - (Optional) One of `global_ipv4` or `public_ipv4`. Defaults to `public_ipv4` for backward
compatibility.
* `facility` - (Optional) Facility where to allocate the public IP address block, makes sense only
if type is `public_ipv4` and must be empty if type is `global_ipv4`. Conflicts with `metro`.
* `metro` - (Optional) Metro where to allocate the public IP address block, makes sense only
if type is `public_ipv4` and must be empty if type is `global_ipv4`. Conflicts with `facility`.
* `description` - (Optional) Arbitrary description.
* `tags` - (Optional) String list of tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the block.
* `cidr_notation` - Address and mask in CIDR notation, e.g. `147.229.15.30/31`.
* `network` - Network IP address portion of the block specification.
* `netmask` - Mask in decimal notation, e.g. `255.255.255.0`.
* `cidr` - length of CIDR prefix of the block as integer.
* `address_family` - Address family as integer. One of `4` or `6`.
* `public` - Boolean flag whether addresses from a block are public.
* `global` - Boolean flag whether addresses from a block are global (i.e. can be assigned in any
facility).

-> **NOTE:** Idempotent reference to a first `/32` address from a reserved block might look
like `join("/", [cidrhost(metal_reserved_ip_block.myblock.cidr_notation,0), "32"])`.

## Import

This resource can be imported using an existing IP reservation ID:

```sh
terraform import equinix_metal_reserved_ip_block {existing_ip_reservation_id}
```
