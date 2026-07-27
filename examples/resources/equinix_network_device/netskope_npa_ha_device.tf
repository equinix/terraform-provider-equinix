# Create NETSKOPE NPA HA device

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

resource "equinix_network_device" "netskope-npa" {
  name            = "NETSKOPE-NPA"
  project_id      = "xxxxxxx"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "NETSKOPE-NPA"
  self_managed    = true
  byol            = true
  interface_count = 1
  package_code    = "STD"
  notifications   = ["test@eq.com"]
  connectivity    = "PRIVATE"
  account_number  = data.equinix_network_account.sv.number
  version         = "R138"
  core_count      = 2
  term_length     = 1
    vendor_configuration = {
      hostname            = "test"
      privateCidrMask     = "24"
      ipAddressType       = "STATIC"
      ipAddress           = "x.x.x.x"
      gatewayIp           = "x.x.x.x",
      primaryNameServer   = "x.x.x.x"
      secondaryNameServer = "x.x.x.x"
      dnsSearchDomain     = "xxxxx"
    }
  secondary_device {
    name           = "NETSKOPE-NPA-Sec"
    metro_code     = data.equinix_network_account.sv.metro_code
    account_number = data.equinix_network_account.sv.number
    notifications  = ["test@eq.com"]
    vendor_configuration = {
      hostname            = "test"
      privateCidrMask     = "24"
      ipAddressType       = "STATIC"
      ipAddress           = "x.x.x.x"
      gatewayIp           = "x.x.x.x",
      primaryNameServer   = "x.x.x.x"
      secondaryNameServer = "x.x.x.x"
      dnsSearchDomain     = "xxxxx"
    }
  }
}
