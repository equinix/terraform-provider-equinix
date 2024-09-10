# Create self configured redundant Arista router with DSA key

data "equinix_network_account" "sv" {
  name = "account-name"
  metro_code = "SV"
}

resource "equinix_network_ssh_key" "test-public-key" {
  name = "key-name"
  public_key = "ssh-dss key-value"
  type = "DSA"
}

resource "equinix_network_device" "arista-ha" {
  name            = "tf-arista-p"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "ARISTA-ROUTER"
  self_managed    = true
  connectivity    = "PRIVATE"
  byol            = true
  package_code    = "CloudEOS"
  notifications   = ["test@equinix.com"]
  hostname        = "arista-p"
  account_number  = data.equinix_network_account.sv.number
  version         = "4.29.0"
  core_count      = 4
  term_length     = 12
  additional_bandwidth = 5
  ssh_key {
    username = "test-username"
    key_name = equinix_network_ssh_key.test-public-key.name
  }
  acl_template_id = "c637a17b-7a6a-4486-924b-30e6c36904b0"
  secondary_device {
    name            = "tf-arista-s"
    metro_code      = data.equinix_network_account.sv.metro_code
    hostname        = "arista-s"
    notifications   = ["test@eq.com"]
    account_number  = data.equinix_network_account.sv.number
    acl_template_id = "fee5e2c0-6198-4ce6-9cbd-bbe6c1dbe138"
  }
}
