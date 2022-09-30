equinix_client_id      = "opYOL7xfD0HLshl2SyKAO4ebLn5uWbhQ"
equinix_client_secret  = "MdK7Qq8IfaSN9yBS"

port_names = ["ops-user100-CX-DC5-NL-Dot1q-BO-10G-PRI-JP-90","ops-user100-CX-SV1-NL-Qinq-STD-1G-PRI-NK-349"]


name = "terra-e2e-generic-sp"
description = " Terra Generic SP"
service_profile_type = "L2_PROFILE"
notifications_type = "BANDWIDTH_ALERT"
notifications_emails = ["example@equinix.com"]
tags = ["Storage", "Compute"]
visibility = "PRIVATE"
access_point_type_configs_type = "COLO"
access_point_type_configs_connection_redundancy_required = false
access_point_type_configs_allow_bandwidth_auto_approval = false
access_point_type_configs_allow_remote_connections = false
access_point_type_configs_connection_label = "Terra Test"
access_point_type_configs_enable_auto_generate_service_key = false
access_point_type_configs_bandwidth_alert_threshold = 10
access_point_type_configs_allow_custom_bandwidth = true
api_config_api_available = false
api_config_equinix_managed_vlan = true
api_config_bandwidth_from_api = false
api_config_integration_id = "test"
api_config_equinix_managed_port = true
authentication_key_required = false
authentication_key_label = "Service Key"
authentication_key_description = "XYZ"
access_point_type_configs_supported_bandwidths = [100, 500]
marketing_info_promotion = false