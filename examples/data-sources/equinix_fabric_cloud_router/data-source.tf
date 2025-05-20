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
  value = [for account in data.equinix_fabric_cloud_router.cloud_router_data_name.account: account.account_number]
}

output "equinix_asn" {
  value = data.equinix_fabric_cloud_router.cloud_router_data_name.equinix_asn
}

output "metro_code" {
  value = [for location in data.equinix_fabric_cloud_router.cloud_router_data_name.location: location.metro_code]
}

output "metro_name" {
  value = [for location in data.equinix_fabric_cloud_router.cloud_router_data_name.location: location.metro_name]
}

output "region" {
  value = [for location in data.equinix_fabric_cloud_router.cloud_router_data_name.location: location.region]
}

output "package_code" {
  value = [for package in data.equinix_fabric_cloud_router.cloud_router_data_name.package: package.code]
}

output "project_id" {
  value = [for project in data.equinix_fabric_cloud_router.cloud_router_data_name.project: project.project_id]
}

output "type" {
  value = data.equinix_fabric_cloud_router.cloud_router_data_name.type
}
