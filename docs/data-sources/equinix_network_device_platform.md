---
subcategory: "Network Edge"
---

# equinix_network_device_platform (Data Source)

Use this data source to get Equinix Network Edge device platform configuration details
for a given device type. For further details, check supported
[Network Edge Vendors and Devices](https://docs.equinix.com/en-us/Content/Interconnection/NE/user-guide/NE-vendors-devices.htm).

## Example Usage

```hcl
# Retrieve platform configuration of a large flavor for a CSR100V device type
# Platform has to support IPBASE software package
data "equinix_network_device_platform" "csrLarge" {
  device_type = "CSR1000V"
  flavor      = "large"
  packages    = ["IPBASE"]
}
```

## Argument Reference

The following arguments are supported:

* `device_type` - (Required) Device type code
* `flavor` - (Optional) Device platform flavor that determines number of CPU cores and memory.
Supported values are: `small`, `medium`, `large`, `xlarge`.
* `core_count` - (Optional) Number of CPU cores used to limit platform search results.
* `packages` - (Optional) List of software package codes to limit platform search results.
* `management_types` - (Optional) List of device management types to limit platform search results.
Supported values are: `EQUINIX-CONFIGURED`, `SELF-CONFIGURED`.
* `license_options` - (Optional) List of device licensing options to limit platform search result.
Supported values are: `BYOL` (for Bring Your Own License), `Sub` (for license subscription).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `memory` - The amount of memory provided by device platform.
* `memory_unit` - Unit of memory provider by device platform.
