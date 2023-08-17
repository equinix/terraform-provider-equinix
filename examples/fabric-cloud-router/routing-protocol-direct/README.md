# ECX FG Connection RP CRUD operations
This example shows how to create Config Direct RP on FG connection .

Note: Each time you need to create a RP resource add-on
make a copy of the base folder - examples/routing-protocol-direct/ and CD into this folder to perform all the CRUD operations.

## Define values for the Fabric Gateway create
At minimum, you must set below variables in `terraform.tfvars` file:
- `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
- `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above
- `rp_name`- Name of routing Protocol
- `rp_type`- Type of routing Protocol entity, "DIRECT" or "BGP"
- connection_uuid = "d557cb4c-9052-4298-b5ca-8a9ed914cf03"
  rp_type = "DIRECT"
  rp_name = "FG-RP"
  equinix_ipv4_ip = "192.1.1.1/30"
  equinix_ipv6_ip = "192::1:1/126"


## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
  In the given example, the folder to perform CRUD operations on an RP resource can be found at examples/routing-protocol-direct/.

- Change directory into - `CD examples/routing-protocol-direct/`
- Initialize Terraform plugins - `terraform init`

## Routing-protocol Direct IP : Create, Read, Update and Delete(CRUD) operations
Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                               Description |
|:----------|:---------------------------------:|--------------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                                                    Creates an FG resource |
| READ      |         `terraform show`          |                          Reads/Shows the current state of the FG resource |
| UPDATE    |    `terraform apply -refresh`     | Updates the FG resource with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                                           Deletes the created FG resource |