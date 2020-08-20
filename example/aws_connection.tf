provider "equinix" {
  client_id     = "your_client_id"
  client_secret = "your_client_secret"
}

data "equinix_ecx_l2_sellerprofile" "aws" {
  name = "AWS Direct Connect"
}

data "equinix_ecx_port" "dot1q-1-pri" {
  name = "sit-001-CX-DC5-NL-Dot1q-BO-10G-PRI-JUN-27"
}

resource "equinix_ecx_l2_connection" "aws-dot1q" {
  name                  = "tf-aws-dot1q"
  profile_uuid          = data.equinix_ecx_l2_sellerprofile.aws.uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = data.equinix_ecx_port.dot1q-1-pri.uuid
  vlan_stag             = 1010
  seller_region         = "us-east-1"
  seller_metro_code     = "DC"
  authorization_key     = "AK123456"
}

//Accepts connection on AWS side
resource "equinix_ecx_l2_connection_accepter" "aws_dot1q" {
  connection_id = equinix_ecx_l2_connection.aws_dot1q.id
  access_key    = "AK123456"
  secret_key    = "SK123456"
}
