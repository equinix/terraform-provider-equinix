data "equinix_fabric_stream_alert_rule" "data_stream_alert_rule" {
  stream_id = "<uuid_of_stream>"
  stream_alert_rule_id = "<uuid_of_stream_alert_rule>"
}

output "stream_alert_rule_state" {
  value = data.equinix_fabric_stream_alert_rule.data_stream_alert_rule.state
}
