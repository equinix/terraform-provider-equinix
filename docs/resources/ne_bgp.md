---
layout: "equinix"
page_title: "Equinix: ne_bgp"
sidebar_current: "docs-equinix-resource-ne-bgp"

description: |-
 Provides Network Edge BGP peering resource.
---

# Resource: ne_bgp

Resource `equinix_ne_bgp` allows creation and management of Network Edge
BGP peering configurations.

## Example Usage

```hcl
# Create BGP peering configuration on a existing connection
# between network device and service provider

resource "equinix_ne_bgp" "test" {
  connection_uuid    = "54014acf-9730-4b55-a791-459283d05fb1"
  local_ip_address   = "10.1.1.1/30"
  local_asn          = 12345
  remote_ip_address  = "10.1.1.2"
  remote_asn         = 66123
  authentication_key = "secret"
}
```

## Argument Reference

* `connection_uuid` - (Required) identifier of a connection established between
network device and remote service provider
* `local_ip_address` - (Required) IP address in CIDR format of a local device
* `local_asn` - (Required) Local ASN number
* `remote_ip_address` - (Required) IP address of remote peer
* `remote_asn` - (Required) Remote ASN number
* `authentication_key` - (Required) shared key used for BGP peer authentication

## Attributes Reference

* `uuid` - BGP peering configuration universally unique identifier
* `device_uuid` - universally unique identifier of a network device that
forms a connection with a given BGP peering
* `state` - BGP peer state, one of:
  * Idle
  * Connect
  * Active
  * OpenSent
  * OpenConfirm
  * Established
* `provisioning_status` - BGP peering configuration provisioning status, one of:
  * PROVISIONING
  * PROVISIONED
  * FAILED
