package connection

import "testing"

func TestAccessPointAuthenticationKeyIsSensitive(t *testing.T) {
	if !accessPointSch().Schema["authentication_key"].Sensitive {
		t.Error("access_point.authentication_key is not sensitive")
	}
}
