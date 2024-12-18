resource "equinix_fabric_service_token" "test"{
  type                 = "VC_TOKEN"
  description          = "Zside COLO Service Token"
  expiration_date_time = "2025-01-18T06:43:49.981Z"
  service_token_connection {
    type = "EVPL_VC"
    supported_bandwidths = [50, 200, 10000]
    z_side {
      access_point_selectors{
        type = "COLO"
        port {
          uuid = "<port_uuid>"
        }
        link_protocol {
          type = "DOT1Q"
          vlan_tag = "2087"
        }
      }
    }
  }
  notifications {
    type   = "ALL"
    emails = ["example@equinix.com", "test1@equinix.com"]
  }
}
