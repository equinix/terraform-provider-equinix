# ECX Fabric Layer2 Redundant Connection to IBM 2

This example shows how to create Layer 2 Connection between ECX Fabric ports and IBM2 Cloud.

## Adjust variables

At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above

`connection_name` - The name of the connection
`connection_type` - Connection type, please refer to OAS schema for enum values.
`notifications_type` - Notification type
`notifications_emails` - List of emails
`bandwidth` - Bandwidth in MBs
`redundancy` - Port redundancy
`purchase_order_number` - Purchase order number applied to billing invoices for this connection.
`aside_ap_type` - Access point type
`aside_link_protocol_type` - Link protocol type
`aside_pri_link_protocol_tag` - Tag number
`zside_ap_type` - Z side access point type
`zside_ap_authentication_key` - IBM authorization key (Account Id), like 6if92b31d921499f903592cd816f6aw1
`zside_ap_profile_type` - Service profile type
`zside_location` - Equinix Metro Code for the Z side access point
`fabric_sp_name` - Service profile name like i.e. AZURE
`equinix_port_name` -  Name of ECX Fabric Port 
`seller_asn` - Seller ASN Number
`seller_region` - Seller Region

## IBM login

Log in to IBM portal with an account that has permission to create necessary resources.

## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
  In the given example, the folder to perform CRUD operations for port2ibm2 redundant connections can be found at examples/fabric/v4/portConnectivity/ibm/ibm2.

- Change directory into - `CD examples/fabric/v4/portConnectivity/ibm/ibm2`
- Initialize Terraform plugins - `terraform init`

## Port to IBM2 connection  : Create, Read, Update and Delete(CRUD) operations
Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                                Description |
|:----------|:---------------------------------:|---------------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                                   Creates a port2ibm2 connection resources |
| READ      |         `terraform show`          |        Reads/Shows the current state of the port2ibm2 connection resources |
| UPDATE    |    `terraform apply -refresh`     |  Updates the connections with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                         Deletes the created port2ibm2 connection resources |
