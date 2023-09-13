# ECX Fabric Layer2 Single Connection from fabric cloud router to Azure

This example shows how create single connection from Fabric Cloud Router to Azure, on ECX Fabric ports.

## Adjust variables
At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above

`fcr_uuid` - UUID of ECX Fabric Cloud Router on a-side
`connection_name` - the name of the connection
`connection_type` - connection type, please refer schema
`notifications_type` - notification type
`notifications_emails` - List of emails
`bandwidth` - bandwidth in MBs
`redundancy` - Port redundancy **Note: in order to use resource dependency with redundancy because it is a set type with
max items of 1; put it into the one() terraform function before attempting to address its children. 
I.e. one(equinix_fabric_connection.connection_name.redundancy).group**
`aside_ap_type` - Fabric Cloud Router type
`zside_ap_type` - Z side access point type
`zside_ap_authentication_key` - Azure authorization key, service key generated from Azure Portal
`zside_ap_profile_type` - Service profile type
`fabric_sp_name` - Service profile name, fetched based on Service Profile get call using Service Profile search schema
`zside_location` - Seller location
`seller_region` - Seller region code

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
| CREATE    |  `terraform apply –auto-approve`  |                                 Creates a fcr2azure connection resource |
| READ      |         `terraform show`          |      Reads/Shows the current state of the fcr2azure connection resource |
| UPDATE    |    `terraform apply -refresh`     | Updates the fcr2azure with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                    Deletes the created fcr2azure connection resource(s) |