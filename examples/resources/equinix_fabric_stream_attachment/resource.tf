resource "equinix_fabric_stream_attachment" "asset" {
  asset_id  = "<id_of_the_asset_being_attached>"
  asset     = "<asset_group>"
  stream_id = "<id_of_the_stream_asset_is_being_attached_to>"
}
