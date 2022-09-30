provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

data "equinix_fabric_service_profiles" "alibaba" {
  filter {
    property = "/name"
    operator = "="
    values = [var.fabric_sp_name]
  }
}

data "equinix_fabric_ports" "qinq-pri" {
  local_var_optionals {
    name = var.equinix_port_name
  }
}

resource "equinix_fabric_connection" "alibaba-qinq" {
  name = var.connection_name
  description = var.description
  type = var.connection_type
  notifications{
    type=var.notifications_type
    emails=var.notifications_emails
  }
  bandwidth = var.bandwidth
  redundancy {priority= var.redundancy}
  order {
    purchase_order_number= var.purchase_order_number
  }
  a_side {
    access_point {
      type= var.aside_ap_type
      port {
        uuid= data.equinix_fabric_ports.qinq-pri.data.0.uuid
      }
      link_protocol {
        type= var.aside_link_protocol_type
        vlan_s_tag= var.aside_link_protocol_stag
        vlan_c_tag= var.aside_link_protocol_ctag
      }
    }
  }
  z_side {
    access_point {
      type= var.zside_ap_type
      authentication_key= var.zside_ap_authentication_key
      seller_region = var.seller_region
      profile {
        type= var.zside_ap_profile_type
        uuid= data.equinix_fabric_service_profiles.alibaba.data.0.uuid
      }
      location {
        metro_code= var.zside_location
      }
    }
  }
}

output "connection_result" {
  value = equinix_fabric_connection.alibaba-qinq.id
}

resource "time_sleep" "wait_for_ingress_alb" {
  destroy_duration = "180s"

  depends_on = [equinix_fabric_connection.alibaba-qinq]
}