provider "equinix" {
  client_id     = "your_client_id"
  client_secret = "your_client_secret"
}

data "equinix_ecx_port" "dot1q-1-pri" {
  name = "sit-001-CX-DC5-NL-Dot1q-BO-10G-PRI-JUN-27"
}

data "equinix_ecx_l2_sellerprofile" "oracle" {
  name = "Oracle Cloud Infrastructure -OCI- FastConnect"
}

resource "equinix_ecx_l2_connection" "oracle-dot1q" {
  name                  = "tf-oracle-dot1q"
  profile_uuid          = var.sp_oracle_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_dot1q_pri_uuid
  vlan_stag             = 1700
  seller_region         = "us-ashburn-1"
  seller_metro_code     = "DC"
  authorization_key     = "123456789"
}
