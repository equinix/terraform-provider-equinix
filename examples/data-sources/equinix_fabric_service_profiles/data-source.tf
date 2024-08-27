data "equinix_fabric_service_profiles" "test" {
  and_filters = true
  filter {
    property = "/type"
    operator = "="
    values 	 = ["L2_PROFILE"]
  }
  filter {
    property = "/name"
    operator = "="
    values   = ["SP_ResourceCreation_PFCR"]
  }
  pagination {
    offset = 0
    limit = 5
  }
  sort {
    direction = "ASC"
    property = "/name"
  }
}

output "number_of_returned_service_profiles" {
  value = length(data.equinix_fabric_service_profiles.test.data)
}

output "first_service_profile_name" {
  value = data.equinix_fabric_service_profiles.test.data.0.name
}

output "first_service_profile_uuid" {
  value = data.equinix_fabric_service_profiles.test.data.0.uuid
}

output "first_service_profile_description" {
  value = data.equinix_fabric_service_profiles.test.data.0.description
}

output "first_service_profile_state" {
  value = data.equinix_fabric_service_profiles.test.data.0.state
}

output "first_service_profile_visibility" {
  value = data.equinix_fabric_service_profiles.test.data.0.visibility
}

output "first_service_profile_metros_code" {
  value = data.equinix_fabric_service_profiles.test.data.0.metros.0.code
}

output "first_service_profile_metros_name" {
  value = data.equinix_fabric_service_profiles.test.data.0.metros.0.name
}

output "first_service_profile_metros_display_name" {
  value = data.equinix_fabric_service_profiles.test.data.0.metros.0.display_name
}

output "first_service_profile_type" {
  value = data.equinix_fabric_service_profiles.test.data.0.type
}
