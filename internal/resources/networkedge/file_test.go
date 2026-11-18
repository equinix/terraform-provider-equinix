package networkedge

import (
	"context"
	"net/http"
	"testing"

	"github.com/equinix/ne-go"
	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/networkedge"
	"github.com/equinix/terraform-provider-equinix/internal/restmock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fileResourceName = "equinix_network_file"

func fileResource(t *testing.T) *schema.Resource {
	t.Helper()
	r, ok := equinix.Provider().ResourcesMap[fileResourceName]
	require.True(t, ok, "provider has resource %s", fileResourceName)
	return r
}

func TestFile_create(t *testing.T) {
	// given
	api := networkedge.NewMockAPI(t)
	d := schema.TestResourceDataRaw(t, fileResource(t).Schema, map[string]any{
		"file_name":        "testFileName",
		"content":          "testContent",
		"metro_code":       "testMetroCode",
		"device_type_code": "testDeviceTypeCode",
		"process_type":     "testProcessType",
		"self_managed":     true,
		"byol":             true,
	})

	// when
	diags := fileResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "create does not return errors: %v", diags)

	posts := api.RequestsFor(http.MethodPost, networkedge.FilesPath)
	require.Len(t, posts, 1, "one upload request was sent")
	assert.Equal(t, map[string]any{
		"metroCode":            "testMetroCode",
		"deviceTypeCode":       "testDeviceTypeCode",
		"processType":          "testProcessType",
		"licenseType":          ne.DeviceLicenseModeBYOL,
		"deviceManagementType": ne.DeviceManagementTypeSelf,
	}, posts[0].Body, "upload form fields match the API contract")
	assert.Equal(t, map[string]restmock.File{
		"file": {Name: "testFileName", Content: "testContent"},
	}, posts[0].Files, "uploaded file matches")
}

func TestFile_read(t *testing.T) {
	// given
	const id = "183b9c9a-251b-4720-85aa-aa80269a0ffe"
	api := networkedge.NewMockAPI(t)
	api.Files.Seed(id, map[string]any{
		"fileName":       "testFileName",
		"metroCode":      "testMetroCode",
		"deviceTypeCode": "testDeviceTypeCode",
		"processType":    "testProcessType",
		"status":         "testStatus",
	})
	d := schema.TestResourceDataRaw(t, fileResource(t).Schema, map[string]any{})
	d.SetId(id)

	// when
	diags := fileResource(t).ReadContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "read does not return errors: %v", diags)
	assert.Equal(t, id, d.Get("uuid"), "UUID matches")
	assert.Equal(t, "testFileName", d.Get("file_name"), "file name matches")
	assert.Equal(t, "testMetroCode", d.Get("metro_code"), "metro code matches")
	assert.Equal(t, "testDeviceTypeCode", d.Get("device_type_code"), "device type code matches")
	assert.Equal(t, "testProcessType", d.Get("process_type"), "process type matches")
	assert.Equal(t, "testStatus", d.Get("status"), "status matches")
}
