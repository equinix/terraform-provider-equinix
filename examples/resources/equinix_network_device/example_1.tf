# Create pair of redundant, managed CSR1000V routers with license subscription
# in two different metro locations

data "equinix_network_account" "dc" {
  metro_code = "DC"
}

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "csr1000v-ha" {
  name            = "tf-csr1000v-p"
  throughput      = 500
  throughput_unit = "Mbps"
  metro_code      = data.equinix_network_account.dc.metro_code
  type_code       = "CSR1000V"
  self_managed    = false
  connectivity    = "INTERNET-ACCESS"
  byol            = false
  package_code    = "SEC"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  hostname        = "csr1000v-p"
  term_length     = 12
  account_number  = data.equinix_network_account.dc.number
  version         = "16.09.05"
  core_count      = 2
  secondary_device {
    name            = "tf-csr1000v-s"
    metro_code      = data.equinix_network_account.sv.metro_code
    hostname        = "csr1000v-s"
    notifications   = ["john@equinix.com", "marry@equinix.com"]
    account_number  = data.equinix_network_account.sv.number
  }
}
