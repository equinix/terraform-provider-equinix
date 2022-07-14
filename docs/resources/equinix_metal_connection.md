---
subcategory: "Metal"
---

# equinix_metal_connection (Resource)

Use this resource to request the creation an Interconnection asset to connect with other parties using [Equinix Fabric - software-defined interconnections](https://metal.equinix.com/developers/docs/networking/fabric/).

~> Equinix Metal connection with service_token_type `a_side` is not generally available and may not be enabled yet for your organization.

## Example Usage

### Shared Connection with a_side token - Redundant Connection from Equinix Metal to a Cloud Service Provider

```hcl
resource "equinix_metal_connection" "example" {
    name               = "tf-metal-to-azure"
    project_id         = local.project_id
    type               = "shared"
    redundancy         = "redundant"
    metro              = "sv"
    speed              = "1000Mbps"
    service_token_type = "a_side"
}

data "equinix_ecx_l2_sellerprofile" "example" {
  name                     = "Azure ExpressRoute"
  organization_global_name = "Microsoft"
}

resource "equinix_ecx_l2_connection" "example" {
  name              = "tf-metal-to-azure"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.example.uuid
  speed             = azurerm_express_route_circuit.example.bandwidth_in_mbps
  speed_unit        = "MB"
  notifications     = ["example@equinix.com"]
  service_token     = equinix_metal_connection.example.service_tokens.0.id
  seller_metro_code = "AM"
  authorization_key = azurerm_express_route_circuit.example.service_key
  named_tag         = "PRIVATE"
  secondary_connection {
    name          = "tf-metal-to-azure"-sec"
    service_token = equinix_metal_connection.example.service_tokens.1.id
  }
}
```

### Shared Connection with z_side token - Non-redundant Connection from your Equinix Fabric Port to Equinix Metal

```hcl
resource "equinix_metal_connection" "example" {
    name               = "tf-port-to-metal"
    project_id         = local.project_id
    type               = "shared"
    redundancy         = "primary"
    metro              = "FR"
    speed              = "200Mbps"
    service_token_type = "z_side"
}

data "equinix_ecx_port" "example" {
  name = "CX-FR5-NL-Dot1q-BO-1G-PRI"
}

resource "equinix_ecx_l2_connection" "example" {
  name                = "tf-port-to-metal"
  zside_service_token = equinix_metal_connection.example.service_tokens.0.id
  speed               = "200"
  speed_unit          = "MB"
  notifications       = ["example@equinix.com"]
  seller_metro_code   = "FR"
  port_uuid           = data.equinix_ecx_port.example.id
  vlan_stag           = 1020
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the connection resource
* `metro` - (Optional) Metro where the connection will be created.
* `facility` - (Optional) Facility where the connection will be created.
* `redundancy` - (Required) Connection redundancy - redundant or primary.
* `type` - (Required) Connection type - dedicated or shared.
* `project_id` - (Optional) ID of the project where the connection is scoped to, must be set for.
* `speed` - (Required) Connection speed - one of 50Mbps, 200Mbps, 500Mbps, 1Gbps, 2Gbps, 5Gbps, 10Gbps.
* `description` - (Optional) Description for the connection resource.
* `mode` - (Optional) Mode for connections in IBX facilities with the dedicated type - standard or tunnel. Default is standard.
* `tags` - (Optional) String list of tags.
* `vlans` - (Optional) Only used with shared connection. Vlans to attach. Pass one vlan for Primary/Single connection and two vlans for Redundant connection.
* `service_token_type` - (Optional) Only used with shared connection. Type of service token to use for the connection, a_side or z_side.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `organization_id` - ID of the organization where the connection is scoped to.
* `status` - Status of the connection resource.
* `ports` - List of connection ports - primary (`ports[0]`) and secondary (`ports[1]`). Schema of
port is described in documentation of the
[equinix_metal_connection datasource](../data-sources/equinix_metal_connection.md).
* `service_tokens` - List of connection service tokens with attributes. Scehma of service_token is described in documentation of the [equinix_metal_connection datasource](../data-sources/equinix_metal_connection.md).
