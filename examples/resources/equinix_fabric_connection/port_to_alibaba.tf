resource "equinix_fabric_connection" "port2alibaba" {
  name = "ConnectionName"
  type = "EVPL_VC"
  notifications {
    type = "ALL"
    emails = ["example@equinix.com", "test1@equinix.com"]
  }
  bandwidth = 50
  redundancy { priority = "PRIMARY" }
  order {
    purchase_order_number = "1-323929"
  }
  a_side {
    access_point {
      type = "COLO"
      port {
        uuid = "<aside_port_uuid>"
      }
      link_protocol {
        type     = "DOT1Q"
        vlan_tag = "2019"
      }
    }
  }
  z_side {
    access_point {
      type               = "SP"
      authentication_key = "<alibaba_account_id>"
      seller_region      = "us-west-1"
      profile {
        type = "L2_PROFILE"
        uuid = "<service_profile_uuid>"
      }
      location {
        metro_code = "SV"
      }
    }
  }
}
