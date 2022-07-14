---
subcategory: "Fabric"
---

# equinix_ecx_l2_connection (Resource)

Resource `equinix_ecx_l2_connection` allows creation and management of Equinix Fabric
layer 2 connections.

## Example Usage

### Non-redundant Connection from own Equinix Fabric Port

```hcl
data "equinix_ecx_l2_sellerprofile" "aws" {
  name = "AWS Direct Connect"
}

data "equinix_ecx_port" "sv-qinq-pri" {
  name = "CX-SV5-NL-Dot1q-BO-10G-PRI"
}

resource "equinix_ecx_l2_connection" "port-2-aws" {
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

### Redundant Connection from own Equinix Fabric Ports

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

resource "equinix_ecx_l2_connection" "ports-2-azure" {
  name              = "tf-azure-pri"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.azure.id
  speed             = 50
  speed_unit        = "MB"
  notifications     = ["john@equinix.com", "marry@equinix.com"]
  port_uuid         = data.equinix_ecx_port.sv-qinq-pri.id
  vlan_stag         = 1482
  vlan_ctag         = 2512
  seller_metro_code = "SV"
  named_tag         = "PRIVATE"
  authorization_key = "c4dff8e8-b52f-4b34-b0d4-c4588f7338f3
  secondary_connection {
    name      = "tf-azure-sec"
    port_uuid = data.equinix_ecx_port.sv-qinq-sec.id
    vlan_stag = 1904
    vlan_ctag = 1631
  }
}
```

### Non-redundant Connection from Network Edge device

```hcl
data "equinix_ecx_l2_sellerprofile" "gcp-1" {
  name = "Google Cloud Partner Interconnect Zone 1"
}

resource "equinix_ecx_l2_connection" "router-to-gcp" {
  name                = "tf-gcp-pri"
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

### Non-redundant Connection from an Equinix customer port using A-Side Service token

```hcl
data "equinix_ecx_l2_sellerprofile" "gcp" {
  name = "Google Cloud Partner Interconnect Zone 1"
}

resource "equinix_ecx_l2_connection" "token-to-gcp" {
  name                = "tf-gcp-pri"
  profile_uuid        = data.equinix_ecx_l2_sellerprofile.gcp-1.id
  service_token       = "e9c22453-d3a7-4d5d-9112-d50173531392"
  speed               = 100
  speed_unit          = "MB"
  notifications       = ["john@equinix.com", "marry@equinix.com"]
  seller_metro_code   = "SV"
  seller_region       = "us-west1"
  authorization_key   = "4d335adc-00fd-4a41-c9f3-782ca31ab3f7/us-west1/1"
}
```

### Non-redundant Connection from own Equinix Fabric Port to an Equinix customer port using Z-Side Service token

```hcl
data "equinix_ecx_port" "sv-qinq-pri" {
  name = "CX-SV5-NL-Dot1q-BO-10G-PRI"
}

resource "equinix_ecx_l2_connection" "port-to-token" {
  name                = "tf-port-token"
  zside_service_token = "e9c22453-d3a7-4d5d-9112-d50173531392"
  speed               = 200
  speed_unit          = "MB"
  notifications       = ["john@equinix.com", "marry@equinix.com"]
  seller_metro_code   = "FR"
  port_uuid           = data.equinix_ecx_port.sv-qinq-pri.id
  vlan_stag           = 1000
}
```

-> **NOTE:** See [Equinix Fabric connecting to the cloud](../guides/equinix_fabric_cloud_providers.md)
guide for more details on how to connect to a CSP.

## Argument Reference

The following arguments are supported:

* `name` - (Required) Connection name. An alpha-numeric 24 characters string which can include only
hyphens and underscores
* `profile_uuid` - (Required) Unique identifier of the service provider's profile.
* `speed` - (Required) Speed/Bandwidth to be allocated to the connection.
* `speed_unit` - (Required) Unit of the speed/bandwidth to be allocated to the connection.
* `notifications` - (Required) A list of email addresses used for sending connection update
notifications.
* `purchase_order_number` - (Optional) Connection's purchase order number to reflect on the invoice
* `port_uuid` - (Required when `device_uuid` or `service_token` are not set) Unique identifier of
the Equinix Fabric Port from which the connection would originate.
* `device_uuid` - (Required when `port_uuid` or `service_token` are not set) Unique identifier of
the Network Edge virtual device from which the connection would originate.
* `device_interface_id` - (Optional) Applicable with `device_uuid`, identifier of network interface
on a given device, used for a connection. If not specified then first available interface will be
selected.
* `service_token`- (Required when `port_uuid` or `device_uuid` are not set) - A-side
service tokens authorize you to create a connection from a customer port, which created the token
for you, to a service profile or your own port.
More details in [A-Side Fabric Service Tokens](https://docs.equinix.com/en-us/Content/Interconnection/Fabric/service%20tokens/Fabric-Service-Tokens.htm#:~:text=the%20service%20token.-,A%2DSide%20Service%20Tokens,-If%20you%20want).
* `zside_service_token`- (Required when `profile_uuid` or `zside_port_uuid` are not set) - Z-side
service tokens authorize you to create a connection from your port or virtual device to a customer
port which created the token for you. `zside_service_token` cannot be used with `secondary_connection`.
More details in [Z-Side Fabric Service Tokens](https://docs.equinix.com/en-us/Content/Interconnection/Fabric/service%20tokens/Fabric-Service-Tokens.htm#:~:text=requirements%20per%20provider.-,Z%2DSide%20Service%20Tokens,-If%20you%20want).

-> **NOTE:** Service tokens can't be reused. To recreate a resource or to create a new one for
another connection even from same interconnection asset, you will need to request another token
from your service provider.

* `vlan_stag` - (Required when port_uuid is set) S-Tag/Outer-Tag of the connection - a numeric
character ranging from 2 - 4094.
* `vlan_ctag` - (Optional) C-Tag/Inner-Tag of the connection - a numeric character ranging from 2
\- 4094.
* `named_tag` - (Optional) The type of peering to set up when connecting to Azure Express Route.
Valid values: `PRIVATE`, `MICROSOFT`, `MANUAL`\*, `PUBLIC`\*.

~> **NOTE:** _"MANUAL"_ peering is deprecated. Use _"PRIVATE"_ or _"MICROSOFT"_ instead. It was
used in cases where `zside_vlan_ctag` was needed. `zside_vlan_ctag` can still be defined without
need to specify `MANUAL` which does not actually correspond to any type of peering supported in
Azure. Check [how to create a connection to Microsoft Azure ExpressRoute with Equinix Fabric](https://docs.equinix.com/en-us/Content/Interconnection/Fabric/connections/Fabric-ms-azure.htm#:~:text=a%20peering%20type.-,(Optional)%C2%A0,-Enter%20a%20vlan)
for more details.

~> **NOTE:** _"PUBLIC"_ peering is deprecated. Use _"MICROSOFT"_ instead. More details in
[Microsoft public peering](https://docs.microsoft.com/en-us/azure/expressroute/about-public-peering)
docs.

* `additional_info` - (Optional) one or more additional information key-value objects
  * `name` - (Required) additional information key
  * `value` - (Required) additional information value
* `zside_port_uuid` - (Optional) Unique identifier of the port on the remote/destination side
(z-side). Allows you to connect between your own ports or virtual devices across your company's
Equinix Fabric deployment, with no need for a private service profile.
* `zside_vlan_stag` - (Optional) S-Tag/Outer-Tag of the connection on the remote/destination
side (z-side) - a numeric character ranging from 2 - 4094.
* `zside_vlan_ctag` - (Optional) C-Tag/Inner-Tag of the connection on the remote/destination
side (z-side) - a numeric character ranging from 2 - 4094.
`secondary_connection` is defined it will internally use same `zside_vlan_ctag` for the secondary
connection.
* `seller_region` - (Optional) The region in which the seller port resides.
* `seller_metro_code` - (Optional) The metro code that denotes the connection’s remote/destination
side (z-side).
* `authorization_key` - (Optional) Unique identifier authorizing Equinix to provision a connection
towards a cloud service provider. At Equinix, an `Authorization Key` is a generic term and is NOT
encrypted on Equinix Fabric. Cloud Service Providers might use a different name to refer to this
key such as `Service Key` or `Authentication Key`. Value depends on a provider service profile,
more information on [Equinix Fabric how to guide](https://developer.equinix.com/docs/ecx-how-to-guide).
* `secondary_connection` - (Optional) Definition of secondary connection for redundant, HA
connectivity. See [Secondary Connection](#secondary-connection) below for more details.

### Secondary Connection

-> **NOTE:** Some service provider do not directly support redundant connections in their service
profiles. However, some of them offer active/active (BGP multipath) or active/passive (failover)
configurations in their platforms and you still achieve that highly resilient network
connections by creating an `equinix_ecx_l2_connection` resource for each connection instead of
defining a `secondary_connection` block.

The `secondary_connection` block supports the following arguments:

* `name` - (Required) secondary connection name
* `speed` - (Optional) Speed/Bandwidth to be allocated to the secondary connection. If not
specified primary `speed` will be used.
* `speed_unit` - (Optional) Unit of the speed/bandwidth to be allocated to the secondary
connection. If not specified primary `speed_unit` will be used.
* `port_uuid` - (Optional) Applicable with primary `port_uuid`. Identifier of the Equinix Fabric Port from
which the secondary connection would originate. If not specified primary `port_uuid` will be used.
* `device_uuid` - (Optional) Applicable with primary `device_uuid`. Identifier of the Network Edge
virtual device from which the secondary connection would originate. If not specified primary
`device_uuid` will be used.
* `device_interface_id` - (Optional) Applicable with `device_uuid`, identifier of network interface
on a given device. If not specified then first available interface will be selected.
* `service_token`- (Optional) Required with primary `service_token`. Unique Equinix Fabric key
given by a provider that grants you authorization to enable connectivity from an Equinix Fabric Port or
virtual device. Each connection (primary and secondary) requires a separate token.
More details in [Fabric Service Tokens](https://docs.equinix.com/en-us/Content/Interconnection/Fabric/service%20tokens/Fabric-Service-Tokens.htm).
* `vlan_stag` - (Required when `port_uuid` is set) S-Tag/Outer-Tag of the secondary connection, a
numeric character ranging from 2 - 4094.
* `vlan_ctag` - (Optional) Applicable with `port_uuid`. C-Tag/Inner-Tag of the secondary
connection, a numeric character ranging from 2 - 4094.
* `seller_metro_code` - (Optional) The metro code that denotes the secondary connection’s
destination (Z side). .
* `seller_region` - (Optional) The region in which the seller port resides. If not specified
primary `seller_region` will be used.
* `authorization_key` - (Optional) Unique identifier authorizing Equinix to provision a connection
towards a cloud service provider. If not specified primary `authorization_key` will be used. However,
some service providers may require different keys for each connection. More information on
[Equinix Fabric how to guide](https://developer.equinix.com/docs/ecx-how-to-guide).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - Unique identifier of the connection.
* `status` - Connection provisioning status on Equinix Fabric side.
* `provider_status` - Connection provisioning status on service provider's side.
* `redundant_uuid` - Unique identifier of the redundant connection, applicable for HA connections.
* `redundancy_type` - Connection redundancy type, applicable for HA connections. Valid values are
`PRIMARY`, `SECONDARY`.
* `redundancy_group` - Unique identifier of group containing a primary and secondary connection.
* `zside_port_uuid` - When not provided as an argument, it is identifier of the z-side port,
assigned by the Fabric.
* `zside_vlan_stag` - When not provided as an argument, it is S-Tag/Outer-Tag of the connection on
the Z side, assigned by the Fabric.
* `actions` - One or more pending actions to complete connection provisioning.
* `secondary_connection`:
  * `zside_port_uuid`
  * `zside_vlan_stag`
  * `zside_vlan_ctag`
  * `redundancy_type`
  * `redundancy_group`

## Update operation behavior

Update of most arguments will force replacement of a connection (including related redundant
connection in HA setup).

Following arguments can be updated. **NOTE** that Equinix Fabric may still forbid updates depending
on current connection state, used service provider or number of updates requested during the day.

* `name`
* `speed` and `speed_unit`

## Timeouts

This resource provides the following [Timeouts configuration](https://www.terraform.io/language/resources/syntax#operation-timeouts)
options:

* create - Default is 5 minutes
* delete - Default is 5 minutes

## Import

Equinix L2 connections can be imported using an existing `id`:

```sh
existing_connection_id='example-uuid-1'
terraform import equinix_ecx_l2_connection.example ${existing_connection_id}
```

**Please Note** that to import a redundant connection you must concatenate `id` of both connections
(primary and secondary) into a single string separated by `:`, e.g.,

```sh
existing_primary_connection_id='example-uuid-1'
existing_secondary_connection_id='example-uuid-2'
terraform import equinix_ecx_l2_connection.example ${existing_primary_connection_id}:${existing_secondary_connection_id}
```
