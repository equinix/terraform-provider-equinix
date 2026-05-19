# add_secondary_device example for C8000V with self-configured BYOL
data "equinix_network_account" "am" {
  name       = "YOUR_ACCOUNT_NAME"
  metro_code = "AM"
  project_id = "YOUR_PROJECT_ID"
}

resource "equinix_network_device" "c8kv_primary" {
  name                      = "C8000V-primary"
  project_id                = "YOUR_PROJECT_ID"
  metro_code                = data.equinix_network_account.am.metro_code
  license_token             = "YOUR_LICENSE_TOKEN"
  type_code                 = "C8000V"
  self_managed              = true
  byol                      = true
  tier                      = 1
  package_code              = "network-advantage"
  generate_default_password = true
  notifications             = ["you@example.com"]
  account_number            = data.equinix_network_account.am.number
  version                   = "17.11.01a"
  hostname                  = "c8kv-primary"
  core_count                = 2
  term_length               = 1
  additional_bandwidth      = 5
  interface_count           = 10
  connectivity              = "INTERNET-ACCESS"
  acl_template_id           = "YOUR_ACL_TEMPLATE_ID"

  ssh_key {
    username = "YOUR_SSH_USERNAME"
    key_name = "YOUR_SSH_KEY_NAME"
  }

  secondary_device {
    name                 = "C8000V-secondary"
    metro_code           = data.equinix_network_account.am.metro_code
    acl_template_id      = "YOUR_ACL_TEMPLATE_ID"
    notifications        = ["you@example.com"]
    hostname             = "c8kv-secondary"
    additional_bandwidth = 5
    account_number       = data.equinix_network_account.am.number
    license_token        = "YOUR_LICENSE_TOKEN"
    vendor_configuration = { hostNamePrefix = "C8KV-secondary" }
    ssh_key {
      username = "YOUR_SSH_USERNAME"
      key_name = "YOUR_SSH_KEY_NAME"
    }
  }
}
