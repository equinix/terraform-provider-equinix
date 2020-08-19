data "equinix_ecx_port" "dot1q-1-pri" {
  name = "sit-001-CX-DC5-NL-Dot1q-BO-10G-PRI-JUN-27"
}

data "equinix_ecx_l2_sellerprofile" "ibm" {
  name = "IBM Cloud Direct Link Exchange"
}

resource "equinix_ecx_l2_connection" "ibm-dot1q" {
  name                  = "tf-ibm-dot1q"
  profile_uuid          = data.equinix_ecx_l2_sellerprofile.ibm.uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = data.equinix_ecx_port.dot1q-1-pri.uuid
  vlan_stag             = 190
  seller_region         = "Washington DC 1"
  seller_metro_code     = "DC"
  authorization_key     = "123456789"
  additional_info {
    name  = "global"
    value = "true"
  }
  additional_info {
    name  = "asn"
    value = "10509"
  }
}
