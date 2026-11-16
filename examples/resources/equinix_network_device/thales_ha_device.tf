# Create Thales HA device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "thales-ha" {
  name            = "Thales"
  project_id      = "xxxxxxx"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "THALES"
  self_managed    = true
  byol            = true
  interface_count = 3
  package_code    = "STD"
  notifications   = ["test@eq.com"]
  account_number  = data.equinix_network_account.sv.number
  version         = "5.6.4.D31681"
  core_count      = 4
  term_length     = 1
  secondary_device {
    name           = "Thales-Sec"
    metro_code     = data.equinix_network_account.sv.metro_code
    account_number = data.equinix_network_account.sv.number
    notifications  = ["test@eq.com"]
  }
}
