# Create Aruba Edgeconnect SDWAN HA device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}


resource "equinix_network_device" "ARUBA-EDGECONNECT-AM" {
  name                 = "TF_Aruba_Edge_Connect"
  project_id           = "XXXXX"
  metro_code           = data.equinix_network_account.sv.metro_code
  type_code            = "EDGECONNECT-SDWAN"
  self_managed         = true
  byol                 = true
  package_code         = "EC-V"
  notifications        = ["test@eq.com"]
  account_number       = data.equinix_network_account.sv.number
  version              = "9.4.2.3"
  core_count           = 2
  term_length          = 1
  additional_bandwidth = 50
  interface_count      = 32
  acl_template_id      = "XXXXXXX"
  vendor_configuration = {
    accountKey : "xxxxx"
    accountName : "xxxx"
    applianceTag : "tests"
    hostname : "test"
  }
  secondary_device {
    name                 = "TF_CHECKPOINT"
    metro_code           = data.equinix_network_account.sv.metro_code
    account_number       = data.equinix_network_account.sv.number
    acl_template_id      = "XXXXXXX"
    notifications        = ["test@eq.com"]
    vendor_configuration = {
      accountKey : "xxxxx"
      accountName : "xxxx"
      applianceTag : "test"
      hostname : "test"
    }
  }
}