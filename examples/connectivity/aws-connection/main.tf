provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

provider "aws" {
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
  region     = var.aws_region
}

data "equinix_ecx_l2_sellerprofile" "aws" {
  name                     = "AWS Direct Connect"
  organization_global_name = "AWS"
}

data "equinix_ecx_port" "dot1q-pri" {
  name = var.equinix_port_name
}

resource "equinix_ecx_l2_connection" "example" {
  name              = "tf-aws-dot1q"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.aws.uuid
  speed             = 50
  speed_unit        = "MB"
  notifications     = ["example@equinix.com"]
  port_uuid         = data.equinix_ecx_port.dot1q-pri.uuid
  vlan_stag         = 1010
  seller_region     = var.aws_region
  seller_metro_code = var.aws_metro_code
  authorization_key = var.aws_account_id
}

locals  {
   aws_connection_id = one([
       for action_data in one(equinix_ecx_l2_connection.example.actions).required_data: action_data["value"]
       if action_data["key"] == "awsConnectionId"
   ])
}

resource "aws_dx_connection_confirmation" "confirmation" {
  connection_id = local.aws_connection_id
}

resource "aws_dx_private_virtual_interface" "example" {
  connection_id    = aws_dx_connection_confirmation.confirmation.id
  name             = "example"
  vlan             = equinix_ecx_l2_connection.example.zside_vlan_stag
  address_family   = "ipv4"
  bgp_asn          = 64999
  amazon_address   = "169.254.0.1/30"
  customer_address = "169.254.0.2/30"
  bgp_auth_key     = "secret"
  vpn_gateway_id   = aws_vpn_gateway.example.id
}

resource "aws_vpc" "example" {
  cidr_block = "10.255.255.0/28"
}

resource "aws_vpn_gateway" "example" {
  vpc_id = aws_vpc.example.id
}
