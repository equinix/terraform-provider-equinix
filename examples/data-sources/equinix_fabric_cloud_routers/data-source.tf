data "equinix_fabric_cloud_routers" "test" {
  filter {
    property = "/name"
    operator = "="
    values 	 = ["Test_PFCR"]
  }
  filter {
    property = "/location/metroCode"
    operator = "="
    values   = ["SV"]
  }
  filter {
    property = "/package/code"
    operator = "="
    values = ["STANDARD"]
    or = true
  }
  filter {
    property = "/state"
    operator = "="
    values = ["ACTIVE"]
    or = true
  }
  pagination {
    offset = 5
    limit = 3
  }
  sort {
    direction = "ASC"
    property = "/name"
  }
}

output "number_of_returned_fcrs" {
  value = length(data.equinix_fabric_cloud_routers.test.data)
}

output "first_fcr_name" {
  value = data.equinix_fabric_cloud_routers.test.data.0.name
}

output "first_fcr_state" {
  value = data.equinix_fabric_cloud_routers.test.data.0.state
}

output "first_fcr_uuid" {
  value = data.equinix_fabric_cloud_routers.test.data.0.uuid
}

output "first_fcr_type" {
  value = data.equinix_fabric_cloud_routers.test.data.0.type
}

output "first_fcr_package_code" {
  value = one(data.equinix_fabric_cloud_routers.test.data.0.package).code
}

output "first_fcr_equinix_asn" {
  value = data.equinix_fabric_cloud_routers.test.data.0.equinix_asn
}

output "first_fcr_location_region" {
  value = one(data.equinix_fabric_cloud_routers.test.data.0.location).region
}

output "first_fcr_location_metro_name" {
  value = one(data.equinix_fabric_cloud_routers.test.data.0.location).metro_name
}

output "first_fcr_location_metro_code" {
  value = one(data.equinix_fabric_cloud_routers.test.data.0.location).metro_code
}

output "first_fcr_project_id" {
  value = one(data.equinix_fabric_cloud_routers.test.data.0.project).project_id
}

output "first_fcr_account_number" {
  value = one(data.equinix_fabric_cloud_routers.test.data.0.account).account_number
}
