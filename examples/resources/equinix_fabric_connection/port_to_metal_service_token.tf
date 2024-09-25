resource "equinix_fabric_connection" "port2metaltoken" {
  name = "Dedicatedp-sharedp-VC"
  type = "EVPL_VC"
  notifications {
    type   = "ALL"
    emails = ["example@equinix.com"]
  }
  bandwidth = 50
  a_side {
    access_point {
      type= "COLO"
      port {
        uuid = "<aside_port_uuid>"
      }
      link_protocol {
        type = "DOT1Q"
        vlan_tag = "1020"
      }
    }
  }
  z_side {
    service_token {
      uuid = "<Metal Service Token>"
    }
  }
}
