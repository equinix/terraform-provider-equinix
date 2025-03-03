# Create VYos Router HA device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "VYOS-AM" {
  name                 = "TF_VYOS"
  project_id           = "XXXXXXX"
  metro_code           = data.equinix_network_account.sv.metro_code
  type_code            = "VYOS-ROUTER"
  self_managed         = true
  byol                 = false
  package_code         = "STD"
  notifications        = ["test@eq.com"]
  account_number       = data.equinix_network_account.sv.number
  version              = "1.4.1-2501"
  hostname             = "test"
  core_count           = 2
  term_length          = 1
  additional_bandwidth = 50
  acl_template_id      = "XXXXXXXX"
  ssh_key {
    username = "test"
    key_name = "xxxxxxxx"
  }
  secondary_device {
    name            = "TF_CHECKPOINT"
    metro_code      = data.equinix_network_account.sv.metro_code
    account_number  = data.equinix_network_account.sv.number
    hostname        = "test"
    acl_template_id = "XXXXXXXXXXX"
    notifications   = ["test@eq.com"]
  }
}