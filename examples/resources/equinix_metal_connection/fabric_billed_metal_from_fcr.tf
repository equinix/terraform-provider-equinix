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
