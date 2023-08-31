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

resource "equinix_ecx_l2_serviceprofile" "seller-profile" {
  name                               = "tf-seller-profile"
  description                        = "Quickly and simply interconnect with our network"
  connection_name_label              = "Software Defined Interconnect"
  bandwidth_alert_threshold          = 20.5
  bandwidth_threshold_notifications  = ["example@example.com"]
  profile_statuschange_notifications = ["example@example.com"]
  vc_statuschange_notifications      = ["example@example.com"]
  features {
    allow_remote_connections = true
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
  speed_band {
    speed      = 1
    speed_unit = "GB"
  }
  speed_band {
    speed      = 500
    speed_unit = "MB"
  }
  speed_band {
    speed      = 100
    speed_unit = "MB"
  }
  speed_band {
    speed      = 50
    speed_unit = "MB"
  }
}
