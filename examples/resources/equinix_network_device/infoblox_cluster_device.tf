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
  connectivity   = "PRIVATE"
  core_count     = 8
  term_length    = 1
  cluster_details {
    cluster_name = "tf-infoblox-cluster"
    node0 {
      vendor_configuration {
        admin_password = "xxxxxxx"
        ip_address     = "X.X.X.X"
        subnet_mask_ip = "X.X.X.X"
        gateway_ip     = "X.X.X.X"
      }
    }
    node1 {
      vendor_configuration {
        admin_password = "xxxxxxx"
        ip_address     = "X.X.X.X"
        subnet_mask_ip = "X.X.X.X"
        gateway_ip     = "X.X.X.X"
      }
    }
  }
}