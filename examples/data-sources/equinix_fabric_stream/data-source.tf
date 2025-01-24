data "equinix_fabric_stream" "data_stream" {
  stream_id = "<uuid_of_stream>"
}

output "stream_state" {
  value = data.equinix_fabric_stream.data_stream.state
}
