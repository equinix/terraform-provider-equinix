terraform {
  required_version = ">= 0.13"

  required_providers {
    equinix = {
      source  = "developer.equinix.com/terraform/equinix"
      version = ">= 9.0.0"
    }
  }
}
