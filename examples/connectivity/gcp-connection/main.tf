provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

provider "google" {
  region  = var.gcp_region
  project = var.gcp_project_name
}

data "equinix_ecx_port" "dot1q-pri" {
  name = var.equinix_port_name
}

data "equinix_ecx_l2_sellerprofile" "gcpi-1" {
  name                     = "Google Cloud Partner Interconnect Zone 1"
  organization_global_name = "GOOGLE"
}

resource "google_compute_network" "test" {
  name = "tf-test"
}

resource "google_compute_router" "test" {
  name    = "tf-test"
  network = google_compute_network.test.name
}

resource "google_compute_interconnect_attachment" "test-1" {
  name                     = "tf-test"
  type                     = "PARTNER"
  router                   = google_compute_router.test.id
  region                   = var.gcp_region
  edge_availability_domain = "AVAILABILITY_DOMAIN_1"
}

resource "equinix_ecx_l2_connection" "gcpi-dot1q" {
  name              = "tf-gcpi-dot1q"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.gcpi-1.uuid
  speed             = 50
  speed_unit        = "MB"
  notifications     = ["example@equinix.com"]
  port_uuid         = data.equinix_ecx_port.dot1q-pri.uuid
  vlan_stag         = 1600
  seller_region     = google_compute_interconnect_attachment.test-1.region
  seller_metro_code = var.gcp_metro_code
  authorization_key = google_compute_interconnect_attachment.test-1.pairing_key
}
