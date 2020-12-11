output "id" {
  description = "Device identifier"
  value       = equinix_network_device.this.uuid
}

output "status" {
  description = "Device provisioning status"
  value       = equinix_network_device.this.status
}

output "license_status" {
  description = "Device license status"
  value       = equinix_network_device.this.license_status
}

output "account_number" {
  description = "Device billing account number"
  value       = equinix_network_device.this.account_number
}

output "cpu_count" {
  description = "Device CPU cores count"
  value       = data.equinix_network_device_platform.this.core_count
}

output "memory" {
  description = "Device memory amount"
  value       = join(" ", [data.equinix_network_device_platform.this.memory, data.equinix_network_device_platform.this.memory_unit])
}

output "software_version" {
  description = "Device software version"
  value       = data.equinix_network_device_software.this.version
}

output "region" {
  description = "Device region"
  value       = equinix_network_device.this.region
}

output "ibx" {
  description = "Device IBX center"
  value       = equinix_network_device.this.ibx
}

output "ssh_ip_address" {
  description = "Device SSH interface IP address"
  value       = equinix_network_device.this.ssh_ip_address
}

output "ssh_ip_fqdn" {
  description = "Device SSH interface FQDN"
  value       = equinix_network_device.this.ssh_ip_fqdn
}

output "interfaces" {
  description = "Device interfaces"
  value       = equinix_network_device.this.interface
}

output "secondary" {
  description = "value"
  value = var.secondary.enabled ? {
    id             = tolist(equinix_network_device.this.secondary_device)[0].uuid
    status         = tolist(equinix_network_device.this.secondary_device)[0].status
    license_status = tolist(equinix_network_device.this.secondary_device)[0].license_status
    account_number = tolist(equinix_network_device.this.secondary_device)[0].account_number
    region         = tolist(equinix_network_device.this.secondary_device)[0].region
    ibx            = tolist(equinix_network_device.this.secondary_device)[0].ibx
    ssh_ip_address = tolist(equinix_network_device.this.secondary_device)[0].ssh_ip_address
    ssh_ip_fqdn    = tolist(equinix_network_device.this.secondary_device)[0].ssh_ip_fqdn
    interfaces     = tolist(equinix_network_device.this.secondary_device)[0].interface
  } : null
}
