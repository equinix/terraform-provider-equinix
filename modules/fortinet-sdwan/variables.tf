variable "metro_code" {
  description = "Metro location code for a network device"
  type        = string
  validation {
    condition     = can(regex("^[A-Z]{2}$", var.metro_code))
    error_message = "Valid metro code consits of two capital leters, i.e. SV, DC."
  }
}

variable "platform" {
  description = "Device platform flavor"
  type        = string
  validation {
    condition     = can(regex("^(small|medium|large)$", var.platform))
    error_message = "One of following platform flavors are supported: small, medium, large."
  }
}

variable "software_package" {
  description = "Device software package"
  type        = string
  validation {
    condition     = can(regex("^(VM02|VM04|VM08)$", var.software_package))
    error_message = "One of following software packages are supported: VM02, VM04, VM08."
  }
}

variable "name" {
  description = "Device name"
  type        = string
  validation {
    condition     = length(var.name) > 2 && length(var.name) < 51
    error_message = "Device name should consist of 3 to 50 characters."
  }
}

variable "hostname" {
  description = "Device hostname"
  type        = string
  validation {
    condition     = length(var.hostname) > 4 && length(var.hostname) < 51
    error_message = "Device name should consist of 5 to 50 characters."
  }
}

variable "term_length" {
  description = "Term length in months"
  type        = number
  validation {
    condition     = can(regex("^(1|12|24|36)$", var.term_length))
    error_message = "One of following term lengths are available: 1, 12, 24, 36 months."
  }
}

variable "notifications" {
  description = "List of email addresses that will be used to send notifications about device"
  type        = list(string)
  validation {
    condition     = length(var.notifications) > 0
    error_message = "Notification list cannot be empty."
  }
}

variable "license_file" {
  description = "Path to device license file"
  type        = string
  default     = ""
}

variable "acl_template_id" {
  description = "Identifier of an ACL template that will be allied on a device"
  type        = string
  validation {
    condition     = length(var.acl_template_id) > 0
    error_message = "ACL template identifier has to be non empty string."
  }
}

variable "additional_bandwidth" {
  description = "Additional internet bandwidth for a device"
  type        = number
  default     = 0
  validation {
    condition     = var.additional_bandwidth == 0 || (var.additional_bandwidth >= 25 && var.additional_bandwidth <= 2001)
    error_message = "Additional internet bandwidth should be between 25 and 2001 Mbps."
  }
}

variable "controller_ip_address" {
  description = "SDWAN controller IP address"
  type        = string
  validation {
    condition     = can(regex("^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$", var.controller_ip_address))
    error_message = "Controller's IP address has to be valid IP address."
  }
}

variable "admin_password" {
  description = "Device admin password"
  type        = string
  validation {
    condition     = length(var.admin_password) > 5 && length(var.admin_password) < 129
    error_message = "Admin password should consist of 6 to 128 characters."
  }
}

variable "secondary" {
  description = "Secondary device attributes"
  type        = map
  default = {
    enabled = false
  }
  validation {
    condition     = can(var.secondary.enabled)
    error_message = "Key enabled has to be defined if secondary block is present."
  }
  validation {
    condition     = ! try(var.secondary.enabled, false) || try(length(var.secondary.hostname) > 4 && length(var.secondary.hostname) < 51, false)
    error_message = "Key hostname should have from 5 to 50 characters."
  }
  validation {
    condition     = ! try(var.secondary.enabled, false) || (can(regex("^[A-Z]{2}$", var.secondary.metro_code)))
    error_message = "Key metro_code has to be defined and consit of two capital leters, i.e. SV, DC."
  }
  validation {
    condition     = ! try(var.secondary.enabled, false) || try(length(var.secondary.license_file) > 0, true)
    error_message = "Key license_file has to be defined as non empty string."
  }
  validation {
    condition     = ! try(var.secondary.enabled, false) || try(length(var.secondary.acl_template_id) > 0, false)
    error_message = "Key acl_template_id has to be defined as non empty string."
  }
  validation {
    condition     = ! try(var.secondary.enabled, false) || try(var.secondary.additional_bandwidth >= 25 && var.secondary.additional_bandwidth <= 2001, true)
    error_message = "Key additional_bandwidth has to be between 25 and 2001 Mbps."
  }
  validation {
    condition     = ! try(var.secondary.enabled, false) || can(regex("^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$", var.secondary.controller_ip_address))
    error_message = "Key controller_ip_address has to be valid IP address."
  }
  validation {
    condition     = ! try(var.secondary.enabled, false) || try(length(var.secondary.admin_password) > 5 && length(var.secondary.admin_password) < 129, false)
    error_message = "Key admin_password should consist of 6 to 128 characters."
  }
}
