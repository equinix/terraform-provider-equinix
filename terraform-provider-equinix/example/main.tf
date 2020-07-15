provider equinix {
  endpoint      = "http://localhost:8080"
  client_id     = "someID"
  client_secret = "someSecret"
}

variable "port_dot1q_pri_uuid" {
  type    = string
  default = "a867f685-41c3-1c37-6de0-320a5c00abdd"
}

variable "port_dot1q_sec_uuid" {
  type    = string
  default = "a867f685-41c4-1c47-6de0-320a5c00abdd"
}

variable "port_qinq_pri_uuid" {
  type    = string
  default = "a867f685-41bd-1bd7-6de0-320a5c00abdd"
}

variable "port_qinq_sec_uuid" {
  type    = string
  default = "a867f685-41be-1be7-6de0-320a5c00abdd"
}

variable "sp_aws_uuid" {
  type    = string
  default = "bb9a0e43-3f9f-478e-8e64-74591fd8fd83"
}

variable "sp_azure_uuid" {
  type    = string
  default = "e685b732-3501-4d4b-b280-0e0049e6f987"
}

variable "sp_gcpi_uuid" {
  type    = string
  default = "5205a692-fb52-43fa-bdf9-9896654ef043"
}

variable "sp_oracle_uuid" {
  type    = string
  default = "c839a858-cd8d-4f2a-b0dc-082623889800"
}

variable "sp_ibm_uuid" {
  type    = string
  default = "2dedda43-16ea-4507-a99e-14459475d18b"
}

variable "sp_alibaba_uuid" {
  type    = string
  default = "22bb9e2e-0004-41ff-b8ff-db504a8c2051"
}

//AWS Direct Connect connection using a Dot1q port
resource "equinix_ecx_l2_connection" "tf-aws-dot1q" {
  name                  = "tf-aws-dot1q"
  profile_uuid          = var.sp_aws_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_dot1q_pri_uuid
  vlan_stag             = 100
  seller_region         = "us-east-1"
  seller_metro_code     = "DC"
  authorization_key     = "357848976964"
}

//AWS Direct Connect connection using a QinQ port
resource "equinix_ecx_l2_connection" "tf-aws-qinq" {
  name                  = "tf-aws-qinq"
  profile_uuid          = var.sp_aws_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_qinq_pri_uuid
  vlan_stag             = 120
  vlan_ctag             = 121
  seller_region         = "us-east-1"
  seller_metro_code     = "DC"
  authorization_key     = "357848976964"
}

//Azure Express Route connection using a Dot1q port - Public
resource "equinix_ecx_l2_connection" "tf-azure-dot1q-pub" {
  name                  = "tf-azure-dot1q-pub-pri"
  profile_uuid          = var.sp_azure_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_dot1q_pri_uuid
  vlan_stag             = 130
  seller_region         = "us-east-1"
  seller_metro_code     = "SV"
  authorization_key     = "c4dff8e8-b52f-4b34-b0d4-c4588f7338f3"
  named_tag             = "Public"
  secondary_connection {
    name      = "tf-azure-dot1q-pub-sec"
    port_uuid = var.port_dot1q_sec_uuid
    vlan_stag = 130
  }
}

//Azure Express Route connection using a Dot1q port - Microsoft
resource "equinix_ecx_l2_connection" "tf-azure-dot1q-mic" {
  name                  = "tf-azure-dot1q-mic-pri"
  profile_uuid          = var.sp_azure_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_dot1q_pri_uuid
  vlan_stag             = 140
  seller_region         = "us-east-1"
  seller_metro_code     = "SV"
  authorization_key     = "c4dff8e8-b52f-4b34-b0d4-c4588f7338f3"
  named_tag             = "Microsoft"
  secondary_connection {
    name      = "tf-azure-dot1q-mic-sec"
    port_uuid = var.port_dot1q_sec_uuid
    vlan_stag = 140
  }
}

resource "equinix_ecx_l2_connection_accepter" "aws_dot1q" {
  connection_id = equinix_ecx_l2_connection.aws_dot1q.id
  access_key    = "AK123456"
  secret_key    = "SK123456"
}

//Azure Express Route connection using a Dot1q port - Manual
resource "equinix_ecx_l2_connection" "tf-azure-dot1q-man" {
  name                  = "tf-azure-dot1q-man-pri"
  profile_uuid          = var.sp_azure_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_dot1q_pri_uuid
  vlan_stag             = 150
  zside_vlan_ctag       = 151
  seller_region         = "us-east-1"
  seller_metro_code     = "SV"
  authorization_key     = "c4dff8e8-b52f-4b34-b0d4-c4588f7338f3"
  named_tag             = "Manual"
  secondary_connection {
    name            = "tf-azure-dot1q-man-sec"
    port_uuid       = var.port_dot1q_sec_uuid
    vlan_stag       = 150
    zside_vlan_ctag = 151
  }
}

//GCPI Connection using a Dot1q port
resource "equinix_ecx_l2_connection" "tf-gcpi-dot1q" {
  name                  = "tf-gcpi-dot1q"
  profile_uuid          = var.sp_gcpi_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_dot1q_pri_uuid
  vlan_stag             = 160
  seller_region         = "us-west1"
  seller_metro_code     = "SV"
  authorization_key     = "33835adc-00fd-4fe1-b9f3-78248e126ef7/us-west1/1"
}

//GCPI Connection using a Dot1q port
resource "equinix_ecx_l2_connection" "tf-gcpi-qinq" {
  name                  = "tf-gcpi-qinq"
  profile_uuid          = var.sp_gcpi_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_qinq_pri_uuid
  vlan_stag             = 160
  vlan_ctag             = 161
  seller_region         = "us-west1"
  seller_metro_code     = "SV"
  authorization_key     = "33835adc-00fd-4fe1-b9f3-78248e126ef8/us-west1/1"
}

//Oracle FastConnect connection using a Dot1q port
resource "equinix_ecx_l2_connection" "tf-oracle-dot1q" {
  name                  = "tf-oracle-dot1q"
  profile_uuid          = var.sp_oracle_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_dot1q_pri_uuid
  vlan_stag             = 170
  seller_region         = "us-ashburn-1"
  seller_metro_code     = "DC"
  authorization_key     = "123456789"
}

//Oracle FastConnect connection using a Dot1q port
resource "equinix_ecx_l2_connection" "tf-oracle-qinq" {
  name                  = "tf-oracle-qinq"
  profile_uuid          = var.sp_oracle_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_qinq_pri_uuid
  vlan_stag             = 180
  vlan_ctag             = 181
  seller_region         = "us-ashburn-1"
  seller_metro_code     = "DC"
  authorization_key     = "ocid1.virtualcircuit.oc1.iad.aaaaaaaafzx4jybgymnyfjfwyxl7f4b6emvp3uast2opcf6z7xzp2lqcnpjq"
}

//IBM Connection using a Dot1q port
resource "equinix_ecx_l2_connection" "tf-ibm-dot1q" {
  name                  = "tf-ibm-dot1q"
  profile_uuid          = var.sp_ibm_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_dot1q_pri_uuid
  vlan_stag             = 190
  seller_region         = "us-east-1"
  seller_metro_code     = "SV"
  authorization_key     = "123456789"
  additional_info {
    name  = "global"
    value = "true"
  }
  additional_info {
    name  = "asn"
    value = "509"
  }
}

//IBM Connection using a Dot1q port
resource "equinix_ecx_l2_connection" "tf-ibm-qinq" {
  name                  = "tf-ibm-qinq"
  profile_uuid          = var.sp_ibm_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_qinq_pri_uuid
  vlan_stag             = 200
  vlan_ctag             = 201
  seller_region         = "us-east-1"
  seller_metro_code     = "SV"
  authorization_key     = "123456789"
  additional_info {
    name  = "global"
    value = "true"
  }
  additional_info {
    name  = "asn"
    value = "509"
  }
}

//Alibaba Express Connect connection using a Dot1q port
resource "equinix_ecx_l2_connection" "tf-alibaba-dot1q" {
  name                  = "tf-alibaba-dot1q"
  profile_uuid          = var.sp_alibaba_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_dot1q_pri_uuid
  vlan_stag             = 210
  seller_region         = "ap-southeast-2"
  seller_metro_code     = "SY"
  authorization_key     = "123456789"
}

//Alibaba Express Connect connection using a Dot1q port
resource "equinix_ecx_l2_connection" "tf-alibaba-qinq" {
  name                  = "tf-alibaba-qinq"
  profile_uuid          = var.sp_alibaba_uuid
  speed                 = 50
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_qinq_pri_uuid
  vlan_stag             = 210
  vlan_ctag             = 211
  seller_region         = "ap-southeast-2"
  seller_metro_code     = "SY"
  authorization_key     = "123456789"
}

//Private profile for myself with dot1q ports
resource "equinix_ecx_l2_serviceprofile" "tf-priv-dot1q" {
  percentage_alert                   = 20.5
  oversubscription_allowed           = false
  api_available                      = false
  connection_name_label              = "Connection"
  equinix_managed_port_vlan          = false
  name                               = "tf-priv"
  bandwidth_threshold_notifications  = ["John.Doe@example.com", "Marry.Doe@example.com"]
  profile_statuschange_notifications = ["John.Doe@example.com", "Marry.Doe@example.com"]
  vc_statuschange_notifications      = ["John.Doe@example.com", "Marry.Doe@example.com"]
  oversubscription                   = "1x"
  private                            = true
  private_user_emails                = ["kkolla@equinix.com"]
  redundancy_required                = false
  tag_type                           = "CTAGED"
  secondary_vlan_from_primary        = false

  features {
    cloud_reach  = true
    test_profile = false
  }
  port {
    uuid       = "a867f685-422f-22f7-6de0-320a5c00abdd"
    metro_code = "NY"
  }
  port {
    uuid       = "a867f685-4231-2317-6de0-320a5c00abdd"
    metro_code = "NY"
  }

  speed_from_api      = false
  customspeed_allowed = false
  speed_band {
    speed      = 1000
    speed_unit = "MB"
  }
  speed_band {
    speed      = 500
    speed_unit = "MB"
  }
  speed_band {
    speed      = 100
    speed_unit = "MB"
  }
}

//Myself: Non-Redundant Connection from a Dot1q port to a Dot1q port 
resource "equinix_ecx_l2_connection" "tf-myself-dot1q-dot1q" {
  name                  = "tf-myself-dot1q-dot1q"
  profile_uuid          = equinix_ecx_l2_serviceprofile.tf-priv-dot1q.uuid
  speed                 = 100
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_dot1q_pri_uuid
  vlan_stag             = 300
  zside_port_uuid       = "a867f685-422f-22f7-6de0-320a5c00abdd"
  zside_vlan_stag       = 500
  seller_region         = "us-east-1"
  seller_metro_code     = "NY"
  authorization_key     = "123456789"
}

//Myself: Redundant Connection from a Dot1q port to a Dot1q port 
resource "equinix_ecx_l2_connection" "tf-myself-dot1q-dot1q-ha" {
  name                  = "tf-myself-dot1q-dot1q-p"
  profile_uuid          = equinix_ecx_l2_serviceprofile.tf-priv-dot1q.uuid
  speed                 = 100
  speed_unit            = "MB"
  notifications         = ["marry@equinix.com", "john@equinix.com"]
  purchase_order_number = "1234567890"
  port_uuid             = var.port_dot1q_pri_uuid
  vlan_stag             = 310
  zside_port_uuid       = "a867f685-422f-22f7-6de0-320a5c00abdd"
  zside_vlan_stag       = 510
  seller_region         = "us-east-1"
  seller_metro_code     = "NY"
  authorization_key     = "123456789"
  secondary_connection {
    name            = "tf-myself-dot1q-dot1q-s"
    port_uuid       = var.port_dot1q_sec_uuid
    vlan_stag       = 310
    zside_port_uuid = "a867f685-4231-2317-6de0-320a5c00abdd"
    zside_vlan_stag = 510
  }
}
