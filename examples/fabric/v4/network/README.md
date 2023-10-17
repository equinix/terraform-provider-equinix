# ECX Fabric network CRUD operations
This example shows how to create Fabric network

Note: Each time you need to create a new Fabric Network resource, 
make a copy of the base folder - examples/fabric/v4/network and CD into this folder to perform all the CRUD operations.

## Define values for the Fabric network create
At minimum, you must set below variables in `terraform.tfvars` file:
  - `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
  - `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above
  - `network_name` - Name of ECX Fabric network 
  - `network_type` - Fabric network type
  - `network_scope` - Fabric network scope
  - `notifications_type` - notification type
  - `notifications_emails` - List of emails

## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
In the given example, the folder to perform CRUD operations on an network resource can be found at examples/fabric-cloud-router/.

- Change directory into - `CD examples/fabric-cloud-router/`
- Initialize Terraform plugins - `terraform init`

## Fabric network : Create, Read, Update and Delete(CRUD) operations
 Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                                Description |
|:----------|:---------------------------------:|---------------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                                                    Creates an network resource |
| READ      |         `terraform show`          |                          Reads/Shows the current state of the network resource |
| UPDATE    |    `terraform apply -refresh`     | Updates the network resource with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                                           Deletes the created network resource |