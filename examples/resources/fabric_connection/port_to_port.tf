resource "equinix_fabric_connection" "port2port" {
  name = "ConnectionName"
  type = "EVPL_VC"
  notifications {
    type = "ALL"
    emails = ["example@equinix.com","test1@equinix.com"]
  }
  bandwidth = 50
  order {
    purchase_order_number= "1-323292"
  }
  a_side {
    access_point {
      type = "COLO"
      port {
        uuid = "<aside_port_uuid>"
      }
      link_protocol {
        type = "QINQ"
        vlan_s_tag = "1976"
        
      }
    }
  }
  z_side {
    access_point {
      type = "COLO"
      port {
        uuid = "<zside_port_uuid>"
      }
      link_protocol {
        type = "QINQ"
        vlan_s_tag = "3711"
      }
      location {
        metro_code= "SV"
      }
    }
  }
}
