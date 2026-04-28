# Create Infoblox Grid Member HA device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "INFOBLOX-SV" {
  name            = "TF_INFOBLOX-NIOS-X"
  project_id      = "xxxxxxx"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "INFOBLOX-NIOSX"
  self_managed    = true
  connectivity    = "INTERNET-ACCESS"
  byol            = true
  package_code    = "STD"
  notifications   = ["test@eq.com"]
  account_number  = data.equinix_network_account.sv.name.number
  version         = "4.0"
  core_count      = 3
  interface_count = 5
  term_length     = 1
  vendor_configuration = {
    hostname = "test"
    token    = "xxxxx"
  }
  secondary_device {
    name           = "TF_INFOBLOX-NIOS-X-Sec"
    metro_code     = data.equinix_network_account.sv.metro_code
    account_number = data.equinix_network_account.sv.name.number
    notifications  = ["test@eq.com"]
    vendor_configuration = {
      hostname = "test"
      token    = "xxxxx"
    }
  }
}
