package equinix

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("NetworkACLTemplate", &resource.Sweeper{
		Name: "NetworkACLTemplate",
		F:    testSweepNetworkACLTemplate,
	})
}

func testSweepNetworkACLTemplate(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	templates, err := config.ne.GetACLTemplates()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching ACL Templates list: %s", err)
		return err
	}
	nonSweepableCount := 0
	for _, template := range templates {
		if !isSweepableTestResource(ne.StringValue(template.Name)) {
			nonSweepableCount++
			continue
		}
		if err := config.ne.DeleteACLTemplate(ne.StringValue(template.UUID)); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkACLTemplate resource %s (%s): %s", ne.StringValue(template.UUID), ne.StringValue(template.Name), err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkACLTemplate resource %s (%s)", ne.StringValue(template.UUID), ne.StringValue(template.Name))
		}
	}
	if nonSweepableCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonSweepableCount)
	}
	return nil
}

func TestAccNetworkACLTemplate(t *testing.T) {
	t.Parallel()
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	context := map[string]interface{}{
		"resourceName":            "test",
		"name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"description":             randString(50),
		"metro_code":              metro.(string),
		"inbound_rule_1_subnets":  []string{"10.0.0.0/16"},
		"inbound_rule_1_protocol": "TCP",
		"inbound_rule_1_src_port": "any",
		"inbound_rule_1_dst_port": "22-23",
		"inbound_rule_2_subnets":  []string{"192.168.16.0/24"},
		"inbound_rule_2_protocol": "UDP",
		"inbound_rule_2_src_port": "any",
		"inbound_rule_2_dst_port": "53",
		"inbound_rule_3_subnets":  []string{"2.2.2.2/32", "5.5.5.5/32"},
		"inbound_rule_3_protocol": "UDP",
		"inbound_rule_3_src_port": "any",
		"inbound_rule_3_dst_port": "any",
	}
	contextWithChanges := copyMap(context)
	contextWithChanges["description"] = randString(50)
	contextWithChanges["inbound_rule_3_subnets"] = []string{"4.4.4.4/32", "16.20.30.0/24"}
	contextWithChanges["inbound_rule_3_protocol"] = "TCP"
	contextWithChanges["inbound_rule_3_dst_port"] = "2048"
	resourceName := "equinix_network_acl_template." + context["resourceName"].(string)
	var template ne.ACLTemplate
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLTemplate(context),
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkACLTemplateExists(resourceName, &template),
					testAccNetworkACLTemplateAttributes(&template, context),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkACLTemplate(contextWithChanges),
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkACLTemplateExists(resourceName, &template),
					testAccNetworkACLTemplateAttributes(&template, contextWithChanges),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
		},
	})
}

func testAccNetworkACLTemplate(ctx map[string]interface{}) string {
	return nprintf(`
resource "equinix_network_acl_template" "%{resourceName}" {
  name          = "%{name}"
  description   = "%{description}"
  metro_code    = "%{metro_code}"

  inbound_rule {
    subnets  = %{inbound_rule_1_subnets}
	protocol = "%{inbound_rule_1_protocol}"
	src_port = "%{inbound_rule_1_src_port}"
	dst_port = "%{inbound_rule_1_dst_port}"
  }

  inbound_rule {
	subnets  = %{inbound_rule_2_subnets}
	protocol = "%{inbound_rule_2_protocol}"
	src_port = "%{inbound_rule_2_src_port}"
	dst_port = "%{inbound_rule_2_dst_port}"
  }

  inbound_rule {
	subnets  = %{inbound_rule_3_subnets}
	protocol = "%{inbound_rule_3_protocol}"
	src_port = "%{inbound_rule_3_src_port}"
	dst_port = "%{inbound_rule_3_dst_port}"
  }
}
`, ctx)
}

func testAccNetworkACLTemplateExists(resourceName string, template *ne.ACLTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		client := testAccProvider.Meta().(*Config).ne
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		resp, err := client.GetACLTemplate(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error when fetching ACL template '%s': %s", rs.Primary.ID, err)
		}
		*template = *resp
		return nil
	}
}

func testAccNetworkACLTemplateAttributes(template *ne.ACLTemplate, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["name"]; ok && ne.StringValue(template.Name) != v.(string) {
			return fmt.Errorf("name does not match %v - %v", ne.StringValue(template.Name), v)
		}
		if v, ok := ctx["description"]; ok && ne.StringValue(template.Description) != v.(string) {
			return fmt.Errorf("name does not match %v - %v", ne.StringValue(template.Description), v)
		}
		if v, ok := ctx["metro_code"]; ok && ne.StringValue(template.MetroCode) != v.(string) {
			return fmt.Errorf("name does not match %v - %v", ne.StringValue(template.MetroCode), v)
		}
		if len(template.InboundRules) != 3 {
			return fmt.Errorf("number of inbound rules does not match %v - %v", len(template.InboundRules), 3)
		}
		for i := 0; i < 3; i++ {
			if ne.IntValue(template.InboundRules[i].SeqNo) != i+1 {
				return fmt.Errorf("inbound_rule %d seqNo does not match %v - %v", i+1, ne.IntValue(template.InboundRules[i].SeqNo), i+1)
			}
			if ne.StringValue(template.InboundRules[i].SrcType) != "SUBNET" {
				return fmt.Errorf("inbound_rule %d srcType does not match %v - %v", i+1, ne.StringValue(template.InboundRules[i].SrcType), "SUBNET")
			}
			if v, ok := ctx[fmt.Sprintf("inbound_rule_%d_subnets", i+1)]; ok && !slicesMatch(template.InboundRules[i].Subnets, v.([]string)) {
				return fmt.Errorf("inbound_rule %d subnets does not match %v - %v", i+1, template.InboundRules[i].Subnets, v)
			}
			if v, ok := ctx[fmt.Sprintf("inbound_rule_%d_protocol", i+1)]; ok && ne.StringValue(template.InboundRules[i].Protocol) != v.(string) {
				return fmt.Errorf("inbound_rule %d protocol does not match %v - %v", i+1, ne.StringValue(template.InboundRules[i].Protocol), v)
			}
			if v, ok := ctx[fmt.Sprintf("inbound_rule_%d_src_port", i+1)]; ok && ne.StringValue(template.InboundRules[i].SrcPort) != v.(string) {
				return fmt.Errorf("inbound_rule %d src_port does not match %v - %v", i+1, ne.StringValue(template.InboundRules[i].SrcPort), v)
			}
			if v, ok := ctx[fmt.Sprintf("inbound_rule_%d_dst_port", i+1)]; ok && ne.StringValue(template.InboundRules[i].DstPort) != v.(string) {
				return fmt.Errorf("inbound_rule %d dst_port does not match %v - %v", i+1, ne.StringValue(template.InboundRules[i].DstPort), v)
			}
		}
		return nil
	}
}
