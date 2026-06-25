package equinix

import (
	"context"
	"fmt"
	"net/http"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_validation "github.com/equinix/terraform-provider-equinix/internal/validation"

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
	"DeviceACLStatus": "device_acl_status",
	"InboundRules":    "inbound_rule",
	"DeviceDetails":   "device_details",
	"ProjectID":       "project_id",
}

var networkACLTemplateDescriptions = map[string]string{
	"UUID":            "Unique identifier of ACL template resource",
	"Name":            "ACL template name",
	"Description":     "ACL template description, up to 200 characters",
	"DeviceACLStatus": "Status of ACL template provisioning process on a device, where template was applied",
	"InboundRules":    "One or more rules to specify allowed inbound traffic. Rules are ordered, matching traffic rule stops processing subsequent ones.",
	"DeviceDetails":   "Device Details to which ACL template is assigned to. ",
	"ProjectID":       "The unique identifier of Project Resource to which ACL template is scoped to",
}

var networkACLTemplateInboundRuleSchemaNames = map[string]string{
	"SeqNo":       "sequence_number",
	"Subnet":      "subnet",
	"Protocol":    "protocol",
	"SrcPort":     "src_port",
	"DstPort":     "dst_port",
	"Description": "description",
}

var networkACLTemplateInboundRuleDescriptions = map[string]string{
	"SeqNo":       "Inbound rule sequence number",
	"Subnet":      "Inbound traffic source IP subnet in CIDR format",
	"Protocol":    "Inbound traffic protocol. One of: `IP`, `TCP`, `UDP`",
	"SrcPort":     "Inbound traffic source ports. Either up to 10, comma separated ports or port range or any word",
	"DstPort":     "Inbound traffic destination ports. Either up to 10, comma separated ports or port range or any word",
	"Description": "Inbound rule description, up to 200 characters",
}

var networkACLTemplateDeviceDetailSchemaNames = map[string]string{
	"UUID":      "uuid",
	"Name":      "name",
	"ACLStatus": "acl_status",
}

func resourceNetworkACLTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkACLTemplateCreate,
		ReadContext:   resourceNetworkACLTemplateRead,
		UpdateContext: resourceNetworkACLTemplateUpdate,
		DeleteContext: resourceNetworkACLTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema:      createNetworkACLTemplateSchema(),
		Description: "Resource allows creation and management of Equinix Network Edge device Access Control List templates",
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
			ValidateFunc: validation.StringLenBetween(1, 200),
			Description:  networkACLTemplateDescriptions["Description"],
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
		networkACLTemplateSchemaNames["DeviceDetails"]: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: networkACLTemplateDeviceDetailsSchema(),
			},
			Description: networkACLTemplateDescriptions["DeviceDetails"],
		},
		networkACLTemplateSchemaNames["ProjectID"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
			Description:  networkACLTemplateDescriptions["ProjectID"],
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
		networkACLTemplateInboundRuleSchemaNames["Subnet"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  networkACLTemplateInboundRuleDescriptions["Subnet"],
			ValidateFunc: validation.IsCIDR,
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
			ValidateFunc: equinix_validation.StringIsPortDefinition,
			Description:  networkACLTemplateInboundRuleDescriptions["SrcPort"],
		},
		networkACLTemplateInboundRuleSchemaNames["DstPort"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: equinix_validation.StringIsPortDefinition,
			Description:  networkACLTemplateInboundRuleDescriptions["DstPort"],
		},
		networkACLTemplateInboundRuleSchemaNames["Description"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 200),
			Description:  networkACLTemplateInboundRuleDescriptions["Description"],
		},
	}
}

func networkACLTemplateDeviceDetailsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		networkACLTemplateDeviceDetailSchemaNames["UUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkACLTemplateInboundRuleDescriptions["SeqNo"],
		},
		networkACLTemplateDeviceDetailSchemaNames["Name"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkACLTemplateInboundRuleDescriptions["Name"],
		},
		networkACLTemplateDeviceDetailSchemaNames["ACLStatus"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkACLTemplateInboundRuleDescriptions["ACLStatus"],
		},
	}
}

func resourceNetworkACLTemplateCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	template := createACLTemplate(d)
	uuid, err := client.CreateACLTemplate(template)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ne.StringValue(uuid))
	diags = append(diags, resourceNetworkACLTemplateRead(ctx, d, m)...)
	return diags
}

func resourceNetworkACLTemplateRead(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	template, err := client.GetACLTemplate(d.Id())
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

func resourceNetworkACLTemplateUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	template := createACLTemplate(d)
	if err := client.ReplaceACLTemplate(d.Id(), template); err != nil {
		return diag.FromErr(err)
	}
	diags = append(diags, resourceNetworkACLTemplateRead(ctx, d, m)...)
	return diags
}

func resourceNetworkACLTemplateDelete(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	if err := client.DeleteACLTemplate(d.Id()); err != nil {
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
	if v, ok := d.GetOk(networkACLTemplateSchemaNames["ProjectID"]); ok {
		template.ProjectID = ne.String(v.(string))
	}
	if v, ok := d.GetOk(networkACLTemplateSchemaNames["InboundRules"]); ok {
		template.InboundRules = expandACLTemplateInboundRules(v.([]any))
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
	if err := d.Set(networkACLTemplateSchemaNames["DeviceACLStatus"], template.DeviceACLStatus); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["DeviceACLStatus"], err)
	}
	if err := d.Set(networkACLTemplateSchemaNames["ProjectID"], template.ProjectID); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["ProjectID"], err)
	}
	if err := d.Set(networkACLTemplateSchemaNames["InboundRules"], flattenACLTemplateInboundRules(template.InboundRules)); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["InboundRules"], err)
	}
	if err := d.Set(networkACLTemplateSchemaNames["DeviceDetails"], flattenACLTemplateDeviceDetails(template.DeviceDetails)); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["DeviceDetails"], err)
	}

	return nil
}

func expandACLTemplateInboundRules(rules []any) []ne.ACLTemplateInboundRule {
	transformed := make([]ne.ACLTemplateInboundRule, len(rules))
	for i := range rules {
		ruleMap := rules[i].(map[string]any)
		rule := ne.ACLTemplateInboundRule{}
		rule.SeqNo = ne.Int(i + 1)
		if v, ok := ruleMap[networkACLTemplateInboundRuleSchemaNames["Subnet"]]; ok {
			rule.Subnet = ne.String(v.(string))
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
		if v, ok := ruleMap[networkACLTemplateInboundRuleSchemaNames["Description"]]; ok {
			rule.Description = ne.String(v.(string))
		}
		transformed[i] = rule
	}
	return transformed
}

func flattenACLTemplateInboundRules(rules []ne.ACLTemplateInboundRule) any {
	transformed := make([]any, len(rules))
	for i := range rules {
		transformedTemplate := make(map[string]any)
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["SeqNo"]] = rules[i].SeqNo
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["Protocol"]] = rules[i].Protocol
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["SrcPort"]] = rules[i].SrcPort
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["DstPort"]] = rules[i].DstPort
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["Subnet"]] = rules[i].Subnet
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["Description"]] = rules[i].Description
		transformed[i] = transformedTemplate
	}
	return transformed
}

func flattenACLTemplateDeviceDetails(rules []ne.ACLTemplateDeviceDetails) any {
	transformed := make([]any, len(rules))
	for i := range rules {
		transformed[i] = map[string]any{
			networkACLTemplateDeviceDetailSchemaNames["UUID"]:      rules[i].UUID,
			networkACLTemplateDeviceDetailSchemaNames["Name"]:      rules[i].Name,
			networkACLTemplateDeviceDetailSchemaNames["ACLStatus"]: rules[i].ACLStatus,
		}
	}
	return transformed
}
