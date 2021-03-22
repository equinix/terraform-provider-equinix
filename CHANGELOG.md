## 1.1.0 (March 22, 2021)

BREAKING CHANGES:

- The environment variable `PACKET_AUTH_TOKEN` is deprecated. Use `METAL_AUTH_TOKEN`.
- `CustomData` fields may be parsed differently with packngo v0.6.0+

FIXES:

- `metal_project_ssh_key` `project_id` is reported correctly
- `metal_ssh_key` `project_id` is reported correctly
- Delete operations now treat 404 and 403 HTTP errors as successful

FEATURES:

- **New Data Source** `metal_facility`
- **New Provider Argument** `max_retries`
- **New Provider Argument** `max_retry_wait_seconds`
- Versioned User-Agent is reported by the HTTP Client
- `metal_device_network_type` field `type` now accepts `hybrid-bonded`

IMPROVEMENTS:

- Depends on packngo [v0.7.0](https://github.com/packethost/packngo/releases/tag/v0.7.0)
- Documentation corrections
- E2E testing has less test jitter and greater overall success
- Removed debug logs from metal\_operating\_system data source
- Removed debug logs from metal\_organization data source
- Removed debug logs from metal\_project data source

## 1.0.0 (December 08, 2020)

BREAKING CHANGES:

- [#1](https://github.com/equinix/terraform-provider-metal/issues/1)
  Users migrating from the Packet provider, please follow the instructions at
  the linked issue. In short, all v3.2.0 Packet provider resources have been
  renamed.

FEATURES:

- **New Resource** `metal_bgp_session`
- **New Resource** `metal_device`
- **New Resource** `metal_device_network_type`
- **New Resource** `metal_ip_attachment`
- **New Resource** `metal_organization`
- **New Resource** `metal_port_vlan_attachment`
- **New Resource** `metal_project`
- **New Resource** `metal_project`
- **New Resource** `metal_project_ssh_key`
- **New Resource** `metal_reserved_ip_block`
- **New Resource** `metal_spot_market_request`
- **New Resource** `metal_ssh_key`
- **New Resource** `metal_volume`
- **New Resource** `metal_vlan`
- **New Resource** `metal_volume_attachment`

- **New Data Resource** `metal_device`
- **New Data Resource** `metal_device_bgp_neighbors`
- **New Data Resource** `metal_ip_block_ranges`
- **New Data Resource** `metal_operating_system`
- **New Data Resource** `metal_organization`
- **New Data Resource** `metal_precreated_ip_block`
- **New Data Resource** `metal_project`
- **New Data Resource** `metal_project_ssh_key`
- **New Data Resource** `metal_spot_market_price`
- **New Data Resource** `metal_spot_market_request`
- **New Data Resource** `metal_volume`
