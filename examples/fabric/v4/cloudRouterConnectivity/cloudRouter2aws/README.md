# ECX Fabric Layer2 Connection from fabric cloud router to AWS

This example shows how create connection from Fabric Cloud Router to AWS, on ECX Fabric ports.

## Adjust variables
At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above

`fcr_uuid` - UUID of ECX Fabric Cloud Router on a-side 
`zside_port_name` -  Name of ECX Fabric z-side port , i.e. ops-user100-CX-SV5-NL-Qinq-BO-10G-SEC-JP-000
`connection_name` - the name of the connection
`connection_type` - connection type, please refer schema
`notifications_type` - notification type
`notifications_emails` - List of emails
`bandwidth` - bandwidth in MBs
`redundancy` - Port redundancy
`aside_ap_type` - Fabric Cloud Router type
`zside_ap_type` - Z side access point type
`zside_ap_authentication_key` - AWS authorization key, account number like 357848912121
`zside_ap_profile_type` - Service profile type
`fabric_sp_name` - Service profile name, fetched based on Service Profile get call using Service Profile search schema
`zside_location` - Seller location
`seller_region` - Seller region code

## AWS login

Log in to AWS portal use account that has permission to create necessary resources.

## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
  In the given example, the folder to perform CRUD operations on a fcr2port connection can be found at examples/fcr2port/.

- Change directory into - `CD fcr2aws/`
- Initialize Terraform plugins - `terraform init`

## Fabric Cloud Router to port connection : Create, Read, Update and Delete(CRUD) operations
Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                            Description |
|:----------|:---------------------------------:|-----------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                                 Creates a fcr2port connection resource |
| READ      |         `terraform show`          |      Reads/Shows the current state of the fcr2port connection resource |
| UPDATE    |    `terraform apply -refresh`     | Updates the fcr2port with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                       Deletes the created fcr2port connection resource |