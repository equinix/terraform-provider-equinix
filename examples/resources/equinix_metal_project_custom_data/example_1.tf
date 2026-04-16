locals {
  project_id = "<UUID_of_your_project>"
}

resource "equinix_metal_project_custom_data" "project_metadata" {
  project_id = local.project_id
  custom_data = jsonencode({
    owner   = "platform-team"
    service = "edge-gateway"
    env     = "prod"
  })
}
