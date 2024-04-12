---
subcategory: "Metal"
---

# equinix_metal_device (Data Source)

The datasource can be used to fetch a single device.

If you need to fetch a list of devices which meet filter criteria, you can use the [equinix_metal_devices](equinix_metal_devices.md) datasource.

~> **Note:** All arguments including the `root_password` and `user_data` will be stored in the raw state as plain-text. [Read more about sensitive data in state](https://developer.hashicorp.com/terraform/language/state/sensitive-data).

## Example Usage

```terraform
# Fetch a device data by hostname and show it's ID

data "equinix_metal_device" "test" {
  project_id = local.project_id
  hostname   = "mydevice"
}

output "id" {
  value = data.equinix_metal_device.test.id
}
```

```terraform
# Fetch a device data by ID and show its public IPv4
data "equinix_metal_device" "test" {
  device_id = "4c641195-25e5-4c3c-b2b7-4cd7a42c7b40"
}

output "ipv4" {
  value = data.equinix_metal_device.test.access_public_ipv4
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Optional) The device name.
* `project_id` - (Optional) The id of the project in which the devices exists.
* `device_id` - (Optional) Device ID.

-> **NOTE:** You should pass either `device_id`, or both `project_id` and `hostname`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_private_ipv4` - The ipv4 private IP assigned to the device.
* `access_public_ipv4` - The ipv4 management IP assigned to the device.
* `access_public_ipv6` - The ipv6 management IP assigned to the device.
* `billing_cycle` - The billing cycle of the device (monthly or hourly).
* `facility` - (**Deprecated**) The facility where the device is deployed. Use metro instead; read the [facility to metro migration guide](https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices)
* `description` - Description string for the device.
* `hardware_reservation_id` - The id of hardware reservation which this device occupies.
* `id` - The ID of the device.
* `metro` - The metro where the device is deployed
* `network` - The device's private and public IP (v4 and v6) network details. See [Network Attribute](#network-attribute) below for more details.
* `network_type` - L2 network type of the device, one of `layer3`, `layer2-bonded`, `layer2-individual`, `hybrid`.
* `operating_system` - The operating system running on the device.
* `plan` - The hardware config of the device.
* `ports` - List of ports assigned to the device. See [Ports Attribute](#ports-attribute) below for more details.
* `root_password` - Root password to the server (if still available).
* `sos_hostname` - The hostname to use for [Serial over SSH](https://deploy.equinix.com/developers/docs/metal/resilience-recovery/serial-over-ssh/) access to the device
* `ssh_key_ids` - List of IDs of SSH keys deployed in the device, can be both user or project SSH keys.
* `state` - The state of the device.
* `tags` - Tags attached to the device.

### Network Attribute

When a device is run without any special network, it will have 3 networks:

* Public IPv4 at `equinix_metal_device.name.network.0`.
* IPv6 at `equinix_metal_device.name.network.1`.
* Private IPv4 at `equinix_metal_device.name.network.2`.

-> **NOTE:** Elastic addresses stack by type. An assigned public IPv4 will go after the management public IPv4 (to index 1), and will then shift the indices of the IPv6 and private IPv4. Assigned private IPv4 will go after the management private IPv4 (to the end of the network list).

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
