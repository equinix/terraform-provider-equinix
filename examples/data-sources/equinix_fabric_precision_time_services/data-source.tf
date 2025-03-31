data "equinix_fabric_precision_time_services" "all" {
  pagination = {
    limit = 2
    offset = 1
  }
  filters = [{
    property = "/type"
    operator = "="
    values = ["PTP"]
  }]
  sort = [{
    direction = "DESC"
    property = "/uuid"
  }]
}


output "ept_service_id" {
  value = data.equinix_fabric_precision_time_services.all.data.0.id
}

output "ept_service_name" {
  value = data.equinix_fabric_precision_time_services.all.data.0.name
}

output "ept_service_state" {
  value = data.equinix_fabric_precision_time_services.all.data.0.state
}

output "ept_service_type" {
  value = data.equinix_fabric_precision_time_services.all.data.0.type
}

output "ept_service_ipv4" {
  value = data.equinix_fabric_precision_time_services.all.data.0.ipv4
}
