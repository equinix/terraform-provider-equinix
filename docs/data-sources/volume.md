---
page_title: "Equinix Metal: metal_volume"
subcategory: ""
description: |-
  Provides an Equinix Metal Block Storage Volume Datasource.
---

# metal\_volume

Provides an Equinix Metal Block Storage Volume datasource to allow you to read existing volumes.

## Example Usage

```hcl
# Read a volume by project ID and name
data "metal_volume" "volume1" {
  name       = "terraform-volume-1"
  project_id = local.project_id
}

output "volume_size" {
  value = data.metal_volume.volume1.size
}
```

## Argument Reference

The following arguments are supported:

* `volume_id` ID of volume for lookup
* `name` - Name of volume for lookup
* `project_id` - The ID the parent Equinix Metal project (for lookup by name)

Either `volume_id` or both `project_id` and `name` must be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the volume
* `name` - The name of the volume
* `project_id` - The project id the volume is in
* `size` - The size in GB of the volume
* `plan` - Performance plan the volume is on
* `billing_cycle` - The billing cycle, defaults to hourly
* `facility` - The facility slug the volume resides in
* `state` - The state of the volume
* `locked` - Whether the volume is locked or not
* `device_ids` - UUIDs of devices to which this volume is attached
