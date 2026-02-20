# Create FG VM Cluster with connectivity PRIVATE and IP Address Type as STATIC

data "equinix_network_account" "sv" {
  metro_code = "SV"
  name       = "account-name"
}

resource "equinix_network_device" "FGVM-SV" {
  name            = "tf-fgvm-cluster-static-znpd"
  metro_code      = "DC"
  type_code       = "FG-VM"
  project_id      = "xxxxxxx"
  self_managed    = true
  connectivity    = "PRIVATE"
  byol            = true
  package_code    = "VM02"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  term_length     = 12
  account_number  = xxxxxx
  version         = "7.6.2"
  interface_count = 10
  core_count      = 2
  ssh_key {
    username = "sanity1"
    key_name = ""
  }
  cluster_details {
    cluster_name = "tf-fgvm--cluster"
    node0 {
      vendor_configuration {
        ip_address              = "x.x.x.x"
        subnet_mask_ip          = "x.x.x.x"
        gateway_ip              = "x.x.x.x"
        management_interface_id = "5"
        hostname                = "test"
        ip_address_type         = "STATIC"
      }
    }
    node1 {
      vendor_configuration {
        ip_address              = "x.x.x.x"
        subnet_mask_ip          = "x.x.x.x"
        gateway_ip              = "x.x.x.x"
        management_interface_id = "5"
        hostname                = "test"
        ip_address_type         = "STATIC"
      }
    }
  }
}
