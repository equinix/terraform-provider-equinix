data "equinix_ecx_port" "dot1q-1-pri" {
  name = "sit-001-CX-DC5-NL-Dot1q-BO-10G-PRI-JUN-27"
}

data "equinix_ecx_port" "qinq-1-pri" {
  name = "sit-001-CX-SY1-NL-QinQ-BO-10G-PRI-JUN-31"
}

data "equinix_ecx_port" "qinq-1-sec" {
  name = "sit-001-CX-NY5-NL-Dot1q-BO-10G-PRI-JP-158"
}

//Private profile with QinQ ports
resource "equinix_ecx_l2_serviceprofile" "priv-qinq" {
  bandwidth_alert_threshold          = 20.5
  oversubscription_allowed           = false
  secondary_vlan_from_primary        = false
  connection_name_label              = "Connection"
  name                               = "tf-priv-qinq"
  bandwidth_threshold_notifications  = ["marry@equinix.com", "john@equinix.com"]]
  profile_statuschange_notifications = ["marry@equinix.com", "john@equinix.com"]]
  vc_statuschange_notifications      = ["marry@equinix.com", "john@equinix.com"]]
  private                            = true
  private_user_emails                = ["marry@equinix.com", "john@equinix.com"]
  tag_type                           = "CTAGED"
  oversubscription                   = "1x"
  features {
    cloud_reach  = true
    test_profile = false
  }
  port {
    uuid       = data.equinix_ecx_port.qinq-1-pri.uuid
    metro_code = data.equinix_ecx_port.qinq-1-pri.metro_code
  }
  port {
    uuid       = data.equinix_ecx_port.qinq-1-sec.uuid
    metro_code = data.equinix_ecx_port.qinq-1-sec.metro_code
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

resource "equinix_ecx_l2_connection" "myself-dot1q-qinq" {
  name                  = "tf-myself-dot1q-qinq"
  profile_uuid          = equinix_ecx_l2_serviceprofile.priv-qinq.uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  port_uuid             = data.equinix_ecx_port.dot1q-1-pri.uuid
  vlan_stag             = 400
  zside_port_uuid       = data.equinix_ecx_port.qinq-1-pri.uuid
  zside_vlan_stag       = 600
  zside_vlan_ctag       = 620
  seller_metro_code     = data.equinix_ecx_port.qinq-1-pri.metro_code
}
