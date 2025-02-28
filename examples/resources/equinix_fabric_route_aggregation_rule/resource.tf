resource "equinix_fabric_route_aggregation_rule" "ra_rule" {
  route_aggregation_id = "<route_aggregation_id>"
  name = "ra-rule-test"
  description = "Route aggregation rule"
  prefix = "192.168.0.0/24"
}

output "route_aggregation_rule_name" {
  value = equinix_fabric_route_aggregation_rule.ra_rule.name
}

output "route_aggregation_rule_description" {
  value = equinix_fabric_route_aggregation_rule.ra_rule.description
}

output "route_aggregation_rule_type" {
  value = equinix_fabric_route_aggregation_rule.ra_rule.type
}

output "route_aggregation_rule_prefix" {
  value = equinix_fabric_route_aggregation_rule.ra_rule.prefix
}

output "route_aggregation_rule_state" {
  value = equinix_fabric_route_aggregation_rule.ra_rule.state
}
