---
page_title: "Migrating from Fabric v3 (ecx) to Fabric v4 (fabric)"
---

# Migrating from Fabric v3 (ecx) to Fabric v4 (fabric)

In December 2023, the Fabric v3 APIs were deprecated and they will reach
End of Life (EOL) in June 2024. The `ecx` Terraform Resources are built on those
APIs and will also reach EOL in June 2024. If you are using them this guide
will help you migrate to the `fabric` Terraform Resources that are built
with the Fabric v4 APIs.

The Fabric v3 Resources are the following data sources and resources:

Data Sources:
* `data "equinix_ecx_l2_sellerprofile"`
* `data "equinix_ecx_l2_sellerprofiles"`
* `data "equinix_ecx_port"`

Resources:
* `resource "equinix_ecx_l2_connection"`
* `resource "equinix_ecx_l2_connection_accepter"`
* `resource "equinix_ecx_l2_service_profile"`

They are being replaced by the family of resources with the 
`equinix_fabric_` prefix.

## Mapping ECX to Fabric

All ECX resources are limited to Layer 2 Connections. This is due to
the limited amount of possibilities that existed at the time of Fabric v3.
The Fabric v4 APIs allow for Layer 2 and Layer 3 Connections.

What this means is that the Fabric resources are more robust so there isn't
an exact one to one mapping, but rather that ECX connection types are now
a subset of a Fabric resource.

### Data Source Mappings:

* Use `data "equinix_fabric_service_profile"` instead of `data "equinix_ecx_sellerprofile"`
* Use `data "equinix_fabric_service_profiles"` instead of `data "equinix_ecx_sellerprofiles"`
* Use `data "equinix_fabric_port"` instead of `data "equinix_ecx_port"`

#### Template Changes Would Be:

Seller Profile:
```hcl
data "equinix_ecx_l2_sellerprofile" "aws" {
  name = "AWS Direct Connect"
}

output "id" {
  value = data.equinix_ecx_l2_sellerprofile.aws.id
}
```
to
```hcl
data "equinix_fabric_service_profile" "aws" {
  uuid = "<uuid_of_aws_service_profile>"
}

output "id" {
  value = data.equinix_fabric_service_profile.aws.id
}
```
or if you still want to search by name you would use the `equinix_fabric_service_profiles` data source
```hcl
data "equinix_fabric_service_profiles" "aws" {
  filter {
    property = "/name"
    operator = "="
    values   = ["AWS Direct Connect"]
  }
}

output "id" {
  value = data.equinix_fabric_service_profile.aws.data.0.id
}
```

Seller Profiles:
```hcl
data "equinix_ecx_l2_sellerprofiles" "aws" {
  organization_global_name = "AWS"
}
```
to 
```hcl
data "equinix_fabric_service_profiles" "aws" {
  filter {
    property = "/name"
    operator = "="
    values   = ["AWS"]
  }
}
```

Port:
```hcl
data "equinix_ecx_port" "tf-pri-dot1q" {
  name = "sit-001-CX-NY5-NL-Dot1q-BO-10G-PRI-JP-157"
}

output "id" {
  value = data.equinix_ecx_port.tf-pri-dot1q.id
}
```
to
```hcl
data "equinix_fabric_ports" "tf-pri-dot1q" {
  filters {
    name = "sit-001-CX-NY5-NL-Dot1q-BO-10G-PRI-JP-157"
  }
}

output "id" {
  value = data.equinix_fabric_ports.tf-pri-dot1q.0.id
}
```

### Resource Mappings

* Use `resource "equinix_fabric_connection"` instead of `resource "equinix_ecx_l2_connection"`
* `resource "equinix_ecx_l2_connection_acceptor` is deprecated.
  * Add your AWS Secret Key and AWS Access Key to `additional_info` property in `resource "equinix_fabric_connection"` for AWS
  * Or use the equivalent resource `aws_dx_connection_confirmation` in the 'AWS' provider instead."
* Use `resource "equinix_fabric_service_profile"` instead of `resource "equinix_ecx_l2_service_profile"`

#### Template Changes

L2 Connection:
```hcl
resource "equinix_ecx_l2_connection" "port-2-aws" {
  name              = "tf-aws"
  profile_uuid      = "<aws_service_profile_uuid>"
  speed             = 200
  speed_unit        = "MB"
  notifications     = ["marry@equinix.com", "john@equinix.com"]
  port_uuid         = "<port_uuid>"
  vlan_stag         = 777
  vlan_ctag         = 1000
  seller_region     = "us-west-1"
  seller_metro_code = "SV"
  authorization_key = "<aws_account_id>"
}
```
to
```hcl
resource "equinix_fabric_connection" "port2aws" {
  name = "tf-aws"
  type = "EVPL_VC" # L2 Connection
  notifications {
    type = "ALL"
    emails = ["marry@equinix.com", "john@equinix.com"]
  }
  bandwidth = 200 # Speed unit is defaulted to MB
  redundancy { priority= "PRIMARY" }
  order {
    purchase_order_number= "1-323929"
  }
  a_side {
    access_point {
      type= "COLO"
      port {
        uuid = "<port_uuid>"
      }
      link_protocol {
        type = "QINQ"
        vlan_s_tag = "777"
        vlan_c_tag = "1000"
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
        uuid = "<aws_service_profile_uuid>"
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
```

Service Profile:
```hcl
resource "equinix_ecx_l2_serviceprofile" "private-profile" {
  name                               = "private-profile"
  description                        = "my private profile"
  connection_name_label              = "Connection"
  bandwidth_threshold_notifications  = ["John.Doe@example.com", "Marry.Doe@example.com"]
  profile_statuschange_notifications = ["John.Doe@example.com", "Marry.Doe@example.com"]
  vc_statuschange_notifications      = ["John.Doe@example.com", "Marry.Doe@example.com"]
  private                            = true
  private_user_emails                = ["John.Doe@example.com", "Marry.Doe@example.com"]
  features {
    allow_remote_connections = true
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
```
to
```hcl
resource "equinix_fabric_service_profile" "private-profile" {
  name = "private-profile"
  description = "my private profile"
  type = "L2_PROFILE" # Need to specify to make it an Layer 2 Profile
  visibility = "PRIVATE"
  notifications = [
    {
      emails = ["John.Doe@example.com", "Marry.Doe@example.com"]
      type = "BANDWIDTH_ALERT"
    },
    {
      emails = ["John.Doe@example.com", "Marry.Doe@example.com"]
      type = "PROFILE_LIFECYCLE"
    },
    {
      emails = ["John.Doe@example.com", "Marry.Doe@example.com"]
      type = "CONNECTION_APPROVAL"
    }
  ]
  allowed_emails = ["John.Doe@example.com", "Marry.Doe@example.com"]
  ports = [
    {
      uuid = "c791f8cb-5cc9-cc90-8ce0-306a5c00a4ee"
      type = "XF_PORT"
    },
    {
      uuid       = "a867f685-4231-2317-6de0-320a5c00abdd"
      type = "XF_PORT"
    }
  ]
  
  access_point_type_configs {
    type = "COLO"
    allow_remote_connections = true
    connection_label = "Connection"
    supported_bandwidths = [ 100, 500, 1000 ]
  }
}
```


## Migrating Terraform State

Once we changed the template accordingly, we can remove the old `equinix_ecx_` resources from Terraform state and import the new ones as `equinix_fabric_` resources by their UUIDs.

In the terraform state and import commands, we use the resource type and name, separated by dot:
```bash
terraform state rm equinix_ecx_l2_connection.example
terraform import equinix_fabric_connection.example <resource_uuid>
```

After that, our templates should be in check with the Terraform state and with the upstream resources in Equinix Fabric. We can verify the migration by running terraform plan, it should show that infrastructure is up to date.


