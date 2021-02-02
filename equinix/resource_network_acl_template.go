package equinix

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/equinix/ne-go"
	"github.com/equinix/rest-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var networkACLTemplateSchemaNames = map[string]string{
	"UUID":            "uuid",
	"Name":            "name",
	"Description":     "description",
	"MetroCode":       "metro_code",
	"DeviceUUID":      "device_id",
	"DeviceACLStatus": "device_acl_status",
	"InboundRules":    "inbound_rule",
}

var networkACLTemplateDescriptions = map[string]string{
	"UUID":            "Unique identifier of ACL template resource",
	"Name":            "ACL template name",
	"Description":     "ACL template description",
	"MetroCode":       "ACL template location metro code",
	"DeviceUUID":      "Identifier of a network device where template was applied",
	"DeviceACLStatus": "Status of ACL template provisioning process on a device, where template was applied",
	"InboundRules":    "One or more rules to specify allowed inbound traffic. Rules are ordered, matching traffic rule stops processing subsequent ones.",
}

var networkACLTemplateInboundRuleSchemaNames = map[string]string{
	"SeqNo":    "sequence_number",
	"SrcType":  "source_type",
	"Subnets":  "subnets",
	"Protocol": "protocol",
	"SrcPort":  "src_port",
	"DstPort":  "dst_port",
}

var networkACLTemplateInboundRuleDescriptions = map[string]string{
	"SeqNo":    "Inbound rule sequence number",
	"SrcType":  "Type of traffic source used in a given innbound rule",
	"Subnets":  "Inbound traffic source IP subnets in CIDR format",
	"Protocol": "Inbound traffic protocol. One of: `IP`, `TCP`, `UDP`",
	"SrcPort":  "Inbound traffic source ports. Either up to 10, comma separated ports or port range or any word",
	"DstPort":  "Inbound traffic destination ports. Either up to 10, comma separated ports or port range or any word",
}

func resourceNetworkACLTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkACLTemplateCreate,
		ReadContext:   resourceNetworkACLTemplateRead,
		UpdateContext: resourceNetworkACLTemplateUpdate,
		DeleteContext: resourceNetworkACLTemplateDelete,
		Schema:        createNetworkACLTemplateSchema(),
		Description:   "Resource allows creation and management of Equinix Network Edge device Access Control List templates",
	}
}

func createNetworkACLTemplateSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		networkACLTemplateSchemaNames["UUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkACLTemplateDescriptions["UUID"],
		},
		networkACLTemplateSchemaNames["Name"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(8, 100),
			Description:  networkACLTemplateDescriptions["Name"],
		},
		networkACLTemplateSchemaNames["Description"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 100),
			Description:  networkACLTemplateDescriptions["Description"],
		},
		networkACLTemplateSchemaNames["MetroCode"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: stringIsMetroCode(),
			Description:  networkACLTemplateDescriptions["MetroCode"],
		},
		networkACLTemplateSchemaNames["DeviceUUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkACLTemplateDescriptions["DeviceUUID"],
		},
		networkACLTemplateSchemaNames["DeviceACLStatus"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkACLTemplateDescriptions["DeviceACLStatus"],
		},
		networkACLTemplateSchemaNames["InboundRules"]: {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: createNetworkACLTemplateInboundRuleSchema(),
			},
			Description: networkACLTemplateDescriptions["InboundRules"],
		},
	}
}

func createNetworkACLTemplateInboundRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		networkACLTemplateInboundRuleSchemaNames["SeqNo"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: networkACLTemplateInboundRuleDescriptions["SeqNo"],
		},
		networkACLTemplateInboundRuleSchemaNames["SrcType"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkACLTemplateInboundRuleDescriptions["SrcType"],
		},
		networkACLTemplateInboundRuleSchemaNames["Subnets"]: {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.IsCIDR,
			},
			Description: networkACLTemplateInboundRuleDescriptions["Subnets"],
		},
		networkACLTemplateInboundRuleSchemaNames["Protocol"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"IP", "TCP", "UDP"}, false),
			Description:  networkACLTemplateInboundRuleDescriptions["Protocol"],
		},
		networkACLTemplateInboundRuleSchemaNames["SrcPort"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: stringIsPortDefinition(),
			Description:  networkACLTemplateInboundRuleDescriptions["SrcPort"],
		},
		networkACLTemplateInboundRuleSchemaNames["DstPort"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: stringIsPortDefinition(),
			Description:  networkACLTemplateInboundRuleDescriptions["DstPort"],
		},
	}
}

func resourceNetworkACLTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*Config)
	var diags diag.Diagnostics
	template := createACLTemplate(d)
	uuid, err := conf.ne.CreateACLTemplate(template)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ne.StringValue(uuid))
	diags = append(diags, resourceNetworkACLTemplateRead(ctx, d, m)...)
	return diags
}

func resourceNetworkACLTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*Config)
	var diags diag.Diagnostics
	template, err := conf.ne.GetACLTemplate(d.Id())
	if err != nil {
		if restErr, ok := err.(rest.Error); ok {
			if restErr.HTTPCode == http.StatusNotFound {
				d.SetId("")
				return diags
			}
		}
		return diag.FromErr(err)
	}
	if err := updateACLTemplateResource(template, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceNetworkACLTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*Config)
	var diags diag.Diagnostics
	template := createACLTemplate(d)
	if err := conf.ne.ReplaceACLTemplate(d.Id(), template); err != nil {
		return diag.FromErr(err)
	}
	diags = append(diags, resourceNetworkACLTemplateRead(ctx, d, m)...)
	return diags
}

func resourceNetworkACLTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*Config)
	var diags diag.Diagnostics
	if devID, ok := d.GetOk(networkACLTemplateSchemaNames["DeviceUUID"]); ok {
		if err := conf.ne.NewDeviceUpdateRequest(devID.(string)).WithACLTemplate("").Execute(); err != nil {
			log.Printf("[WARN] could not unassign ACL template %q from device %q: %s", d.Id(), devID, err)
		}
	}
	if err := conf.ne.DeleteACLTemplate(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createACLTemplate(d *schema.ResourceData) ne.ACLTemplate {
	template := ne.ACLTemplate{}
	if v, ok := d.GetOk(networkACLTemplateSchemaNames["Name"]); ok {
		template.Name = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkACLTemplateSchemaNames["Description"]); ok {
		template.Description = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkACLTemplateSchemaNames["MetroCode"]); ok {
		template.MetroCode = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkACLTemplateSchemaNames["InboundRules"]); ok {
		template.InboundRules = expandACLTemplateInboundRules(v.([]interface{}))
	}
	return template
}

func updateACLTemplateResource(template *ne.ACLTemplate, d *schema.ResourceData) error {
	if err := d.Set(networkACLTemplateSchemaNames["UUID"], template.UUID); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["UUID"], err)
	}
	if err := d.Set(networkACLTemplateSchemaNames["Name"], template.Name); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["Name"], err)
	}
	if err := d.Set(networkACLTemplateSchemaNames["Description"], template.Description); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["Description"], err)
	}
	if err := d.Set(networkACLTemplateSchemaNames["MetroCode"], template.MetroCode); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["MetroCode"], err)
	}
	if err := d.Set(networkACLTemplateSchemaNames["DeviceUUID"], template.DeviceUUID); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["DeviceUUID"], err)
	}
	if err := d.Set(networkACLTemplateSchemaNames["DeviceACLStatus"], template.DeviceACLStatus); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["DeviceACLStatus"], err)
	}
	if err := d.Set(networkACLTemplateSchemaNames["InboundRules"], flattenACLTemplateInboundRules(template.InboundRules)); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["InboundRules"], err)
	}
	return nil
}

func expandACLTemplateInboundRules(rules []interface{}) []ne.ACLTemplateInboundRule {
	transformed := make([]ne.ACLTemplateInboundRule, len(rules))
	for i := range rules {
		ruleMap := rules[i].(map[string]interface{})
		rule := ne.ACLTemplateInboundRule{}
		rule.SeqNo = ne.Int(i + 1)
		rule.SrcType = ne.String("SUBNET")
		if v, ok := ruleMap[networkACLTemplateInboundRuleSchemaNames["Subnets"]]; ok {
			rule.Subnets = expandListToStringList(v.([]interface{}))
		}
		if v, ok := ruleMap[networkACLTemplateInboundRuleSchemaNames["Protocol"]]; ok {
			rule.Protocol = ne.String(v.(string))
		}
		if v, ok := ruleMap[networkACLTemplateInboundRuleSchemaNames["SrcPort"]]; ok {
			rule.SrcPort = ne.String(v.(string))
		}
		if v, ok := ruleMap[networkACLTemplateInboundRuleSchemaNames["DstPort"]]; ok {
			rule.DstPort = ne.String(v.(string))
		}
		transformed[i] = rule
	}
	return transformed
}

func flattenACLTemplateInboundRules(rules []ne.ACLTemplateInboundRule) interface{} {
	transformed := make([]interface{}, len(rules))
	for i := range rules {
		transformed[i] = map[string]interface{}{
			networkACLTemplateInboundRuleSchemaNames["SeqNo"]:    rules[i].SeqNo,
			networkACLTemplateInboundRuleSchemaNames["SrcType"]:  rules[i].SrcType,
			networkACLTemplateInboundRuleSchemaNames["Subnets"]:  rules[i].Subnets,
			networkACLTemplateInboundRuleSchemaNames["Protocol"]: rules[i].Protocol,
			networkACLTemplateInboundRuleSchemaNames["SrcPort"]:  rules[i].SrcPort,
			networkACLTemplateInboundRuleSchemaNames["DstPort"]:  rules[i].DstPort,
		}
	}
	return transformed
}
