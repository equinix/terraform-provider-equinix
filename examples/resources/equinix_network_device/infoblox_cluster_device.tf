# Create Infoblox Grid Member HA device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "INFOBLOX-SV" {
  name           = "TF_INFOBLOX"
  project_id     = "XXXXXXXXXX"
  metro_code     = data.equinix_network_account.sv.metro_code
  type_code      = "INFOBLOX-GRID-MEMBER"
  self_managed   = true
  byol           = true
  package_code   = "STD"
  notifications  = ["test@eq.com"]
  account_number = data.equinix_network_account.sv.number
  version        = "9.0.5"
  hostname       = "test"
  connectivity   = "PRIVATE"
  core_count     = 8
  term_length    = 1
  cluster_details {
    cluster_name = "tf-infoblox-cluster"
    node0 {
      vendor_configuration {
        admin_password = "Welcome@1"
        ip_address     = "192.168.1.35"
        subnet_mask_ip = "255.255.255.0"
        gateway_ip     = "192.168.1.1"
        hostname       = "test"
      }
    }
    node1 {
      vendor_configuration {
        admin_password = "Welcome@1"
        ip_address     = "192.168.1.35"
        subnet_mask_ip = "255.255.255.0"
        gateway_ip     = "192.168.1.1"
        hostname       = "test"
      }
    }
  }
}