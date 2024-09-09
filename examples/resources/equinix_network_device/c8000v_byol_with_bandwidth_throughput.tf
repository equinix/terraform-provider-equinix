# Create C8000V BYOL device with numeric bandwidth throughput information

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "c8000v-byol-throughput" {
  name            = "tf-c8000v-byol"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "C8000V"
  self_managed    = true
  byol            = true
  package_code    = "VM100"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  term_length     = 12
  account_number  = data.equinix_network_account.sv.number
  version         = "17.11.01a"
  interface_count = 10
  core_count      = 2
  throughput 	  = "100"
  throughput_unit = "Mbps"
  ssh_key {
    username = "test"
    key_name = "test-key"
  }
  acl_template_id = "0bff6e05-f0e7-44cd-804a-25b92b835f8b"
}