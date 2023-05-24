# ECX Fabric Layer2 Connection between two own ports

This example shows how create layer 2 connection between two, own ECX Fabric ports.

## Adjust variables
At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above

`aside_port_name` - Name of ECX Fabric a-side port i.e. ops-user100-CX-SV5-NL-Qinq-STD-1G-SEC-JP-111
`zside_port_name` -  Name of ECX Fabric z-side port , i.e. ops-user100-CX-SV5-NL-Qinq-BO-10G-SEC-JP-000

`connection_name` - the name of the connection
`connection_type` - connection type, please refer schema
`notifications_type` - notification type
`notifications_emails` - List of emails
`bandwidth` - bandwidth in MBs
`redundancy` - Port redundancy
`aside_ap_type` - Access point type
`aside_port_uuid` - Port uuid, fetched based on port call using Port resource
`aside_link_protocol_type` - link protocol type
`aside_link_protocol_stag` - a-side s tag number
`zside_ap_type` - Z side access point type
`aside_link_protocol_stag` - z-side s tag number
`zside_location` - Seller location

## Initialize

Change directory to project root to run terra test or change to example directory and initialize Terraform plugins
by running `terraform init`.

## Deploy template

Apply changes by running `terraform apply`, then **inspect proposed plan**
and approve it.
