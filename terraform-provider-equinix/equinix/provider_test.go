package equinix

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"equinix": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(endpointEnvVar); v == "" {
		t.Fatalf("%s env variable has to be set", endpointEnvVar)
	}
	if v := os.Getenv(clientIDEnvVar); v == "" {
		t.Fatalf("%s env variable has to be set", clientIDEnvVar)
	}
	if v := os.Getenv(clientSecretEnvVar); v == "" {
		t.Fatalf("%s env variable has to be set", clientSecretEnvVar)
	}
}

func nprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		var strVal string
		switch val.(type) {
		case []string:
			r := regexp.MustCompile(`" "`)
			strVal = r.ReplaceAllString(fmt.Sprintf("%q", val), `", "`)
		default:
			strVal = fmt.Sprintf("%v", val)
		}
		format = strings.Replace(format, "%{"+key+"}", strVal, -1)
	}
	return format
}
