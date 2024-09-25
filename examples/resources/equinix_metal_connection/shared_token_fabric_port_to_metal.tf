resource "equinix_metal_vlan" "example" {
    project_id      = local.my_project_id
    metro           = "FR"
}

resource "equinix_metal_connection" "example" {
    name               = "tf-port-to-metal"
    project_id         = local.project_id
    type               = "shared"
    redundancy         = "primary"
    metro              = "FR"
    speed              = "200Mbps"
    service_token_type = "z_side"
    contact_email      = "username@example.com"
    vlans              = [
      equinix_metal_vlan.example.vxlan
    ]
}

data "equinix_fabric_ports" "example" {
  filters {
    name = "CX-FR5-NL-Dot1q-BO-1G-PRI"
  }
}

resource "equinix_fabric_connection" "example" {
  name = "port-2-shared-metal-token"
  type = "EVPL_VC"
  notifications {
    type   = "ALL"
    emails = ["example@equinix.com"]
  }
  bandwidth = 50
  a_side {
    access_point {
      type= "COLO"
      port {
        uuid = data.equinix_fabric_ports.example.id
      }
      link_protocol {
        type = "DOT1Q"
        vlan_tag = "1020"
      }
    }
  }
  z_side {
    service_token {
      uuid = equinix_metal_connection.example.service_tokens.0.id
    }
  }
}
