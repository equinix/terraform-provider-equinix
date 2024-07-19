resource "equinix_metal_vlan" "example1" {
  project_id      = local.my_project_id
  metro           = "SV"
}
resource "equinix_metal_connection" "example" {
  name            = "tf-port-to-metal-legacy"
  project_id      = local.my_project_id
  metro           = "SV"
  redundancy      = "redundant"
  type            = "shared"
  contact_email   = "username@example.com"
  vlans              = [
    equinix_metal_vlan.example1.vxlan,
    equinix_metal_vlan.example2.vxlan
  ]
}
data "equinix_fabric_service_profiles" "zside" {
  count = var.zside_ap_type == "SP" ? 1 : 0
  filter {
    property = "/name"
    operator = "="
    values   = ["AWS Direct Connect"]
  }
}
resource "equinix_fabric_connection" "example" {
  name = "tf-NIMF-metal-2-aws-legacy"
  type = "EVPL_VC"
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
      type               = "METAL_NETWORK"
      authentication_key = equinix_metal_connection.metal-connection.authorization_code
    }
  }
  z_side {
    access_point {
      type               = "SP"
      authentication_key = local.aws_account_id
      seller_region      = "us-west-1"
      profile {
        type = "L2_PROFILE"
        uuid = data.equinix_fabric_service_profiles.zside[0].id
      }
      location {
        metro_code ="SV"
      }
    }
  }
}
