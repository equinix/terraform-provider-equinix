# Create self configured single Aviatrix device with cloud init file

data "equinix_network_account" "sv" {
  metro_code = "SV"
}

variable "filepath" { default = "cloudInitFileFolder/TF-AVX-cloud-init-file.txt" }

resource "equinix_network_file" "aviatrix-cloudinit-file" {
  file_name = "TF-AVX-cloud-init-file.txt"
  content = file("${path.module}/${var.filepath}")
  metro_code = data.equinix_network_account.sv.metro_code
  device_type_code = "AVIATRIX_EDGE"
  process_type = "CLOUD_INIT"
  self_managed = true
  byol = true
}

resource "equinix_network_device" "aviatrix-single" {
  name            = "tf-aviatrix"
  metro_code      = data.equinix_network_account.sv.metro_code
  type_code       = "AVIATRIX_EDGE"
  self_managed    = true
  byol            = true
  package_code    = "STD"
  notifications   = ["john@equinix.com"]
  term_length     = 12
  account_number  = data.equinix_network_account.sv.number
  version         = "6.9"
  core_count      = 2
  cloud_init_file_id = equinix_network_file.aviatrix-cloudinit-file.uuid
  acl_template_id = "c06150ea-b604-4ad1-832a-d63936e9b938"
}
