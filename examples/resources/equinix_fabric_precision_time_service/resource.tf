resource "equinix_fabric_precision_time_service" "ptp" {
  type = "PTP"
  name = "tf_acc_eptptp_PFCR"
  package = {
    code = "PTP_STANDARD"
  }
  connections = [
    {
      uuid = "<connection_id>"
    }
  ]
  ipv4 = {
    primary = "191.168.254.241"
    secondary = "191.168.254.242"
    network_mask = "255.255.255.240"
    default_gateway = "191.168.254.254"
  }
}

output "ept_service_id" {
  value = equinix_fabric_precision_time_service.ptp.id
}

output "ept_service_name" {
  value = equinix_fabric_precision_time_service.ptp.name
}

output "ept_service_state" {
  value = equinix_fabric_precision_time_service.ptp.state
}

output "ept_service_type" {
  value = equinix_fabric_precision_time_service.ptp.type
}

output "ept_service_connection" {
  value = equinix_fabric_precision_time_service.ptp.connections
}

output "ept_service_ipv4" {
  value = equinix_fabric_precision_time_service.ptp.ipv4
}
