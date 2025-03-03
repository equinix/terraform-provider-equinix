# Create Cisco FTD Cluster with Connectivity- PRIVATE

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "Cisco-FTD-SV" {
  name            = "TF_Cisco_NGFW_CLUSTER_ZNPD"
  project_id      = "XXXXXXX"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "Cisco_NGFW"
  self_managed    = true
  connectivity    = "PRIVATE"
  byol            = true
  package_code    = "FTDv10"
  notifications   = ["test@eq.com"]
  account_number  = data.equinix_network_account.sv.number
  version         = "7.0.4-55"
  hostname        = "test"
  core_count      = 4
  term_length     = 1
  interface_count = 10
  cluster_details {
    cluster_name = "tf-ftd-cluster"
    node0 {
      vendor_configuration {
        hostname        = "test"
        activation_key  = "XXXXX"
        controller1     = "X.X.X.X"
        management_type = "FMC"
      }
    }
    node1 {
      vendor_configuration {
        hostname        = "test"
        management_type = "FMC"
      }
    }
  }
}