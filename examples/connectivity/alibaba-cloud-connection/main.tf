provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

data "equinix_ecx_l2_sellerprofile" "alibaba" {
  name                     = "Alibaba Cloud Express Connect"
  organization_global_name = "Alibaba"
}

data "equinix_ecx_port" "dot1q-pri" {
  name = var.equinix_port_name
}

resource "equinix_ecx_l2_connection" "example" {
  name              = "tf-alibaba-dot1q"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.alibaba.id
  speed             = 50
  speed_unit        = "MB"
  notifications     = ["example@equinix.com"]
  port_uuid         = data.equinix_ecx_port.dot1q-pri.id
  vlan_stag         = 2100
  seller_region     = var.alibaba_region
  seller_metro_code = var.alibaba_metro_code
  authorization_key = var.alibaba_account_id
}
