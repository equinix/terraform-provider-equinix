---
layout: "equinix"
page_title: "Equinix: equinix_network_ssh_user"
subcategory: ""
description: |-
 Provides Equinix Network Edge SSH user resource
---

# Resource: equinix_network_ssh_user

Resource `equinix_network_ssh_user` allows creation and management of Equinix Network
Edge SSH users.

## Example Usage

```hcl
# Create SSH user with password auth method and associate it with
# two virtual network devices

resource "equinix_network_ssh_user" "john" {
  username = "john"
  password = "secret"
  device_ids = [
    equinix_ne_device.csr1000v-ha.uuid,
    equinix_ne_device.csr1000v-ha.redundant_uuid
  ]
}
```

## Argument Reference

* `username` - (Required) SSH user login name
* `password` - (Required) SSH user password
* `device_ids` - (Required) list of device identifiers to which user will have access

## Attributes Reference

* `uuid` - SSH user unique identifier

## Import

This resource can be imported using an existing ID:

```sh
terraform import equinix_network_ssh_user.example {existing_id}
```
