equinix_client_id      = "MyEquinixClientId"
equinix_client_secret  = "MyEquinixSecret"

connection_name = "terra_e2e_ibm"
description = "Test Connection"
connection_type = "EVPL_VC"
notifications_type = "ALL"
notifications_emails = ["example@equinix.com"]
bandwidth = 50
redundancy = "PRIMARY"
purchase_order_number = "1-323292"
aside_ap_type = "COLO"
aside_link_protocol_type = "QINQ"
aside_link_protocol_stag = "2019"
aside_link_protocol_ctag = "1921"
zside_ap_type = "SP"
zside_ap_authentication_key = "5bf92b31d921499f963592cd816f6be7"
zside_ap_profile_type = "L2_PROFILE"
zside_location = "SV"
seller_region = "San Jose 2"
fabric_sp_name = "IBM Cloud Direct Link 2"
equinix_port_name = "ops-user100-CX-SV1-NL-Qinq-STD-10G-PRI-NK-333"
additional_info = [{"name":"ASN","value":"1232"},{"name":"CER IPv4 CIDR","value":"10.254.0.0/16"},{"name":"IBM IPv4 CIDR","value":"172.16.0.0/12"}]

