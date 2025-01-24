resource "equinix_fabric_stream" "new_stream" {
  type = "TELEMETRY_STREAM"
  name = "<name_of_stream_resource>"
  description = "<description_of_stream_resource>"
  project = {
    project_id = "<destination_project_id_for_stream"
  }
}

output "stream_state" {
  value = equinix_fabric_stream.new_stream.state
}
