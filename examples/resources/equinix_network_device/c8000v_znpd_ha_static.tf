# Create C8000V HA - BYOL device with connectivity PRIVATE with static IP address type

data "equinix_network_account" "sv" {
  metro_code = "SV"
  name       = "account-name"
}

resource "equinix_network_device" "c8000v-byol" {
  name            = "tf-c8000v-byol"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "C8000V"
  self_managed    = true
  byol            = true
  package_code    = "network-essentials"
  connectivity    = "PRIVATE"
  notifications   = ["john@equinix.com", "marry@equinix.com", "fred@equinix.com"]
  term_length     = 12
  account_number  = data.equinix_network_account.sv.number
  version         = "17.11.01a"
  interface_count = 10
  core_count      = 2
  tier            = 1
  ssh_key {
    username = "test"
    key_name = "test-key"
  }
  vendor_configuration = {
    restApiSupportRequirement = "true", ipAddressType = "STATIC", ipAddress = "x.x.x.x", gatewayIp = "x.x.x.x",
    subnetMaskIp              = "x.x.x.x", managementInterfaceId= "6"
  }
  secondary_device {
    name                 = "tf-c8000v-byol-secondary"
    metro_code           = data.equinix_network_account.sv.metro_code
    hostname             = "csr8000v-s"
    notifications        = ["john@equinix.com", "marry@equinix.com"]
    account_number       = data.equinix_network_account.sv.number
    vendor_configuration = {
      restApiSupportRequirement = "true", ipAddressType = "STATIC", ipAddress = "x.x.x.x", gatewayIp = "x.x.x.x",
      subnetMaskIp              = "x.x.x.x", managementInterfaceId= "6"
    }
  }
}