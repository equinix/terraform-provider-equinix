# Create Fortinet SDWAN single device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "FTNT-SDWAN-SV" {
  name                 = "TF_FTNT-SDWAN"
  project_id           = "XXXXXXXXXX"
  metro_code           = data.equinix_network_account.sv.metro_code
  type_code            = "FG-SDWAN"
  self_managed         = true
  byol                 = true
  package_code         = "VM02"
  notifications        = ["test@eq.com"]
  account_number       = data.equinix_network_account.sv.number
  version              = "7.0.14"
  hostname             = "test"
  core_count           = 2
  term_length          = 1
  additional_bandwidth = 50
  acl_template_id      = "XXXXXXXX"
  vendor_configuration = {
    adminPassword = "XXXXX"
    controller1 = "X.X.X.X"
  }
}