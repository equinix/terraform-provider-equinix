---
layout: "equinix"
page_title: "Equinix: ne_sshuser"
sidebar_current: "docs-equinix-resource-ne-sshuser"
description: |-
 Provides Network Edge SSH user resource.
---

# Resource: ne_sshuser

Resource `equinix_ne_sshuser` allows creation and management of Network Edge
SSH users.

## Example Usage

```hcl
# Create SSH user with password auth method and associate it with
# two virtual network devices

resource "equinix_ne_sshuser" "john" {
  username = "john"
  password = "secret"
  devices = [
    equinix_ne_device.csr1000v-ha.uuid,
    equinix_ne_device.csr1000v-ha.redundant_uuid
  ]
}
```

## Argument Reference

* `username` - (Required) SSH user login name
* `password` - (Required) SSH user password
* `devices` - (Required) list of device identifiers to which user will have access

## Attributes Reference

* `uuid` - SSH user universally unique identifier
