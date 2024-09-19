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