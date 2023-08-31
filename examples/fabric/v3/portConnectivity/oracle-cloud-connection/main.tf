provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

provider "oci" {
  tenancy_ocid     = var.oci_tenancy_ocid
  user_ocid        = var.oci_user_ocid
  fingerprint      = var.oci_fingerprint
  private_key_path = var.oci_private_key_path
  region           = var.oci_region
}

data "oci_core_fast_connect_provider_services" "fc_provider_services" {
  compartment_id = var.oci_compartment_id
}

locals {
  fc_provider_services_id = element(
    data.oci_core_fast_connect_provider_services.fc_provider_services.fast_connect_provider_services,
    index(
      data.oci_core_fast_connect_provider_services.fc_provider_services.fast_connect_provider_services.*.provider_name,
      var.oci_fastconnect_provider
    )
  ).id
}

resource "oci_core_drg" "oci_drg" {
  display_name   = "TestOCIDRG"
  compartment_id = var.oci_compartment_id
}

resource "oci_core_virtual_circuit" "oci_virtual_circuit" {
  type                 = "PRIVATE"
  compartment_id       = var.oci_compartment_id
  bandwidth_shape_name = "1 Gbps"

  cross_connect_mappings {
    customer_bgp_peering_ip = "10.1.0.50/30"
    oracle_bgp_peering_ip   = "10.1.0.49/30"
  }
  customer_asn        = "12234"
  display_name        = "TestOCIVC"
  gateway_id          = oci_core_drg.oci_drg.id //The OCID of the dynamic routing gateway (DRG) that this virtual circuit uses
  provider_service_id = local.fc_provider_services_id
  region              = var.oci_region
}

data "equinix_ecx_l2_sellerprofile" "oracle" {
  name                     = "Oracle Cloud Infrastructure -OCI- FastConnect"
  organization_global_name = "ORACLE"
}

data "equinix_ecx_port" "dot1q-1-pri" {
  name = var.equinix_port_name
}

resource "equinix_ecx_l2_connection" "oci-dot1q" {
  name              = "tf-oci-dot1q"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.oci.uuid
  speed             = 1
  speed_unit        = "GB"
  notifications     = ["example@equinix.com"]
  port_uuid         = data.equinix_ecx_port.dot1q-1-pri.uuid
  vlan_stag         = 3030
  seller_region     = var.oci_region
  seller_metro_code = var.oci_metro_code
  authorization_key = oci_core_virtual_circuit.oci_virtual_circuit.id //Virtual Circuit OCID
}
