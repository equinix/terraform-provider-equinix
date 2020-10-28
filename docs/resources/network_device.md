---
layout: "equinix"
page_title: "Equinix: equinix_network_device"
subcategory: ""
description: |-
 Provides Network Edge device resource
---

# Resource: equinix_network_device

Resource `equinix_network_device` allows creation and management of Network Edge
virtual network devices.

Network Edge virtual network devices can be created in two modes:

* **managed** (default) where Equinix manages connectivity and services in the
device and customer gets limited access to the device
* **self-configured** where customer provisions and manages own services in the device
with less restricted access. Some device types are offered only in this mode

In addition to management modes, there are two software license modes available:

* **subscription**  where Equinix provides software license, including end-to-end
support, and bills for the service respectively.
* **BYOL** [bring your own license] where customer brings his own, already procured
device software license. There are no charges associated with such license.
It is the only licensing mode for *self-configured* devices

## Example Usage

```hcl
# Create pair of redundant, managed CSR1000V routers with license subscription
# in two different metro locations

data "equinix_network_account" "dc" {
  metro_code = "DC"
}

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "csr1000v-ha" {
  name            = "tf-csr1000v-p"
  throughput      = 500
  throughput_unit = "Mbps"
  metro_code      = data.equinix_network_account.dc.metro_code
  type_code       = "CSR1000V"
  package_code    = "SEC"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  hostname        = "csr1000v-p"
  acls            = ["10.0.0.0/24", "192.168.0.0/24", "1.1.1.1/32"]
  term_length     = 6
  account_number  = data.equinix_network_account.dc.number
  version         = "16.09.05"
  core_count      = 2
  secondary_device {
    name           = "tf-csr1000v-s"
    metro_code     = data.equinix_network_account.sv.metro_code
    hostname       = "csr1000v-s"
    acls           = ["1.1.1.1/32", "2.2.2.2/32", "4.4.4.4/32", "5.5.5.5/32"]
    notifications  = ["john@equinix.com", "marry@equinix.com"]
    account_number = data.equinix_network_account.sv.number
  }
}
```

## Argument Reference

* `name` - (Required) Device name
* `type_code` - (Required) Device type code
* `metro_code` - (Required) Metro location of a device
* `throughput` - (Required) License throughput for a device
* `throughput_unit` - (Required) License throughput unit (Mbps or Gbps)
* `hostname` - (Required) Device hostname
* `package_code` - (Required) Code of a software package used for a device
* `version` - (Required) Software version for a device
* `byol` - (Optional) Boolean value determining if device licensing mode will be
*bring your own license* or *subscription* (default)
* `license_token` - (Optional) License Token can be provided for some device types
in BYOL licensing mode
* `acls` - (Optional) List of IP address subnets that will be loaded as an access
control list for users accessing a device
* `account_number` - (Required) Billing account number for a device
* `notifications` - (Required) List of email addresses that will receive device
status notifications
* `purchase_order_number` - (Optional) Purchase order number associated
with a device order
* `term_length` - (Required) Term length
* `additional_bandwidth` - (Optional) Additional Internet bandwidth, in Mbps,
that will be added for the device in addition to 15Mbps included by default
* `order_reference` - (Optional) Name/number used to identify device order on
the invoice
* `interface_count` - (Optional) Number of network interfaces on a device. If not
specified then default number for a given device type will be used.
* `core_count` - (Required) Number of CPU cores for a device
* `self_managed` - (Optional) Boolean value determining if device will be self-managed
or Equinix managed (default)
* `vendor_configuration` - (Optional) map of device parameters and values that
are vendor specific and can or have to be provided for some device types
* `secondary_device` - (Optional) Definition of secondary device for redundant
device configurations
  * `name` - (Required) Device name
  * `metro_code` - (Required) Metro location of a device
  * `hostname` - (Required) Device hostname
  * `license_token` - (Optional) License Token can be provided for some device types
  * `acls` - (Optional) List of IP address subnets that will be loaded as an access
  * `account_number` - (Required) Billing account number for a device
  * `notifications` - (Required) List of email addresses that will receive device
  * `additional_bandwidth` - (Optional) Additional Internet bandwidth, in Mbps,
  * `vendor_configuration` - (Optional) map of device parameters and values that

## Attributes Reference

* `uuid` - Device universally unique identifier
* `status` - Device provisioning status
* `license_status` - Device license registration status
* `acls_status` - Device ACL provisioning
* `ibx` - Name of Equinix exchange
* `region` - Region in which device metro is located
* `ssh_ip_address` - IP address to use for SSH connectivity with the device
* `ssh_ip_fqdn` - FQDN to use for SSH connectivity with the device
* `redundancy_type` - Indicates if device is primary or secondary
(in HA configuration)
* `redundant_id` - Universally unique identifier for a redundant device
(in HA configuration)
* `interface` - List of device interfaces
  * `id` - interface identifier
  * `name` - interface name
  * `status` -  interface status (AVAILABLE, RESERVED, ASSIGNED)
  * `operational_status` - interface operation status (up or down)
  * `mac_address` - interface MAC address
  * `ip_address` - interface IP address
  * `assigned_type` - interface management type (Equinix Managed or empty)
  * `type` - interface type
