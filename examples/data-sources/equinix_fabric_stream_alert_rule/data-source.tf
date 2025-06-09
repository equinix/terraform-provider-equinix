data "equinix_fabric_stream_alert_rule" "data_stream_alert_rule" {
  stream_id = "<uuid_of_stream>"
  stream_alert_rule_id = "<uuid_of_stream_alert_rule>"
}

output "stream_alert_rule_type" {
  value = data.equinix_fabric_stream_alert_rule.alert_rule.type
}

output "stream_alert_rule_id" {
  value = data.equinix_fabric_stream_alert_rule.alert_rule.uuid
}

output "stream_alert_rule_state" {
  value = data.equinix_fabric_stream_alert_rule.alert_rule.state
}

output "stream_alert_rule_stream_id" {
  value = data.equinix_fabric_stream_alert_rule.alert_rule.stream_id
}