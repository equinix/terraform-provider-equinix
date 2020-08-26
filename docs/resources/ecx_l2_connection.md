---
layout: "equinix"
page_title: "Equinix: ecx_l2_connection"
sidebar_current: "docs-equinix-datasource-ecx-l2-connection"
description: |-
  Get an ECX L2 connection resource.
---

# equinix_ecx_l2_connection

Resource `equinix_ecx_l2_connection` is used to manage layer 2 connections in Equinix Cloud Exchange (ECX) Fabric.

## Example Usage

### Non-redundant Connection

```hcl
resource "equinix_ecx_l2_connection" "aws_dot1q" {
 name = "tf-single-aws"
 profile_uuid = "2a4f7e27-dff8-4f15-aeda-a11ffe9ccf73"
 speed = 200
 speed_unit = "MB"
 notifications = ["marry@equinix.com", "john@equinix.com"]
 port_uuid = "febc9d80-11e0-4dc8-8eb8-c41b6b378df2"
 vlan_stag = 777
 vlan_ctag = 1000
 seller_region = "us-east-1"
 seller_metro_code = "SV"
 authorization_key = "1234456"
}
```

### Redundant Connection

```hcl
resource "equinix_ecx_l2_connection" "redundant_self" {
  name = "tf-redundant-self"
  profile_uuid = "2a4f7e27-dff8-4f15-aeda-a11ffe9ccf73"
  speed = 50
  speed_unit = "MB"
  notifications = ["john@equinix.com", "marry@equinix.com"]
  port_uuid = "febc9d80-11e0-4dc8-8eb8-c41b6b378df2"
  vlan_stag = 800
  zside_port_uuid = "03a969b5-9cea-486d-ada0-2a4496ed72fb"
  zside_vlan_stag = 1010
  seller_region = "us-east-1"
  seller_metro_code = "SV"
  secondary_connection {
    name = "tf-redundant-self-sec"
    port_uuid = "86872ae5-ca19-452b-8e69-bb1dd5f93bd1"
    vlan_stag = 999
    vlan_ctag = 1000
    zside_port_uuid = "393b2f6e-9c66-4a39-adac-820120555420"
    zside_vlan_stag = 1022
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - _(Required)_ Name of the primary connection - An alpha-numeric 24 characters string which can include only hyphens and underscores ('-' & '\_').
- `profile_uuid` - _(Required)_ Unique identifier of the provider's service profile.
- `speed` - _(Required)_ Speed/Bandwidth to be allocated to the connection.
- `speed_unit` - _(Required)_ Unit of the speed/bandwidth to be allocated to the connection.
- `notifications` - _(Required)_ A list of email addresses that would be notified when there are any updates on this connection.
- `purchase_order_number` - _(Optional)_ Test field to link the purchase order numbers to the connection on Equinix which would be reflected on the invoice.
- `port_uuid` - _(Required when device_uuid is not set)_ Unique identifier of the buyer's port from which the connection would originate.
- `device_uuid` - _(Required when port_uuid is not set)_ Unique identifier of the Network Edge virtual device from which the connection would originate.
- `vlan_stag` - _(Required when port_uuid is set)_ S-Tag/Outer-Tag of the connection - a numeric character ranging from 2 - 4094.
- `vlan_ctag` - _(Optional)_ C-Tag/Inner-Tag of the connection - a numeric character ranging from 2 - 4094.
- `named_tag` - _(Optional)_ The type of peering to set up in case when connecting to Azure Express Route. One of _"Public"_, _"Private"_, _"Microsoft"_, _"Manual"_
- `additional_info` - _(Optional)_ one or more additional information key-value objects
  - `name` - _(Required)_ additional information key
  - `value` - _(Required)_ additional information value
- `zside_port_uuid` - _(Optional)_ Unique identifier of the port on the Z side.
- `zside_vlan_stag` - _(Optional)_ S-Tag/Outer-Tag of the connection on the Z side.
- `zside_vlan_ctag` - _(Optional)_ C-Tag/Inner-Tag of the connection on the Z side.
- `seller_region` - _(Optional)_ The region in which the seller port resides.
- `seller_metro_code` - _(Optional)_ The metro code that denotes the connection’s destination (Z side).
- `authorization_key` - _(Optional)_ Text field based on the service profile you want to connect to.
- `secondary_connection` - _(Optional)_ Definition of secondary connection for redundant connectivity. Most attributes are derived from primary connection, except below:
  - `name` - _(Required)_
  - `port_uuid` - _(Required when device_uuid is not set)_
  - `device_uuid` - _(Required when port_uuid is not set)_
  - `vlan_stag` - _(Required when port_uuid is set)_
  - `vlan_ctag` - _(Optional)_
  - `zside_port_uuid` - _(Optional)_
  - `zside_vlan_stag` - _(Optional)_
  - `zside_vlan_ctag` - _(Optional)_

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

- `uuid` - Unique identifier of the connection
- `status` - Status of the connection
- `redundant_uuid` - Unique identifier of the redundant connection (i.e. secondary connection)

## Update operation behavior

As for now, update of ECXF L2 connection implies removal of old connection (in redundant scenario - both primary and secondary connections), and creation of new one, with required set of attributes.
