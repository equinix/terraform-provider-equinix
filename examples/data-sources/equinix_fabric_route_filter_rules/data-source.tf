data "equinix_fabric_route_filter_rules" "rf_rules" {
  route_filter_id = "<route_filter_policy_id"
  limit = 100
  offset = 5
}

output "first_route_filter_rule_name" {
  value = data.equinix_fabric_route_filter_rules.rf_rules.data.0.name
}

output "first_route_filter_rule_description" {
  value = data.equinix_fabric_route_filter_rules.rf_rules.data.0.description
}

output "first_route_filter_rule_prefix" {
  value = data.equinix_fabric_route_filter_rules.rf_rules.data.0.prefix
}

output "first_route_filter_rule_prefix_match" {
  value = data.equinix_fabric_route_filter_rules.rf_rules.data.0.prefix_match
}
