# Equinix Network Edge: Silver Peak Unity EdgeConnect

A Terraform module to create Silver Peak Unity EdgeConnect SD-WAN network device
on the Equinix platform.

Device is created as self-managed device with *bring your own license*
licensing mode.

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 0.12.0 |
| equinix | >= 1.1.0 |

## Providers

| Name | Version |
|---------|----------|
| equinix | >= 1.1.0 |

## Assumptions

* there is one billing account in `Active` status for every used
`metro_code`
* most recent, stable version of a software for a given `software_package` is used
* secondary device name is same as primary with `-secondary` suffix added
* secondary device notification list is same as for primary

## Example usage

```hcl
provider equinix {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

resource "equinix_network_acl_template" "pri-acl" {
  name        = "acl-pri"
  description = "Primary ACL template"
  metro_code  = "SV"
  inbound_rule {
    subnets  = ["10.141.30.0/24"]
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
}

resource "equinix_network_acl_template" "sec-acl" {
  name        = "acl-sec"
  description = "Secondary ACL template"
  metro_code  = "SY"
  inbound_rule {
    subnets  = ["10.141.50.0/24"]
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
}

module "silverpeak-sdwan" {
  source               = "../modules/silverpeak-sdwan"
  metro_code           = "SV"
  platform             = "small"
  software_package     = "EC-M"
  name                 = "tf-tst-silverpeak"
  hostname             = "tf-sp-pri"
  term_length          = 1
  notifications        = ["test@test.com"]
  acl_template_id      = equinix_network_acl_template.pri-acl.id
  account_name         = "myAccountName"
  account_key          = "myAccountKey"
  appliance_tag        = "myApplianceTag"
  additional_bandwidth = 25
  interface_count      = 32
  secondary = {
    enabled              = true
    hostname             = "tf-sp-sec"
    metro_code           = "SY"
    acl_template_id      = equinix_network_acl_template.sec-acl.id
    account_name         = "myAccountName"
    account_key          = "myAccountKey"
    additional_bandwidth = 25
    appliance_tag        = "myApplianceTagSecondary"
  }
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
|metro_code|Two-letter device location's metro code|`string`|`""`|yes|
|platform|Device hardware platform flavor: small, medium, large|`string`|`""`|yes|
|software_package|Device software package: EC-M, EC-L, EC-XL|`string`|`""`|yes|
|name|Device name|`string`|`""`|yes|
|hostname|Device hostname|`string`|`""`|yes|
|term_length|Term length in months: 1, 12, 24, 36|`number`|`0`|yes|
|notifications|List of email addresses that will receive notifications about device|`list(string)`|n/a|yes|
|interface_count|Number of network interfaces on a device: 10 or 32|`number`|`10`|no|
|additional_bandwidth|Amount of additional internet bandwidth for a device, in Mbps|`number`|`0`|no|
|acl_template_id|Identifier of a network ACL template that will be applied on a device|`string`|`""`|yes|
|account_name|Device account name|`string`|`""`|yes|
|account_key|Device account key|`string`|`""`|yes|
|appliance_tag|Device appliance tag|`string`|`""`|no|
|secondary|Map of secondary device attributes in redundant setup|`map`|N/A|no|

Secondary device map attributes:

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
|enabled|Value that determines if secondary device will be created|`bool`|`false`|no|
|metro_code|Two-letter secondary device location's metro code|`string`|`""`|yes|
|hostname|Device hostname|`string`|`""`|yes|
|additional_bandwidth|Amount of additional internet bandwidth for a secondary device, in Mbps|`number`|`0`|no|
|acl_template_id|Identifier of a network ACL template that will be applied on a secondary device|`string`|`""`|yes|
|account_name|Secondary device account name|`string`|`""`|yes|
|account_key|Secondary device account key|`string`|`""`|yes|
|appliance_tag|Secondary device appliance tag|`string`|`""`|no|

## Outputs

| Name | Description |
|------|-------------|
|id|Device identifier|
|status|Device provisioning status|
|license_status|Device license status|
|account_number|Device billing account number|
|cpu_count|Number of device CPU cores|
|memory|Amount of device memory|
|software_version|Device software version|
|region|Device region|
|ibx|Device IBX center code|
|ssh_ip_address|Device SSH interface IP address|
|ssh_ip_fqdn|Device SSH interface FQDN|
|interfaces|List of network interfaces present on a device|
|secondary|Secondary device outputs (same as for primary). Present when secondary device was enabled|
