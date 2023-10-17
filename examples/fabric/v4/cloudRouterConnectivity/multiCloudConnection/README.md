# ECX Fabric Layer2 MultiCloud Connection: FCR 2 AWS and Azure

This example shows how to create Layer 2 Connection between FCR to AWS and Azure.

## Adjust variables

At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above

`fcr_name` - Name of ECX Fabric Cloud Router on a-side , i.e. amcrh007-fcr
`fcr_type` - Fabric Cloud Router type
`fcr_location` - Fabric Cloud Router location
`fcr_project` - Fabric Cloud Router project
`fcr_account` - Fabric Cloud Router account
`fcr_package` - Fabric Cloud Router package type, i.e. PRO
`notifications_type` - notification type
`notifications_emails` - List of emails


`azure_connection_name` - The name of the Azure connection
`azure_connection_type` - Connection type, please refer to OAS schema for enum values.
`azure_notifications_type` - Notification type
`azure_notifications_emails` - List of emails
`azure_bandwidth` - Bandwidth in MBs
`azure_redundancy` - Port redundancy PRIMARY or SECONDARY
`azure_purchase_order_number` - Purchase order number applied to billing invoices for this connection.
`azure_peering_type` - Peering Type
`azure_aside_ap_type` - Access point type

`azure_zside_ap_type` - Z side access point type
`azure_zside_ap_authentication_key` - AZURE authorization key, like c620477c-3f30-41e8-a0b9-cf324a12121d
`azure_zside_ap_profile_type` - Service profile type
`azure_zside_location` - Equinix Metro Code for the Z side access point
`azure_fabric_sp_name` - Service profile name like i.e. AZURE

`azure_rp_name`- Name of Direct routing Protocol
`azure_rp_type`- Type of Direct routing Protocol entity, "DIRECT"
`azure_equinix_ipv4_ip` = Equinix Side IpV4 Address
`azure_equinix_ipv6_ip` = Equinix Side IpV6 Address

`azure_bgp_rp_name` - Name of BGP routing Protocol
`azure_bgp_rp_type` - Type of BGP routing Protocol entity, "BGP"
`azure_bgp_customer_peer_ipv4` - Customer Side IpV4 Address
`azure_bgp_customer_peer_ipv6` - Customer Side IpV6 Address
`azure_bgp_enabled_ipv4` - Enable BGP IpV4 session from customer side
`azure_bgp_enabled_ipv6` - Enable BGP IpV6 session from customer side
`azure_bgp_customer_asn` - Customer ASN Number

`aws_connection_name` - The name of the AWS connection
`aws_connection_type` - connection type, please refer schema
`aws_notifications_type` - notification type
`aws_notifications_emails` - List of emails
`aws_bandwidth` - bandwidth in MBs
`aws_redundancy` - Port redundancy
`aws_aside_ap_type` - Fabric Cloud Router type
`aws_zside_ap_type` - Z side access point type
`aws_zside_ap_authentication_key` - AWS authorization key, account number like 357848912121
`aws_access_key` - AWS access key, like BQR12AHQKSYUTPBGHPIJ
`aws_secret_key` - AWS secret key, like 2qwrbYTUUIQWOOEIHDJSKbhikjhalpe
`aws_zside_ap_profile_type` - Service profile type
`aws_fabric_sp_name` - Service profile name, fetched based on Service Profile get call using Service Profile search schema
`aws_zside_location` - Seller location
`aws_seller_region` - Seller region code

`aws_rp_name`- Name of Direct routing Protocol
`aws_rp_type`- Type of Direct routing Protocol entity, "DIRECT"
`aws_equinix_ipv4_ip` = Equinix Side IpV4 Address
`aws_equinix_ipv6_ip` = Equinix Side IpV6 Address

`aws_bgp_rp_name` - Name of BGP routing Protocol
`aws_bgp_rp_type` - Type of BGP routing Protocol entity, "BGP"
`aws_bgp_customer_peer_ipv4` - Customer Side IpV4 Address
`aws_bgp_customer_peer_ipv6` - Customer Side IpV6 Address
`aws_bgp_enabled_ipv4` - Enable BGP IpV4 session from customer side
`aws_bgp_enabled_ipv6` - Enable BGP IpV6 session from customer side
`aws_bgp_customer_asn` - Customer ASN Number

## Azure login

Log in to Azure portal with an account that has permission to create necessary resources.

Create an Azure ExpressRoute Circuit and use its Service Key as the Authentication Key in the examples.

Bandwidth in Terraform must match the bandwidth of the ExpressRoute Circuit created in Azure.

## AWS login

Log in to AWS portal use account that has permission to create necessary resources.

## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
  In the given example, the folder to perform CRUD operations for multi cloud connections can be found at examples/fabric/v4/cloudRouterConnectivity/MutliCloudConnection

- Change directory into - `CD examples/fabric/v4/cloudRouterConnectivity/MutliCloudConnection`
- Initialize Terraform plugins - `terraform init`

## Multi Cloud connection  : Create, Read, Update and Delete(CRUD) operations
Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                               Description |
|:----------|:---------------------------------:|--------------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                                  Creates multi-cloud connection resources |
| READ      |         `terraform show`          |      Reads/Shows the current state of the multi-cloud connection resources |
| UPDATE    |    `terraform apply -refresh`     | Updates the connections with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                       Deletes the created multi-cloud connection resources |
