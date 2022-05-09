---
subcategory: "Network Edge"
---

# equinix_network_device (Data Source)

Use this data source to get Equinix Network Edge device details.

## Example Usage

```hcl
# Retrieve data for an existing Equinix Network Edge device with UUID "f0b5c553-cdeb-4bc3-95b8-23db9ccfd5ee"
data "equinix_network_device" "by_uuid" {
  uuid = "f0b5c553-cdeb-4bc3-95b8-23db9ccfd5ee"
}

# Retrieve data for an existing Equinix Network Edge device named "Arcus-Gateway-A1"
data "equinix_network_device" "by_name" {
  name = "Arcus-Gateway-A1"
}
```

## Argument Reference

* `uuid` - (Optional) UUID of an existing Equinix Network Edge device
* `name` - (Optional) Name of an existing Equinix Network Edge device
* `valid_status_list` - (Optional) Device states to be considered valid when searching for a device by name

NOTE: Exactly one of either `uuid` or `name` must be specified.

## Attributes Reference

* `uuid` - Device unique identifier
* `status` - Device provisioning status
  * INITIALIZING
  * PROVISIONING
  * PROVISIONED  (**NOTE: By default data source will only return devices in this state.  To include other states see `valid_state_list`**)
  * WAITING_FOR_PRIMARY
  * WAITING_FOR_SECONDARY
  * WAITING_FOR_REPLICA_CLUSTER_NODES 
  * CLUSTER_SETUP_IN_PROGRESS 
  * FAILED
  * DEPROVISIONING
  * DEPROVISIONED
* `valid_status_list` - Comma separated list of device states (from see `status` for full list) to be considered valid. Default is 'PROVISIONED'.  Case insensitive. 
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
