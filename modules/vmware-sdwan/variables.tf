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
    condition     = can(regex("^(VeloCloud-2|VeloCloud-4|VeloCloud-8)$", var.software_package))
    error_message = "One of following software packages are supported: VeloCloud-2, VeloCloud-4, VeloCloud-8."
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

variable "activation_key" {
  description = "Device activation key"
  type        = string
  validation {
    condition     = can(regex("^[A-Z]{4}-[A-Z]{4}-[A-Z]{4}-[A-Z]{4}", var.activation_key))
    error_message = "Activation key should consist of four groups of four uppercase letters sepparated with dashes, i.e. AVWD-MMWE-WLEP-KKWM."
  }
}

variable "controller_fqdn" {
  description = "FQDN of a primary controller"
  type        = string
}

variable "root_password" {
  description = "Device root password"
  type        = string
  default     = ""
  validation {
    condition     = can(regex("^()|(.{8,128})$", var.root_password))
    error_message = "Device roout password should consist of at least 8 characters and must include an uppercase letter, a number and a special character."
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
    condition     = ! try(var.secondary.enabled, false) || (can(regex("^[A-Z]{2}$", var.secondary.metro_code)))
    error_message = "Key metro_code has to be defined and consit of two capital leters, i.e. SV, DC."
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
    condition     = ! try(var.secondary.enabled, false) || can(regex("^[A-Z]{4}-[A-Z]{4}-[A-Z]{4}-[A-Z]{4}", var.secondary.activation_key))
    error_message = "Key activation_key has to be declared and consist of four groups of four uppercase letters sepparated with dashes, i.e. AVWD-MMWE-WLEP-KKWM."
  }
  validation {
    condition     = ! try(var.secondary.enabled, false) || try(length(var.secondary.controller_fqdn) > 0, false)
    error_message = "Key controller_fqdn has to be defined as non empty string."
  }
  validation {
    condition     = ! try(var.secondary.enabled, false) || try(length(var.secondary.root_password) > 7 && length(var.secondary.root_password) < 129, true)
    error_message = "Key root_password has to be from 8 to 128 characters long."
  }
}
