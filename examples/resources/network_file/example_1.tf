variable "filepath" { default = "fileFolder/fileName.txt" }

resource "equinix_network_file" "test-file" {
  file_name = "fileName.txt"
  content = file("${path.module}/${var.filepath}")
  metro_code = "SV"
  device_type_code = "AVIATRIX_EDGE"
  process_type = "CLOUD_INIT"
  self_managed = true
  byol = true
}
