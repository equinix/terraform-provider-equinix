
# Create a new read-only user API key

resource "equinix_metal_user_api_key" "test" {
  description = "Read-only user key"
  read_only   = true
}
