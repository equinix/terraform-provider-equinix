# ECX Fabric Cloud Router CRUD operations
This example shows how to create Fabric Cloud Router.

Note: Each time you need to create a new Fabric Cloud Router resource, 
make a copy of the base folder - examples/fabric-cloud-router/ and CD into this folder to perform all the CRUD operations.

## Define values for the Fabric Cloud Router create
At minimum, you must set below variables in `terraform.tfvars` file:
  - `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
  - `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above
  - `fcr_name` - Name of ECX Fabric Cloud Router on a-side , i.e. amcrh007-fcr
  - `fcr_type` - Fabric Cloud Router type
  - `fcr_location` - Fabric Cloud Router location
  - `fcr_project` - Fabric Cloud Router project
  - `fcr_account` - Fabric Cloud Router account
  - `fcr_package` - Fabric Cloud Router package type, i.e. PRO
  - `notifications_type` - notification type
  - `notifications_emails` - List of emails

## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
In the given example, the folder to perform CRUD operations on an FCR resource can be found at examples/fabric-cloud-router/.

- Change directory into - `CD examples/fabric-cloud-router/`
- Initialize Terraform plugins - `terraform init`

## Fabric Cloud Router : Create, Read, Update and Delete(CRUD) operations
 Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                                Description |
|:----------|:---------------------------------:|---------------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                                                    Creates an FCR resource |
| READ      |         `terraform show`          |                          Reads/Shows the current state of the FCR resource |
| UPDATE    |    `terraform apply -refresh`     | Updates the FCR resource with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                                           Deletes the created FCR resource |