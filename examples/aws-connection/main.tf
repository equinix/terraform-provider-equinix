provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

data "equinix_ecx_l2_sellerprofile" "aws" {
  name = "AWS Direct Connect"
}

data "equinix_ecx_port" "dot1q-pri" {
  name = var.equinix_port_name
}

resource "equinix_ecx_l2_connection" "aws-dot1q" {
  name              = "tf-aws-dot1q"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.aws.uuid
  speed             = 50
  speed_unit        = "MB"
  notifications     = ["example@equinix.com"]
  port_uuid         = data.equinix_ecx_port.dot1q-pri.uuid
  vlan_stag         = 1010
  seller_region     = "us-east-1"
  seller_metro_code = "DC"
  authorization_key = var.aws_account_id
}

//Accepts connection on AWS side
resource "equinix_ecx_l2_connection_accepter" "aws-dot1q" {
  connection_id = equinix_ecx_l2_connection.aws-dot1q.id
  access_key    = var.aws_access_key
  secret_key    = var.aws_secret_key
}
