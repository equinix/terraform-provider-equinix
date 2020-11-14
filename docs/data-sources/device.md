---
page_title: "Equinix Metal: metal_device"
subcategory: ""
description: |-
  Provides an Equinix Metal device datasource. This can be used to read existing devices.
---

# metal_device

Provides an Equinix Metal device datasource.

~> **Note:** All arguments including the `root_password` and `user_data` will be stored in
 the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
# Fetch a device data by hostname and show it's ID

data "metal_device" "test" {
  project_id = local.project_id
  hostname   = "mydevice"
}

output "id" {
  value = data.metal_device.test.id
}
```

```hcl
# Fetch a device data by ID and show its public IPv4
data "metal_device" "test" {}

output "ipv4" {
  value = data.metal_device.test.access_public_ipv4
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - The device name
* `project_id` - The id of the project in which the devices exists
* `device_id` - Device ID

User can lookup devices either by `device_id` or `project_id` and `hostname`.

## Attributes Reference

The following attributes are exported:

* `access_private_ipv4` - The ipv4 private IP assigned to the device
* `access_public_ipv4` - The ipv4 management IP assigned to the device
* `access_public_ipv6` - The ipv6 management IP assigned to the device
* `billing_cycle` - The billing cycle of the device (monthly or hourly)
* `facility` - The facility where the device is deployed.
* `description` - Description string for the device
* `hardware_reservation_id` - The id of hardware reservation which this device occupies
* `id` - The ID of the device
* `network` - The device's private and public IP (v4 and v6) network details. When a device is run without any special network configuration, it will have 3 networks:
  * Public IPv4 at `metal_device.name.network.0`
  * IPv6 at `metal_device.name.network.1`
  * Private IPv4 at `metal_device.name.network.2`
  Elastic addresses then stack by type - an assigned public IPv4 will go after the management public IPv4 (to index 1), and will then shift the indices of the IPv6 and private IPv4. Assigned private IPv4 will go after the management private IPv4 (to the end of the network list).
  The fields of the network attributes are:
  * `address` - IPv4 or IPv6 address string
  * `cidr` - Bit length of the network mask of the address
  * `gateway` - Address of router
  * `public` - Whether the address is routable from the Internet
  * `family` - IP version - "4" or "6"
* `network_type` - L2 network type of the device, one of "layer3", "layer2-bonded", "layer2-individual", "hybrid"
* `operating_system` - The operating system running on the device
* `plan` - The hardware config of the device
* `ports` - Ports assigned to the device
  * `name` - Name of the port (e.g. `eth0`, or `bond0`)
  * `id` - ID of the port
  * `type` - Type of the port (e.g. `NetworkPort` or `NetworkBondPort`)
  * `mac` - MAC address assigned to the port
  * `bonded` - Whether this port is part of a bond in bonded network setup
* `root_password` - Root password to the server (if still available)
* `ssh_key_ids` - List of IDs of SSH keys deployed in the device, can be both user or project SSH keys
* `state` - The state of the device
* `tags` - Tags attached to the device
