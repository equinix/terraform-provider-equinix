# ECX Fabric Gateway CRUD operations
This example shows how to create Fabric Gateway.

Note: Each time you need to create a new Fabric Gateway resource, 
make a copy of the base folder - examples/fabric-gateway/ and CD into this folder to perform all the CRUD operations.

## Define values for the Fabric Gateway create
At minimum, you must set below variables in `terraform.tfvars` file:
  - `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
  - `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above
  - `fg_name` - Name of ECX Fabric Gateway on a-side , i.e. amcrh007-fg
  - `fg_type` - Fabric Gateway type
  - `fg_location` - Fabric Gateway location
  - `fg_project` - Fabric Gateway project
  - `fg_account` - Fabric Gateway account
  - `fg_package` - Fabric Gateway package type, i.e. PRO
  - `notifications_type` - notification type
  - `notifications_emails` - List of emails

## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
In the given example, the folder to perform CRUD operations on an FG resource can be found at examples/fabric-gateway/.

- Change directory into - `CD examples/fabric-gateway/`
- Initialize Terraform plugins - `terraform init`

## Fabric Gateway : Create, Read, Update and Delete(CRUD) operations
 Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                               Description |
|:----------|:---------------------------------:|--------------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                                                    Creates an FG resource |
| READ      |         `terraform show`          |                          Reads/Shows the current state of the FG resource |
| UPDATE    |    `terraform apply -refresh`     | Updates the FG resource with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                                           Deletes the created FG resource |