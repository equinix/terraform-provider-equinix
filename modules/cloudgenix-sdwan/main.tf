locals {
  metro_codes = ! var.secondary.enabled || (var.metro_code == try(var.secondary.metro_code, "")) ? toset([var.metro_code]) : toset([var.metro_code, var.secondary.metro_code])
}

data "equinix_network_account" "this" {
  for_each   = local.metro_codes
  metro_code = each.key
  status     = "Active"
}

data "equinix_network_device_type" "this" {
  category    = "SDWAN"
  vendor      = "CloudGenix"
  metro_codes = local.metro_codes
}

data "equinix_network_device_platform" "this" {
  device_type = data.equinix_network_device_type.this.code
  flavor      = var.platform
}

data "equinix_network_device_software" "this" {
  device_type = data.equinix_network_device_type.this.code
  packages    = [var.software_package]
  stable      = true
  most_recent = true
}

resource "equinix_network_device" "this" {
  self_managed         = true
  byol                 = true
  name                 = var.name
  type_code            = data.equinix_network_device_type.this.code
  package_code         = var.software_package
  version              = data.equinix_network_device_software.this.version
  core_count           = data.equinix_network_device_platform.this.core_count
  metro_code           = var.metro_code
  account_number       = data.equinix_network_account.this[var.metro_code].number
  term_length          = var.term_length
  notifications        = var.notifications
  acl_template_id      = var.acl_template_id
  additional_bandwidth = var.additional_bandwidth > 0 ? var.additional_bandwidth : null
  vendor_configuration = {
    licenseKey    = var.license_key
    licenseSecret = var.license_secret
  }
  dynamic "secondary_device" {
    for_each = try(var.secondary.enabled, false) ? [1] : []
    content {
      name                 = "${var.name}-secondary"
      metro_code           = var.secondary.metro_code
      account_number       = data.equinix_network_account.this[var.secondary.metro_code].number
      notifications        = var.notifications
      acl_template_id      = var.secondary.acl_template_id
      additional_bandwidth = try(var.secondary.additional_bandwidth, null)
      vendor_configuration = {
        licenseKey    = var.secondary.license_key
        licenseSecret = var.secondary.license_secret
      }
    }
  }
}


