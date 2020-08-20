provider "equinix" {
  client_id     = "your_client_id"
  client_secret = "your_client_secret"
}

provider "google" {
  credentials = file("account.json")
  project     = "my-project-id"
}

data "equinix_ecx_port" "dot1q-1-pri" {
  name = "sit-001-CX-DC5-NL-Dot1q-BO-10G-PRI-JUN-27"
}

data "equinix_ecx_l2_sellerprofile" "gcpi-1" {
  name = "Google Cloud Partner Interconnect Zone 1"
}

resource "google_compute_router" "test" {
  name    = "tf-test"
  network = google_compute_network.foobar.name
}

resource "google_compute_interconnect_attachment" "test-1" {
  name                     = "tf-test"
  type                     = "PARTNER"
  router                   = google_compute_router.test.id
  region                   = "us-west1"
  edge_availability_domain = "AVAILABILITY_DOMAIN_1"
}

resource "equinix_ecx_l2_connection" "gcpi-dot1q" {
  name                  = "tf-gcpi-dot1q"
  profile_uuid          = data.equinix_ecx_l2_sellerprofile.gcpi-1.uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = data.equinix_ecx_port.dot1q-1-pri.uuid
  vlan_stag             = 1600
  seller_region         = google_compute_interconnect_attachment.test-1.region
  seller_metro_code     = "SV"
  authorization_key     = google_compute_interconnect_attachment.test-1.pairing_key
}
