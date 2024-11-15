data "equinix_fabric_service_token" "service-token" {
  uuid = "<uuid_of_service_token>"
}

output "id" {
  value = data.equinix_fabric_service_token.service-token.id
}

output "type" {
  value = data.equinix_fabric_service_token.service-token.type
}

output "expiration_date_time" {
  value = data.equinix_fabric_service_token.service-token.expiration_date_time
}

output "supported_bandwidths" {
  value = data.equinix_fabric_service_token.service-token.service_token_connection.0.supported_bandwidths
}

output "virtual_device_type" {
  value = data.equinix_fabric_service_token.service-token.service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.type
}

output "virtual_device_uuid" {
  value = data.equinix_fabric_service_token.service-token.service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.uuid
}

output "interface_type" {
  value = data.equinix_fabric_service_token.service-token.service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.type
}

output "interface_uuid" {
  value = data.equinix_fabric_service_token.service-token.service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.id
}
