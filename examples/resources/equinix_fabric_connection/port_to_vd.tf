resource "equinix_fabric_connection" "vd2port" {
  name = "ConnectionName"
  type = "EVPL_VC"
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
      type = "VD"
      virtual_device {
        type = "EDGE"
        uuid = "<device_uuid>"
      }
      interface {
        type = "NETWORK"
        id = 7
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
        type       = "DOT1Q"
        vlan_s_tag = "3711"
      }
      location {
        metro_code = "SV"
      }
    }
  }
}
