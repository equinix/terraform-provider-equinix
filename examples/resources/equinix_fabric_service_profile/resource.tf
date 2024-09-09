resource "equinix_fabric_service_profile" "new_service_profile" {
  description = "Service Profile for Receiving Connections"
  name = "Name Of Business + Use Case Tag"
  type = "L2_PROFILE"
  visibility = "PUBLIC"
  notifications = [
    {
      emails = ["someone@sample.com"]
      type = "BANDWIDTH_ALERT"
    }
  ]
  allowed_emails = ["test@equinix.com", "testagain@equinix.com"]
  ports = [
    {
      uuid = "c791f8cb-5cc9-cc90-8ce0-306a5c00a4ee"
      type = "XF_PORT"
    }
  ]
  
  access_point_type_configs {
    type = "COLO"
    allow_remote_connections = true
    allow_custom_bandwidth = true
    allow_bandwidth_auto_approval = false
    connection_redundancy_required = false
    connection_label = "Service Profile Tag1"
    bandwidth_alert_threshold = 10
    supported_bandwidths = [ 100, 500 ]
  }
}
