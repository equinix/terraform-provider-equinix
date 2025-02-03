data "equinix_fabric_streams" "data_streams" {
  pagination = {
    limit = 2
    offset = 1
  }
}

output "number_of_returned_streams" {
  value = length(data.equinix_fabric_streams.data_streams.data)
}
