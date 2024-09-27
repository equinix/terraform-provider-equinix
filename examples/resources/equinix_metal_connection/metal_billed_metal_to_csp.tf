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