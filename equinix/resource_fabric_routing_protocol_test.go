package equinix

import "testing"

func TestFabricRoutingProtocolBgpAuthKeyIsSensitive(t *testing.T) {
	if !createFabricRoutingProtocolResourceSchema()["bgp_auth_key"].Sensitive {
		t.Error("bgp_auth_key is not sensitive")
	}
	if !readFabricRoutingProtocolResourceSchema()["bgp_auth_key"].Sensitive {
		t.Error("data source bgp_auth_key is not sensitive")
	}
}
