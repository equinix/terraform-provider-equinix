# Create F5XC device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "f5xc-single" {
  name            = "tf-f5xc"
  project_id      = "XXXXXX"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "F5XC"
  self_managed    = true
  byol            = true
  connectivity    = "INTERNET-ACCESS"
  package_code    = "STD"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  term_length     = 1
  account_number  = data.equinix_network_account.sv.number
  acl_template_id = "xxxx"
  version         = "9.2025.17"
  interface_count = 8
  core_count      = 8
  vendor_configuration   = {"token" = "XXXXXXXXXX", "hostname" = "XXXX"}
}