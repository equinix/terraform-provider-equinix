# Create Checkpoint single device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "CHECKPOINT-SV" {
  name                 = "TF_CHECKPOINT"
  project_id           = "XXXX"
  metro_code           = data.equinix_network_account.sv.metro_code
  type_code            = "CGUARD"
  self_managed         = true
  byol                 = true
  package_code         = "STD"
  notifications        = ["test@eq.com"]
  account_number       = data.equinix_network_account.sv.number
  version              = "R81.20"
  hostname             = "test"
  core_count           = 2
  term_length          = 1
  additional_bandwidth = 5
  acl_template_id      = "XXXXXXX"
  ssh_key {
    username = "XXXXX"
    key_name = "XXXXXX"
  }
}