provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

provider "metal" {
  auth_token = var.metal_auth_token
}

data "metal_project" "this" {
    name = var.metal_project_name
}

resource "random_string" "random" {
  length  = 3
  special = false
}

locals {
  connection_name  = format("%s-%s",var.connection_name, random_string.random.result)
  metal_speed_unit = var.connection_speed_unit == "GB" ? "Gbps" : "Mbps"
}

resource "metal_connection" "this" {
    name               = local.connection_name
    project_id         = data.metal_project.this.project_id
    metro              = var.connection_metro
    redundancy         = "primary"
    type               = "shared"
    service_token_type = "a_side"
    description        = var.connection_description
    tags               = ["terraform"]
    speed              = format("%d%s", var.connection_speed, local.metal_speed_unit)
}

data "equinix_ecx_l2_sellerprofile" "this" {
    name = var.seller_profile_name
}

resource "equinix_ecx_l2_connection" "this" {
    name              = local.connection_name
    profile_uuid      = data.equinix_ecx_l2_sellerprofile.this.uuid
    speed             = var.connection_speed
    speed_unit        = var.connection_speed_unit
    notifications     = var.connection_notification_users
    seller_metro_code = var.connection_metro
    seller_region     = var.seller_region
    authorization_key = var.seller_authorization_key
    service_token     = metal_connection.this.service_tokens.0.id
}

output "fabric-connection" {
  value = {
      "name"   = local.connection_name,
      "id"     = equinix_ecx_l2_connection.this.id
      "status" = equinix_ecx_l2_connection.this.status
  }
}

output "metal-connection" {
  value = {
      "name"   = local.connection_name,
      "id"     = metal_connection.this.id
      "status" = metal_connection.this.status
  }
}
