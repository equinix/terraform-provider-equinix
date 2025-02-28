data "equinix_fabric_stream_subscriptions" "all" {
  stream_id = "<stream_id>"
  pagination = {
    limit  = 10
    offset = 0
  }
}
