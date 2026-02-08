package acl_template

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
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
	"MetroCode":       "metro_code",
	"DeviceUUID":      "device_id",
	"DeviceACLStatus": "device_acl_status",
	"InboundRules":    "inbound_rule",
	"DeviceDetails":   "device_details",
	"ProjectID":       "project_id",
}

var networkACLTemplateDescriptions = map[string]string{
	"UUID":            "Unique identifier of ACL template resource",
	"Name":            "ACL template name",
	"Description":     "ACL template description, up to 200 characters",
	"MetroCode":       "ACL template location metro code",
	"DeviceUUID":      "Identifier of a network device where template was applied",
	"DeviceACLStatus": "Status of ACL template provisioning process on a device, where template was applied",
	"InboundRules":    "One or more rules to specify allowed inbound traffic. Rules are ordered, matching traffic rule stops processing subsequent ones.",
	"DeviceDetails":   "Device Details to which ACL template is assigned to. ",
	"ProjectID":       "The unique identifier of Project Resource to which ACL template is scoped to",
}

var networkACLTemplateInboundRuleSchemaNames = map[string]string{
	"SeqNo":       "sequence_number",
	"SrcType":     "source_type",
	"Subnets":     "subnets",
	"Subnet":      "subnet",
	"Protocol":    "protocol",
	"SrcPort":     "src_port",
	"DstPort":     "dst_port",
	"Description": "description",
}

var networkACLTemplateInboundRuleDescriptions = map[string]string{
	"SeqNo":       "Inbound rule sequence number",
	"SrcType":     "Type of traffic source used in a given inbound rule",
	"Subnets":     "Inbound traffic source IP subnets in CIDR format",
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

var networkACLTemplateDeviceDetailDescription = map[string]string{
	"UUID":      "Unique Identifier for the device",
	"Name":      "Device Name",
	"ACLStatus": "Device ACL Provisioning status",
}

var networkACLTemplateDeprecateDescriptions = map[string]string{
	"DeviceUUID": "Refer to device details get device information",
	"MetroCode":  "Metro Code is no longer required",
	"Subnets":    "Use Subnet instead",
}

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkACLTemplateCreate,
		ReadContext:   resourceNetworkACLTemplateRead,
		UpdateContext: resourceNetworkACLTemplateUpdate,
		DeleteContext: resourceNetworkACLTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
		networkACLTemplateSchemaNames["MetroCode"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Deprecated:   networkACLTemplateDeprecateDescriptions["MetroCode"],
			ValidateFunc: equinix_validation.StringIsMetroCode,
			Description:  networkACLTemplateDescriptions["MetroCode"],
		},
		networkACLTemplateSchemaNames["DeviceUUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Deprecated:  networkACLTemplateDeprecateDescriptions["DeviceUUID"],
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
		networkACLTemplateInboundRuleSchemaNames["SrcType"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkACLTemplateInboundRuleDescriptions["SrcType"],
			Deprecated:  "Source Type will not be returned",
		},
		networkACLTemplateInboundRuleSchemaNames["Subnets"]: {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.IsCIDR,
			},
			Description: networkACLTemplateInboundRuleDescriptions["Subnets"],
			Deprecated:  networkACLTemplateDeprecateDescriptions["Subnets"],
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

func resourceNetworkACLTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceNetworkACLTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceNetworkACLTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceNetworkACLTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	if devID, ok := d.GetOk(networkACLTemplateSchemaNames["DeviceUUID"]); ok {
		if err := client.NewDeviceUpdateRequest(devID.(string)).WithACLTemplate("").Execute(); err != nil {
			log.Printf("[WARN] could not unassign ACL template %q from device %q: %s", d.Id(), devID, err)
		}
	}
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
	if err := d.Set(networkACLTemplateSchemaNames["DeviceACLStatus"], template.DeviceACLStatus); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["DeviceACLStatus"], err)
	}
	if err := d.Set(networkACLTemplateSchemaNames["ProjectID"], template.ProjectID); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["ProjectID"], err)
	}
	var inboundRules []ne.ACLTemplateInboundRule
	if v, ok := d.GetOk(networkACLTemplateSchemaNames["InboundRules"]); ok {
		inboundRules = expandACLTemplateInboundRules(v.([]interface{}))
	}
	if err := d.Set(networkACLTemplateSchemaNames["InboundRules"], flattenACLTemplateInboundRules(inboundRules, template.InboundRules)); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["InboundRules"], err)
	}
	if err := d.Set(networkACLTemplateSchemaNames["DeviceDetails"], flattenACLTemplateDeviceDetails(template.DeviceDetails)); err != nil {
		return fmt.Errorf("error reading %s: %s", networkACLTemplateSchemaNames["DeviceDetails"], err)
	}

	return nil
}

func expandACLTemplateInboundRules(rules []interface{}) []ne.ACLTemplateInboundRule {
	transformed := make([]ne.ACLTemplateInboundRule, len(rules))
	for i := range rules {
		ruleMap := rules[i].(map[string]interface{})
		rule := ne.ACLTemplateInboundRule{}
		rule.SeqNo = ne.Int(i + 1)
		if v, ok := ruleMap[networkACLTemplateInboundRuleSchemaNames["Subnets"]]; ok {
			rule.Subnets = converters.IfArrToStringArr(v.([]interface{}))
		}
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

func flattenACLTemplateInboundRules(existingRules []ne.ACLTemplateInboundRule, rules []ne.ACLTemplateInboundRule) interface{} {
	setSubnets := checkExistingSubnets(existingRules)
	transformed := make([]interface{}, len(rules))
	for i := range rules {
		transformedTemplate := make(map[string]interface{})
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["SeqNo"]] = rules[i].SeqNo
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["SrcType"]] = rules[i].SrcType
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["Protocol"]] = rules[i].Protocol
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["SrcPort"]] = rules[i].SrcPort
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["DstPort"]] = rules[i].DstPort
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["Subnet"]] = rules[i].Subnet
		transformedTemplate[networkACLTemplateInboundRuleSchemaNames["Description"]] = rules[i].Description
		if setSubnets {
			transformedTemplate[networkACLTemplateInboundRuleSchemaNames["Subnets"]] = rules[i].Subnets
		}
		transformed[i] = transformedTemplate
	}
	return transformed
}

func checkExistingSubnets(existingRules []ne.ACLTemplateInboundRule) bool {
	for i := range existingRules {
		if existingRules[i].Subnets != nil && len(existingRules[i].Subnets) > 0 {
			return true
		}
	}
	return false
}

func flattenACLTemplateDeviceDetails(rules []ne.ACLTemplateDeviceDetails) interface{} {
	transformed := make([]interface{}, len(rules))
	for i := range rules {
		transformed[i] = map[string]interface{}{
			networkACLTemplateDeviceDetailSchemaNames["UUID"]:      rules[i].UUID,
			networkACLTemplateDeviceDetailSchemaNames["Name"]:      rules[i].Name,
			networkACLTemplateDeviceDetailSchemaNames["ACLStatus"]: rules[i].ACLStatus,
		}
	}
	return transformed
}
