---
subcategory: "Network Edge"
---

# equinix_network_device (Resource)

Resource `equinix_network_device` allows creation and management of Equinix Network
Edge virtual network devices.

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
```hcl
# Create self configured PANW cluster with BYOL license

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "panw-cluster" {
  name            = "tf-panw"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "PA-VM"
  self_managed    = true
  byol            = true
  package_code    = "VM100"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  term_length     = 6
  account_number  = data.equinix_network_account.sv.number
  version         = "10.1.3"
  interface_count = 10
  core_count      = 2
  ssh_key {
    username = "test"
    key_name = "test-key"
  }
  acl_template_id = "0bff6e05-f0e7-44cd-804a-25b92b835f8b"
  cluster_details {
    cluster_name    = "tf-panw-cluster"
    node0 {
      vendor_configuration {
        hostname = "panw-node0"
      }
      license_token = "licenseToken"
    }
    node1 {
      vendor_configuration {
        hostname = "panw-node1"
      }
      license_token = "licenseToken"
    }
  }
}
```

## Argument Reference

* `name` - (Required) Device name
* `type_code` - (Required) Device type code
* `metro_code` - (Required) Device location metro code
* `hostname` - (Optional) Device hostname prefix
* `package_code` - (Required) Device software package code
* `version` - (Required) Device software software version
* `core_count` - (Required) Number of CPU cores used by device,
* `term_length` - (Required) Device term length
* `self_managed` - (Optional) Boolean value that determines device management mode:
*self-managed* or *Equinix managed* (default)
* `byol` - (Optional) Boolean value that determines device licensing mode:
*bring your own license* or *subscription* (default)
* `license_token` - (Optional) License Token applicable for some device types
in BYOL licensing mode
* `license_file` - (Optional) Path to the license file that will be uploaded and
applied on a device. Applicable for some devices types in BYOL licensing mode
* `throughput` - (Optional) Device license throughput
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
that will be allocated to the device (in addition to default 15Mbps)
* `interface_count` - (Optional) Number of network interfaces on a device. If not
specified, default number for a given device type will be used
* `wan_interafce_id` - (Optional) Specify the WAN/SSH interface id. If not
  specified, default WAN/SSH interface for a given device type will be used
* `vendor_configuration` - (Optional) Map of vendor specific configuration parameters
for a device (controller1, activationKey, managementType, siteId, systemIpAddress)
* `ssh-key` - (Optional) Definition of SSH key that will be provisioned
on a device (max one key)
* `secondary_device` - (Optional) Definition of secondary device for redundant
device configurations
* `cluster_details` - (Optional) An object that has the cluster details

The `secondary_device` block supports the following arguments:

* `name` - (Required) Secondary device name
* `metro_code` - (Required) Metro location of a secondary device
* `hostname` - (Optional) Secondary device hostname
* `license_token` - (Optional) License Token can be provided
for some device types o the device
* `license_file` - (Optional) Path to the license file that
will be uploaded and applied on a secondary device. Applicable for some devices
types in BYOL licensing mode
* `account_number` - (Required) Billing account number for
secondary device
* `notifications` - (Required) List of email addresses that
will receive notifications about secondary device
* `additional_bandwidth` - (Optional) Additional Internet
bandwidth, in Mbps, for a secondary device
* `vendor_configuration` - (Optional) Map of vendor specific
configuration parameters for a secondary device (controller1, activationKey, managementType, siteId, systemIpAddress)
* `acl_template_id` - Identifier of an ACL template that will
be applied on a secondary device
* `ssh-key` - (Optional) up to one definition of SSH key that will be provisioned
on a secondary device

The `ssh_key` block supports the following arguments:

* `username` - (Required) username associated with given key
* `name` - (Required) reference by name to previously provisioned public SSH key

The `cluster_details` block supports the following arguments:

* `cluster_name` - (Required) The name of the cluster device
* `node0` - (Required) An object that has node0 details
* `node1` - (Required) An object that has node1 details

The `node0` block supports the following arguments:

* `vendor_configuration` - (Optional) An object that has fields relevant to the vendor of the cluster device
* `license_file_id` - (Optional) License file id. This is necessary for Fortinet and Juniper clusters
* `license_token` - (Optional) License token. This is necessary for Palo Alto clusters

The `node1` block supports the following arguments:

* `vendor_configuration` - (Optional) An object that has fields relevant to the vendor of the cluster device
* `license_file_id` - (Optional) License file id. This is necessary for Fortinet and Juniper clusters
* `license_token` - (Optional) License token. This is necessary for Palo Alto clusters

The `vendor_configuration` block supports the following arguments:

* `hostname` - (Optional) Hostname. This is necessary for Palo Alto, Juniper, and Fortinet clusters
* `admin_password` - (Optional) The administrative password of the device. You can use it to log in to the console. This field is not available for all device types
* `controller1` - (Optional) System IP Address. Mandatory for the Fortinet SDWAN cluster device
* `activation_key` - (Optional) Activation key. This is required for Velocloud clusters
* `controller_fqdn` - (Optional) Controller fqdn. This is required for Velocloud clusters
* `root_password` - (Optional) The CLI password of the device. This field is relevant only for the Velocloud SDWAN cluster

## Attributes Reference

* `uuid` - Device unique identifier
* `status` - Device provisioning status
  * INITIALIZING
  * PROVISIONING
  * WAITING_FOR_PRIMARY
  * WAITING_FOR_SECONDARY
  * WAITING_FOR_REPLICA_CLUSTER_NODES 
  * CLUSTER_SETUP_IN_PROGRESS 
  * FAILED
  * PROVISIONED
  * DEPROVISIONING
  * DEPROVISIONED
* `license_status` - Device license registration status
  * APPLYING_LICENSE
  * REGISTERED
  * APPLIED
  * WAITING_FOR_CLUSTER_SETUP
  * REGISTRATION_FAILED
* `license_file_id` - Unique identifier of applied license file
* `ibx` - Device location Equinix Business Exchange name
* `region` - Device location region
* `acl_template_id` - Unique identifier of applied ACL template
* `ssh_ip_address` - IP address of SSH enabled interface on the device
* `ssh_ip_fqdn` - FQDN of SSH enabled interface on the device
* `redundancy_type` - Device redundancy type applicable for HA devices, either
primary or secondary
* `redundant_id` - Unique identifier for a redundant device applicable for HA devices
* `interface` - List of device interfaces
  * `interface.#.id` - interface identifier
  * `interface.#.name` - interface name
  * `interface.#.status` -  interface status (AVAILABLE, RESERVED, ASSIGNED)
  * `interface.#.operational_status` - interface operational status (up or down)
  * `interface.#.mac_address` - interface MAC address
  * `interface.#.ip_address` - interface IP address
  * `interface.#.assigned_type` - interface management type (Equinix Managed or empty)
  * `interface.#.type` - interface type
* `asn` - Autonomous system number
* `zone_code` - Device location zone code
* `cluster_id` - The id of the cluster
* `num_of_nodes` - The number of nodes in the cluster

## Timeouts

This resource provides the following [Timeouts configuration](https://www.terraform.io/language/resources/syntax#operation-timeouts)
options:

* create - Default is 90 minutes
* update - Default is 30 minutes
* delete - Default is 30 minutes

## Import

This resource can be imported using an existing ID:

```sh
terraform import equinix_network_device.example {existing_id}
```
