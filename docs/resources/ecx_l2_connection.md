---
layout: "equinix"
page_title: "Equinix: equinix_ecx_l2_connection"
subcategory: ""
description: |-
  Provides ECX Fabric Layer 2 connection resource
---

# Resource: equinix_ecx_l2_connection

Resource `equinix_ecx_l2_connection` is used to manage layer 2 connections in
Equinix Cloud Exchange (ECX) Fabric.

## Example Usage

### Non-redundant Connection

```hcl
data "equinix_ecx_l2_sellerprofile" "aws" {
  name = "AWS Direct Connect"
}

data "equinix_ecx_port" "sv-qinq-pri" {
  name = "CX-SV5-NL-Dot1q-BO-10G-PRI"
}

resource "equinix_ecx_l2_connection" "aws" {
  name              = "tf-aws"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.aws.id
  speed             = 200
  speed_unit        = "MB"
  notifications     = ["marry@equinix.com", "john@equinix.com"]
  port_uuid         = data.equinix_ecx_port.sv-qinq-pri.id
  vlan_stag         = 777
  vlan_ctag         = 1000
  seller_region     = "us-west-1"
  seller_metro_code = "SV"
  authorization_key = "345742915919"
}
```

### Redundant Connection

```hcl
data "equinix_ecx_l2_sellerprofile" "azure" {
  name = "Azure Express Route"
}

data "equinix_ecx_port" "sv-qinq-pri" {
  name = "CX-SV5-NL-Dot1q-BO-10G-PRI"
}

data "equinix_ecx_port" "sv-qinq-sec" {
  name = "CX-SV1-NL-Dot1q-BO-10G-SEC"
}

resource "equinix_ecx_l2_connection" "azure" {
  name              = "tf-azure-pri"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.azure.id
  speed             = 50
  speed_unit        = "MB"
  notifications     = ["john@equinix.com", "marry@equinix.com"]
  port_uuid         = data.equinix_ecx_port.sv-qinq-pri.id
  vlan_stag         = 1482
  vlan_ctag         = 2512
  seller_metro_code = "SV"
  named_tag         = "Public"
  authorization_key = "c4dff8e8-b52f-4b34-b0d4-c4588f7338f3
  secondary_connection {
    name      = "tf-azure-sec"
    port_uuid = data.equinix_ecx_port.sv-qinq-sec.id
    vlan_stag = 1904
    vlan_ctag = 1631
  }
}
```

### Connection from Network Edge device

```hcl
data "equinix_ecx_l2_sellerprofile" "gcp-1" {
  name = "Google Cloud Partner Interconnect Zone 1"
}

resource "equinix_ecx_l2_connection" "router-gcp" {
  name                = "tf-azure-pri"
  profile_uuid        = data.equinix_ecx_l2_sellerprofile.gcp-1.id
  device_uuid         = equinix_network_device.myrouter.id
  device_interface_id = 5
  speed               = 100
  speed_unit          = "MB"
  notifications       = ["john@equinix.com", "marry@equinix.com"]
  seller_metro_code   = "SV"
  seller_region       = "us-west1"
  authorization_key   = "4d335adc-00fd-4a41-c9f3-782ca31ab3f7/us-west1/1"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the primary connection - An alpha-numeric 24 characters
string which can include only hyphens and underscores ('-' & '\_').
- `profile_uuid` - (Required) Unique identifier of the provider's service profile.
- `speed` - (Required) Speed/Bandwidth to be allocated to the connection.
- `speed_unit` - (Required) Unit of the speed/bandwidth to be allocated
to the connection.
- `notifications` - (Required) A list of email addresses that would be notified
when there are any updates on this connection.
- `purchase_order_number` - (Optional) Test field to link the purchase order
numbers to the connection on Equinix which would be reflected on the invoice.
- `port_uuid` - (Required when device_uuid is not set) Unique identifier of
the buyer's port from which the connection would originate.
- `device_uuid` - (Required when port_uuid is not set) Unique identifier of
the Network Edge virtual device from which the connection would originate.
- `device_interface_id` - (Optional) Applicable with `device_uuid`, identifier of
 network interface on a given device. If not specified then first available interface
 will be selected,
- `vlan_stag` - (Required when port_uuid is set) S-Tag/Outer-Tag of the connection
\- a numeric character ranging from 2 - 4094.
- `vlan_ctag` - (Optional) C-Tag/Inner-Tag of the connection - a numeric
character ranging from 2 - 4094.
- `named_tag` - (Optional) The type of peering to set up in case when connecting
to Azure Express Route. One of _"Public"_, _"Private"_, _"Microsoft"_, _"Manual"_
- `additional_info` - (Optional) one or more additional information key-value objects
  - `name` - (Required) additional information key
  - `value` - (Required) additional information value
- `zside_port_uuid` - (Optional) Unique identifier of the port on the Z side.
- `zside_vlan_stag` - (Optional) S-Tag/Outer-Tag of the connection on the Z side.
- `zside_vlan_ctag` - (Optional) C-Tag/Inner-Tag of the connection on the Z side.
- `seller_region` - (Optional) The region in which the seller port resides.
- `seller_metro_code` - (Optional) The metro code that denotes the connection’s
destination (Z side).
- `authorization_key` - (Optional) Text field based on the service profile
you want to connect to.
- `secondary_connection` - (Optional) Definition of secondary connection for
 redundant connectivity. Most attributes are derived from primary connection,
 except below:
  - `name` - (Required)
  - `port_uuid` - (Required when device_uuid is not set)
  - `device_uuid` - (Required when port_uuid is not set)
  - `vlan_stag` - (Required when port_uuid is set)
  - `vlan_ctag` - (Optional, can be set with port_uuid)

## Attributes Reference

In addition to the arguments listed above, the following computed attributes
are exported:

- `uuid` - Unique identifier of the connection
- `status` - Status of the connection on ECXF side
- `provider_status` - status of the connection at the service provider's end
- `redundant_uuid` - Unique identifier of the redundant connection
(i.e. secondary connection)
- `redundancy_type` - type of connection, either primary or secondary
- `zside_port_uuid` - when not provided as an argument, it is identifier of the
z-side port, assigned by the Fabric
- `zside_vlan_stag` - when not provided as an argument, it is S-Tag/Outer-Tag of
 the connection on the Z side, assigned by the Fabric
- `zside_vlan_ctag` - when not provided as an argument, it is C-Tag/Inner-Tag of
 the connection on the Z side, assigned by the Fabric
- `secondary`
  - `zside_port_uuid` - as for primary connection
  - `zside_vlan_stag` - as for primary connection
  - `zside_vlan_ctag` - as for primary connection

## Update operation behavior

Update of most arguments will force replacement of a connection (including related
redundant connection in HA setup).

Following arguments can be updated. **NOTE** that ECXF may still forbid updates depending
on current connection state, used service provider or number of updates requested
during the day.

- `name`
- `speed` and `speed_unit`
