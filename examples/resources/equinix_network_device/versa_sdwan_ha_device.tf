# Create Fortinet SDWAN single device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "FTNT-SDWAN-SV" {
  name                 = "TF_VERSA-SDWAN"
  project_id           = "XXXXXXXXX"
  metro_code           = data.equinix_network_account.sv.metro_code
  type_code            = "VERSA_SDWAN"
  self_managed         = true
  byol                 = true
  package_code         = "FLEX_VNF_2"
  notifications        = ["test@eq.com"]
  account_number       = data.equinix_network_account.sv.number
  version              = "21.2.3"
  core_count           = 2
  term_length          = 1
  additional_bandwidth = 50
  acl_template_id      = "XXXXXXXXX"
  vendor_configuration = {
    controller1 = "X.X.X.X"
    controller2 = "X.X.X.X"
    localId     = "test@test.com"
    remoteId    = "test@test.com"
    serialNumber = "4"
  }
  secondary_device {
    name                 = "Praveena_TF_VERSA"
    metro_code           = data.equinix_network_account.sv.metro_code
    account_number       = data.equinix_network_account.sv.number
    acl_template_id      = "XXXXXXXX"
    notifications        = ["test@eq.com"]
    vendor_configuration = {
      controller1 = "X.X.X.X"
      controller2 = "X.X.X.X"
      localId     = "test@test.com"
      remoteId    = "test@test.com"
      serialNumber = "4"
    }
  }
}