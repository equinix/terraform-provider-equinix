---
subcategory: "Network Edge"
---

# equinix_network_bgp (Resource)

Resource `equinix_network_bgp` allows creation and management of Equinix Network
Edge BGP peering configurations.

## Example Usage

```hcl
# Create BGP peering configuration on a existing connection
# between network device and service provider

resource "equinix_network_bgp" "test" {
  connection_id      = "54014acf-9730-4b55-a791-459283d05fb1"
  local_ip_address   = "10.1.1.1/30"
  local_asn          = 12345
  remote_ip_address  = "10.1.1.2"
  remote_asn         = 66123
  authentication_key = "secret"
}
```

## Argument Reference

The following arguments are supported:

* `connection_id` - (Required) identifier of a connection established between.
network device and remote service provider that will be used for peering.
* `local_ip_address` - (Required) IP address in CIDR format of a local device.
* `local_asn` - (Required) Local ASN number.
* `remote_ip_address` - (Required) IP address of remote peer.
* `remote_asn` - (Required) Remote ASN number.
* `authentication_key` - (Optional) shared key used for BGP peer authentication.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - BGP peering configuration unique identifier.
* `device_id` - unique identifier of a network device that is a local peer in a given BGP peering
configuration.
* `state` - BGP peer state, one of `Idle`, `Connect`, `Active`, `OpenSent`, `OpenConfirm`,
`Established`.
* `provisioning_status` - BGP peering configuration provisioning status, one of `PROVISIONING`,
`PENDING_UPDATE`, `PROVISIONED`, `FAILED`.

## Import

This resource can be imported using an existing ID:

```sh
terraform import equinix_network_bgp.example {existing_id}
```
