# Equinix Network Edge: VMWare SD-WAN Edge module

A Terraform module to create VMWare SD-WAN Edge network device
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

module "vmware-sdwan" {
  source           = "../modules/vmware-sdwan"
  metro_code       = "SV"
  platform         = "small"
  software_package = "VeloCloud-2"
  name             = "tf-tst-velocloud"
  term_length      = 1
  notifications    = ["test@test.com"]
  acl_template_id  = equinix_network_acl_template.pri-acl.id
  activation_key   = "AAAA-BBBB-CCCC-DDDD"
  controller_fqdn  = "test.com"
  root_password    = "Gy3ehr38R5"
  secondary = {
    enabled         = true
    metro_code      = "SY"
    acl_template_id = equinix_network_acl_template.sec-acl.id
    activation_key  = "DDDD-EEEE-FFFF-GGGG"
    controller_fqdn = "test2.com"
    root_password   = "Gy3ehr38R5"
  }
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
|metro_code|Two-letter device location's metro code|`string`|`""`|yes|
|platform|Device hardware platform flavor: small, medium, large|`string`|`""`|yes|
|software_package|Device software package: VeloCloud-2, VeloCloud-4, VeloCloud-8|`string`|`""`|yes|
|name|Device name|`string`|`""`|yes|
|term_length|Term length in months: 1, 12, 24, 36|`number`|`0`|yes|
|notifications|List of email addresses that will receive notifications about device|`list(string)`|n/a|yes|
|additional_bandwidth|Amount of additional internet bandwidth for a device, in Mbps|`number`|`0`|no|
|acl_template_id|Identifier of a network ACL template that will be applied on a device|`string`|`""`|yes|
|activation_key|Device activation key|`string`|`""`|yes|
|controller_fqdn|FQDN of a primary controller|`string`|`""`|yes|
|root_password|Device root password|`string`|`""`|no|
|secondary|Map of secondary device attributes in redundant setup|`map`|N/A|no|

Secondary device map attributes:

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
|enabled|Value that determines if secondary device will be created|`bool`|`false`|no|
|metro_code|Two-letter secondary device location's metro code|`string`|`""`|yes|
|additional_bandwidth|Amount of additional internet bandwidth for a secondary device, in Mbps|`number`|`0`|no|
|acl_template_id|Identifier of a network ACL template that will be applied on a secondary device|`string`|`""`|yes|
|activation_key|Secondary device activation key|`string`|`""`|yes|
|controller_fqdn|FQDN of a secondary controller|`string`|`""`|yes|
|root_password|Secondary device root password|`string`|`""`|no|

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
