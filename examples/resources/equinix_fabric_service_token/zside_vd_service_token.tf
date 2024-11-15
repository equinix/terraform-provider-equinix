resource "equinix_fabric_service_token" "test" {
  type                 = "VC_TOKEN"
  description = "Zside VD Service Token"
  expiration_date_time = "2025-01-18T06:43:49.986Z"
  service_token_connection {
    type                 = "EVPL_VC"
    supported_bandwidths = [50, 200, 10000]
    z_side {
      access_point_selectors {
        type = "VD"
        virtual_device {
          type = "EDGE"
          uuid = "<device_uuid>"
        }
        interface {
          type = "NETWORK"
        }
      }
    }
  }
  notifications {
    type   = "ALL"
    emails = ["example@equinix.com"]
  }
}
