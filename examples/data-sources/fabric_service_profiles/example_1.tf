data "equinix_fabric_service_profiles" "service_profiles_data_name" {
  filter {
    property = "/name"
    operator = "="
    values   = ["<list_of_profiles_to_return>"]
  }
}

output "id" {
  value = data.equinix_fabric_service_profile.service_profiles_data_name.data.0.id
}

output "name" {
  value = data.equinix_fabric_service_profile.service_profiles_data_name.data.0.name
}

output "type" {
  value = data.equinix_fabric_service_profile.service_profiles_data_name.data.0.type
}

output "visibility" {
  value = data.equinix_fabric_service_profile.service_profiles_data_name.data.0.visibility
}

output "org_name" {
  value = data.equinix_fabric_service_profile.service_profiles_data_name.data.0.account.0.organization_name
}

output "access_point_type_configs_type" {
  value = data.equinix_fabric_service_profile.service_profiles_data_name.data.0.access_point_type_configs.0.type
}

output "allow_remote_connections" {
  value = data.equinix_fabric_service_profile.service_profiles_data_name.data.0.access_point_type_configs.0.allow_remote_connections
}

output "supported_bandwidth_0" {
  value = data.equinix_fabric_service_profile.service_profiles_data_name.data.0.access_point_type_configs.0.supported_bandwidths.0
}

output "supported_bandwidth_1" {
  value = data.equinix_fabric_service_profile.service_profiles_data_name.data.0.access_point_type_configs.0.supported_bandwidths.1
}

output "redundandy_required" {
  value = data.equinix_fabric_service_profile.service_profiles_data_name.data.0.access_point_type_configs.0.connection_redundancy_required
}

output "allow_over_subscription" {
  value = data.equinix_fabric_service_profile.service_profiles_data_name.data.0.access_point_type_configs.0.api_config.0.allow_over_subscription
}
