# Get Project SSH Key by name
data "equinix_metal_project_ssh_key" "my_key" {
  search     = "username@hostname"
  project_id = local.project_id
}
