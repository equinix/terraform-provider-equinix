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
  connection_id = equinix_ecx_l2_connection.awsConn.id
  access_key    = "AKIAIXKQARIFBC3QJKYQ"
  secret_key    = "ARIFW1lWbqNSOqSkCAOXAhep22UGyLJvkDBAIG/6"
}
```

## Argument Reference

* `connection_id`
* `access_key`
* `secret_key`

## Attribute Reference

* `aws_connection_id` - the ID of a hosted Direct Connect connection on AWS side,
applicable for accepter resource with connections to AWS only
