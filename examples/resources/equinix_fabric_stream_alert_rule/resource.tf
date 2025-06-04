resource "equinix_fabric_stream_alert_rule" "new_stream_alert_rule" {
  stream_id          = "<stream_id>"
  name               = "<name>"
  type               = "METRIC_ALERT"
  description        = "<description>"
  enabled            = true
  operand            = "ABOVE"
  window_size        = "<window_size>"
  warning_threshold  = "<warning_threshold>"
  critical_threshold = "<critical_threshold>"
  metric_name        = "equinix.fabric.connection.bandwidth_tx.usage"
  resource_selector = {
    include = ["*/connections/<connection_id>"]
  }
}