variable "equinix_client_id" {
  type        = string
  description = "API Consumer Key available under 'My Apps' in developer portal."
}

variable "equinix_client_secret" {
  type        = string
  description = "API Consumer secret available under 'My Apps' in developer portal."
}

variable "metal_auth_token" {
  type        = string
  description = "This is your Equinix Metal API Auth token."
}

variable "metal_project_name" {
  type        = string
  description = "Name of the project where the connection is scoped to."
}

variable "connection_name" {
  type        = string
  description = "Name of the connection resource that will be created in both Equinix Metal and Equinix Fabric services."
  default     = "metal-to-sp"
}

variable "connection_metro" {
  type        = string
  description = "Metro where the connection will be created."
  default     = "SV"

  validation {
    condition     = can(regex("^[A-Z]{2}$", var.connection_metro))
    error_message = "Valid metro code consits of two capital leters, i.e. SV, DC."
  }
}

variable "connection_description" {
  type        = string
  description = "Description for the connection resource."
  default     = "Connect from Equinix Metal to Service provider using a-side token."
}

variable "connection_speed" {
  type        = number
  description = "Speed/Bandwidth to be allocated to the connection - (MB or GB). "
  default     = 50
}

variable "connection_speed_unit" {
  type        = string
  description = "Unit of the speed/bandwidth to be allocated to the connection."
  default     = "MB"

  validation {
    condition = contains(["MB", "GB"], var.connection_speed_unit)
    error_message = "Invalid speed unit. Required one of:  MB, GB."
  }
}

variable "connection_notification_users" {
  type        = list(string)
  description = "A list of email addresses used for sending connection update notifications."
  validation {
    condition     = length(var.connection_notification_users) > 0
    error_message = "Notification list cannot be empty."
  }
}

variable "seller_profile_name" {
  type        = string
  description = "Name of the service provider to connect with, i.e. 'AWS Direct Connect'"
}

variable "seller_authorization_key" {
  type        = string
  description = "Text field used to authorize connection on the provider side. Value depends on a provider service profile used for connection"
}

variable "seller_region" {
  type        = string
  description = "The region code in which the seller port resides. i.e. 'us-west-1'"
}
