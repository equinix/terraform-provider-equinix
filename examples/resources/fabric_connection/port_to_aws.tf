resource "equinix_fabric_connection" "port2aws" {
  name = "ConnectionName"
  type = "EVPL_VC"
  notifications {
    type = "ALL"
    emails = ["example@equinix.com","test1@equinix.com"]
  }
  bandwidth = 50
  redundancy { priority= "PRIMARY" }
  order {
    purchase_order_number= "1-323929"
  }
  a_side {
    access_point {
      type= "COLO"
      port {
        uuid = "<aside_port_uuid>"
      }
      link_protocol {
        type = "QINQ"
        vlan_s_tag = "2019"
        vlan_c_tag = "2112"
      }
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
  
  additional_info = [
    { key = "accessKey", value = "<aws_access_key>" },
    { key = "secretKey", value = "<aws_secret_key>" }
  ]
}
