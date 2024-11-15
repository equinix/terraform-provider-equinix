data "equinix_fabric_service_tokens" "service-tokens" {
  filter {
    property = "/type"
    operator = "="
    values 	 = "EVPL_VC"
  }
  filter {
    property = "/state"
    operator = "="
    values 	 = ["INACTIVE"]
  }
  pagination {
    offset = 0
    limit = 5
    total = 25
  }
}

output "number_of_returned_service_tokens" {
  value = length(data.equinix_fabric_service_tokens.service-tokens.data)
}

output "first_service_token_id" {
  value = data.equinix_fabric_service_tokens.service-tokens.data.0.id
}

output "first_service_token_type" {
  value = data.equinix_fabric_service_tokens.service-tokens.data.0.type
}

output "first_service_token_expiration_date_time" {
  value = data.equinix_fabric_service_tokens.service-tokens.data.0.expiration_date_time
}

output "first_service_token_supported_bandwidths" {
  value = data.equinix_fabric_service_tokens.service-tokens.data.0.service_token_connection.0.supported_bandwidths
}

output "first_service_token_virtual_device_type" {
  value = data.equinix_fabric_service_tokens.service-tokens.data.0.service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.type
}

output "first_service_token_virtual_device_uuid" {
  value = data.equinix_fabric_service_tokens.service-tokens.data.0.service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.uuid
}

output "first_service_token_interface_type" {
  value = data.equinix_fabric_service_tokens.service-tokens.data.0.service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.type
}

output "first_service_token_interface_uuid" {
  value = data.equinix_fabric_service_tokens.service-tokens.data.0.service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.id
}
