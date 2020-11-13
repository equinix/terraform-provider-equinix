---
layout: "equinix"
page_title: "Equinix: equinix_ecx_l2_connection_accepter"
subcategory: ""
description: |-
  Provides ECX Fabric Layer 2 connection accepter resource
---

# Resource: equinix_ecx_l2_connection_accepter

Resource `equinix_ecx_l2_connection_accepter` is used to accept layer2 connection
on provider side.

Resource leverages ECX Fabric integration with service providers.
Currently supported providers are:

* `AWS` (AWS Direct Connect)

## Example Usage

```hcl
resource "equinix_ecx_l2_connection_accepter" "accepter" {
  connection_id = equinix_ecx_l2_connection.awsConn.id
  access_key    = "AKIAIXKQARIFBC3QJKYQ"
  secret_key    = "ARIFW1lWbqNSOqSkCAOXAhep22UGyLJvkDBAIG/6"
}
```

## Argument Reference

* `connection_id` - (Required) Identifier of Layer 2 connection that will be accepted
* `access_key` - (Required) Access Key used to accept connection on provider side
* `secret_key` - (Required) Secret Key used to accept connection on provider side

## Attribute Reference

* `aws_connection_id` - the ID of a hosted Direct Connect connection on AWS side,
applicable for accepter resource with connections to AWS only
