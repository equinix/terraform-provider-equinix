# ECX Fabric Layer2 Connection to Azure

This example shows how to create layer 2 connection between ECX Fabric port and AZURE Cloud.
Example covers **provisioning of both sides** of the connection.

## Adjust variables

At minimum, you must set below variables in `terraform.tfvars` file:

- `equinix_client_id` - Equinix client ID (consumer key), obtained after registering app in the developer platform
- `equinix_client_secret` - Equinix client secret ID (consumer secret), obtained same way as above
- `fabric_sp_name` - Service profile name like i.e. AZURE
- `equinix_port_name` -  Name of ECX Fabric port that should be connected to AZURE, i.e. ops-user100-CX-SV5-NL-Qinq-BO-10G-SEC-JP-199
- `connection_name` - the name of the connection
- `connection_type` - connection type, please refer schema
- `notifications_type` - notification type
- `notifications_emails` - List of emails
- `bandwidth` - bandwidth in MBs
- `redundancy` - Port redundancy
- `aside_ap_type` - Access point type
- `peering_type` - Peering type
- `aside_port_uuid` - Port uuid, fetched based on port call using Port resource
- `aside_link_protocol_type` - link protocol type
- `aside_link_protocol_stag` - s tag number
- `aside_link_protocol_ctag` - c tag number
- `zside_ap_type` - Z side access point type
- `zside_ap_authentication_key` - AZURE authorization key, like c620477c-3f30-41e8-a0b9-cfdb4a12121d
- `zside_ap_profile_type` - Service profile type
- `zside_ap_profile_uuid` - Service profile uuid, fetched based on Service Profile get call using Service Profile search schema
- `zside_location` - Seller location

## Azure login

Log in to Azure portal use account that has permission to create necessary resources.

## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
  In the given example, the folder to perform CRUD operations for port2azure connection can be found at examples/fabric/v4/portConnectivity/azure/.

- Change directory into - `CD examples/fabric/v4/portConnectivity/azure/`
- Initialize Terraform plugins - `terraform init`

## Port to Azure connection  : Create, Read, Update and Delete(CRUD) operations
Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                              Description |
|:----------|:---------------------------------:|-------------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                                 Creates a port2azure connection resource |
| READ      |         `terraform show`          |      Reads/Shows the current state of the port2azure connection resource |
| UPDATE    |    `terraform apply -refresh`     | Updates the port2azure with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                       Deletes the created port2azure connection resource |
