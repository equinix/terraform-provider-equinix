---
layout: "packet"
page_title: "Packet: operating_system"
sidebar_current: "docs-packet-datasource-operating-system"
description: |-
  Get a Packet operating system image
---

# packet\_operating\_system

Use this data source to get Packet Operating System image.

## Example Usage

```hcl
data "packet_operating_system" "example" {
  name             = "Container Linux"
  distro           = "coreos"
  version          = "alpha"
  provisionable_on = "baremetal_1"
}

resource "packet_device" "server" {
  hostname         = "tf.coreos2"
  plan             = "baremetal_1"
  facility         = "ewr1"
  operating_system = "${data.packet_operating_system.example.id}"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.cool_project.id}"
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

