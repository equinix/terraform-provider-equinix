resource "equinix_fabric_route_filter" "rf_policy" {
  name = "RF_Policy_Name",
  project {
    projectId = "<project_id>"
  },
  type = "BGP_IPv4_PREFIX_FILTER",
  description = "Route Filter Policy for X Purpose",
}

output "id" {
  value = equinix_fabric_route_filter.rf_policy.id
}

output "type" {
  value = equinix_fabric_route_filter.rf_policy.type
}

output "state" {
  value = equinix_fabric_route_filter.rf_policy.state
}

output "not_matched_rules_action" {
  value = equinix_fabric_route_filter.rf_policy.not_matched_rule_action
}

output "connections_count" {
  value = equinix_fabric_route_filter.rf_policy.connections_count
}

output "rules_count" {
  value = equinix_fabric_route_filter.rf_policy.rules_count
}
