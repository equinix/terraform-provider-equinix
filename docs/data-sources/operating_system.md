---
page_title: "Equinix Metal: operating_system"
subcategory: ""
description: |-
  Get an Equinix Metal operating system image
---

# metal\_operating\_system

Use this data source to get Equinix Metal Operating System image.

## Example Usage

```hcl
data "metal_operating_system" "example" {
  name             = "Container Linux"
  distro           = "coreos"
  version          = "alpha"
  provisionable_on = "c1.small.x86"
}

resource "metal_device" "server" {
  hostname         = "tf.coreos2"
  plan             = "c1.small.x86"
  facilities       = ["ewr1"]
  operating_system = data.metal_operating_system.example.id
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
