data "equinix_fabric_stream_alert_rules" "data_stream_alert_rules" {
  stream_id = "<uuid_of_stream>"
  pagination = {
    limit = 5
    offset = 1
  }
}

output "stream_alert_rules_type" {
  value = data.equinix_fabric_stream_alert_rules.alert_rules.data[0].type
}

output "stream_alert_rules_id" {
  value = data.equinix_fabric_stream_alert_rules.alert_rules.data[0].uuid
}

output "stream_alert_rules_state" {
  value = data.equinix_fabric_stream_alert_rules.alert_rules.data[0].state
}

output "stream_alert_rules_stream_id" {
  value = data.equinix_fabric_stream_alert_rules.alert_rules.data[0].stream_id
}
