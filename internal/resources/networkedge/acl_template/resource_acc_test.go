package acl_template

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/nprintf"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	tstResourcePrefix = "tfacc"
)

func copyMap(source map[string]interface{}) map[string]interface{} {
	target := make(map[string]interface{})
	for k, v := range source {
		target[k] = v
	}
	return target
}

func TestAccNetworkACLTemplate(t *testing.T) {
	context := map[string]interface{}{
		"resourceName":               "test",
		"name":                       fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"description":                acctest.RandString(50),
		"inbound_rule_1_subnet":      "10.0.0.0/16",
		"inbound_rule_1_protocol":    "TCP",
		"inbound_rule_1_src_port":    "any",
		"inbound_rule_1_dst_port":    "22-23",
		"inbound_rule_1_description": acctest.RandString(50),
		"inbound_rule_2_subnet":      "192.168.16.0/24",
		"inbound_rule_2_protocol":    "UDP",
		"inbound_rule_2_src_port":    "any",
		"inbound_rule_2_dst_port":    "53",
		"inbound_rule_3_subnet":      "2.2.2.2/32",
		"inbound_rule_3_protocol":    "UDP",
		"inbound_rule_3_src_port":    "any",
		"inbound_rule_3_dst_port":    "any",
	}
	contextWithChanges := copyMap(context)
	contextWithChanges["description"] = acctest.RandString(50)
	contextWithChanges["inbound_rule_1_description"] = acctest.RandString(50)
	contextWithChanges["inbound_rule_3_subnet"] = "4.4.4.4/32"
	contextWithChanges["inbound_rule_3_protocol"] = "TCP"
	contextWithChanges["inbound_rule_3_dst_port"] = "2048"
	resourceName := "equinix_network_acl_template." + context["resourceName"].(string)
	var template ne.ACLTemplate
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
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
	return nprintf.Nprintf(`
resource "equinix_network_acl_template" "%{resourceName}" {
  name          = "%{name}"
  description   = "%{description}"

  inbound_rule {
    subnet   = "%{inbound_rule_1_subnet}"
	protocol = "%{inbound_rule_1_protocol}"
	src_port = "%{inbound_rule_1_src_port}"
	dst_port = "%{inbound_rule_1_dst_port}"
	description = "%{inbound_rule_1_description}"
  }

  inbound_rule {
	subnet   = "%{inbound_rule_2_subnet}"
	protocol = "%{inbound_rule_2_protocol}"
	src_port = "%{inbound_rule_2_src_port}"
	dst_port = "%{inbound_rule_2_dst_port}"
  }

  inbound_rule {
	subnet   = "%{inbound_rule_3_subnet}"
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
		client := acceptance.TestAccProvider.Meta().(*config.Config).Ne
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
			return fmt.Errorf("description does not match %v - %v", ne.StringValue(template.Description), v)
		}
		if len(template.InboundRules) != 3 {
			return fmt.Errorf("number of inbound rules does not match %v - %v", len(template.InboundRules), 3)
		}
		for i := 0; i < 3; i++ {
			if ne.IntValue(template.InboundRules[i].SeqNo) != i+1 {
				return fmt.Errorf("inbound_rule %d seqNo does not match %v - %v", i+1, ne.IntValue(template.InboundRules[i].SeqNo), i+1)
			}
			if v, ok := ctx[fmt.Sprintf("inbound_rule_%d_subnet", i+1)]; ok && ne.StringValue(template.InboundRules[i].Subnet) != v.(string) {
				return fmt.Errorf("inbound_rule %d subnet does not match %v - %v", i+1, ne.StringValue(template.InboundRules[i].Subnet), v)
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
			if v, ok := ctx[fmt.Sprintf("inbound_rule_%d_description", i+1)]; ok && ne.StringValue(template.InboundRules[i].Description) != v.(string) {
				return fmt.Errorf("inbound_rule %d description does not match %v - %v", i+1, ne.StringValue(template.InboundRules[i].Description), v)
			}
		}
		return nil
	}
}
