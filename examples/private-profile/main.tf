provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

data "equinix_ecx_port" "dot1q-1-pri" {
  name = var.equinix_pri_port_name
}

data "equinix_ecx_port" "dot1q-1-sec" {
  name = var.equinix_sec_port_name
}

//Private profile for myself with qinq ports
resource "equinix_ecx_l2_serviceprofile" "private-profile" {
  name                               = "tf-private-profile"
  description                        = "Private profile for limited use only"
  connection_name_label              = "Software Defined Interconnect"
  bandwidth_alert_threshold          = 20.5
  bandwidth_threshold_notifications  = ["example@example.com"]
  profile_statuschange_notifications = ["example@example.com"]
  vc_statuschange_notifications      = ["example@example.com"]
  private                            = true
  private_user_emails                = ["John@equinix.com", "Marry@equinix.com"]
  features {
    allow_remote_connections = false
    test_profile             = false
  }
  port {
    uuid       = data.equinix_ecx_port.dot1q-1-pri.uuid
    metro_code = data.equinix_ecx_port.dot1q-1-pri.metro_code
  }
  port {
    uuid       = data.equinix_ecx_port.dot1q-1-sec.uuid
    metro_code = data.equinix_ecx_port.dot1q-1-sec.metro_code
  }
  api_integration = true
  integration_id  = "Example-company-CrossConnect-01"
  tag_type        = "CTAGED"
  speed_from_api  = true
}
