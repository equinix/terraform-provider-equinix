resource "equinix_fabric_connection" "epl" {
  name = "ConnectionName"
  type = "EVPLAN_VC"
  notifications {
    type   = "ALL"
    emails = ["example@equinix.com", "test1@equinix.com"]
  }
  bandwidth = 50
  order {
    purchase_order_number = "1-323292"
  }
  a_side {
    access_point {
      type = "COLO"
      port {
        uuid = "<aside_port_uuid>"
      }
      link_protocol {
        type = "DOT1Q"
        vlan_s_tag = "1976"

      }
    }
  }
  z_side {
    access_point {
      type = "NETWORK"
      network {
        uuid = "<network_uuid>"
      }
    }
  }
}
