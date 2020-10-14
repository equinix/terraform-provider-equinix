---
layout: "equinix"
page_title: "Equinix: equinix_network_device_platform"
subcategory: ""
description: |-
 Get information on Network Edge device platform configuration
---

# Data Source: equinix_network_device_platform

Use this data source to get Network Edge device platform configuration details
for a given device type.

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

* `device_type` - (Required) Device type code
* `flavor` - (Optional) Device platform flavor that determines number of CPU cores
and memory. Supported values:
  * small
  * medium
  * large
  * xlarge
* `core_count` - (Optional) Limits platforms to those that provide given number
of CPU cores
* `packages` - (Optional) Limits platforms to those that support provided software
package codes
* `management_types` - (Optional) Limits platforms to those that support given
device management types. Supported values:
  * EQUINIX-CONFIGURED
  * SELF-CONFIGURED
* `license_options` - (Optional) Limits platforms to those that support given
licensing options. Supported values:
  * BYOL (for Bring Your Own License)
  * Sub (for license subscription)

## Attributes Reference

* `memory` - The amount of memory
* `memory_unit` - The unit of memory
