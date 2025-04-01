data "equinix_fabric_precision_time_service" "ept-test" {
  ept_service_id = "<ept_service_id"
}

output "ept_service_id" {
  value = data.equinix_fabric_precision_time_service.ept-test.id
}

output "ept_service_name" {
  value = data.equinix_fabric_precision_time_service.ept-test.name
}

output "ept_service_state" {
  value = data.equinix_fabric_precision_time_service.ept-test.state
}

output "ept_service_type" {
  value = data.equinix_fabric_precision_time_service.ept-test.type
}

output "ept_service_ipv4" {
  value = data.equinix_fabric_precision_time_service.ept-test.ipv4
}

output "ept_service_connection" {
  value = equinix_fabric_precision_time_service.ptp.connections
}
