provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

data "equinix_ecx_l2_sellerprofile" "alibaba" {
  name = "Alibaba Cloud Express Connect"
}

data "equinix_ecx_port" "dot1q-pri" {
  name = var.equinix_port_name
}

resource "equinix_ecx_l2_connection" "alibaba-dot1q" {
  name                  = "tf-alibaba-dot1q"
  profile_uuid          = data.equinix_ecx_l2_sellerprofile.alibaba.uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["example@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = data.equinix_ecx_port.dot1q-pri.uuid
  vlan_stag             = 2100
  seller_region         = "ap-southeast-2"
  seller_metro_code     = "SY"
  authorization_key     = var.alibaba_account_id
}
