package equinix

import (
	"fmt"
	"testing"

	"github.com/equinix/ecx-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccECXL2ConnectionAccepter(t *testing.T) {
	t.Parallel()
	portName, _ := schema.EnvDefaultFunc(priPortEnvVar, "sit-001-CX-SV1-NL-Dot1q-BO-10G-PRI-JUN-33")()
	spName, _ := schema.EnvDefaultFunc(priSpEnvVar, "AWS Direct Connect")()
	context := map[string]interface{}{
		"connectionResourceName": "conn",
		"name":                   fmt.Sprintf("tf-tst-%s", randString(6)),
		"profile_name":           spName.(string),
		"speed":                  50,
		"speed_unit":             "MB",
		"notifications":          []string{"marry@equinix.com", "john@equinix.com"},
		"port_name":              portName.(string),
		"vlan_stag":              randInt(2000),
		"seller_region":          "us-west-2",
		"seller_metro_code":      "SV",
		"authorization_key":      "123456789012",
		"accepterResourceName":   "accepter",
		"access_key":             "AKIAGGJKJU7BC3QJKYQ",
		"secret_key":             "CXGJW1lWbqENEqSkBAK",
	}
	connectionResourceName := fmt.Sprintf("equinix_ecx_l2_connection.%s", context["connectionResourceName"].(string))
	accepterResourceName := fmt.Sprintf("equinix_ecx_l2_connection_accepter.%s", context["accepterResourceName"].(string))
	var testConn ecx.L2Connection
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccECXL2ConnectionAccepter(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(accepterResourceName, "connection_id"),
					resource.TestCheckResourceAttrPair(accepterResourceName, "connection_id", connectionResourceName, "id"),
					resource.TestCheckResourceAttrSet(accepterResourceName, "aws_connection_id"),
					testAccECXL2ConnectionExists(connectionResourceName, &testConn),
					testAccECXL2ConnectionAccepterStatus(&testConn),
				),
			},
		},
	})
}

func testAccECXL2ConnectionAccepter(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_ecx_l2_sellerprofile" "aws" {
  name = "%{profile_name}"
}
	  
data "equinix_ecx_port" "port" {
  name = "%{port_name}"
}

locals {
  aws_first_metro        = tolist(data.equinix_ecx_l2_sellerprofile.aws.metro)[0]
  aws_first_metro_region = keys(local.aws_first_metro.regions)[0]
}

resource "equinix_ecx_l2_connection" "%{connectionResourceName}" {
  name              = "%{name}"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.aws.id
  speed             = %{speed}
  speed_unit        = "%{speed_unit}"
  notifications     = %{notifications}
  port_uuid         = data.equinix_ecx_port.port.id
  vlan_stag         = %{vlan_stag}
  seller_region     = local.aws_first_metro_region
  seller_metro_code = local.aws_first_metro.code
  authorization_key = "%{authorization_key}"
}

resource "equinix_ecx_l2_connection_accepter" "%{accepterResourceName}" {
   connection_id = equinix_ecx_l2_connection.%{connectionResourceName}.id
   access_key    = "%{access_key}"
   secret_key    = "%{secret_key}"
}
`, ctx)
}

func testAccECXL2ConnectionAccepterStatus(conn *ecx.L2Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if conn.ProviderStatus != ecx.ConnectionStatusProvisioned {
			return fmt.Errorf("provider_status does not match %v - %v", ecx.ConnectionStatusProvisioned, conn.ProviderStatus)
		}
		return nil
	}
}
