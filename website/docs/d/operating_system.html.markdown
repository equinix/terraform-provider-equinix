---
layout: "packet"
page_title: "Packet: opearating_system"
sidebar_current: "docs-packet-datasource-operating-system"
description: |-
  Get an operationg on a Packet Operating System image
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

```

## Argument Reference

 * `distro` - (Optional) Name of the OS distribution.
 * `name` - (Optional) Name or part of the name of the distribution. Case insensitive.
 * `provisionable_on` - (Optional) Plan name.
 * `version` - (Optional) Version of the distribution

## Attributes Reference

 * `operating_system` - Opearting system of a device.

