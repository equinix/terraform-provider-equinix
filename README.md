<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Provider for Equinix Platform

The Terraform Equinix provider is a plugin for Terraform that allows for lifecycle management of Equinix Platform resources.

[![Build Status](https://travis-ci.com/equinix/terraform-provider-equinix.svg?branch=master)](https://travis-ci.com/github/equinix/terraform-provider-equinix)
[![Go Report Card](https://goreportcard.com/badge/github.com/equinix/terraform-provider-equinix)](https://goreportcard.com/report/github.com/equinix/terraform-provider-equinix)
[![GoDoc](https://godoc.org/github.com/go-resty/resty?status.svg)](https://godoc.org/github.com/equinix/terraform-provider-equinix)
![GitHub](https://img.shields.io/github/license/equinix/terraform-provider-equinix)

---

## Requirements
* [Terraform](https://www.terraform.io/downloads.html) 0.12+
* [Go](https://golang.org/doc/install) 1.14+ (to build provider plugin)

## Quick starts
* [Using the provider](#using-the-provider)
  * [Provider configuration](#provider-configuration)
  * [Data source: ECXF port](#ecx-port-data-source)
  * [Data source: ECXF layer 2 seller profile](#ecx-layer2-seller-profile-data-source)
  * [Resource: ECXF layer 2 connection](#ecx-l2-connection-resource)
  * [Resource: ECXF layer 2 service profile](#ecx-l2-service-profile-resource)
* [Developing the provider](#developing-the-provider)

## Using the provider

### Provider configuration
Equinix provider requires few basic configuration parameters to operate:
- *API endpoint* - Equinix Platform API base URL
- *client identifier* - used for API endpoint authorization with oAuth client credentials grant
- *client secret* - just as above

Above parameters can be provided in terraform file or as environment variables. Nevertheless, please note that it is not recommended to keep sensitive data in plain text files.

Example provider configuration in `main.tf` file:
```
provider equinix {
  endpoint = "https://api.equinix.com"
  client_id = "someID"
  client_secret = "someSecret"
}
```

Example provider configuration using `environment variables`:
```
export EQUINIX_API_ENDPOINT=https://api.equinix.com"
export EQUINIX_API_CLIENTID=someID
export EQUINIX_API_CLIENTSECRET=someSecret
```

### ECX Port data source
Data source `equinix_ecx_port` is used to fetch attributes of ECX port (like UUID) with given port name.

Example usage:
```
data "equinix_ecx_port" "tf-pri-dot1q" {
  name = "sit-001-CX-NY5-NL-Dot1q-BO-10G-PRI-JP-157"
}
```

#### Argument Reference
* `name` - *(Required)* Name of the port

#### Attributes Reference
* `uuid` - Unique identifier of the port
* `status` - Status of the connection
* `region` - Region in which the port resides
*	`ibx` - Equinix IBX where the port resides.
* `metro_code` - The metro code of the metro where the port resides
* `priority` - The priority of the device (primary / secondary) where the port resides
* `encapsulation` - The VLAN encapsulation of the port (Dot1q or QinQ)
* `buyout` - Indicates whether the port supports unlimited connections. If  "false", the port is a standard port with limited connections. If  "true", the port is an "unlimited connections" port that allows multiple connections at no additional charge. 
* `bandwidth` - Port Bandwidth in bytes.
* `status` - Port status that indicates whether a port has been assigned or is ready for connection.

### ECX Layer2 Seller Profile data source
Data source `equinix_ecx_l2_sellerprofile` is used to fetch attributes of ECX Layer2 Seller Profile (like UUID) with a given profile name.

Example usage:
```
data "equinix_ecx_l2_sellerprofile" "tf-aws" {
  name = "AWS Direct Connect"
}
```

#### Argument Reference
* `name` - *(Optional)* Name of seller profile
* `organization_name` - *(Optional)* Name of seller's organization
* `organization_global_name` - *(Optional)* Name of seller's global organization

#### Attributes Reference
* `uuid` - Unique identifier of seller profile
* `speed_from_api` - information if seller is deriving connection speed from an API call
* `speed_customization_allowed` - information if seller allows customer to enter a custom connection speed
* `redundancy_required` - information if seller requires connections to be redundant
* `encapsulation` - seller profile's encapsulation (Dot1q or QinQ)
* `speed_band` - one or more specifications of speed/bandwidth supported by seller profile
  * `speed` - speed/bandwidth supported by this profile
  * `speed_unit` - unit of the speed/bandwidth supported by this profile
* `metro` - one or more specifications of metro locations supported by seller profile
  * `code` - metro code
  * `name` - metro name
  * `ibxes` - list of IBXes supported within given metro
  * `regions` -  list of regions supported within given metro
* `additional_info` - one or more specifications of additional buyer information attrubutes that can be provided in connection definition that uses given seller profile
  * `name` - name of an attribute
  * `description` - textual description of an attribute
  * `data_type` - data type of an attribute *(BOOLEAN / INTEGER / STRING)*
  * `mandatory` - specifies if attribute is mandatory to create connection

### ECX L2 connection resource
Resource `equinix_ecx_l2_connection` is used to manage layer 2 connections in Equinix Cloud Exchange (ECX) Fabric.

Example usage - non redundant connection:
```
resource "equinix_ecx_l2_connection" "aws_dot1q" {
 name = "tf-single-aws"
 profile_uuid = "2a4f7e27-dff8-4f15-aeda-a11ffe9ccf73"
 speed = 200
 speed_unit = "MB"
 notifications = ["marry@equinix.com", "john@equinix.com"]
 port_uuid = "febc9d80-11e0-4dc8-8eb8-c41b6b378df2"
 vlan_stag = 777
 vlan_ctag = 1000
 seller_region = "us-east-1"
 seller_metro_code = "SV"
 authorization_key = "1234456"
}
```

Example usage - redundant connection:
```
resource "equinix_ecx_l2_connection" "redundant_self" {
  name = "tf-redundant-self"
  profile_uuid = "2a4f7e27-dff8-4f15-aeda-a11ffe9ccf73"
  speed = 50
  speed_unit = "MB"
  notifications = ["john@equinix.com", "marry@equinix.com"]
  port_uuid = "febc9d80-11e0-4dc8-8eb8-c41b6b378df2"
  vlan_stag = 800
  zside_port_uuid = "03a969b5-9cea-486d-ada0-2a4496ed72fb"
  zside_vlan_stag = 1010
  seller_region = "us-east-1"
  seller_metro_code = "SV"
  secondary_connection {
    name = "tf-redundant-self-sec"
    port_uuid = "86872ae5-ca19-452b-8e69-bb1dd5f93bd1"
    vlan_stag = 999
    vlan_ctag = 1000
    zside_port_uuid = "393b2f6e-9c66-4a39-adac-820120555420"
    zside_vlan_stag = 1022
  }
}
```

#### Argument Reference
The following arguments are supported:
* `name` - *(Required)* Name of the primary connection - An alpha-numeric 24 characters string which can include only hyphens and underscores ('-' & '\_').
* `profile_uuid` - *(Required)* Unique identifier of the provider's service profile.
* `speed` - *(Required)* Speed/Bandwidth to be allocated to the connection.
* `speed_unit` - *(Required)* Unit of the speed/bandwidth to be allocated to the connection.
* `notifications` - *(Required)* A list of email addresses that would be notified when there are any updates on this connection.
* `purchase_order_number` - *(Optional)* Test field to link the purchase order numbers to the connection on Equinix which would be reflected on the invoice.
* `port_uuid` - *(Required when device_uuid is not set)* Unique identifier of the buyer's port from which the connection would originate.
* `device_uuid` - *(Required when port_uuid is not set)* Unique identifier of the  Network Edge virtual device from which the connection would originate.
* `vlan_stag` - *(Required when port_uuid is set)* S-Tag/Outer-Tag of the connection - a numeric character ranging from 2 - 4094.
* `vlan_ctag` - *(Optional)* C-Tag/Inner-Tag of the connection - a numeric character ranging from 2 - 4094.
* `named_tag` - *(Optional)* The type of peering to set up in case when connecting to Azure Express Route. One of _"Public"_, _"Private"_, _"Microsoft"_, _"Manual"_
* `additional_info` - *(Optional)* one or more additional information key-value objects
  * `name` - *(Required)* additional information key
  * `value` - *(Required)* additional information value
* `zside_port_uuid` - *(Optional)* Unique identifier of the port on the Z side.
* `zside_vlan_stag` - *(Optional)* S-Tag/Outer-Tag of the connection on the Z side.
* `zside_vlan_ctag` - *(Optional)* C-Tag/Inner-Tag of the connection on the Z side.
* `seller_region` - *(Optional)* The region in which the seller port resides.
* `seller_metro_code` - *(Optional)* The metro code that denotes the connection’s destination (Z side).
* `authorization_key` - *(Optional)* Text field based on the service profile you want to connect to.
* `secondary_connection` - *(Optional)* Definition of secondary connection for redundant connectivity. Most attributes are derived from primary connection, except below:
  * `name` - *(Required)*
  * `port_uuid` - *(Required when device_uuid is not set)*
  * `device_uuid` - *(Required when port_uuid is not set)*
  * `vlan_stag` - *(Required when port_uuid is set)* 
  * `vlan_ctag` - *(Optional)*
  * `zside_port_uuid` - *(Optional)*
  * `zside_vlan_stag` - *(Optional)*
  * `zside_vlan_ctag` - *(Optional)*

#### Attributes Reference
In addition to the arguments listed above, the following computed attributes are exported:
* `uuid` - Unique identifier of the connection
* `status` - Status of the connection
* `redundant_uuid` - Unique identifier of the redundant connection (i.e. secondary connection)

#### Update operation behavior
As for now, update of ECXF L2 connection implies removal of old connection (in redundant scenario - both primary and secondary connections), and creation of new one, with required set of attributes.

### ECX L2 service profile resource
Resource `equinix_ecx_l2_serviceprofile` is used to manage layer 2 service profiles in Equinix Cloud Exchange (ECX) Fabric.

Example usage:
```
resource "equinix_ecx_l2_serviceprofile" "private-profile" {
  bandwidth_alert_threshold          = 20.5
  oversubscription_allowed           = false
  connection_name_label              = "Connection"
  name                               = "private-profile"
  bandwidth_threshold_notifications  = ["John.Doe@example.com", "Marry.Doe@example.com"]
  profile_statuschange_notifications = ["John.Doe@example.com", "Marry.Doe@example.com"]
  vc_statuschange_notifications      = ["John.Doe@example.com", "Marry.Doe@example.com"]
  oversubscription                   = "1x"
  private                            = true
  private_user_emails                = ["John.Doe@example.com", "Marry.Doe@example.com"]
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

#### Argument Reference
The following arguments are supported by `equinix_ecx_l2_serviceprofile` resource:
* `bandwidth_alert_threshold` - *(Required)* specifies the port bandwidth threshold percentage. If the bandwidth limit is met or exceeded, an alert is sent to the seller
* `speed_customization_allowed` - *(Required)* allow customer to enter a custom speed
* `oversubscription_allowed` - *(Optional)* regardless of the utilization, the Equinix service will continue to add connections to your links until we reach the oversubscription limit. By selecting this service, you acknowledge that you will manage decisions on when to increase capacity on these links
* `api_integration` - *(Required)* API integration allows you to complete connection provisioning in less than five minutes. Without API Integration, additional manual steps will be required and the provisioning will likely take longer
* `authkey_label` - *(Optional)* the Authentication Key service allows Service Providers with QinQ ports to accept groups of connections or VLANs from Dot1q port customers. This is similar to S-Tag/C-Tag capabilities
* `connection_name_label` - *(Required)* name of the connection
* `ctag_label` - *(Optional)* C-Tag/Inner-Tag of the connection - A numeric character ranging from 2 to 4094
* `servicekey_autogenerated` - *(Optional)* indicates whether multiple connections can be created with the same authorization key to connect to this service profile after the first connection has been approved by the seller
* `equinix_managed_port_vlan` - *(Required)* only applicable if API available is set true. It indicates whether the port and VLAN details are managed by Equinix. 
* `integration_id` - *(Optional)* specifies the API integration ID that was provided to the customer during onboarding. You can validate your API integration ID using the validateIntegrationId API.
* `name` - *(Required)* name of the service profile - An alpha-numeric 50 characters string which can include only hyphens and underscores ('-' & '\_').
* `bandwidth_threshold_notifications` - *(Required)* an array of email ids you would like to notify when there are any updates on your connection
* `profile_statuschange_notifications` - *(Required)* an array of email ids you would like to notify when there are any updates on your connection
* `vc_statuschange_notifications` - *(Required)* an array of email ids you would like to notify when there are any updates on your connection
* `oversubscription` - *(Optional)* you can set an alert for when a percentage of your profile has been sold. Service providers like to use this functionality to alert them when they need to add more ports or when they need to create a new service profile
* `private` - *(Required)* indicates whether or not this is a private profile. If private, it can only be available for creating connections if correct permissions are granted (i.e. not public like AWS/Azure/Oracle/Google, etc.
* `private_user_emails` - *(Optional)* an array of users email ids who have permission to access this service profile
* `redundancy_required` - *(Required)* specify if your connections require redundancy. If yes, then users need to create a secondary redundant connection
* `speed_from_api` - *(Required) derive speed from an API call
* `tag_type` - *(Optional)* specifies additional tagging information required by the seller profile
* `secondary_vlan_from_primary` - *(Required)* indicates whether the VLAN ID of the secondary connection is the same as the primary connection
*  `features` - *(Required)* contains feature-related information
  * `cloud_reach` - *(Required)* indicates whether or not connections to this profile can be created from remote metros
  * `test_profile` - *(Required)* indicates whether or not this profile can be used for test connections.
* `port` - *(Required)* one or more definitions of ports residing in locations, from which your customers will be able to access services using given profile
  * `uuid` - *(Required)* unique identifier of the port
  * `metro_code` - *(Required)* the metro where the port resides
* `speed_band` - *(Required)* one or more definitions of supported speed/bandwidth
  * `speed` - *(Required)* The speed/bandwidth supported by this profile
  * `speed_unit` - *(Required)* unit of the speed/bandwidth supported by this profile

#### Attributes Reference
In addition to the arguments listed above, the following computed attributes are exported:
* `uuid` - Unique identifier of the 
* `status` - Status of the service profile

Building the provider
------------------
1. Clone Equinix Terraform SDK repository

  *Equinix Terraform SDK contains provider source code along with code of required libraries and corresponding tools.*

  **NOTE**: in future, Equinix Go repositories may be released under open source license and moved to Github

   ```
   $ git clone https://oauth2:ACCESS_TOKEN@git.equinix.com/developer-platform/equinix-terraform-sdk.git
   ```

2. Build the provider

   Enter the provider directory and build the provider:
   ```
   $ cd equinix-terraform-sdk/terraform-provider-equinix
   $ make build
   ```

3. Install the provider

  Provider binary can be installed in terraform plugins directory `~/.terraform.d/plugins` by running make with *install* target:
  ```
  $ make install
  ```

## Developing the provider
* use Go programming best practices, *gofmt, go_vet, golint, ineffassign*, etc.
* enter the provider directory
  ```
  $ cd terraform-provider-equinix
  ```
* to build, use make `build` target
  ```
  $ make build
  ```
* to run unit tests, use make `test` target
  ```
  $ make test
  ```
* to run acceptance tests, use make `testacc` target
  ```
  $ make testacc
  ```
  Check "Running acceptance tests" section for more details.
  

### Running acceptance tests
  **NOTE**: acceptance tests create resources on real infrastructure, thus may be subject for costs. In order to run acceptance tests, you must set necessary provider configuration attributes.

  ```
  $ export EQUINIX_API_ENDPOINT=https://api.equinix.com"
  $ export EQUINIX_API_CLIENTID=someID
  $ export EQUINIX_API_CLIENTSECRET=someSecret
  $ make testacc
  ```

#### ECX Port acceptance tests
ECX Port data source acceptance tests use below parameters, that can be set to match with desired tesing environment. If not set, defaults values, **from Sandbox enviroment** are used.
* **TF_ACC_ECX_PORT_NAME** - sets name of the port used in data source

#### ECX L2 connection acceptance tests
ECX Layer 2 connection acceptance tests use below parameters, that can be set to match with desired tesing environment. If not set, defaults values, **from Sandbox enviroment** are used.

* **TF_ACC_ECX_L2_AWS_SP_ID** - sets UUID of Layer2 service profile for AWS
* **TF_ACC_ECX_L2_AZURE_SP_ID** - sets UUID of Layer2 service profile for Azure 
* **TF_ACC_ECX_PRI_DOT1Q_PORT_ID** - sets UUID of Dot1Q encapsulated port on primary device
* **TF_ACC_ECX_SEC_DOT1Q_PORT_ID** - sets UUID of Dot1Q encapsulated port on secondary device

Example - running tests on Sandbox environment but with defined ports:
```
  $ export EQUINIX_API_ENDPOINT=https://sandboxapi.equinix.com"
  $ export EQUINIX_API_CLIENTID=someID
  $ export EQUINIX_API_CLIENTSECRET=someSecret
  $ export TF_ACC_ECX_PRI_DOT1Q_PORT_ID="6ca3704b-c660-4c6f-9e66-3282f8de787b"
  $ export TF_ACC_ECX_SEC_DOT1Q_PORT_ID="7a80ab13-4e04-455c-82e3-79d962d0c0c3"
  $ make testacc
```
