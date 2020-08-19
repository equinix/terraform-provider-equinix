data "equinix_ecx_l2_sellerprofile" "alibaba" {
  name = "Alibaba Express Connect"
}

data "equinix_ecx_port" "dot1q-1-pri" {
  name = "sit-001-CX-DC5-NL-Dot1q-BO-10G-PRI-JUN-27"
}

resource "equinix_ecx_l2_connection" "alibaba-dot1q" {
  name                  = "tf-alibaba-dot1q"
  profile_uuid          = data.equinix_ecx_l2_sellerprofile.alibaba.uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = data.equinix_ecx_port.dot1q-1-pri.uui
  vlan_stag             = 2100
  seller_region         = "ap-southeast-2"
  seller_metro_code     = "SY"
  authorization_key     = "123456789"
}
