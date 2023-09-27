# ECX Fabric Layer2 Connection from fabric cloud router to ipwan

This example shows how create connection from Fabric Cloud Router to ipwan, on ECX Fabric ports.

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
`zside_ap_type` - Z side access point type, ipwan

## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
  In the given example, the folder to perform CRUD operations on a fcr2wan connection can be found at examples/fcr2port/.

- Change directory into - `CD examples/cloudRouter2wan/`
- Initialize Terraform plugins - `terraform init`

## Fabric Cloud Router to wan connection : Create, Read, Update and Delete(CRUD) operations
Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                           Description |
|:----------|:---------------------------------:|----------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                                 Creates a fcr2wan connection resource |
| READ      |         `terraform show`          |      Reads/Shows the current state of the fcr2wan connection resource |
| UPDATE    |    `terraform apply -refresh`     | Updates the fcr2wan with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                       Deletes the created fcr2wan connection resource |