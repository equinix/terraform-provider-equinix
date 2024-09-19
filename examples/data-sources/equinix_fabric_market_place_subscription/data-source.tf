data "equinix_fabric_market_place_subscription" "subscription-test" {
  uuid = "<uuid_of_marketplace_subscription>"
}
output "id" {
  value = data.equinix_fabric_market_place_subscription.subscription-test.id
}
output "status" {
  value = data.equinix_fabric_market_place_subscription.subscription-test.status
}
output "marketplace" {
  value = data.equinix_fabric_market_place_subscription.subscription-test.marketplace
}
output "offer_type" {
  value = data.equinix_fabric_market_place_subscription.subscription-test.offer_type
}
output "is_auto_renew" {
  value = data.equinix_fabric_market_place_subscription.subscription-test.is_auto_renew
}
