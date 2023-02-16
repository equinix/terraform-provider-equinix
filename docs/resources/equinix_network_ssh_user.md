---
subcategory: "Network Edge"
---

# equinix_network_ssh_user (Resource)

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
    equinix_network_device.csr1000v-ha.uuid,
    equinix_network_device.csr1000v-ha.redundant_uuid
  ]
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) SSH user login name.
* `password` - (Required) SSH user password.
* `device_ids` - (Required) list of device identifiers to which user will have access.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - SSH user unique identifier.

## Import

This resource can be imported using an existing ID:

```sh
terraform import equinix_network_ssh_user.example {existing_id}
```
