provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

data "equinix_fabric_service_profiles" "azure" {
  filter {
    property = "/name"
    operator = "="
    values   = [var.fabric_sp_name]
  }
}

data "equinix_fabric_ports" "qinq-pri" {
  filters {
    name = var.equinix_pri_port_name
  }
}

data "equinix_fabric_ports" "qinq-sec" {
  filters {
    name = var.equinix_sec_port_name
  }
}

resource "equinix_fabric_connection" "azure-qinq" {
  name = var.pri_connection_name
  type = var.connection_type
  notifications {
    type   = var.notifications_type
    emails = var.notifications_emails
  }
  bandwidth = var.bandwidth
  redundancy {
    priority = "PRIMARY"
  }
  order {
    purchase_order_number = var.purchase_order_number
  }
  a_side {
    access_point {
      type = var.aside_ap_type
      port {
        uuid = data.equinix_fabric_ports.qinq-pri.data.0.uuid
      }
      link_protocol {
        type       = var.aside_link_protocol_type
        vlan_s_tag = var.aside_pri_link_protocol_stag
      }
    }
  }
  z_side {
    access_point {
      type               = var.zside_ap_type
      authentication_key = var.zside_ap_authentication_key
      profile {
        type = var.zside_ap_profile_type
        uuid = data.equinix_fabric_service_profiles.azure.data.0.uuid
      }
      location {
        metro_code = var.zside_location
      }
    }
  }
}

resource "equinix_fabric_connection" "azure-qinq-second-connection" {
  name = var.sec_connection_name
  type = var.connection_type
  notifications {
    type   = var.notifications_type
    emails = var.notifications_emails
  }
  bandwidth = var.bandwidth
  redundancy {
    priority = "SECONDARY"
    group    = one(equinix_fabric_connection.azure-qinq.redundancy).group
  }
  order {
    purchase_order_number = var.purchase_order_number
  }
  a_side {
    access_point {
      type = var.aside_ap_type
      port {
        uuid = data.equinix_fabric_ports.qinq-sec.data.0.uuid
      }
      link_protocol {
        type       = var.aside_link_protocol_type
        vlan_s_tag = var.aside_sec_link_protocol_stag
      }
    }
  }
  z_side {
    access_point {
      type               = var.zside_ap_type
      authentication_key = var.zside_ap_authentication_key
      profile {
        type = var.zside_ap_profile_type
        uuid = data.equinix_fabric_service_profiles.azure.data.0.uuid
      }
      location {
        metro_code = var.zside_location
      }
    }
  }
}

output "connection_result" {
  value = equinix_fabric_connection.azure-qinq.id
}

output "second_connection_result" {
  value = equinix_fabric_connection.azure-qinq-second-connection.id
}
