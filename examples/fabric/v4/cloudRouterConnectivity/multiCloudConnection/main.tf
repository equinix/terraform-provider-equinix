provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

resource "equinix_fabric_cloud_router" "test" {
  name = var.fcr_name
  type = var.fcr_type
  notifications {
    type   = var.notifications_type
    emails = var.notifications_emails
  }
  order {
    purchase_order_number = var.purchase_order_number
  }
  location {
    metro_code = var.fcr_location
  }
  package {
    code = var.fcr_package
  }
  project {
    project_id = var.fcr_project
  }
  account {
    account_number = var.fcr_account
  }
}

output "fcr_result" {
  value = equinix_fabric_cloud_router.test.id
}

data "equinix_fabric_service_profiles" "azure" {
  filter {
    property = "/name"
    operator = "="
    values   = [var.azure_fabric_sp_name]
  }
}

resource "equinix_fabric_connection" "fcr2azure" {
  name = var.azure_connection_name
  type = var.azure_connection_type

  notifications {
    type   = var.azure_notifications_type
    emails = var.azure_notifications_emails
  }
  bandwidth = var.azure_bandwidth
  redundancy {
    priority = var.azure_redundancy
  }
  order {
    purchase_order_number = var.azure_purchase_order_number
  }
  a_side {
    access_point {
      type = var.azure_aside_ap_type
      router {
        uuid = equinix_fabric_cloud_router.test.id
      }
    }
  }

  z_side {
    access_point {
      type               = var.azure_zside_ap_type
      authentication_key = var.azure_zside_ap_authentication_key
      peering_type       = var.azure_peering_type
      profile {
        type = var.azure_zside_ap_profile_type
        uuid = data.equinix_fabric_service_profiles.azure.id
      }
      location {
        metro_code = var.azure_zside_location
      }
    }
  }
}
output "azure_connection_name" {
  value = equinix_fabric_connection.fcr2azure.name
}
output "azure_connection_id" {
  value = equinix_fabric_connection.fcr2azure.id
}

resource "equinix_fabric_routing_protocol" "azure-direct-protocol" {
  connection_uuid = equinix_fabric_connection.fcr2azure.id
  type            = var.azure_rp_type
  name            = var.azure_rp_name
  direct_ipv4 {
    equinix_iface_ip = var.azure_equinix_ipv4_ip
  }
  direct_ipv6 {
    equinix_iface_ip = var.azure_equinix_ipv6_ip
  }
}

output "azure_rp_direct_id" {
  value = equinix_fabric_routing_protocol.azure-direct-protocol.id
}

resource "equinix_fabric_routing_protocol" "azure-bgp-protocol" {
  connection_uuid = equinix_fabric_connection.fcr2azure.id
  type            = var.azure_bgp_rp_type
  name            = var.azure_bgp_rp_name
  bgp_ipv4 {
    customer_peer_ip = var.azure_bgp_customer_peer_ipv4
    enabled          = var.azure_bgp_enabled_ipv4
  }
  bgp_ipv6 {
    customer_peer_ip = var.azure_bgp_customer_peer_ipv6
    enabled          = var.azure_bgp_enabled_ipv6
  }
  customer_asn = var.azure_bgp_customer_asn
  depends_on   = [equinix_fabric_routing_protocol.azure-direct-protocol]
}

output "azure_rp_bgp_id" {
  value = equinix_fabric_routing_protocol.azure-bgp-protocol.id
}


data "equinix_fabric_service_profiles" "aws" {
  filter {
    property = "/name"
    operator = "="
    values   = [var.aws_fabric_sp_name]
  }
}

resource "equinix_fabric_connection" "fcr2aws" {
  name = var.aws_connection_name
  type = var.aws_connection_type
  notifications {
    type   = var.aws_notifications_type
    emails = var.aws_notifications_emails
  }
  additional_info = [{ "key" = "accessKey", "value" = var.aws_access_key }, { "key" = "secretKey", "value" = var.aws_secret_key }]
  bandwidth       = var.aws_bandwidth
  redundancy { priority = var.aws_redundancy }
  order {
    purchase_order_number = var.aws_purchase_order_number
  }
  a_side {
    access_point {
      type = var.aws_aside_ap_type
      router {
        uuid = equinix_fabric_cloud_router.test.id
      }
    }
  }
  z_side {
    access_point {
      type               = var.aws_zside_ap_type
      authentication_key = var.aws_zside_ap_authentication_key
      seller_region      = var.aws_seller_region
      profile {
        type = var.aws_zside_ap_profile_type
        uuid = data.equinix_fabric_service_profiles.aws.id
      }
      location {
        metro_code = var.aws_zside_location
      }
    }
  }
}

output "aws_connection_name" {
  value = equinix_fabric_connection.fcr2aws.name
}

output "aws_connection_id" {
  value = equinix_fabric_connection.fcr2aws.id
}

resource "equinix_fabric_routing_protocol" "aws-direct-protocol" {
  connection_uuid = equinix_fabric_connection.fcr2aws.id
  type            = var.aws_rp_type
  name            = var.aws_rp_name
  direct_ipv4 {
    equinix_iface_ip = var.aws_equinix_ipv4_ip
  }
  direct_ipv6 {
    equinix_iface_ip = var.aws_equinix_ipv6_ip
  }
}

output "aws_rp_direct_id" {
  value = equinix_fabric_routing_protocol.aws-direct-protocol.id
}

resource "equinix_fabric_routing_protocol" "aws-bgp-protocol" {
  connection_uuid = equinix_fabric_connection.fcr2aws.id
  type            = var.aws_bgp_rp_type
  name            = var.aws_bgp_rp_name
  bgp_ipv4 {
    customer_peer_ip = var.aws_bgp_customer_peer_ipv4
    enabled          = var.aws_bgp_enabled_ipv4
  }
  bgp_ipv6 {
    customer_peer_ip = var.aws_bgp_customer_peer_ipv6
    enabled          = var.aws_bgp_enabled_ipv6
  }
  customer_asn = var.aws_bgp_customer_asn

  depends_on = [equinix_fabric_routing_protocol.aws-direct-protocol]
}

output "aws_rp_bgp_id" {
  value = equinix_fabric_routing_protocol.aws-bgp-protocol.id
}
