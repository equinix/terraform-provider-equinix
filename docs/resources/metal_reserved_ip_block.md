---
subcategory: "Metal"
---

# equinix_metal_reserved_ip_block (Resource)

Provides a resource to create and manage blocks of reserved IP addresses in a project.

When a user provisions first device in a metro, Equinix Metal API automatically allocates IPv6/56 and private IPv4/25 blocks. The new device then gets IPv6 and private IPv4 addresses from those block. It also gets a public IPv4/31 address. Every new device in the project and metro will automatically get IPv6 and private IPv4 addresses from these pre-allocated blocks. The IPv6 and private IPv4 blocks can't be created, only imported. With this resource, it's possible to create either public IPv4 blocks or global IPv4 blocks.

Public blocks are allocated in a metro. Addresses from public blocks can only be assigned to devices in the metro. Public blocks can have mask from /24 (256 addresses) to /32 (1 address). If you create public block with this resource, you must fill the metro argument.

Addresses from global blocks can be assigned in any metro. Global blocks can have mask from /30 (4 addresses), to /32 (1 address). If you create global block with this resource, you must specify type = "global_ipv4" and you must omit the metro argument.

Once IP block is allocated or imported, an address from it can be assigned to device with the `equinix_metal_ip_attachment` resource.

See the [Virtual Routing and Forwarding documentation](https://deploy.equinix.com/developers/docs/metal/layer2-networking/vrf/) for product details and API reference material.

## Example Usage

Allocate reserved IP blocks:

```terraform
# Allocate /31 block of max 2 public IPv4 addresses in Silicon Valley (sv) metro for myproject

resource "equinix_metal_reserved_ip_block" "two_elastic_addresses" {
  project_id = local.project_id
  metro      = "sv"
  quantity   = 2
}

# Allocate 1 floating IP in Silicon Valley (sv) metro

resource "equinix_metal_reserved_ip_block" "test" {
  project_id = local.project_id
  type       = "public_ipv4"
  metro      = "sv"
  quantity   = 1
}

# Allocate 1 global floating IP, which can be assigned to device in any metro

resource "equinix_metal_reserved_ip_block" "test" {
  project_id = local.project_id
  type       = "global_ipv4"
  quantity   = 1
}
```

Allocate a block and run a device with public IPv4 from the block

```terraform
# Allocate /31 block of max 2 public IPv4 addresses in Silicon Valley (sv) metro
resource "equinix_metal_reserved_ip_block" "example" {
  project_id = local.project_id
  metro      = "sv"
  quantity   = 2
}

# Run a device with both public IPv4 from the block assigned

resource "equinix_metal_device" "nodes" {
  project_id       = local.project_id
  metro            = "sv"
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
* `quantity` - (Optional) The number of allocated `/32` addresses, a power of 2. Required when `type` is not `vrf`.
* `type` - (Optional) One of `global_ipv4`, `public_ipv4`, or `vrf`. Defaults to `public_ipv4` for backward compatibility.
* `facility` - (**Deprecated**) Facility where to allocate the public IP address block, makes sense only if type is `public_ipv4` and must be empty if type is `global_ipv4`. Conflicts with `metro`. Use metro instead; read the [facility to metro migration guide](https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices)
* `metro` - (Optional) Metro where to allocate the public IP address block, makes sense only if type is `public_ipv4` and must be empty if type is `global_ipv4`. Conflicts with `facility`.
* `description` - (Optional) Arbitrary description.
* `tags` - (Optional) String list of tags.
* `vrf_id` - (Optional) Only valid and required when `type` is `vrf`. VRF ID for type=vrf reservations.
* `wait_for_state` - (Optional) Wait for the IP reservation block to reach a desired state on resource creation. One of: `pending`, `created`. The `created` state is default and recommended if the addresses are needed within the configuration. An error will be returned if a timeout or the `denied` state is encountered.
* `custom_data` - (Optional) Custom Data is an arbitrary object (submitted in Terraform as serialized JSON) to assign to the IP Reservation. This may be helpful for self-managed IPAM. The object must be valid JSON.
* `network` - (Optional) Only valid as an argument and required when `type` is `vrf`. An unreserved network address from an existing `ip_range` in the specified VRF.
* `cidr` - (Optional) Only valid as an argument and required when `type` is `vrf`. The size of the network to reserve from an existing VRF ip_range. `cidr` can only be specified with `vrf_id`. Range is 22-31. Virtual Circuits require 30-31. Other VRF resources must use a CIDR in the 22-29 range.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the block.
* `cidr_notation` - Address and mask in CIDR notation, e.g. `147.229.15.30/31`.
* `network` - Network IP address portion of the block specification.
* `netmask` - Mask in decimal notation, e.g. `255.255.255.0`.
* `cidr` - length of CIDR prefix of the block as integer.
* `address_family` - Address family as integer. One of `4` or `6`.
* `public` - Boolean flag whether addresses from a block are public.
* `global` - Boolean flag whether addresses from a block are global (i.e. can be assigned in any metro).
* `vrf_id` - VRF ID of the block when type=vrf

-> **NOTE:** Idempotent reference to a first `/32` address from a reserved block might look like `join("/", [cidrhost(metal_reserved_ip_block.myblock.cidr_notation,0), "32"])`.

## Import

This resource can be imported using an existing IP reservation ID:

```sh
terraform import equinix_metal_reserved_ip_block {existing_ip_reservation_id}
```
