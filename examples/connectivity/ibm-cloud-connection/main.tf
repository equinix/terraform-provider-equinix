provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

data "equinix_ecx_port" "dot1q-pri" {
  name = var.equinix_port_name
}

data "equinix_ecx_l2_sellerprofile" "ibm" {
  name                     = "IBM Cloud Direct Link Exchange"
  organization_global_name = "IBM"
}

resource "equinix_ecx_l2_connection" "ibm-dot1q" {
  name              = "tf-ibm-dot1q"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.ibm.uuid
  speed             = 50
  speed_unit        = "MB"
  notifications     = ["example@equinix.com"]
  port_uuid         = data.equinix_ecx_port.dot1q-pri.uuid
  vlan_stag         = 1900
  seller_region     = var.ibm_region
  seller_metro_code = var.ibm_metro_code
  authorization_key = var.ibm_account_id
  additional_info {
    name  = "global"
    value = "true"
  }
  additional_info {
    name  = "asn"
    value = "10509"
  }
}
