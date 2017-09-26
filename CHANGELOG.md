## 0.1.1 (Unreleased)

INTERNAL:

* provider: Add logging transport for HTTP client (`DEBUG` log now contains all HTTP requests & responses) [GH-25]

IMPROVEMENTS:

* resource/packet_device: Add `public_ipv4_subnet_size` field [GH-7]
* resource/packet_device: Add `hardware_reservation_id` field [GH-14]
* resource/packet_device: Add `root_password` attribute [GH-15]
* resource/packet_volume_attachment: Add resource for attaching volumes [GH-9]
* resource/packet_reserved_ip_block: Add resource for reserving blocks of IP addresses [GH-21]
* resource/packet_ip_attachment: Add resource for attaching floating IP addresses to devices [GH-21]

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
