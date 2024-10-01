resource "equinix_fabric_route_filter_rule" "rf_rule" {
  route_filter_id = "<route_filter_policy_id>"
  name = "Route Filter Rule Name"
  prefix = "192.168.0.0/24"
  prefix_match = "exact"
  description = "Route Filter Rule for X Purpose"
}

output "route_filter_rule_id" {
  value = equinix_fabric_route_filter_rule.rf_rule.id
}

output "route_filter_id" {
  value = equinix_fabric_route_filter_rule.rf_rule.route_filter_id
}

output "route_filter_rule_prefix" {
  value = equinix_fabric_route_filter_rule.rf_rule.prefix
}

output "route_filter_rule_prefix_match" {
  value = equinix_fabric_route_filter_rule.rf_rule.prefix_match
}
