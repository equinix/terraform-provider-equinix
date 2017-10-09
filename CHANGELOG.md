## 1.1.0 (October 09, 2017)

INTERNAL:

* Capture breaking change in the Packet API ((https://github.com/packethost/packngo/pull/47))

## 1.0.0 (September 28, 2017)

INTERNAL:

* provider: Add logging transport for HTTP client (`DEBUG` log now contains all HTTP requests & responses) ([#25](https://github.com/terraform-providers/terraform-provider-packet/issues/25))

IMPROVEMENTS:

* resource/packet_device: Add `public_ipv4_subnet_size` field ([#7](https://github.com/terraform-providers/terraform-provider-packet/issues/7))
* resource/packet_device: Add `hardware_reservation_id` field ([#14](https://github.com/terraform-providers/terraform-provider-packet/issues/14))
* resource/packet_device: Add `root_password` attribute ([#15](https://github.com/terraform-providers/terraform-provider-packet/issues/15))
* resource/packet_volume_attachment: Add resource for attaching volumes ([#9](https://github.com/terraform-providers/terraform-provider-packet/issues/9))
* resource/packet_reserved_ip_block: Add resource for reserving blocks of IP addresses ([#21](https://github.com/terraform-providers/terraform-provider-packet/issues/21))
* resource/packet_ip_attachment: Add resource for attaching floating IP addresses to devices ([#21](https://github.com/terraform-providers/terraform-provider-packet/issues/21))
* resource/packet_devce: Allow to use next-available hardware reservation ([#26](https://github.com/terraform-providers/terraform-provider-packet/issues/26))
* resource/packet_reserved_ip_block: Make reserved IP block resource importable ([#29](https://github.com/terraform-providers/terraform-provider-packet/issues/29))


## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
