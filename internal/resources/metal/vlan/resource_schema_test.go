package vlan

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"
)

func TestResourceSchema(t *testing.T) {
	s := resourceSchema(context.Background())

	assert.NotNil(t, s)
	assert.Contains(t, s.Attributes, "id")
	assert.Contains(t, s.Attributes, "project_id")
	assert.Contains(t, s.Attributes, "description")
	assert.Contains(t, s.Attributes, "facility")
	assert.Contains(t, s.Attributes, "metro")
	assert.Contains(t, s.Attributes, "vxlan")

	idAttr := s.Attributes["id"].(schema.StringAttribute)
	assert.True(t, idAttr.Computed)
	assert.Equal(t, "The unique identifier for this Metal Vlan", idAttr.Description)

	projectIDAttr := s.Attributes["project_id"].(schema.StringAttribute)
	assert.True(t, projectIDAttr.Required)
	assert.Equal(t, "ID of parent project", projectIDAttr.Description)

	facilityAttr := s.Attributes["facility"].(schema.StringAttribute)
	assert.True(t, facilityAttr.Optional)
	assert.True(t, facilityAttr.Computed)
	assert.NotEmpty(t, facilityAttr.DeprecationMessage)

	metroAttr := s.Attributes["metro"].(schema.StringAttribute)
	assert.True(t, metroAttr.Optional)
	assert.True(t, metroAttr.Computed)
	assert.Equal(t, "Metro in which to create the VLAN", metroAttr.Description)
}
