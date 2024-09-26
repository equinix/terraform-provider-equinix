data "equinix_fabric_route_filter_rule" "rf_rule" {
  route_filter_id = "<route_filter_policy_id>"
  uuid = "<route_filter_rule_uuid>"
}

output "route_filter_rule_name" {
  value = data.equinix_fabric_route_filter_rule.rf_rule.name
}


output "route_filter_rule_description" {
  value = data.equinix_fabric_route_filter_rule.rf_rule.description
}

output "route_filter_rule_prefix" {
  value = data.equinix_fabric_route_filter_rule.rf_rule.prefix
}

output "route_filter_rule_prefix_match" {
  value = data.equinix_fabric_route_filter_rule.rf_rule.prefix_match
}
