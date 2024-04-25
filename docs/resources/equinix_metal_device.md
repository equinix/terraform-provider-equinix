---
subcategory: "Metal"
---

# equinix_metal_device (Resource)

Provides an Equinix Metal device resource. This can be used to create,
modify, and delete devices.

~> **NOTE:** All arguments including the `root_password` and `user_data` will be stored in
 the raw state as plain-text.
[Read more about sensitive data in state](https://developer.hashicorp.com/terraform/language/state/sensitive-data).

## Example Usage

Create a device and add it to cool_project

```hcl
resource "equinix_metal_device" "web1" {
  hostname         = "tf.coreos2"
  plan             = "c3.small.x86"
  metro            = "sv"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}
```

Same as above, but boot via iPXE initially, using the Ignition Provider for provisioning

```hcl
resource "equinix_metal_device" "pxe1" {
  hostname         = "tf.coreos2-pxe"
  plan             = "c3.small.x86"
  metro            = "sv"
  operating_system = "custom_ipxe"
  billing_cycle    = "hourly"
  project_id       = local.project_id
  ipxe_script_url  = "https://rawgit.com/cloudnativelabs/pxe/master/packet/coreos-stable-metal.ipxe"
  always_pxe       = "false"
  user_data        = data.ignition_config.example.rendered
}
```

Create a device without a public IP address in metro ny, with only a /30 private IPv4 subnet (4 IP addresses)

```hcl
resource "equinix_metal_device" "web1" {
  hostname         = "tf.coreos2"
  plan             = "c3.small.x86"
  metro            = "ny"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
  ip_address {
    type = "private_ipv4"
    cidr = 30
  }
}
```

Deploy device on next-available reserved hardware and do custom partitioning.

```hcl
resource "equinix_metal_device" "web1" {
  hostname                = "tftest"
  plan                    = "c3.small.x86"
  metro                   = "ny"
  operating_system        = "ubuntu_20_04"
  billing_cycle           = "hourly"
  project_id              = local.project_id
  hardware_reservation_id = "next-available"
  storage                 = <<EOS
{
  "disks": [
    {
      "device": "/dev/sda",
      "wipeTable": true,
      "partitions": [
        {
          "label": "BIOS",
          "number": 1,
          "size": "4096"
        },
        {
          "label": "SWAP",
          "number": 2,
          "size": "3993600"
        },
        {
          "label": "ROOT",
          "number": 3,
          "size": "0"
        }
      ]
    }
  ],
  "filesystems": [
    {
      "mount": {
        "device": "/dev/sda3",
        "format": "ext4",
        "point": "/",
        "create": {
          "options": [
            "-L",
            "ROOT"
          ]
        }
      }
    },
    {
      "mount": {
        "device": "/dev/sda2",
        "format": "swap",
        "point": "none",
        "create": {
          "options": [
            "-L",
            "SWAP"
          ]
        }
      }
    }
  ]
}
EOS
}
```

Create a device and allow the `user_data` and `custom_data` attributes to change in-place (i.e., without destroying and recreating the device):

```hcl
resource "equinix_metal_device" "pxe1" {
  hostname         = "tf.coreos2-pxe"
  plan             = "c3.small.x86"
  metro            = "sv"
  operating_system = "custom_ipxe"
  billing_cycle    = "hourly"
  project_id       = local.project_id
  ipxe_script_url  = "https://rawgit.com/cloudnativelabs/pxe/master/packet/coreos-stable-metal.ipxe"
  always_pxe       = "false"
  user_data        = local.user_data
  custom_data      = local.custom_data

  behavior {
    allow_changes = [
      "custom_data",
      "user_data"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `always_pxe` - (Optional) If true, a device with OS `custom_ipxe` will continue to boot via iPXE
on reboots.
* `behavior` - (Optional) Behavioral overrides that change how the resource handles certain attribute updates. See [Behavior](#behavior) below for more details.
* `billing_cycle` - (Optional) monthly or hourly
* `custom_data` - (Optional) A string of the desired Custom Data for the device.  By default, changing this attribute will cause the provider to destroy and recreate your device.  If `reinstall` is specified or `behavior.allow_changes` includes `"custom_data"`, the device will be updated in-place instead of recreated.
* `description` - (Optional) The device description.
* `facilities` - (**Deprecated**) List of facility codes with deployment preferences. Equinix Metal API will go
through the list and will deploy your device to first facility with free capacity. List items must
be facility codes or `any` (a wildcard). To find the facility code, visit
[Facilities API docs](https://metal.equinix.com/developers/api/facilities/), set your API auth
token in the top of the page and see JSON from the API response. Conflicts with `metro`.  Use metro instead; read the [facility to metro migration guide](https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices)
* `force_detach_volumes` - (Optional) Delete device even if it has volumes attached. Only applies
for destroy action.
* `hardware_reservation_id` - (Optional) The UUID of the hardware reservation where you want this
device deployed, or `next-available` if you want to pick your next available reservation
automatically. Changing this from a reservation UUID to `next-available` will re-create the device
in another reservation. Please be careful when using hardware reservation UUID and `next-available`
together for the same pool of reservations. It might happen that the reservation which Equinix
Metal API will pick as `next-available` is the reservation which you refer with UUID in another
equinix_metal_device resource. If that happens, and the equinix_metal_device with the UUID is
created later, resource creation will fail because the reservation is already in use (by the
resource created with `next-available`). To workaround this, have the `next-available` resource
[explicitly depend_on](https://learn.hashicorp.com/terraform/getting-started/dependencies.html#implicit-and-explicit-dependencies)
the resource with hardware reservation UUID, so that the latter is created first. For more details,
see [issue #176](https://github.com/packethost/terraform-provider-packet/issues/176).
* `hostname` - (Optional) The device hostname used in deployments taking advantage of Layer3 DHCP
or metadata service configuration.
* `ip_address` - (Optional) A list of IP address types for the device. See
[IP address](#ip-address) below for more details.
* `ipxe_script_url` - (Optional) URL pointing to a hosted iPXE script. More information is in the
[Custom iPXE](https://metal.equinix.com/developers/docs/servers/custom-ipxe/) doc.
* `metro` - (Optional) Metro area for the new device. Conflicts with `facilities`.
* `operating_system` - (Required) The operating system slug. To find the slug, or visit
[Operating Systems API docs](https://metal.equinix.com/developers/api/operatingsystems), set your
API auth token in the top of the page and see JSON from the API response.
* `plan` - (Required) The device plan slug. To find the plan slug, visit the
[bare-metal server](https://deploy.equinix.com/product/bare-metal/servers/) and [plan documentation](https://deploy.equinix.com/developers/docs/metal/hardware/standard-servers/).
* `project_id` - (Required) The ID of the project in which to create the device
* `project_ssh_key_ids` - (Optional) Array of IDs of the project SSH keys which should be added to the device. If you specify this array, only the listed project SSH keys (and any SSH keys for the users specified in user_ssh_key_ids) will be added. If no SSH keys are specified (both user_ssh_keys_ids and project_ssh_key_ids are empty lists or omitted), all parent project keys, parent project members keys and organization members keys will be included.  Project SSH keys can be created with the [equinix_metal_project_ssh_key](equinix_metal_project_ssh_key.md) resource.
* `user_ssh_key_ids` - (Optional) Array of IDs of the users whose SSH keys should be added to the device. If you specify this array, only the listed users' SSH keys (and any project SSH keys specified in project_ssh_key_ids) will be added. If no SSH keys are specified (both user_ssh_keys_ids and project_ssh_key_ids are empty lists or omitted), all parent project keys, parent project members keys and organization members keys will be included. User SSH keys can be created with the [equinix_metal_ssh_key](equinix_metal_ssh_key.md) resource.
* `reinstall` - (Optional) Whether the device should be reinstalled instead of destroyed when
modifying user_data, custom_data, or operating system. See [Reinstall](#reinstall) below for more
details.
* `storage` - (Optional) JSON for custom partitioning. Only usable on reserved hardware. More
information in in the
[Custom Partitioning and RAID](https://metal.equinix.com/developers/docs/servers/custom-partitioning-raid/)
doc. Please note that the disks.partitions.size attribute must be a string, not an integer. It can
be a number string, or size notation string, e.g. "4G" or "8M" (for gigabytes and megabytes).
* `tags` - (Optional) Tags attached to the device.
* `termination_time` - (Optional) Timestamp for device termination. For example `2021-09-03T16:32:00+03:00`.
If you don't supply timezone info, timestamp is assumed to be in UTC.
* `user_data` - (Optional) A string of the desired User Data for the device.  By default, changing this attribute will cause the provider to destroy and recreate your device.  If `reinstall` is specified or `behavior.allow_changes` includes `"user_data"`, the device will be updated in-place instead of recreated.
* `wait_for_reservation_deprovision` - (Optional) Only used for devices in reserved hardware. If
set, the deletion of this device will block until the hardware reservation is marked provisionable
(about 4 minutes in August 2019).

### Behavior

The `behavior` block has below fields:

* `allow_changes` - (Optional) List of attributes that are allowed to change without recreating the instance. Supported attributes: `custom_data`, `user_data`"

### IP address

The `ip_address` block has below fields:

* `type` - (Required) One of `private_ipv4`, `public_ipv4`, `public_ipv6`.
* `cidr` - (Optional) CIDR suffix for IP address block to be assigned, i.e. amount of addresses.
* `reservation_ids` - (Optional) List of UUIDs of [IP block reservations](metal_reserved_ip_block.md)
from which the public IPv4 address should be taken.

You can supply one `ip_address` block per IP address type. If you use the `ip_address` you must
always pass a block for `private_ipv4`.

To learn more about using the reserved IP addresses for new devices, see the examples in the
[equinix_metal_reserved_ip_block](metal_reserved_ip_block.md) documentation.

### Reinstall

The `reinstall` block has below fields:

* `enabled` - (Optional) Whether the provider should favour reinstall over destroy and create. Defaults to
`false`.
* `preserve_data` - (Optional) Whether the non-OS disks should be kept or wiped during reinstall.
Defaults to `false`.
* `deprovision_fast` - (Optional) Whether the OS disk should be filled with `00h` bytes before reinstall.
Defaults to `false`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/configuration/resources#operation-timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when creating the Device. This includes the time to provision the OS.
* `update` - (Defaults to 20 mins) Used when updating the Device. This includes the time needed to reprovision instances when `reinstall` arguments are used.
* `delete` - (Defaults to 20 mins) Used when deleting the Device. This includes the time to deprovision a hardware reservation when `wait_for_reservation_deprovision` is enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_private_ipv4` - The ipv4 private IP assigned to the device.
* `access_public_ipv4` - The ipv4 maintenance IP assigned to the device.
* `access_public_ipv6` - The ipv6 maintenance IP assigned to the device.
* `billing_cycle` - The billing cycle of the device (monthly or hourly).
* `created` - The timestamp for when the device was created.
* `deployed_facility` - (**Deprecated**) The facility where the device is deployed. Use metro instead; read the [facility to metro migration guide](https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices)
* `deployed_hardware_reservation_id` - ID of hardware reservation where this device was deployed.
It is useful when using the `next-available` hardware reservation.
* `description` - Description string for the device.
* `hostname` - The hostname of the device.
* `id` - The ID of the device.
* `locked` - Whether the device is locked or unlocked. Locking a device prevents you from deleting or reinstalling the device or performing a firmware update on the device, and it prevents an instance with a termination time set from being reclaimed, even if the termination time was reached
* `metro` - The metro area where the device is deployed.
* `network` - The device's private and public IP (v4 and v6) network details. See
[Network Attribute](#network-attribute) below for more details.
* `network_type` - (Deprecated) Network type of a device, used in
[Layer 2 networking](https://metal.equinix.com/developers/docs/networking/layer2/). Since this
attribute is deprecated you should handle Network Type with one of
[equinix_metal_port](equinix_metal_port.md),
[equinix_metal_device_network_type](equinix_metal_device_network_type.md) resources or
[equinix_metal_port](../data-sources/equinix_metal_port.md) datasource.
See [network_types guide](../guides/network_types.md) for more info.
* `operating_system` - The operating system running on the device.
* `plan` - The hardware config of the device.
* `ports` - List of ports assigned to the device. See [Ports Attribute](#ports-attribute) below for
more details.
* `project_id` - The ID of the project the device belongs to.
* `root_password` - Root password to the server (disabled after 24 hours).
* `sos_hostname` - The hostname to use for [Serial over SSH](https://deploy.equinix.com/developers/docs/metal/resilience-recovery/serial-over-ssh/) access to the device
* `ssh_key_ids` - List of IDs of SSH keys deployed in the device, can be both user and project SSH keys.
* `state` - The status of the device.
* `tags` - Tags attached to the device.
* `updated` - The timestamp for the last time the device was updated.

### Network Attribute

When a device is run without any special network, it will have 3 networks:

* Public IPv4 at `equinix_metal_device.name.network.0`.
* IPv6 at `equinix_metal_device.name.network.1`.
* Private IPv4 at `equinix_metal_device.name.network.2`.

-> **NOTE:** Elastic addresses stack by type. An assigned public IPv4 will go after the management
public IPv4 (to index 1), and will then shift the indices of the IPv6 and private IPv4. Assigned
private IPv4 will go after the management private IPv4 (to the end of the network list).

Each element in the `network` list exports:

* `address` - IPv4 or IPv6 address string.
* `cidr` - Bit length of the network mask of the address.
* `gateway` - Address of router.
* `public` - Whether the address is routable from the Internet.
* `family` - IP version. One of `4`, `6`.

### Ports Attribute

Each element in the `ports` list exports:

* `name` - Name of the port (e.g. `eth0`, or `bond0`).
* `id` - ID of the port.
* `type` - Type of the port (e.g. `NetworkPort` or `NetworkBondPort`).
* `mac` - MAC address assigned to the port.
* `bonded` - Whether this port is part of a bond in bonded network setup.

## Import

This resource can be imported using an existing device ID:

```sh
terraform import equinix_metal_device {existing_device_id}
```
