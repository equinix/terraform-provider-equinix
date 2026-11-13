package equinix_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("equinix_fabric_cloud_router_PFCR", &resource.Sweeper{
		Name: "equinix_fabric_cloud_router",
		F:    testSweepCloudRouters,
	})
}

func testSweepCloudRouters(_ string) error {
	return nil
}

func TestAccCloudRouterCreateOnlyRequiredParameters_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkCloudRouterDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRouterCreateOnlyRequiredParameterConfig_PFCR("fcr_acctest_PFCR"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "name", "fcr_acctest_PFCR"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "type", "XF_ROUTER"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "order.0.purchase_order_number", "1-234567"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "location.0.metro_code", "SV"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "package.0.code", "STANDARD"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "project.0.project_id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "account.0.account_number"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "href"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "state", "PROVISIONED"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "equinix_asn"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.test", "connections_count", "0"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "change_log.0.updated_date_time"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.test", "change_log.0.deleted_date_time"),
				),
			},
			{
				Config: testAccCloudRouterCreateOnlyRequiredParameterConfig_PFCR("fcr_update_PFCR"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_cloud_router.test", "name", "fcr_update_PFCR"),
				),
			},
		},
	})
}

func testAccCloudRouterCreateOnlyRequiredParameterConfig_PFCR(name string) string {
	return fmt.Sprintf(`resource "equinix_fabric_cloud_router" "test"{
		type = "XF_ROUTER"
		name = "%s"
		location{
			metro_code  = "SV"
		}
		package{
			code = "STANDARD"
		}
		order{
			purchase_order_number = "1-234567"
		}
		notifications{
			type = "ALL"
			emails = [
				"test@equinix.com",
				"test1@equinix.com"
			]
		}
		project{
			project_id = "33ec651f-cc99-48e0-94d3-47466899cdc7"
		}
		account {
			account_number = 201257
		}
	}`, name)
}

func TestAccCloudRouterCreateMixedParameters_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkCloudRouterDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRouterCreateMixedParameterConfig_PFCR(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_cloud_router.example", "name", "fcr_acc_test_PFCR"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "type", "XF_ROUTER"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "notifications.0.type", "ALL"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "notifications.0.emails.0", "test@equinix.com"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "order.0.purchase_order_number", "1-234567"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "location.0.metro_code", "SV"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "location.0.region", "AMER"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "location.0.metro_name", "Silicon Valley"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "package.0.code", "STANDARD"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "project.0.project_id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "account.0.account_number"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "href"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "state", "PROVISIONED"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "equinix_asn"),
					resource.TestCheckResourceAttr("equinix_fabric_cloud_router.example", "connections_count", "0"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "change_log.0.created_by"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "change_log.0.created_date_time"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "change_log.0.updated_date_time"),
					resource.TestCheckResourceAttrSet("equinix_fabric_cloud_router.example", "change_log.0.deleted_date_time"),
				),
			},
		},
	})
}
func testAccCloudRouterCreateMixedParameterConfig_PFCR() string {
	return `
		resource "equinix_fabric_cloud_router" "example"{
		type = "XF_ROUTER"
		name = "fcr_acc_test_PFCR"
		location{
			region      = "AMER"
			metro_code  = "SV"
			metro_name = "Silicon Valley"
		}
		package{
			code = "STANDARD"
		}
		order{
			purchase_order_number = "1-234567"
		}
		notifications{
			type = "ALL"
			emails = [
				"test@equinix.com",
				"test1@equinix.com"
					]
		}
		project{
			project_id = "33ec651f-cc99-48e0-94d3-47466899cdc7"
		}
		account {
			account_number = 201257
		}
	}`
}

func checkCloudRouterDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_cloud_router" {
			continue
		}
		err := equinix.WaitUntilCloudRouterDeprovisioned(ctx, rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for cloud router deletion. ID: %s, Err: %s", rs.Primary.ID, err)
		}
	}
	return nil
}

func testCloudRouterReadConfig(t *testing.T, status int) *config.Config {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`[{"errorCode":"EQ-2040012","errorMessage":"Cloud Router does not exist"}]`))
	}))
	t.Cleanup(srv.Close)
	return &config.Config{BaseURL: srv.URL, Token: "test-token"}
}

func TestCloudRouterRead_notFoundRemovesFromState(t *testing.T) {
	r := equinix.Provider().ResourcesMap["equinix_fabric_cloud_router"]
	d := r.TestResourceData()
	d.SetId("00000000-0000-0000-0000-000000000000")

	diags := r.ReadContext(context.Background(), d, testCloudRouterReadConfig(t, http.StatusNotFound))

	if diags.HasError() {
		t.Fatalf("expected no error on 404, got %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("expected cloud router to be removed from state, got id %q", d.Id())
	}
}

func TestCloudRouterRead_otherErrorsAreReturned(t *testing.T) {
	r := equinix.Provider().ResourcesMap["equinix_fabric_cloud_router"]
	d := r.TestResourceData()
	d.SetId("00000000-0000-0000-0000-000000000000")

	diags := r.ReadContext(context.Background(), d, testCloudRouterReadConfig(t, http.StatusForbidden))

	if !diags.HasError() {
		t.Fatal("expected an error on 403")
	}
}

func TestCloudRouterDataSourceRead_notFoundIsAnError(t *testing.T) {
	ds := equinix.Provider().DataSourcesMap["equinix_fabric_cloud_router"]
	d := ds.TestResourceData()
	if err := d.Set("uuid", "00000000-0000-0000-0000-000000000000"); err != nil {
		t.Fatal(err)
	}

	diags := ds.ReadContext(context.Background(), d, testCloudRouterReadConfig(t, http.StatusNotFound))

	if !diags.HasError() {
		t.Fatal("expected the data source to fail when the cloud router does not exist")
	}
}
