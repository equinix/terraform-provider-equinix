provider "equinix" {
  client_id     = "opYOL7xfD0HLshl2SyKAO4ebLn5uWbhQ"
  client_secret = "MdK7Qq8IfaSN9yBS"
}

resource "equinix_fabric_service_profile" "generic" {
  name        = "terra-e2e-generic-sp"
  description = "Generic SP"
  type        = "L2_PROFILE"
  notifications {
    emails        = ["opsuser100@equinix.com"]
    type          = "BANDWIDTH_ALERT"
    send_interval = ""
  }
  tags       = ["Storage", "Compute"]
  visibility = "PRIVATE"
  ports {
    uuid = "c4d9350e-77c5-7c5d-1ce0-306a5c00a600"
    type = "XF_PORT"
    location {
      metro_code = "SV"
    }
    cross_connect_id          = ""
    seller_region             = ""
    seller_region_description = ""
  }
  access_point_type_configs {
    type                             = "COLO"
    connection_redundancy_required   = false
    allow_bandwidth_auto_approval    = false
    allow_remote_connections         = false
    connection_label                 = "test"
    enable_auto_generate_service_key = false
    bandwidth_alert_threshold        = 10
    allow_custom_bandwidth           = true
    api_config {
      api_available        = false
      equinix_managed_vlan = true
      bandwidth_from_api   = false
      integration_id       = "test"
      equinix_managed_port = true
    }
    authentication_key {
      required    = false
      label       = "Service Key"
      description = "XYZ"
    }
    supported_bandwidths = [100, 500]
  }
  marketing_info {
    promotion = false
  }
}

output "connection_result" {
  value = "equinix_fabric_service_profile.generic.uuid"
}