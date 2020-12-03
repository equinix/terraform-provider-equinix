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
  term_length     = 6
  account_number  = data.equinix_network_account.dc.number
  version         = "16.09.05"
  core_count      = 2
  secondary_device {
    name            = "tf-csr1000v-s"
    metro_code      = data.equinix_network_account.sv.metro_code
    hostname        = "csr1000v-s"
    notifications   = ["john@equinix.com", "marry@equinix.com"]
    account_number  = data.equinix_network_account.sv.number
  }
}
```

## Argument Reference

* `name` - (Required) Device name
* `type_code` - (Required) Device type code
* `metro_code` - (Required) Metro location of a device
* `hostname` - (Optional) Device hostname
* `package_code` - (Required) Code of a software package used for a device
* `version` - (Required) Software version for a device
* `core_count` - (Required) Number of CPU cores for a device
* `term_length` - (Required) Term length
* `self_managed` - (Optional) Boolean value determining if device will be self-managed
or Equinix managed (default)
* `byol` - (Optional) Boolean value determining if device licensing mode will be
*bring your own license* or *subscription* (default)
* `license_token` - (Optional) License Token can be provided for some device types
in BYOL licensing mode
* `license_file` - (Optional) Path to the license file that will be uploaded and
applied on a device. Applicable for some devices types in BYOL licensing mode
* `throughput` - (Optional) License throughput for a device
* `throughput_unit` - (Optional) License throughput unit (Mbps or Gbps)
* `account_number` - (Required) Billing account number for a device
* `notifications` - (Required) List of email addresses that will receive device
status notifications
* `purchase_order_number` - (Optional) Purchase order number associated
with a device order
* `order_reference` - (Optional) Name/number used to identify device order on
the invoice
* `acl_template_id` - (Optional) Identifier of an ACL template that
will be applied on the device
* `additional_bandwidth` - (Optional) Additional Internet bandwidth, in Mbps,
that will be added for the device in addition to 15Mbps included by default
* `interface_count` - (Optional) Number of network interfaces on a device. If not
specified then default number for a given device type will be used.
* `vendor_configuration` - (Optional) map of vendor specific configuration parameters
for a device
* `ssh-key` - (Optional) up to one definition of SSH key that will be provisioned
on a device
  * `ssh-key.#.username` - (Required) username associated with given key
  * `ssh-key.#.name` - (Required) name of SSH key as defined in
`equinix_network_ssh_key` resource
* `secondary_device` - (Optional) Definition of secondary device for redundant
device configurations
  * `secondary_device.#.name` - (Required) Secondary device name
  * `secondary_device.#.metro_code` - (Required) Metro location of a secondary device
  * `secondary_device.#.hostname` - (Optional) Secondary device hostname
  * `secondary_device.#.license_token` - (Optional) License Token can be provided
 for some device types o the device
  * `secondary_device.#.license_file` - (Optional) Path to the license file that
  will be uploaded and applied on a secondary device. Applicable for some devices
  types in BYOL licensing mode
  * `secondary_device.#.account_number` - (Required) Billing account number for
  secondary device
  * `secondary_device.#.notifications` - (Required) List of email addresses that
  will receive notifications about secondary device
  * `secondary_device.#.additional_bandwidth` - (Optional) Additional Internet
 bandwidth, in Mbps, for a secondary device
  * `secondary_device.#.vendor_configuration` - (Optional) map of vendor specific
 configuration parameters for a secondary device
  * `secondary_device.#.acl_template_id` - Identifier of an ACL template that will
  be applied on a secondary device
  * `ssh-key` - (Optional) up to one definition of SSH key that will be provisioned
on a secondary device
    * `ssh-key.#.username` - (Required) username associated with given key
    * `ssh-key.#.name` - (Required) name of SSH key as defined in
      `equinix_network_ssh_key` resource

## Attributes Reference

* `uuid` - Device universally unique identifier
* `status` - Device provisioning status
  * INITIALIZING
  * PROVISIONING
  * WAITING_FOR_PRIMARY
  * WAITING_FOR_SECONDARY
  * FAILED
  * PROVISIONED
  * DEPROVISIONING
  * DEPROVISIONED
* `license_status` - Device license registration status
  * APPLYING_LICENSE
  * REGISTERED
  * APPLIED
  * REGISTRATION_FAILED
* `ibx` - Name of Equinix exchange
* `region` - Region in which device metro is located
* `ssh_ip_address` - IP address to use for SSH connectivity with the device
* `ssh_ip_fqdn` - FQDN to use for SSH connectivity with the device
* `redundancy_type` - Indicates if device is primary or secondary
(in HA configuration)
* `redundant_id` - Universally unique identifier for a redundant device
(in HA configuration)
* `interface` - List of device interfaces
  * `interface.#.id` - interface identifier
  * `interface.#.name` - interface name
  * `interface.#.status` -  interface status (AVAILABLE, RESERVED, ASSIGNED)
  * `interface.#.operational_status` - interface operation status (up or down)
  * `interface.#.mac_address` - interface MAC address
  * `interface.#.ip_address` - interface IP address
  * `interface.#.assigned_type` - interface management type (Equinix Managed or empty)
  * `interface.#.type` - interface type

## Timeouts

This resource provides the following [Timeouts configuration](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)
options:

* create - Default is 60 minutes
* update - Default is 10 minutes
