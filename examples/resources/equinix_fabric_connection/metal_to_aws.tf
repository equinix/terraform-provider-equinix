resource "equinix_fabric_connection" "metal2aws" {
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
      type               = "METAL_NETWORK"
      authentication_key = "<metal_authorization_code>"
    }
  }
  z_side {
    access_point {
      type = "SP"
      authentication_key = "<aws_account_id>"
      seller_region = "us-west-1"
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
