data "equinix_fabric_stream_attachments" "all" {
  pagination = {
    limit  = 100
    offset = 0
  }
  filters = [{
    property = "<filter_property>"
    operator = "="
    values   = ["<list_of_values_to_filter>"]
  }]
  sort = [{
    direction = "<DESC|ASC>"
    property  = "/uuid"
  }]
}
