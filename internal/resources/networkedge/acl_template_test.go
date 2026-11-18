package networkedge

import (
	"context"
	"net/http"
	"testing"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/networkedge"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const aclTemplateResourceName = "equinix_network_acl_template"

func aclTemplateResource(t *testing.T) *schema.Resource {
	t.Helper()
	r, ok := equinix.Provider().ResourcesMap[aclTemplateResourceName]
	require.True(t, ok, "provider has resource %s", aclTemplateResourceName)
	return r
}

func TestACLTemplate_create(t *testing.T) {
	// given
	api := networkedge.NewMockAPI(t)
	d := schema.TestResourceDataRaw(t, aclTemplateResource(t).Schema, map[string]any{
		"name":        "test",
		"description": "testTemplate",
		"project_id":  "68ccfd49-39b1-478e-957a-67c72f719d7a",
		"inbound_rule": []any{
			map[string]any{
				"subnet":      "10.0.0.0/24",
				"protocol":    "TCP",
				"src_port":    "any",
				"dst_port":    "8080",
				"description": "test inbound rule",
			},
			map[string]any{
				"protocol": "IP",
				"src_port": "any",
				"dst_port": "any",
			},
		},
	})

	// when
	diags := aclTemplateResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "create does not return errors: %v", diags)

	posts := api.RequestsFor(http.MethodPost, networkedge.ACLTemplatesPath)
	require.Len(t, posts, 1, "one create request was sent")
	body := posts[0].Body
	assert.Equal(t, "test", body["name"], "name matches")
	assert.Equal(t, "testTemplate", body["description"], "description matches")
	assert.Equal(t, "68ccfd49-39b1-478e-957a-67c72f719d7a", body["projectId"], "project ID matches")
	// Terraform gives every nested attribute a value, so an unset subnet or
	// description reaches the API as an empty string. Sequence numbers follow the rule order.
	assert.Equal(t, []any{
		map[string]any{
			"seqNo":       float64(1),
			"subnet":      "10.0.0.0/24",
			"protocol":    "TCP",
			"srcPort":     "any",
			"dstPort":     "8080",
			"description": "test inbound rule",
		},
		map[string]any{
			"seqNo":       float64(2),
			"subnet":      "",
			"protocol":    "IP",
			"srcPort":     "any",
			"dstPort":     "any",
			"description": "",
		},
	}, body["inboundRules"], "inbound rules match the API contract")
}

func TestACLTemplate_read(t *testing.T) {
	// given
	const id = "9a2f9d62-6fbc-4e2c-b6a5-4f5a3a1cbbc1"
	api := networkedge.NewMockAPI(t)
	api.ACLTemplates.Seed(id, map[string]any{
		"name":        "test",
		"description": "testTemplate",
		"inboundRules": []any{
			map[string]any{
				"seqNo":       1,
				"subnet":      "10.0.0.0/24",
				"protocol":    "TCP",
				"srcPort":     "any",
				"dstPort":     "8080",
				"description": "test inbound rule",
			},
		},
		"virtualDeviceDetails": []any{
			map[string]any{"uuid": "6e0c2b09-37e9-4fd2-b9e4-e68b08d8b29d", "name": "test-device1", "aclStatus": "Provisioned"},
			map[string]any{"uuid": "6e0c2b09-37e9-4fd2-b9e4-e68b08d8b29e", "name": "test-device2", "aclStatus": "Provisioning"},
		},
	})
	d := schema.TestResourceDataRaw(t, aclTemplateResource(t).Schema, map[string]any{})
	d.SetId(id)

	// when
	diags := aclTemplateResource(t).ReadContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "read does not return errors: %v", diags)
	assert.Equal(t, "test", d.Get("name"), "name matches")
	assert.Equal(t, "testTemplate", d.Get("description"), "description matches")
	assert.Equal(t, []any{
		map[string]any{
			"sequence_number": 1,
			"subnet":          "10.0.0.0/24",
			"protocol":        "TCP",
			"src_port":        "any",
			"dst_port":        "8080",
			"description":     "test inbound rule",
		},
	}, d.Get("inbound_rule"), "inbound rules match")
	assert.Equal(t, []any{
		map[string]any{"uuid": "6e0c2b09-37e9-4fd2-b9e4-e68b08d8b29d", "name": "test-device1", "acl_status": "Provisioned"},
		map[string]any{"uuid": "6e0c2b09-37e9-4fd2-b9e4-e68b08d8b29e", "name": "test-device2", "acl_status": "Provisioning"},
	}, d.Get("device_details"), "device details match")
}
