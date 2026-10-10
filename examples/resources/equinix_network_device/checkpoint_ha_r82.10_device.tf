# Create Checkpoint ha device - R82.10 version

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "CHECKPOINT-SV" {
  name            = "Praveena_TF_CHECKPOINT"
  project_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "CGUARD"
  self_managed    = true
  byol            = true
  package_code    = "STD"
  notifications   = ["test@eq.com"]
  account_number  = data.equinix_network_account.sv.number
  version         = "R82.10_JHF_24"
  interface_count = 10 # 10 interfaces for version - R82.10
  hostname        = "test"
  core_count      = 2
  term_length     = 1
  vendor_configuration = {
    clusterMember = false # default value false if not present
    sicKey        = "vpnxxxxxxxxx" # mandatory for R82.10 version
  }
  acl_template_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  ssh_key {
    username = "admin"
    key_name = "xxxxx"
  }
  secondary_device {
    name            = "TF_CHECKPOINT-secondary"
    metro_code      = data.equinix_network_account.sv.metro_code
    hostname        = "test"
    notifications   = ["test@eq.com"]
    account_number  = data.equinix_network_account.sv.number
    acl_template_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    vendor_configuration = {
      clusterMember = false # default value false if not present
      sicKey        = "vpnxxxxxxxxx" # mandatory for R82.10 version
    }
  }
}
