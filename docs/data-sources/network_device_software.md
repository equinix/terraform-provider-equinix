---
subcategory: "Network Edge"
---

# equinix_network_device_software (Data Source)

Use this data source to get Equinix Network Edge device software details for a given
device type.

## Example Usage

```hcl
# Retrieve details for CSR1000V device software with latest path of 16.09 version
# that supports IPBASE package
data "equinix_network_device_software" "csrLatest1609" {
  device_type   = "CSR1000V"
  version_regex = "^16.09.+"
  packages      = ["IPBASE"]
  most_recent   = true
}
```

## Argument Reference

* `device_type` - (Required) Code of a device type
* `version_regex` - (Optional) A regex string to apply on returned versions
and filter search results
* `stable` - (Optional) Boolean value to limit query results to stable versions only
* `packages` - (Optional) Limits returned versions to those that are supported
by given software package codes
* `most_recent` - (Optional) Boolean value to indicate that most recent version
should be used *(in case when more than one result is returned)*

## Attributes Reference

* `version` - Version number
* `image_name` - Software image name
* `date` - Version release date
* `status` - Version status
* `release_notes_link` - Link to version release notes
