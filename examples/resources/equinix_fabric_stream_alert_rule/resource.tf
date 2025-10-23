resource "equinix_fabric_stream_alert_rule" "new_stream_alert_rule" {
  stream_id          = "<stream_id>"
  name               = "<name>"
  type               = "METRIC_ALERT"
  description        = "<description>"
  enabled            = true
  metric_selector = {
    include = ["equinix.fabric.connection.bandwidth_tx.usage"]
  }
  detection_method = {
    operand            = "ABOVE"
    window_size        = "<window_size>"
    warning_threshold  = "<warning_threshold>"
    critical_threshold = "<critical_threshold>"
  }
  resource_selector = {
    include = ["*/connections/<connection_id>"]
  }
}

output "stream_alert_rule_type" {
  value = equinix_fabric_stream_alert_rule.new_stream_alert_rule.type
}

output "stream_alert_rule_id" {
  value = equinix_fabric_stream_alert_rule.new_stream_alert_rule.uuid
}

output "stream_alert_rule_stream_id" {
  value = equinix_fabric_stream_alert_rule.new_stream_alert_rule.stream_id
}

output "stream_alert_rule_state" {
  value = equinix_fabric_stream_alert_rule.new_stream_alert_rule.state
}
