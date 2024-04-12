data "equinix_fabric_cloud_router" "cloud_router_data_name" {
  uuid = "<uuid_of_cloud_router>"
}

output "id" {
  value = data.equinix_fabric_cloud_router.cloud_router_data_name.id
}

output "name" {
  value = data.equinix_fabric_cloud_router.cloud_router_data_name.name
}

output "account_number" {
  value = data.equinix_fabric_cloud_router.cloud_router_data_name.account.0.account_number
}

output "equinix_asn" {
  value = data.equinix_fabric_cloud_router.cloud_router_data_name.equinix_asn
}

output "metro_code" {
  value = data.equinix_fabric_cloud_router.cloud_router_data_name.location.0.metro_code
}

output "metro_name" {
  value = data.equinix_fabric_cloud_router.cloud_router_data_name.location.0.metro_name
}

output "region" {
  value = data.equinix_fabric_cloud_router.cloud_router_data_name.location.0.region
}

output "package_code" {
  value = data.equinix_fabric_cloud_router.cloud_router_data_name.package.0.code
}

output "project_id" {
  value = data.equinix_fabric_cloud_router.cloud_router_data_name.project.0.project_id
}

output "type" {
  value = data.equinix_fabric_cloud_router.cloud_router_data_name.type
}
