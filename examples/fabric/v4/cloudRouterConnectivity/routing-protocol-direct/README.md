# ECX Fabric Cloud Router Connection RP CRUD operations
This example shows how to Config Direct Routing Protocol details for an existing FCR connection .

Note: Each time you need to create a Direct RP resource add-on
make a copy of the base folder - examples/fabric/v4/cloudRouterConnectivity/routing-protocol-direct/ and CD into this folder to perform all the CRUD operations.

## Define values for the Fabric Cloud Router create
At minimum, you must set below variables in `terraform.tfvars` file:
- `equinix_client_id` - Equinix client ID (consumer key), obtained after registering app in the developer platform
- `equinix_client_secret` - Equinix client secret ID (consumer secret),obtained same way as above
- `rp_name`- Name of routing Protocol
- `rp_type`- Type of routing Protocol entity, "DIRECT" or "BGP"
- `connection_uuid` = FCR Connection UUID
- `equinix_ipv4_ip` = Equinix Side IpV4 Address
- `equinix_ipv6_ip` = Equinix Side IpV6 Address


## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
  In the given example, the folder to perform CRUD operations on an RP resource can be found at examples/fabric/v4/cloudRouterConnectivity/routing-protocol-direct/.

- Change directory into - `CD examples/fabric/v4/cloudRouterConnectivity/routing-protocol-direct/`
- Initialize Terraform plugins - `terraform init`

## Routing-protocol Direct IP : Create, Read, Update and Delete(CRUD) operations
Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                                                                                            Description |
|:----------|:---------------------------------:|---------------------------------------------------------------------------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                                                                                             Creates a Direct Routing Protocol resource |
| READ      |         `terraform show`          |                                                                  Reads/Shows the current state of the Direct Routing Protocol resource |
| UPDATE    |    `terraform apply -refresh`     |                                         Updates the Direct Routing Protocol resource with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                                                                                   Deletes the created Direct Routing Protocol resource |