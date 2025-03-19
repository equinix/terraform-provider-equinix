# Create Infoblox Grid Member HA device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "INFOBLOX-SV" {
  name                 = "TF_INFOBLOX"
  project_id           = "XXXXXXXXXX"
  metro_code           = data.equinix_network_account.sv.metro_code
  type_code            = "INFOBLOX-GRID-MEMBER"
  self_managed         = true
  connectivity         = "PRIVATE"
  byol                 = true
  package_code         = "STD"
  notifications        = ["test@eq.com"]
  account_number       = data.equinix_network_account.sv.number
  version              = "9.0.5"
  core_count           = 8
  term_length          = 1
  vendor_configuration = {
    adminPassword = "X.X.X.X"
    ipAddress     = "X.X.X.X"
    subnetMaskIp  = "X.X.X.X"
    gatewayIp     = "X.X.X.X"
  }
  secondary_device {
    name                 = "TF_INFOBLOX-Sec"
    metro_code           = data.equinix_network_account.sv.metro_code
    account_number       = data.equinix_network_account.sv.number
    notifications        = ["test@eq.com"]
    vendor_configuration = {
      adminPassword = "X.X.X.X"
      ipAddress     = "X.X.X.X"
      subnetMaskIp  = "X.X.X.X"
      gatewayIp     = "X.X.X.X"
    }
  }
}