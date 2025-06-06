data "equinix_fabric_stream_alert_rules" "data_stream_alert_rules" {
  stream_id = "<uuid_of_stream>"
  pagination = {
    limit = 5
    offset = 1
  }
}

output "number_of_returned_stream_alert_rules" {
  value = length(data.equinix_fabric_stream_alert_rules.data_stream_alert_rules.data)
}
