# ECX Fabric Cloud Router Connection Routing Protocols CRUD operations
This example shows how to create Routing Protocols on FCR connection .

## Define values for the Fabric Cloud Router create
At minimum, you must set below variables in `terraform.tfvars` file:
- `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
- `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above
- `connection_uuid`- Connection uuid to apply the routing details to
- `direct_rp_name`- Name of Direct Routing Protocol instance 
- `equinix_ipv4_ip` - IPv4 for Direct Routing Protocol
- `equinix_ipv6_ip` - IPv6 for Direct Routing Protocol
- `bgp_rp_name` - Name of BGP Routing Protocol instance 
- `customer_peer_ipv4` - Customer Peering IPv4 for BGP Routing Protocol
- `customer_peer_ipv6` - Customer Peering IPv6 for BGP Routing Protocol
- `bgp_enabled_ipv4` - Enable flag for IPv4 on BGP Routing Protocol instance
- `bgp_enabled_ipv6` - Enable flag for IPv6 on BGP Routing Protocol instance
- `customer_asn` - Customer BGP ASN Number
- `equinix_asn` - Equinix BGP ASN Number (Will default if not supplied)

## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
  In the given example, the folder to perform CRUD operations on an RP resource can be found at examples/routing-protocol-bgp/.

- Change directory into - `CD examples/fabric/v4/cloudRouterConnectivity/routing-protocols/`
- Initialize Terraform plugins - `terraform init`

## Routing-protocol BGP : Create, Read, Update and Delete(CRUD) operations
Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                                              Description |
|:----------|:---------------------------------:|-----------------------------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                                                     Creates a routing_protocols resource |
| READ      |         `terraform show`          |                          Reads/Shows the current state of the routing_protocols resource |
| UPDATE    |    `terraform apply -refresh`     | Updates the routing_protocols resource with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                                           Deletes the created routing_protocols resource |