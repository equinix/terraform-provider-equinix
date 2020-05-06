Terraform Provider for Equinix Platform
==================
* Contact us : https://developer.equinix.com/contact-us

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------------
* [Terraform](https://www.terraform.io/downloads.html) 0.12+
* [Go](https://golang.org/doc/install) 1.14+ (to build provider plugin)

Using the provider
------------------
**NOTE**: dedicated documentation page will be created in future

The Equinix provider is used to manage Equinix Platform infrastructure using Terraform.

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
* `name` - *(Required)* Name of the primary connection - An alpha-numeric 24 characters string which can include only hyphens and underscores ('-' & '_').
* `profile_uuid` - *(Required)* Unique identifier of the provider's service profile.
* `speed` - *(Required)* Speed/Bandwidth to be allocated to the connection.
* `speed_unit` - *(Required)* Unit of the speed/bandwidth to be allocated to the connection.
* `notifications` - *(Required)* A list of email addresses that would be notified when there are any updates on this connection.
* `purchase_order_number` - *(Optional)* Test field to link the purchase order numbers to the connection on Equinix which would be reflected on the invoice.
* `port_uuid` - *(Required)* Unique identifier of the buyer's port from which the connection would originate.
* `vlan_stag` - *(Required)* S-Tag/Outer-Tag of the connection - a numeric character ranging from 2 - 4094.
* `vlan_ctag` - *(Optional)* C-Tag/Inner-Tag of the connection - a numeric character ranging from 2 - 4094.
* `zside_port_uuid` - *(Optional)* Unique identifier of the port on the Z side.
* `zside_vlan_stag` - *(Optional)* S-Tag/Outer-Tag of the connection on the Z side.
* `zside_vlan_ctag` - *(Optional)* C-Tag/Inner-Tag of the connection on the Z side.
* `seller_region` - *(Required)* The region in which the seller port resides.
* `seller_metro_code` - *(Required)* The metro code that denotes the connection’s destination (Z side).
* `authorization_key` - *(Optional)* Text field based on the service profile you want to connect to.
* `secondary_connection` - *(Optional)* Definition of secondary connection for redundant connectivity. Most attributes are derived from primary connection, except below:
  * `name` - *(Required)*
  * `port_uuid` - *(Required)*
  * `vlan_stag` - *(Required)*
  * `vlan_ctag` - *(Optional)*
  * `zside_port_uuid` - *(Optional)*
  * `zside_vlan_stag` - *(Optional)*
  * `zside_vlan_ctag` - *(Optional)*

#### Attributes Reference
In addition to the arguments listed above, the following computed attributes are exported:
* `uuid` - Unique identifier of the connection
* `status` - Status of the connection
* `redundant_uuid` - Unique identifier of the redundant connection (i.e. secondary connection)

#### Update operation behaviour
As for now, update of ECXF L2 connection implies removal of old connection (in redundant scenario - both primary and secondary connections), and creation of new one, with required set of attributes.


Building the provider
------------------
1. Clone Equinix Terraform SDK repository

  *Equinix Terraform SDK contains provider source code along with code of required librarires and corresponding tools.*

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

Developing the provider
------------------
* use Go programming best practices, *gofmt, go_vet, golint, ineffassign*, etc.
* enter the provider directory
   ```
   $ cd equinix-terraform-sdk/terraform-provider-equinix
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

   **NOTE**: acceptance tests create resources on real infrastructure, thus may be subject for costs. In order to run acceptance tests, you must set necessary provider configuration attributes.
   ```
   $ export EQUINIX_API_ENDPOINT=https://api.equinix.com"
   $ export EQUINIX_API_CLIENTID=someID
   $ export EQUINIX_API_CLIENTSECRET=someSecret
   $ make testacc
   ```
