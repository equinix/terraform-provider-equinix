# List CIDR expressions of all the allocated IP block in you project.

# Declare your project ID
locals {
  project_id = "<UUID_of_your_project>"
}

data "equinix_metal_ip_block_ranges" "test" {
  project_id = local.project_id
}

output "out" {
  value = data.equinix_metal_ip_block_ranges.test
}
