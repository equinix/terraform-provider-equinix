---
page_title: "Equinix: equinix_metal_operating_system"
subcategory: ""
description: |-
  Get an Equinix Metal operating system image
---

# Data Source: equinix\_metal\_operating\_system

Use this data source to get Equinix Metal Operating System image.

## Example Usage

```hcl
data "equinix_metal_operating_system" "example" {
  distro           = "ubuntu"
  version          = "20.04"
  provisionable_on = "c3.medium.x86"
}

resource "equinix_metal_device" "server" {
  hostname         = "tf.ubuntu"
  plan             = "c3.medium.x86"
  facilities       = ["ny5"]
  operating_system = data.equinix_metal_operating_system.example.id
  billing_cycle    = "hourly"
  project_id       = local.project_id
}
```

## Argument Reference

* `distro` - (Optional) Name of the OS distribution.
* `name` - (Optional) Name or part of the name of the distribution. Case insensitive.
* `provisionable_on` - (Optional) Plan name.
* `version` - (Optional) Version of the distribution

## Attributes Reference

* `id` - Operating system slug
* `slug` - Operating system slug (same as `id`)
