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