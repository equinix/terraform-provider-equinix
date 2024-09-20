---
subcategory: "Metal"
---

# equinix_metal_connection (Resource)

Use this resource to request the creation an Interconnection asset to connect with other parties using [Equinix Fabric - software-defined interconnections](https://metal.equinix.com/developers/docs/networking/fabric/).

## Example Usage

### Fabric Billed Shared Virtual Connection - Non-redundant connection from your own Equinix Fabric Port to Equinix Metal

```terraform
resource "equinix_metal_vlan" "example" {
  project_id = "<metal_project_id>"
  metro      = "FR"
}

resource "equinix_metal_connection" "example" {
  name               = "tf-metal-from-port"
  project_id         = "<metal_project_id>"
  type               = "shared"
  redundancy         = "primary"
  metro              = "FR"
  speed              = "200Mbps"
  service_token_type = "z_side"
  contact_email      = "username@example.com"
  vlans              = [equinix_metal_vlan.example.vxlan]
}

data "equinix_fabric_ports" "a_side" {
  filters {
    name = "<name_of_port||port_prefix>"
  }
}

resource "equinix_fabric_connection" "example" {
  name      = "tf-metal-from-port"
  type      = "EVPL_VC"
  bandwidth = "200"
  notifications {
    type   = "ALL"
    emails = ["username@example.com"]
  }
  order { purchase_order_number = "1-323292" }
  project { project_id = "<fabric_project_id>" }
  a_side {
    access_point {
      type = "COLO"
      port {
        uuid = data.equinix_fabric_ports.a_side.data.0.uuid
      }
      link_protocol {
        type     = "DOT1Q"
        vlan_tag = 1234
      }
    }
  }
  z_side {
    service_token {
      uuid = equinix_metal_connection.example.service_tokens.0.id
    }
  }
}
```

-> NOTE: There is an [Equinix Fabric L2 Connection To Equinix Metal Terraform module](https://registry.terraform.io/modules/equinix-labs/fabric-connection-metal/equinix/latest) available with full-fledged examples of connections from Fabric Ports, Network Edge Devices or Service Tokens. Check out the [example for shared connection with Z-side Service Token](https://registry.terraform.io/modules/equinix-labs/fabric-connection-metal/equinix/0.2.0/examples/fabric-port-connection-with-zside-token).

### Fabric Billed Shared Virtual Connection - Non-redundant connection from your own Network Edge device to Equinix Metal

```terraform
resource "equinix_metal_vrf" "example" {
  name       = "tf-metal-from-ne"
  metro      = "FR"
  local_asn  = "65001"
  ip_ranges  = ["10.99.1.0/24"]
  project_id = equinix_metal_project.test.id
}

resource "equinix_metal_connection" "example" {
  name               = "tf-metal-from-ne"
  project_id         = "<metal_project_id>"
  type               = "shared"
  redundancy         = "primary"
  metro              = "FR"
  speed              = "200Mbps"
  service_token_type = "z_side"
  contact_email      = "username@example.com"
  vrfs               = [equinix_metal_vrf.example.vxlan]
}

resource "equinix_fabric_connection" "example" {
  name      = "tf-metal-from-ne"
  type      = "EVPL_VC"
  bandwidth = "200"
  notifications {
    type   = "ALL"
    emails = ["username@example.com"]
  }
  order { purchase_order_number = "1-323292" }
  project { project_id = "<fabric_project_id>" }
  a_side {
    access_point {
      type = "VD"
      virtual_device {
        type = "EDGE"
        uuid = equinix_network_device.example.id
      }
    }
  }
  z_side {
    service_token {
      uuid = equinix_metal_connection.example.service_tokens.0.id
    }
  }
}
```

### Fabric Billed Shared Virtual Connection- Non-redundant connection from Equinix Fabric Cloud Router to Equinix Metal

```terraform
resource "equinix_metal_vlan" "example1" {
  project_id = "<metal_project_id>"
  metro      = "SV"
}

resource "equinix_metal_connection" "example" {
  name          = "tf-metal-from-fcr"
  project_id    = "<metal_project_id>"
  metro         = "SV"
  redundancy    = "primary"
  type          = "shared_port_vlan"
  contact_email = "username@example.com"
  speed         = "200Mbps"
  vlans         = [equinix_metal_vlan.example1.vxlan]
}

resource "equinix_fabric_connection" "example" {
  name      = "tf-metal-from-fcr"
  type      = "IP_VC"
  bandwidth = "200"
  notifications {
    type   = "ALL"
    emails = ["username@example.com"]
  }
  project { project_id = "<fabric_project_id>" }
  a_side {
    access_point {
      type = "CLOUD_ROUTER"
      router {
        uuid = equinix_fabric_cloud_router.example.id
      }
    }
  }
  z_side {
    access_point {
      type               = "METAL_NETWORK"
      authentication_key = equinix_metal_connection.example.authorization_code
    }
  }
}
```

### Metal Billed Shared Virtual Connection - Redundant connection from Equinix Metal to a Cloud Service Provider

```terraform
resource "equinix_metal_connection" "example" {
  name               = "tf-metal-2-azure"
  project_id         = "<metal_project_id>"
  type               = "shared"
  redundancy         = "redundant"
  metro              = "SV"
  speed              = "1Gbps"
  service_token_type = "a_side"
  contact_email      = "username@example.com"
}

data "equinix_fabric_service_profiles" "zside" {
  filter {
    property = "/name"
    operator = "="
    values   = ["Azure ExpressRoute"]
  }
}

resource "equinix_fabric_connection" "example_primary" {
  name      = "tf-metal-2-azure-pri"
  type      = "EVPL_VC"
  bandwidth = azurerm_express_route_circuit.example.bandwidth_in_mbps
  redundancy { priority = "PRIMARY" }
  notifications {
    type   = "ALL"
    emails = ["username@example.com"]
  }
  project { project_id = "<fabric_project_id>" }
  a_side {
    service_token {
      uuid = equinix_metal_connection.example.service_tokens.0.id
    }
  }
  z_side {
    access_point {
      type               = "SP"
      authentication_key = azurerm_express_route_circuit.example.service_key
      profile {
        type = "L2_PROFILE"
        uuid = data.equinix_fabric_service_profiles.zside.id
      }
      location {
        metro_code = "SV"
      }
    }
  }
}

resource "equinix_fabric_connection" "example_secondary" {
  name      = "tf-metal-2-azure-sec"
  type      = "EVPL_VC"
  bandwidth = azurerm_express_route_circuit.example.bandwidth_in_mbps
  redundancy {
    priority = "SECONDARY"
    group    = one(equinix_fabric_connection.example_primary.redundancy).group
  }
  notifications {
    type   = "ALL"
    emails = ["username@example.com"]
  }
  project { project_id = "<fabric_project_id>" }
  a_side {
    service_token {
      uuid = equinix_metal_connection.example.service_tokens.1.id
    }
  }
  z_side {
    access_point {
      type               = "SP"
      authentication_key = azurerm_express_route_circuit.example.service_key
      profile {
        type = "L2_PROFILE"
        uuid = data.equinix_fabric_service_profiles.zside.id
      }
      location {
        metro_code = "SV"
      }
    }
  }
}
```

-> NOTE: There are multiple [Equinix Fabric L2 Connection Terraform modules](https://registry.terraform.io/search/modules?namespace=equinix-labs&q=fabric-connection) available with full-fledged examples of connections from Fabric Ports, Network Edge Devices or Service Token to most popular Cloud Service Providers. Check out the examples for Equinix Metal shared connection with A-side Service Token included in each of them: [AWS](https://registry.terraform.io/modules/equinix-labs/fabric-connection-aws/equinix/latest/examples/service-token-metal-to-aws-connection), [Azure](https://registry.terraform.io/modules/equinix-labs/fabric-connection-azure/equinix/latest/examples/service-token-metal-to-azure-connection), [Google Cloud](https://registry.terraform.io/modules/equinix-labs/fabric-connection-gcp/equinix/latest/examples/service-token-metal-to-gcp-connection), [IBM Cloud](https://registry.terraform.io/modules/equinix-labs/fabric-connection-ibm/equinix/latest/examples/service-token-metal-to-ibm-connection), [Oracle Cloud](https://registry.terraform.io/modules/equinix-labs/fabric-connection-oci/equinix/latest/examples/service-token-metal-to-oci-connection), [Alibaba Cloud](https://registry.terraform.io/modules/equinix-labs/fabric-connection-alibaba/equinix/latest/examples/service-token-metal-to-alibaba-connection).

### Metal Billed Shared Virtual Connection - Non-Redundant connection from Equinix Metal to your own Equinix Fabric Port

```terraform
resource "equinix_metal_connection" "example" {
  name               = "tf-metal-2-port"
  project_id         = "<metal_project_id>"
  type               = "shared"
  redundancy         = "redundant"
  metro              = "FR"
  speed              = "1Gbps"
  service_token_type = "a_side"
  contact_email      = "username@example.com"
}

data "equinix_fabric_ports" "a_side" {
  filters {
    name = "<name_of_port||port_prefix>"
  }
}

resource "equinix_fabric_connection" "example" {
  name = "tf-metal-2-port"
  type = "EVPL_VC"
  notifications {
    type   = "ALL"
    emails = ["username@example.com"]
  }
  project {
    project_id = "<fabric_project_id>"
  }
  bandwidth = "100"
  order {
    purchase_order_number = "1-323292"
  }
  a_side {
    service_token {
      uuid = equinix_metal_connection.example.service_tokens.0.id
    }
  }
  z_side {
    access_point {
      type = "COLO"
      port {
        uuid = data.equinix_fabric_ports.a_side.data.0.uuid
      }
      link_protocol {
        type     = "DOT1Q"
        vlan_tag = 1234
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the connection resource
- `redundancy` (String) Connection redundancy - redundant or primary
- `type` (String) Connection type - dedicated, shared or shared_port_vlan

### Optional

- `contact_email` (String) The preferred email used for communication and notifications about the Equinix Fabric interconnection
- `description` (String) Description of the connection resource
- `facility` (String, Deprecated) Facility where the connection will be created
- `metro` (String) Metro where the connection will be created
- `mode` (String) Mode for connections in IBX facilities with the dedicated type - standard or tunnel
- `organization_id` (String) ID of the organization responsible for the connection. Applicable with type "dedicated"
- `project_id` (String) ID of the project where the connection is scoped to. Required with type "shared"
- `service_token_type` (String) Only used with shared connection. Type of service token to use for the connection, a_side or z_side
- `speed` (String) Connection speed -  Values must be in the format '<number>Mbps' or '<number>Gpbs', for example '100Mbps' or '50Gbps'.  Actual supported values will depend on the connection type and whether the connection uses VLANs or VRF.
- `tags` (List of String) Tags attached to the connection
- `vlans` (List of Number) Only used with shared connection. VLANs to attach. Pass one vlan for Primary/Single connection and two vlans for Redundant connection
- `vrfs` (List of String) Only used with shared connection. VRFs to attach. Pass one VRF for Primary/Single connection and two VRFs for Redundant connection

### Read-Only

- `authorization_code` (String) Only used with Fabric Shared connection. Fabric uses this token to be able to give more detailed information about the Metal end of the network, when viewing resources from within Fabric.
- `id` (String) The unique identifier of the resource
- `ports` (List of Object) List of connection ports - primary (`ports[0]`) and secondary (`ports[1]`) (see [below for nested schema](#nestedatt--ports))
- `service_tokens` (List of Object) Only used with shared connection. List of service tokens required to continue the setup process with [equinix_fabric_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/fabric_connection) or from the [Equinix Fabric Portal](https://fabric.equinix.com/dashboard) (see [below for nested schema](#nestedatt--service_tokens))
- `status` (String) Status of the connection resource
- `token` (String, Deprecated) Only used with shared connection. Fabric Token required to continue the setup process with [equinix_fabric_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/fabric_connection) or from the [Equinix Fabric Portal](https://fabric.equinix.com/dashboard)

<a id="nestedatt--ports"></a>
### Nested Schema for `ports`

Read-Only:

- `id` (String)
- `link_status` (String)
- `name` (String)
- `role` (String)
- `speed` (Number)
- `status` (String)
- `virtual_circuit_ids` (List of String)


<a id="nestedatt--service_tokens"></a>
### Nested Schema for `service_tokens`

Read-Only:

- `expires_at` (String)
- `id` (String)
- `max_allowed_speed` (String)
- `role` (String)
- `state` (String)
- `type` (String)
