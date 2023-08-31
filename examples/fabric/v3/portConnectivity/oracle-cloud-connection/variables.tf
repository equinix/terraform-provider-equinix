variable "equinix_client_id" {}

variable "equinix_client_secret" {}

variable "equinix_port_name" {}

variable "oci_tenancy_ocid" {}

variable "oci_user_ocid" {}

variable "oci_fingerprint" {}

variable "oci_private_key_path" {}


variable "oci_compartment_id" {}

variable "oci_fastconnect_provider" {
  default = "Equinix"
}

variable "oci_region" {}

variable "oci_metro_code" {}
