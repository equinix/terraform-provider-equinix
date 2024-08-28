# Create self configured redundant BlueCat DNS and DHCP Server
data "equinix_network_account" "sv" {
  name       = "account-name"
  metro_code = "SV"
}

resource "equinix_network_ssh_key" "test-public-key" {
  name       = "key-name"
  public_key = "ssh-dss key-value"
  type       = "DSA"
}

resource "equinix_network_device" "bluecat-bdds-ha" {
  name                 = "tf-bluecat-bdds-p"
  metro_code           = data.equinix_network_account.sv.metro_code
  type_code            = "BLUECAT"
  self_managed         = true
  connectivity         = "PRIVATE"
  byol                 = true
  package_code         = "STD"
  notifications        = ["test@equinix.com"]
  account_number       = data.equinix_network_account.sv.number
  version              = "9.6.0"
  core_count           = 2
  term_length          = 12
  vendor_configuration = {
    "hostname" = "test"
    "privateAddress" : "x.x.x.x"
    "privateCidrMask" : "24"
    "privateGateway" : "x.x.x.x"
    "licenseKey" : "xxxxx-xxxxx-xxxxx-xxxxx-xxxxx"
    "licenseId" : "xxxxxxxxxxxxxxx"
  }
  ssh_key {
    username = "test-username"
    key_name = equinix_network_ssh_key.test-public-key.name
  }
  secondary_device {
    name                 = "tf-bluecat-bdds-s"
    metro_code           = data.equinix_network_account.sv.metro_code
    notifications        = ["test@eq.com"]
    account_number       = data.equinix_network_account.sv.number
    vendor_configuration = {
      "hostname" = "test"
      "privateAddress" : "x.x.x.x"
      "privateCidrMask" : "24"
      "privateGateway" : "x.x.x.x"
      "licenseKey" : "xxxxx-xxxxx-xxxxx-xxxxx-xxxxx"
      "licenseId" : "xxxxxxxxxxxxxxx"
    }
  }
}
