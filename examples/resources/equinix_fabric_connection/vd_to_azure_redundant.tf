resource "equinix_fabric_connection" "vd2azure_primary" {
  name = "ConnectionName"
  type = "EVPL_VC"
  redundancy { priority = "PRIMARY" }
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
        type = "CLOUD"
        id = 7
      }
    }
  }
  z_side {
    access_point {
      type = "SP"
      authentication_key = "<Azure_ExpressRouter_Auth_Key>"
      peering_type = "PRIVATE"
      profile {
        type = "L2_PROFILE"
        uuid = "<Azure_Service_Profile_UUID>"
      }
      location {
        metro_code = "SV"
      }
    }
  }
}

resource "equinix_fabric_connection" "vd2azure_secondary" {
  name = "ConnectionName"
  type = "EVPL_VC"
  redundancy {
    priority = "SECONDARY"
    group = one(equinix_fabric_connection.vd2azure_primary.redundancy).group
  }
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
        type = "CLOUD"
        id = 5
      }
    }
  }
  z_side {
    access_point {
      type = "SP"
      authentication_key = "<Azure_ExpressRouter_Auth_Key>"
      peering_type = "PRIVATE"
      profile {
        type = "L2_PROFILE"
        uuid = "<Azure_Service_Profile_UUID>"
      }
      location {
        metro_code = "SV"
      }
    }
  }
}
