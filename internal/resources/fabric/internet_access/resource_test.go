package internetaccess_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const eiaResource = "equinix_fabric_internet_access_service.test"

func TestAccFabricInternetAccessServiceCreate_PFCR(t *testing.T) {
	testData := testinghelpers.GetFabricInternetAccessTestData(t)
	pfcr, ok := testData["pfcr"]
	if !ok {
		t.Skip("skipping: no 'pfcr' key in TF_ACC_FABRIC_INTERNET_ACCESS_TEST_DATA")
	}

	vlan, err := testinghelpers.RandomVlan(pfcr["port_uuid"])
	if err != nil {
		t.Fatalf("unable to get an available VLAN: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInternetAccessServiceDeleted,
		Steps: []resource.TestStep{
			{
				Config: testAccInternetAccessServiceConfig(pfcr, vlan),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(eiaResource, tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(eiaResource, tfjsonpath.New("state"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(eiaResource, tfjsonpath.New("type"), knownvalue.StringExact("SINGLE_IA")),
					statecheck.ExpectKnownValue(eiaResource, tfjsonpath.New("name"), knownvalue.StringExact("tf_acc_eia_PFCR")),
					statecheck.ExpectKnownValue(eiaResource, tfjsonpath.New("bandwidth"), knownvalue.Int32Exact(50)),
					statecheck.ExpectKnownValue(eiaResource, tfjsonpath.New("routing_protocol"),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"type": knownvalue.StringExact("STATIC"),
						}),
					),
				},
			},
			{
				ResourceName:      eiaResource,
				ImportState:       true,
				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"timeouts", "state"},
			},
		},
	})
}

func testAccInternetAccessServiceConfig(data map[string]string, vlan int) string {
	return fmt.Sprintf(`
resource "equinix_fabric_connection" "ia" {
  name = "tf_acc_eia_vc_PFCR"
  type = "IA_VC"
  notifications {
    type   = "ALL"
    emails = ["example@equinix.com"]
  }
  bandwidth = 50
  project {
    project_id = %[1]q
  }
  a_side {
    access_point {
      type = "COLO"
      port {
        uuid = %[2]q
      }
      link_protocol {
        type     = "DOT1Q"
        vlan_tag = %[3]d
      }
    }
  }
  z_side {
    access_point {
      type = "SP"
      profile {
        type = "IA_PROFILE"
        uuid = %[4]q
      }
      location {
        metro_code = %[5]q
      }
    }
  }
}

resource "equinix_fabric_ip_block" "route" {
  type = "IPV4_IP_BLOCK"
  location = {
    metro_code = %[5]q
  }
  prefix_length = 30
  project = {
    project_id = %[1]q
  }
  account = {
    account_number = %[6]q
  }
}

resource "equinix_fabric_internet_access_service" "test" {
  name      = "tf_acc_eia_PFCR"
  type      = "SINGLE_IA"
  bandwidth = 50
  project = {
    project_id = %[1]q
  }
  account = {
    account_number = %[6]q
  }
  billing = {
    type = "FIXED"
  }
  routing_protocol = {
    type = "STATIC"
    customer_routes = [{
      ip_block = {
        uuid = equinix_fabric_ip_block.route.uuid
      }
    }]
    connections = [{
      uuid = equinix_fabric_connection.ia.id
    }]
  }
}
`, data["project_id"], data["port_uuid"], vlan, data["service_profile_uuid"], data["metro_code"], data["account_number"])
}

func checkInternetAccessServiceDeleted(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_internet_access_service" {
			continue
		}
		svc, resp, err := client.InternetAccessServicesApi.GetEiaService(ctx, rs.Primary.ID).Execute()
		if err != nil {
			if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 403) {
				continue
			}
			return fmt.Errorf("error checking Internet Access service %s: %w", rs.Primary.ID, err)
		}
		if svc.GetState() != fabricv4.INTERNETACCESSSERVICESTATE_DEPROVISIONED {
			return fmt.Errorf("Internet Access service %s still exists in state %s", rs.Primary.ID, svc.GetState())
		}
	}
	return nil
}
