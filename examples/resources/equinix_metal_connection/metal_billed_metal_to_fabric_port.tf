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


    