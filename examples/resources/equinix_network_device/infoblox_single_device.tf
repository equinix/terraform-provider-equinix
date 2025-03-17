# Create Infoblox Grid Member Single device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "INFOBLOX-SV" {
  name                 = "TF_INFOBLOX"
  project_id           = "XXXXXXXXXX"
  metro_code           = data.equinix_network_account.sv.metro_code
  type_code            = "INFOBLOX-GRID-MEMBER"
  self_managed         = true
  byol                 = true
  connectivity         = "PRIVATE"
  package_code         = "STD"
  notifications        = ["test@eq.com"]
  account_number       = data.equinix_network_account.sv.number
  version              = "9.0.5"
  hostname             = "test"
  core_count           = 8
  term_length          = 1
  vendor_configuration = {
    adminPassword = "X.X.X.X"
    ipAddress     = "X.X.X.X"
    subnetMaskIp = "X.X.X.X"
    gatewayIp     = "X.X.X.X"
  }
}