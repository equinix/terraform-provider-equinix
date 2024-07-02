# Create a new read-only API key in existing project
resource "equinix_metal_project_api_key" "test" {
  project_id  = local.existing_project_id
  description = "Read-only key scoped to a projct"
  read_only   = true
}
