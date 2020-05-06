provider equinix {
  endpoint = "http://localhost:8080"
  client_id = "someID"
  client_secret = "someSecret"
}

resource "equinix_ecx_l2_connection" "aws_dot1q" {
 name = "tf-single-aws"
 profile_uuid = "2a4f7e27-dff8-4f15-aeda-a11ffe9ccf73"
 speed = 200
 speed_unit = "MB"
 notifications = ["marry@equinix.com", "john@equinix.com"]
 purchase_order_number = "1234567890"
 port_uuid = "febc9d80-11e0-4dc8-8eb8-c41b6b378df2"
 vlan_stag = 777
 vlan_ctag = 1000
 seller_region = "us-east-1"
 seller_metro_code = "SV"
 authorization_key = "1234456"
}

resource "equinix_ecx_l2_connection" "redundant_self" {
  name = "tf-redundant-self"
  profile_uuid = "2a4f7e27-dff8-4f15-aeda-a11ffe9ccf73"
  speed = 50
  speed_unit = "MB"
  notifications = ["john@equinix.com", "marry@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid = "febc9d80-11e0-4dc8-8eb8-c41b6b378df2"
  vlan_stag = 800
  zside_port_uuid = "03a969b5-9cea-486d-ada0-2a4496ed72fb"
  zside_vlan_stag = 1010
  seller_region = "us-east-1"
  seller_metro_code = "SV"
  secondary_connection {
    name = "tf-redundant-self-sec"
    port_uuid = "86872ae5-ca19-452b-8e69-bb1dd5f93bd1"
    vlan_stag = 999
    vlan_ctag = 1000
    zside_port_uuid = "393b2f6e-9c66-4a39-adac-820120555420"
    zside_vlan_stag = 1022
  }
}
