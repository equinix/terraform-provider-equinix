# ECX Fabric Layer2 Connection to Oracle

This example shows how to create layer 2 connection between ECX Fabric port
and Oracle Cloud.
Example covers **provisioning of both sides** of the connection.

## Adjust variables

At minimum, you must set below variables in `terraform.tfvars` file:

* `equinix_client_id` - Equinix client ID (consumer key), obtained after
  registering app in the developer platform
* `equinix_client_secret` - Equinix client secret ID (consumer secret),
  obtained same way as above

`fabric_sp_name` - Service profile name like i.e. **Oracle Cloud Infrastructure -OCI- FastConnect**
`equinix_port_name` -  Name of ECX Fabric port that should be
connected to Oracle, i.e. ops-user100-CX-SV5-NL-Qinq-BO-10G-SEC-JP-000
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
`zside_ap_authentication_key` - Oracle authorization key following a pattern, like **ocid1.virtualcircuit.oc1.phx.aaaaaaaa62hgxtczqwpaour5gnaq4qn2wupjwpqtsknr2mvxblrcdtusda1a**
`zside_ap_profile_type` - Service profile type
`zside_ap_profile_uuid` - Service profile uuid, fetched based on Service Profile get call using Service Profile search schema
`zside_location` - Seller location
`seller_region` - Seller region code

## Oracle login

Log in to Oracle portal use account that has permission to create necessary resources.

## Initialize

Change directory to project root to run terra test or change to example directory and initialize Terraform plugins
by running `terraform init`.

## Deploy template

Apply changes by running `terraform apply`, then **inspect proposed plan**
and approve it.
