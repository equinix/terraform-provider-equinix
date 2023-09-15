# ECX Fabric Layer2 Two Redundant Connections from fabric cloud router to Azure

This example shows how create two redundant connections from Fabric Cloud Router to Azure, on ECX Fabric ports.

## Adjust variables
At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above

`pri_connection_name` - the name of the primary connection
`sec_connection_name` - the name of the secondary connection
`connection_type` - connection type, please refer schema
`notifications_type` - notification type
`notifications_emails` - List of emails
`bandwidth` - bandwidth in MBs
`aside_ap_type` - Fabric Cloud Router type
`peering_type` - Peering type for the ECX Fabric Cloud Router on the a-side; typically PRIVATE
**Note: You can use one Cloud Router for both connections if you would like**
`cloud_router_primary_uuid` - UUID of ECX Fabric Cloud Router on a-side
`cloud_router_secondary_uuid` - UUID of ECX Fabric Cloud Router on a-side for secondary connection
`zside_ap_type` - Z side access point type
`zside_ap_authentication_key` - Azure authorization key, service key generated from Azure Portal
`zside_ap_profile_type` - Service profile type
`zside_ap_profile_uuid` - Service profile UUID
`zside_location` - Seller location
`fabric_sp_name` - Service profile name, fetched based on Service Profile get call using Service Profile search schema

## Azure login

Log in to Azure portal use account that has permission to create necessary resources.

## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
  In the given example, the folder to perform CRUD operations on a fcr2port connection can be found at examples/fcr2port/.

- Change directory into - `CD cloudRouter2azure/two-connections`
- Initialize Terraform plugins - `terraform init`

## Fabric Cloud Router to port connection : Create, Read, Update and Delete(CRUD) operations
Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                             Description |
|:----------|:---------------------------------:|------------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                      Creates a fcr2azure redundant connection resources |
| READ      |         `terraform show`          |     Reads/Shows the current state of the fcr2azure connection resources |
| UPDATE    |    `terraform apply -refresh`     | Updates the fcr2azure with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                      Deletes the created fcr2azure connection resources |