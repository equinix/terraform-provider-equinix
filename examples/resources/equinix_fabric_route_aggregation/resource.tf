resource "equinix_fabric_route_aggregation" "new-ra" {
  type = "BGP_IPv4_PREFIX_AGGREGATION"
  name = "new-ra"
  description = "Test aggregation"
  project = {
    project_id = "776847000642406"
  }
}
