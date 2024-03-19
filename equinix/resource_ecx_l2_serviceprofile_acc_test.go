package equinix

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/comparisons"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/nprintf"

	"github.com/equinix/ecx-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFabricL2ServiceProfile_Private(t *testing.T) {
	priPortName, _ := schema.EnvDefaultFunc(priPortEnvVar, "sit-001-CX-SV1-NL-Dot1q-BO-10G-PRI-JUN-33")()
	secPortName, _ := schema.EnvDefaultFunc(secPortEnvVar, "sit-001-CX-SV5-NL-Dot1q-BO-10G-SEC-JUN-36")()
	context := map[string]interface{}{
		"resourceName":                       "test",
		"name":                               fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"description":                        acctest.RandString(100),
		"bandwidth_threshold_notifications":  []string{"John.Doe@example.com", "Marry.Doe@example.com"},
		"profile_statuschange_notifications": []string{"John.Doe@example.com", "Marry.Doe@example.com"},
		"vc_statuschange_notifications":      []string{"John.Doe@example.com", "Marry.Doe@example.com"},
		"private":                            true,
		"private_user_emails":                []string{"John.Doe@example.com", "Marry.Doe@example.com"},
		"features_cloud_reach":               true,
		"features_test_profile":              true,
		"port1_name":                         priPortName.(string),
		"port2_name":                         secPortName.(string),
		"speedband1_speed":                   500,
		"speedband1_speed_unit":              "MB",
		"speedband2_speed":                   200,
		"speedband2_speed_unit":              "MB",
	}
	resourceName := fmt.Sprintf("equinix_ecx_l2_serviceprofile.%s", context["resourceName"].(string))
	var testProfile ecx.L2ServiceProfile
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccECXL2PrivateServiceProfile(context),
				Check: resource.ComposeTestCheckFunc(
					testAccECXL2ServiceProfileExists(resourceName, &testProfile),
					testAccECXL2ServiceProfileAttributes(&testProfile, context),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_user_emails", "features.0.test_profile"},
			},
		},
	})
}

func testAccECXL2ServiceProfileExists(resourceName string, profile *ecx.L2ServiceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		client := testAccProvider.Meta().(*config.Config).Ecx
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		resp, err := client.GetL2ServiceProfile(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error when fetching L2 service profile %v", err)
		}
		if ecx.StringValue(resp.UUID) != rs.Primary.ID {
			return fmt.Errorf("resource ID does not match %v - %v", rs.Primary.ID, resp.UUID)
		}
		*profile = *resp
		return nil
	}
}

func testAccECXL2ServiceProfileAttributes(profile *ecx.L2ServiceProfile, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["bandwidth_alert_threshold"]; ok && ecx.Float64Value(profile.AlertPercentage) != v.(float64) {
			return fmt.Errorf("bandwidth_alert_threshold does not match %v - %v", ecx.Float64Value(profile.AlertPercentage), v)
		}
		if v, ok := ctx["speed_customization_allowed"]; ok && ecx.BoolValue(profile.AllowCustomSpeed) != v.(bool) {
			return fmt.Errorf("speed_customization_allowed does not match %v - %v", ecx.BoolValue(profile.AllowCustomSpeed), v)
		}
		if v, ok := ctx["oversubscription_allowed"]; ok && ecx.BoolValue(profile.AllowOverSubscription) != v.(bool) {
			return fmt.Errorf("oversubscription_allowed does not match %v - %v", ecx.BoolValue(profile.AllowOverSubscription), v)
		}
		if v, ok := ctx["api_integration"]; ok && ecx.BoolValue(profile.APIAvailable) != v.(bool) {
			return fmt.Errorf("api_integration does not match %v - %v", profile.APIAvailable, v)
		}
		if v, ok := ctx["authkey_label"]; ok && ecx.StringValue(profile.AuthKeyLabel) != v.(string) {
			return fmt.Errorf("authkey_label does not match %v - %v", profile.AuthKeyLabel, v)
		}
		if v, ok := ctx["connection_name_label"]; ok && ecx.StringValue(profile.ConnectionNameLabel) != v.(string) {
			return fmt.Errorf("connection_name_label does not match %v - %v", ecx.StringValue(profile.ConnectionNameLabel), v)
		}
		if v, ok := ctx["ctag_label"]; ok && ecx.StringValue(profile.CTagLabel) != v.(string) {
			return fmt.Errorf("ctag_label does not match %v - %v", ecx.StringValue(profile.CTagLabel), v)
		}
		if v, ok := ctx["description"]; ok && ecx.StringValue(profile.Description) != v.(string) {
			return fmt.Errorf("description does not match %v - %v", ecx.StringValue(profile.Description), v)
		}
		if v, ok := ctx["servicekey_autogenerated"]; ok && ecx.BoolValue(profile.EnableAutoGenerateServiceKey) != v.(bool) {
			return fmt.Errorf("servicekey_autogenerated does not match %v - %v", ecx.BoolValue(profile.EnableAutoGenerateServiceKey), v)
		}
		if v, ok := ctx["equinix_managed_port_vlan"]; ok && ecx.BoolValue(profile.EquinixManagedPortAndVlan) != v.(bool) {
			return fmt.Errorf("equinix_managed_port_vlan does not match %v - %v", ecx.BoolValue(profile.EquinixManagedPortAndVlan), v)
		}
		if v, ok := ctx["integration_id"]; ok && ecx.StringValue(profile.IntegrationID) != v.(string) {
			return fmt.Errorf("integration_id does not match %v - %v", ecx.StringValue(profile.IntegrationID), v)
		}
		if v, ok := ctx["name"]; ok && ecx.StringValue(profile.Name) != v.(string) {
			return fmt.Errorf("name does not match %v - %v", ecx.StringValue(profile.Name), v)
		}
		if v, ok := ctx["bandwidth_threshold_notifications"]; ok && !comparisons.SlicesMatch(profile.OnBandwidthThresholdNotification, v.([]string)) {
			return fmt.Errorf("bandwidth_threshold_notifications does not match %v - %v", profile.OnBandwidthThresholdNotification, v)
		}
		if v, ok := ctx["profile_statuschange_notifications"]; ok && !comparisons.SlicesMatch(profile.OnProfileApprovalRejectNotification, v.([]string)) {
			return fmt.Errorf("profile_statuschange_notifications does not match %v - %v", profile.OnProfileApprovalRejectNotification, v)
		}
		if v, ok := ctx["vc_statuschange_notifications"]; ok && !comparisons.SlicesMatch(profile.OnVcApprovalRejectionNotification, v.([]string)) {
			return fmt.Errorf("vc_statuschange_notifications does not match %v - %v", profile.OnVcApprovalRejectionNotification, v)
		}
		if v, ok := ctx["oversubscription"]; ok && ecx.StringValue(profile.OverSubscription) != v.(string) {
			return fmt.Errorf("oversubscription does not match %v - %v", ecx.StringValue(profile.OverSubscription), v)
		}
		if v, ok := ctx["private"]; ok && ecx.BoolValue(profile.Private) != v.(bool) {
			return fmt.Errorf("private does not match %v - %v", ecx.BoolValue(profile.Private), v)
		}
		if v, ok := ctx["private_user_emails"]; ok && !comparisons.SlicesMatchCaseInsensitive(profile.PrivateUserEmails, v.([]string)) {
			return fmt.Errorf("private_user_emails does not match %v - %v", profile.PrivateUserEmails, v)
		}
		if v, ok := ctx["redundancy_required"]; ok && ecx.BoolValue(profile.RequiredRedundancy) != v.(bool) {
			return fmt.Errorf("redundancy_required does not match %v - %v", ecx.BoolValue(profile.RequiredRedundancy), v)
		}
		if v, ok := ctx["speed_from_api"]; ok && ecx.BoolValue(profile.SpeedFromAPI) != v.(bool) {
			return fmt.Errorf("speed_from_api does not match %v - %v", ecx.BoolValue(profile.SpeedFromAPI), v)
		}
		if v, ok := ctx["tag_type"]; ok && ecx.StringValue(profile.TagType) != v.(string) {
			return fmt.Errorf("tag_type does not match %v - %v", ecx.StringValue(profile.TagType), v)
		}
		if v, ok := ctx["secondary_vlan_from_primary"]; ok && ecx.BoolValue(profile.VlanSameAsPrimary) != v.(bool) {
			return fmt.Errorf("secondary_vlan_from_primary does not match %v - %v", ecx.BoolValue(profile.VlanSameAsPrimary), v)
		}
		if v, ok := ctx["features_cloud_reach"]; ok && ecx.BoolValue(profile.Features.CloudReach) != v.(bool) {
			return fmt.Errorf("features.cloud_reach does not match %v - %v", ecx.BoolValue(profile.Features.CloudReach), v)
		}
		if len(profile.Ports) != 2 {
			return fmt.Errorf("ports.# length does not match %v - %v", len(profile.Ports), 2)
		}
		return nil
	}
}

func testAccECXL2PrivateServiceProfile(ctx map[string]interface{}) string {
	return nprintf.NPrintf(`
data "equinix_ecx_port" "port1" {
    name = "%{port1_name}"
}
	  
data "equinix_ecx_port" "port2" {
    name = "%{port2_name}"
}

resource "equinix_ecx_l2_serviceprofile" "%{resourceName}" {
	name                               = "%{name}"
	description                        = "%{description}"
	bandwidth_threshold_notifications  = %{bandwidth_threshold_notifications}
	profile_statuschange_notifications = %{profile_statuschange_notifications}
	vc_statuschange_notifications      = %{vc_statuschange_notifications}
	private                            = %{private}
	private_user_emails                = %{private_user_emails}
	features {
	  allow_remote_connections  = %{features_cloud_reach}
	  test_profile              = %{features_test_profile}
	}
	port {
	  uuid       = data.equinix_ecx_port.port1.id
	  metro_code = data.equinix_ecx_port.port1.metro_code
	}
	port {
	  uuid       = data.equinix_ecx_port.port2.id
	  metro_code = data.equinix_ecx_port.port2.metro_code
	}
	speed_band {
	  speed      = %{speedband1_speed}
	  speed_unit = "%{speedband1_speed_unit}"
	}
	speed_band {
	  speed      = %{speedband2_speed}
	  speed_unit = "%{speedband2_speed_unit}"
	}
 }
`, ctx)
}
