resource "equinix_metal_connection" "example" {
    name               = "tf-metal-to-azure"
    project_id         = local.project_id
    type               = "shared"
    redundancy         = "redundant"
    metro              = "sv"
    speed              = "1000Mbps"
    service_token_type = "a_side"
    contact_email      = "username@example.com"
}

data "equinix_fabric_service_profiles" "example" {
  filter {
    property = "/name"
    operator = "="
    values   = ["Azure ExpressRoute"]
  }
}

resource "equinix_fabric_connection" "example" {
  name = "shared-metal-token-2-azure"
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
    service_token {
      uuid = "<service_token_uuid>"
    }
  }
  z_side {
    access_point {
      type = "SP"
      authentication_key = "<Azure_ExpressRouter_Auth_Key>"
      peering_type = "PRIVATE"
      profile {
        type = "L2_PROFILE"
        uuid = data.equinix_fabric_service_profiles.example[0].data.0.uuid
      }
      location {
        metro_code = "SV"
      }
    }
  }
}
