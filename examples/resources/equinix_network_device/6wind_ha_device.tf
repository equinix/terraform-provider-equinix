# Create 6WIND VSR HA device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "6wind-vsr" {
  name            = "6WIND-VSR"
  project_id      = "df3f3b30-48b8-4dac-bcea-d7e719a7f436"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "6WIND-VSR"
  self_managed    = true
  byol            = true
  interface_count = 10
  package_code    = "STD"
  notifications   = ["test@eq.com"]
  account_number  = data.equinix_network_account.sv.number
  version         = "3.10.8"
  core_count      = 2
  term_length     = 1
  vendor_configuration = {
    hostname = "test"
    token    = "xxxx"
  }
  ssh_key {
    username = "xxxx"
    key_name = "xxxxx"
  }
  secondary_device {
    name           = "6WIND-VSR-Sec"
    metro_code     = data.equinix_network_account.sv.metro_code
    account_number = data.equinix_network_account.sv.number
    notifications  = ["test@eq.com"]
    vendor_configuration = {
      hostname = "test"
      token    = "xxxx"
    }
  }
}
