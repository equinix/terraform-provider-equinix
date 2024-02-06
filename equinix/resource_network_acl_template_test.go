package equinix

import (
	"testing"

	"github.com/equinix/ne-go"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestNetworkACLTemplate_createFromResourceData(t *testing.T) {
	expected := ne.ACLTemplate{
		Name:        ne.String("test"),
		Description: ne.String("testTemplate"),
		MetroCode:   ne.String("SV"),
		ProjectID:   ne.String("68ccfd49-39b1-478e-957a-67c72f719d7a"),
		InboundRules: []ne.ACLTemplateInboundRule{
			{
				SeqNo:       ne.Int(1),
				Subnets:     []string{"10.0.0.0/24", "1.1.1.1/32"},
				Subnet:      ne.String("10.0.0.0/24"),
				Protocol:    ne.String("TCP"),
				SrcPort:     ne.String("any"),
				DstPort:     ne.String("8080"),
				Description: ne.String("test inbound rule"),
			},
		},
	}
	rawData := map[string]interface{}{
		networkACLTemplateSchemaNames["Name"]:        ne.StringValue(expected.Name),
		networkACLTemplateSchemaNames["Description"]: ne.StringValue(expected.Description),
		networkACLTemplateSchemaNames["MetroCode"]:   ne.StringValue(expected.MetroCode),
		networkACLTemplateSchemaNames["ProjectID"]:   ne.StringValue(expected.ProjectID),
	}
	d := schema.TestResourceDataRaw(t, createNetworkACLTemplateSchema(), rawData)
	d.Set(networkACLTemplateSchemaNames["InboundRules"], flattenACLTemplateInboundRules(expected.InboundRules, expected.InboundRules))
	// when
	result := createACLTemplate(d)
	// then
	assert.Equal(t, expected, result, "Created ACL Template matches expected result")
}

func TestNetworkACLTemplate_updateResourceData(t *testing.T) {
	initial := ne.ACLTemplate{
		Name:        ne.String("test"),
		Description: ne.String("testTemplate"),
		MetroCode:   ne.String("SV"),
		InboundRules: []ne.ACLTemplateInboundRule{
			{
				SeqNo:    ne.Int(1),
				Subnets:  []string{"10.0.0.0/24", "1.1.1.1/32"},
				Subnet:   ne.String("10.0.0.0/24"),
				Protocol: ne.String("TCP"),
				SrcPort:  ne.String("any"),
				DstPort:  ne.String("8080"),
			},
		},
	}
	rawData := map[string]interface{}{
		networkACLTemplateSchemaNames["Name"]:        ne.StringValue(initial.Name),
		networkACLTemplateSchemaNames["Description"]: ne.StringValue(initial.Description),
		networkACLTemplateSchemaNames["MetroCode"]:   ne.StringValue(initial.MetroCode),
	}
	d := schema.TestResourceDataRaw(t, createNetworkACLTemplateSchema(), rawData)
	d.Set(networkACLTemplateSchemaNames["InboundRules"], flattenACLTemplateInboundRules(initial.InboundRules, initial.InboundRules))

	input := &ne.ACLTemplate{
		Name:        ne.String("test"),
		Description: ne.String("testTemplate"),
		MetroCode:   ne.String("SV"),
		InboundRules: []ne.ACLTemplateInboundRule{
			{
				SeqNo:       ne.Int(1),
				Subnets:     []string{"10.0.0.0/24", "1.1.1.1/32"},
				Subnet:      ne.String("10.0.0.0/24"),
				Protocol:    ne.String("TCP"),
				SrcPort:     ne.String("any"),
				DstPort:     ne.String("8080"),
				Description: ne.String("test inbound rule"),
			},
		},
		DeviceDetails: []ne.ACLTemplateDeviceDetails{
			{
				UUID:      ne.String("6e0c2b09-37e9-4fd2-b9e4-e68b08d8b29d"),
				Name:      ne.String("testDevice"),
				ACLStatus: ne.String("Provisioned"),
			},
		},
	}
	// when
	err := updateACLTemplateResource(input, d)
	// then
	assert.Nil(t, err, "Update of resource data does not return error")
	assert.Equal(t, ne.StringValue(input.Name), d.Get(networkACLTemplateSchemaNames["Name"]), "Name matches")
	assert.Equal(t, ne.StringValue(input.Description), d.Get(networkACLTemplateSchemaNames["Description"]), "Description matches")
	assert.Equal(t, ne.StringValue(input.MetroCode), d.Get(networkACLTemplateSchemaNames["MetroCode"]), "MetroCode matches")
	assert.Equal(t, input.InboundRules, expandACLTemplateInboundRules(d.Get(networkACLTemplateSchemaNames["InboundRules"]).([]interface{})), "InboundRules matches")
}

func TestNetworkACLTemplate_expandInboundRules(t *testing.T) {
	// given
	input := []interface{}{
		map[string]interface{}{
			networkACLTemplateInboundRuleSchemaNames["Subnets"]:     []interface{}{"10.0.0.0/24", "1.1.1.1/32"},
			networkACLTemplateInboundRuleSchemaNames["Protocol"]:    "TCP",
			networkACLTemplateInboundRuleSchemaNames["SrcPort"]:     "any",
			networkACLTemplateInboundRuleSchemaNames["DstPort"]:     "8080",
			networkACLTemplateInboundRuleSchemaNames["Description"]: "description of inbound rule",
		},
		map[string]interface{}{
			networkACLTemplateInboundRuleSchemaNames["Subnet"]:   "3.3.3.3/32",
			networkACLTemplateInboundRuleSchemaNames["Protocol"]: "IP",
			networkACLTemplateInboundRuleSchemaNames["SrcPort"]:  "any",
			networkACLTemplateInboundRuleSchemaNames["DstPort"]:  "any",
		},
	}

	var nilSubnet *string = nil
	var nilSubnets []string = nil

	expected := []ne.ACLTemplateInboundRule{
		{
			SeqNo:       ne.Int(1),
			Subnets:     converters.IfArrToStringArr(input[0].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["Subnets"]].([]interface{})),
			Subnet:      nilSubnet,
			Protocol:    ne.String(input[0].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["Protocol"]].(string)),
			SrcPort:     ne.String(input[0].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["SrcPort"]].(string)),
			DstPort:     ne.String(input[0].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["DstPort"]].(string)),
			Description: ne.String(input[0].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["Description"]].(string)),
		},
		{
			SeqNo:    ne.Int(2),
			Subnets:  nilSubnets,
			Subnet:   ne.String(input[1].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["Subnet"]].(string)),
			Protocol: ne.String(input[1].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["Protocol"]].(string)),
			SrcPort:  ne.String(input[1].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["SrcPort"]].(string)),
			DstPort:  ne.String(input[1].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["DstPort"]].(string)),
		},
	}
	// when
	result := expandACLTemplateInboundRules(input)
	// then
	assert.Equal(t, expected, result, "Expanded ACL template inbound rules matches expected result")
}

func TestNetworkACLTemplate_flattenInboundRules(t *testing.T) {
	input := []ne.ACLTemplateInboundRule{
		{
			SeqNo:       ne.Int(1),
			SrcType:     ne.String("SUBNET"),
			Subnets:     []string{"10.0.0.0/24", "1.1.1.1/32"},
			Protocol:    ne.String("TCP"),
			SrcPort:     ne.String("any"),
			DstPort:     ne.String("8080"),
			Description: ne.String("description of inbound rule"),
		},
		{
			SeqNo:    ne.Int(2),
			SrcType:  ne.String("SUBNET"),
			Subnet:   ne.String("3.3.3.3/32"),
			Protocol: ne.String("IP"),
			SrcPort:  ne.String("any"),
			DstPort:  ne.String("any"),
		},
	}
	initial := input
	expected := []interface{}{
		map[string]interface{}{
			networkACLTemplateInboundRuleSchemaNames["SeqNo"]:       input[0].SeqNo,
			networkACLTemplateInboundRuleSchemaNames["SrcType"]:     input[0].SrcType,
			networkACLTemplateInboundRuleSchemaNames["Subnets"]:     input[0].Subnets,
			networkACLTemplateInboundRuleSchemaNames["Subnet"]:      input[0].Subnet,
			networkACLTemplateInboundRuleSchemaNames["Protocol"]:    input[0].Protocol,
			networkACLTemplateInboundRuleSchemaNames["SrcPort"]:     input[0].SrcPort,
			networkACLTemplateInboundRuleSchemaNames["DstPort"]:     input[0].DstPort,
			networkACLTemplateInboundRuleSchemaNames["Description"]: input[0].Description,
		},
		map[string]interface{}{
			networkACLTemplateInboundRuleSchemaNames["SeqNo"]:       input[1].SeqNo,
			networkACLTemplateInboundRuleSchemaNames["SrcType"]:     input[1].SrcType,
			networkACLTemplateInboundRuleSchemaNames["Subnets"]:     input[1].Subnets,
			networkACLTemplateInboundRuleSchemaNames["Subnet"]:      input[1].Subnet,
			networkACLTemplateInboundRuleSchemaNames["Protocol"]:    input[1].Protocol,
			networkACLTemplateInboundRuleSchemaNames["SrcPort"]:     input[1].SrcPort,
			networkACLTemplateInboundRuleSchemaNames["DstPort"]:     input[1].DstPort,
			networkACLTemplateInboundRuleSchemaNames["Description"]: input[1].Description,
		},
	}
	// when
	result := flattenACLTemplateInboundRules(initial, input)
	// then
	assert.Equal(t, expected, result, "Flattened ACL template inbound rules match expected result")
}

func TestNetworkACLTemplate_flattenDeviceDetails(t *testing.T) {
	input := []ne.ACLTemplateDeviceDetails{
		{
			UUID:      ne.String("6e0c2b09-37e9-4fd2-b9e4-e68b08d8b29d"),
			Name:      ne.String("test-device1"),
			ACLStatus: ne.String("Provisioned"),
		},
		{
			UUID:      ne.String("6e0c2b09-37e9-4fd2-b9e4-e68b08d8b29e"),
			Name:      ne.String("test-device2"),
			ACLStatus: ne.String("Provisioning"),
		},
	}
	expected := []interface{}{
		map[string]interface{}{
			networkACLTemplateDeviceDetailSchemaNames["UUID"]:      input[0].UUID,
			networkACLTemplateDeviceDetailSchemaNames["Name"]:      input[0].Name,
			networkACLTemplateDeviceDetailSchemaNames["ACLStatus"]: input[0].ACLStatus,
		},
		map[string]interface{}{
			networkACLTemplateDeviceDetailSchemaNames["UUID"]:      input[1].UUID,
			networkACLTemplateDeviceDetailSchemaNames["Name"]:      input[1].Name,
			networkACLTemplateDeviceDetailSchemaNames["ACLStatus"]: input[1].ACLStatus,
		},
	}
	// when
	result := flattenACLTemplateDeviceDetails(input)
	// then
	assert.Equal(t, expected, result, "Flattened ACL template Device Details match expected result")
}
