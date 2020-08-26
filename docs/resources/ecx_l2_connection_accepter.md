---
layout: "equinix"
page_title: "Equinix: equinix_ecx_l2_connection_accepter"
sidebar_current: "docs-equinix-resource-ecx-l2-connection-accepter"
description: |-
  Provides a Resource for Attaching IP Subnets from a Reserved Block to a Device
---

# equinix\_ecx\_l2\_connection\_accepter

Resource to approve hosted Layer 2 connections.

The resource relies on the Equinix Cloud Exchange Fabric API. The parameters and
attributes available map to the fields described at
<https://developer.equinix.com/catalog/buyerv3#operation/performUserActionUsingPATCH>

## Example Usage

```hcl
resource "equinix_ecx_l2_connection_accepter" "accepter" {
  connection_id = "xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx"
  access_key = "AKIAIXKQARIFBC3QJKYQ"
  secret_key = "ARIFW1lWbqNSOqSkCAOXAhep22UGyLJvkDBAIG/6"
}
```

## Argument Reference

* `connection_id`
* `access_key`
* `secret_key`

## Attribute Reference

This resource exports no additional attributes.
