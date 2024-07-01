---
subcategory: "Network Edge"
---

# equinix_network_device_link (Resource)

Resource `equinix_network_device_link` allows creation and management of Equinix Network Edge virtual network device links.

## Example Usage

```terraform
# Example of device link with HA device pair
# where each device is in different metro
resource "equinix_network_device_link" "test" {
  name   = "test-link"
  subnet = "192.168.40.64/27"
  project_id  = "a86d7112-d740-4758-9c9c-31e66373746b"
  device {
    id           = equinix_network_device.test.uuid
    asn          = equinix_network_device.test.asn > 0 ? equinix_network_device.test.asn : 22111
    interface_id = 6
  }
  device {
    id           = equinix_network_device.test.secondary_device[0].uuid
    asn          = equinix_network_device.test.secondary_device[0].asn > 0 ? equinix_network_device.test.secondary_device[0].asn : 22333
    interface_id = 7
  }
  link {
    account_number  = equinix_network_device.test.account_number
    src_metro_code  = equinix_network_device.test.metro_code
    dst_metro_code  = equinix_network_device.test.secondary_device[0].metro_code
    throughput      = "50"
    throughput_unit = "Mbps"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) device link name.
* `subnet` - (Optional) device link subnet in CIDR format. Not required for link between self configured devices.
* `device` - (Required) definition of one or more devices belonging to the device link. See [Device](#device) section below for more details.
* `link` - (Deprecated) definition of one or more, inter metro, connections belonging to the device link. See [Link](#link) section below for more details.
* `metro_link` - (Optional) definition of one or more, inter metro, connections belonging to the device link. See [Metro Link](#Metro_Link) section below for more details.
* `redundancy_type` - (Optional) Whether the connection should be created through Fabric's primary or secondary port. Supported values: `PRIMARY` (Default), `SECONDARY`, `HYBRID`
* `project_id` - (Optional) Unique Identifier for the project resource where the device link is scoped to.If you leave it out, the device link will be created under the default project id of your organization.

### Device

The `device` block supports the following arguments:

* `id` - (Required) Device identifier.
* `asn` - (Optional) Device ASN number. Not required for self configured devices.
* `interface_id` - (Optional) Device network interface identifier to use for device link connection.

### Link

The `link` block supports the following arguments:

* `account_number` - (Required) billing account number to be used for connection charges
* `throughput` - (Required) connection throughput.
* `throughput_unit` - (Required) connection throughput unit (Mbps or Gbps).
* `src_metro_code` - (Required) connection source metro code.
* `dst_metro_code` - (Required) connection destination metro code.
* `src_zone_code` - (Deprecated) connection source zone code is not required.
* `dst_zone_code` - (Deprecated) connection destination zone code is not required.

### Metro_Link

The `Metro link` block supports the following arguments:

* `account_number` - (Required) billing account number to be used for connection charges
* `throughput` - (Required) connection throughput.
* `throughput_unit` - (Required) connection throughput unit (Mbps or Gbps).
* `metro_code` - (Required) connection metro code.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - Device link unique identifier.
* `status` - Device link provisioning status. One of `PROVISIONING`, `PROVISIONED`, `DEPROVISIONING`, `DEPROVISIONED`, `FAILED`.

The `device` block attributes:

* `ip_address` - IP address from device link subnet that was assigned to the device
* `status` - device link provisioning status on a given device. One of `PROVISIONING`, `PROVISIONED`, `DEPROVISIONING`, `DEPROVISIONED`, `FAILED`.

## Timeouts

This resource provides the following [Timeouts configuration](https://www.terraform.io/language/resources/syntax#operation-timeouts) options:

* create - Default is 10 minutes
* update - Default is 10 minutes
* delete - Default is 10 minutes

## Import

This resource can be imported using an existing ID:

```sh
terraform import equinix_network_device_link.example {existing_id}
```
