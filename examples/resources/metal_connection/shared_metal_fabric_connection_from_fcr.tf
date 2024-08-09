resource "equinix_metal_vlan" "example1" {
  project_id      = local.my_project_id
  metro           = "SV"
}
resource "equinix_metal_connection" "example" {
  name            = "tf-port-to-metal-legacy"
  project_id      = local.my_project_id
  metro           = "SV"
  redundancy      = "primary"
  type            = "shared_port_vlan"
  contact_email   = "username@example.com"
  vlans              = [ equinix_metal_vlan.example1.vxlan ]
}
resource "equinix_fabric_connection" "example" {
  name = "tf-NIMF-metal-2-aws-legacy"
  type = "IP_VC"
  notifications {
    type   = "ALL"
    emails = "sername@example.com"
  }
  project {
    project_id = local.fabric_project_id
  }
  bandwidth       = "200"
  order {
    purchase_order_number = "1-323292"
  }
  a_side {
      access_point {
        type = "CLOUD_ROUTER"
        router {
          uuid = local.cloud_router_uuid
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
