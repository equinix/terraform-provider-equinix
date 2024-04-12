data "equinix_fabric_service_profile" "service_profile_data_name" {
  uuid = "<uuid_of_service_profile>"
}

output "id" {
  value = data.equinix_fabric_service_profile.service_profile_data_name.id
}

output "name" {
  value = data.equinix_fabric_service_profile.service_profile_data_name.name
}

output "type" {
  value = data.equinix_fabric_service_profile.service_profile_data_name.type
}

output "visibility" {
  value = data.equinix_fabric_service_profile.service_profile_data_name.visibility
}

output "org_name" {
  value = data.equinix_fabric_service_profile.service_profile_data_name.account.0.organization_name
}

output "access_point_type_configs_type" {
  value = data.equinix_fabric_service_profile.service_profile_data_name.access_point_type_configs.0.type
}

output "allow_remote_connections" {
  value = data.equinix_fabric_service_profile.service_profile_data_name.access_point_type_configs.0.allow_remote_connections
}

output "supported_bandwidth_0" {
  value = data.equinix_fabric_service_profile.service_profile_data_name.access_point_type_configs.0.supported_bandwidths.0
}

output "supported_bandwidth_1" {
  value = data.equinix_fabric_service_profile.service_profile_data_name.access_point_type_configs.0.supported_bandwidths.1
}

output "redundandy_required" {
  value = data.equinix_fabric_service_profile.service_profile_data_name.access_point_type_configs.0.connection_redundancy_required
}

output "allow_over_subscription" {
  value = data.equinix_fabric_service_profile.service_profile_data_name.access_point_type_configs.0.api_config.0.allow_over_subscription
}
