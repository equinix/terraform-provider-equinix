# Create FG VM ha device with connectivity PRIVATE and IP Address Type as STATIC

data "equinix_network_account" "sv" {
  metro_code = "SV"
  name       = "account-name"
}

resource "equinix_network_device" "FTNT-FIREWALL-SV" {
  name                 = "TF_FTNT-FIREWALL"
  project_id           = "XXXXXXXXXX"
  metro_code           = data.equinix_network_account.sv.metro_code
  interface_count      = 10
  type_code            = "FG-VM"
  self_managed         = true
  byol                 = true
  connectivity         = "PRIVATE"
  package_code         = "VM02"
  notifications        = ["test@eq.com"]
  account_number       = data.equinix_network_account.sv.number
  version              = "7.6.3"
  hostname             = "test"
  core_count           = 2
  term_length          = 1
  vendor_configuration = {
    gatewayIp     = "X.X.X.X"
    ipAddress     = "X.X.X.X"
    ipAddressType = "STATIC"
    subnetMaskIp  = "x.x.x.x"
  }
  secondary_device {
    name                 = "TF_FTNT-FIREWALL-secondary"
    metro_code           = data.equinix_network_account.sv.metro_code
    hostname             = "fg-vm-znpd"
    notifications        = ["john@equinix.com", "marry@equinix.com"]
    account_number       = data.equinix_network_account.sv.number
    vendor_configuration = {
      ipAddressType = "STATIC", ipAddress = "x.x.x.x", gatewayIp = "x.x.x.x",
      subnetMaskIp  = "x.x.x.x", managementInterfaceId = "6"
    }
  }
}