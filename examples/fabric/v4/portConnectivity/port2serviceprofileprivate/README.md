# ECX Fabric Layer2 private seller profile

This example shows how to create layer2 public seller profile.

## Adjust variables
* `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above

`fabric_sp_name` - Service profile name like i.e. **Equinix Direct Connect - Private**
`aside_port_name` -  Name of ECX Fabric port that should be connected used for the connection, i.e. ops-user100-CX-SV5-NL-Qinq-BO-10G-SEC-JP-000

`connection_name` - the name of the connection
`connection_type` - connection type, please refer schema
`notifications_type` - notification type
`notifications_emails` - List of emails
`bandwidth` - bandwidth in MBs
`redundancy` - Port redundancy
`aside_ap_type` - Access point type
`aside_port_uuid` - Port uuid, fetched based on port call using Port resource
`aside_link_protocol_type` - link protocol type
`aside_link_protocol_tag` - tag number
`zside_ap_type` - Z side access point type
`zside_ap_profile_type` - Service profile type
`zside_ap_profile_uuid` - Service profile uuid, fetched based on Service Profile get call using Service Profile search schema
`zside_location` - Seller location

## Initialize

Change directory to example directory and initialize Terraform plugins
by running `terraform init`.

## Deploy template

Apply changes by running `terraform apply`, then **inspect proposed plan**
and approve it.
