# Create ZSCALER APPC device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "zscaler-appc-single" {
  name            = "tf-zscaler-appc"
  project_id      = "XXXXXX"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "ZSCALER-APPC"
  self_managed    = true
  byol            = true
  connectivity    = "PRIVATE"
  package_code    = "STD"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  term_length     = 12
  account_number  = data.equinix_network_account.sv.number
  version         = "23.395.1"
  interface_count = 1
  core_count      = 4
  vendor_configuration   = {"provisioningKey" = "XXXXXXXXXX", "hostname" = "XXXX"}
  ssh_key {
    username = "test"
    key_name = "test-key"
  }
}