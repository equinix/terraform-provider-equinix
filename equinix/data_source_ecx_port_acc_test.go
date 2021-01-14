package equinix

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestAccFabricPort(t *testing.T) {
	t.Parallel()
	portName, _ := schema.EnvDefaultFunc(priPortEnvVar, "smandalika@equinix.com1-SV1-Dot1q-L-Primary-161350")()
	context := map[string]interface{}{
		"port-resourceName": "test",
		"port-name":         portName,
	}
	resourceName := fmt.Sprintf("data.equinix_ecx_port.%s", context["port-resourceName"].(string))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withPort().build(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestMatchResourceAttr(resourceName, "ibx", regexp.MustCompile(`^[A-Z]{2}[0-9]+$`)),
					resource.TestMatchResourceAttr(resourceName, "metro_code", regexp.MustCompile(`^[A-Z]{2}$`)),
					resource.TestMatchResourceAttr(resourceName, "priority", regexp.MustCompile(`^(Primary|Secondary)$`)),
					resource.TestMatchResourceAttr(resourceName, "encapsulation", regexp.MustCompile(`^(Dot1q|Qinq)$`)),
					resource.TestMatchResourceAttr(resourceName, "buyout", regexp.MustCompile(`^(true|false)$`)),
					resource.TestMatchResourceAttr(resourceName, "bandwidth", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

func testAccECXPortt(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_ecx_port" "%{resourceName}" {
  name = "%{name}"
}
`, ctx)
}
