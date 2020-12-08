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
