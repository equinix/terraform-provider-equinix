# Create self configured redundant BlueCat Edge Service Point
data "equinix_network_account" "sv" {
  name = "account-name"
  metro_code = "SV"
}

resource "equinix_network_file" "bluecat-edge-service-point-cloudinit-primary-file" {
  file_name = "TF-BLUECAT-ESP-cloud-init-file.txt"
  content = file("${path.module}/${var.filepath}")
  metro_code = data.equinix_network_account.sv.metro_code
  device_type_code = "BLUECAT-EDGE-SERVICE-POINT"
  process_type = "CLOUD_INIT"
  self_managed = true
  byol = true
}

resource "equinix_network_file" "bluecat-edge-service-point-cloudinit-secondary-file" {
  file_name = "TF-BLUECAT-ESP-cloud-init-file.txt"
  content = file("${path.module}/${var.filepath}")
  metro_code = data.equinix_network_account.sv.metro_code
  device_type_code = "BLUECAT-EDGE-SERVICE-POINT"
  process_type = "CLOUD_INIT"
  self_managed = true
  byol = true
}

resource "equinix_network_device" "bluecat-edge-service-point-ha" {
  name            = "tf-bluecat-edge-service-point-p"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "BLUECAT-EDGE-SERVICE-POINT"
  self_managed    = true
  connectivity    = "PRIVATE"
  byol            = true
  package_code    = "STD"
  notifications   = ["test@equinix.com"]
  account_number  = data.equinix_network_account.sv.number
  cloud_init_file_id = equinix_network_file.bluecat-edge-service-point-cloudinit-primary-file.uuid
  version         = "4.6.3"
  core_count      = 4
  term_length     = 12
  secondary_device {
    name            = "tf-bluecat-edge-service-point-s"
    metro_code      = data.equinix_network_account.sv.metro_code
    notifications   = ["test@eq.com"]
    account_number  = data.equinix_network_account.sv.number
    cloud_init_file_id = equinix_network_file.bluecat-edge-service-point-cloudinit-secondary-file.uuid
  }
}
