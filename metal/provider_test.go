package metal

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"metal": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

/*
func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}
*/

func testAccPreCheck(t *testing.T) {
	v := os.Getenv("METAL_AUTH_TOKEN")

	if v == "" {
		v = os.Getenv("PACKET_AUTH_TOKEN")
	}

	if v == "" {
		t.Fatal("METAL_AUTH_TOKEN must be set for acceptance tests")
	}
}
