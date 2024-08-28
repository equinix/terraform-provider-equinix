# Create self configured single Catalyst 8000V (Autonomous Mode) router with license token

data "equinix_network_account" "sv" {
  name = "account-name"
  metro_code = "SV"
}

resource "equinix_network_device" "c8kv-single" {
  name            = "tf-c8kv"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "C8000V"
  self_managed    = true
  byol            = true
  package_code    = "network-essentials"
  notifications   = ["test@equinix.com"]
  hostname        = "C8KV"
  account_number  = data.equinix_network_account.sv.number
  version         = "17.06.01a"
  core_count      = 2
  term_length     = 12
  license_token = "valid-license-token"
  additional_bandwidth = 5
  ssh_key {
    username = "test-username"
    key_name = "valid-key-name"
  }
  acl_template_id = "3e548c02-9164-4197-aa23-05b1f644883c"
}
