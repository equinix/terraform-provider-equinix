# ECX Fabric Layer2 Connection between two own ports

This example shows how create layer 2 connection between two, own ECX Fabric ports.

## Adjust variables
At minimum, you must set below variables in `terraform.tfvars` file:

- `equinix_client_id` - Equinix client ID (consumer key), obtained after registering app in the developer platform
- `equinix_client_secret` - Equinix client secret ID (consumer secret), obtained same way as above
- `aside_port_name` - Name of ECX Fabric a-side port i.e. ops-user100-CX-SV5-NL-Qinq-STD-1G-SEC-JP-111
- `zside_port_name` -  Name of ECX Fabric z-side port , i.e. ops-user100-CX-SV5-NL-Qinq-BO-10G-SEC-JP-000
- `connection_name` - the name of the connection
- `connection_type` - connection type, please refer schema
- `notifications_type` - notification type
- `notifications_emails` - List of emails
- `bandwidth` - bandwidth in MBs
- `redundancy` - Port redundancy
- `aside_ap_type` - Access point type
- `aside_port_uuid` - Port uuid, fetched based on port call using Port resource
- `aside_link_protocol_type` - link protocol type
- `aside_link_protocol_stag` - a-side s tag number
- `zside_ap_type` - Z side access point type
- `aside_link_protocol_stag` - z-side s tag number
- `zside_location` - Seller location

## Initialize
- First step is to initialize the terraform directory/resource we are going to work on.
  In the given example, the folder to perform CRUD operations for port2portself connection can be found at examples/fabric/v4/portConnectivity/port2portself/.

- Change directory into - `CD examples/fabric/v4/portConnectivity/port2portself/`
- Initialize Terraform plugins - `terraform init`

## Port to Port2Portself connection  : Create, Read, Update and Delete(CRUD) operations
Note: `–auto-approve` command does not prompt the user for validating the applying config. Remove it to get a prompt to confirm the operation.

| Operation |              Command              |                                                                 Description |
|:----------|:---------------------------------:|----------------------------------------------------------------------------:|
| CREATE    |  `terraform apply –auto-approve`  |                                 Creates a port2portself connection resource |
| READ      |         `terraform show`          |      Reads/Shows the current state of the port2portself connection resource |
| UPDATE    |    `terraform apply -refresh`     | Updates the port2portself with values provided in the terraform.tfvars file |
| DELETE    | `terraform destroy –auto-approve` |                       Deletes the created port2portself connection resource |