data "equinix_fabric_route_filter" "rf_policy" {
  uuid = "<uuid_of_route_filter"
}

output "id" {
  value = data.equinix_fabric_route_filter.rf_policy.id
}

output "type" {
  value = data.equinix_fabric_route_filter.rf_policy.type
}

output "state" {
  value = data.equinix_fabric_route_filter.rf_policy.state
}

output "not_matched_rules_action" {
  value = data.equinix_fabric_route_filter.rf_policy.not_matched_rule_action
}

output "connections_count" {
  value = data.equinix_fabric_route_filter.rf_policy.connections_count
}

output "rules_count" {
  value = data.equinix_fabric_route_filter.rf_policy.rules_count
}
