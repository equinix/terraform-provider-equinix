package networkedge

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/equinix/terraform-provider-equinix/internal/networkedge"
	"github.com/equinix/terraform-provider-equinix/internal/restmock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const deviceResourceName = "equinix_network_device"

func deviceResource(t *testing.T) *schema.Resource {
	t.Helper()
	r, ok := equinix.Provider().ResourcesMap[deviceResourceName]
	require.True(t, ok, "provider has resource %s", deviceResourceName)
	return r
}

// writeLicenseFile creates a license file in a temporary directory and returns its path.
func writeLicenseFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestDevice_create(t *testing.T) {
	// given
	api := networkedge.NewMockAPI(t)
	licensePath := writeLicenseFile(t, "test.lic", "license")
	d := schema.TestResourceDataRaw(t, deviceResource(t).Schema, map[string]any{
		"name":                      "device",
		"type_code":                 "CSR1000V",
		"metro_code":                "SV",
		"throughput":                100,
		"tier":                      1,
		"throughput_unit":           "Mbps",
		"hostname":                  "test",
		"package_code":              "SEC",
		"version":                   "9.0.1",
		"byol":                      false,
		"license_token":             "sWf3df4gaAvbbexw45ga4f",
		"license_file":              licensePath,
		"acl_template_id":           "a624178c-6d59-4798-9a7f-2ddf2c7c5881",
		"account_number":            "123456",
		"notifications":             []any{"bla@bla.com"},
		"purchase_order_number":     "1234567890",
		"term_length":               1,
		"additional_bandwidth":      50,
		"order_reference":           "12312121sddsf1231",
		"interface_count":           10,
		"wan_interface_id":          "5",
		"core_count":                2,
		"self_managed":              false,
		"vendor_configuration":      map[string]any{"key": "value"},
		"ssh_key":                   []any{map[string]any{"username": "user", "key_name": "key"}},
		"connectivity":              "INTERNET-ACCESS",
		"project_id":                "68ccfd49-39b1-478e-957a-67c72f719d7a",
		"diverse_device_id":         "ed7891bd-15b4-4f72-ac56-d96cfdacddcc",
		"generate_default_password": false,
	})

	// when
	diags := deviceResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "create does not return errors: %v", diags)

	// The license file is uploaded first, then its ID is sent in the create request.
	uploads := api.RequestsFor(http.MethodPost, networkedge.LicenseFilesPath)
	require.Len(t, uploads, 1, "one license file upload request was sent")
	assert.Equal(t, map[string]any{
		"metroCode":            "SV",
		"deviceTypeCode":       "CSR1000V",
		"licenseType":          "BYOL",
		"deviceManagementType": "SELF-CONFIGURED",
	}, uploads[0].Body, "license file upload form fields match the API contract")
	assert.Equal(t, map[string]restmock.File{"file": {Name: "test.lic", Content: "license"}}, uploads[0].Files, "uploaded file matches")
	licenseFileIDs := api.LicenseFiles.IDs()
	require.Len(t, licenseFileIDs, 1, "API holds one license file")

	posts := api.RequestsFor(http.MethodPost, networkedge.DevicesPath)
	require.Len(t, posts, 1, "one create request was sent")
	assert.Equal(t, map[string]any{
		"virtualDeviceName":       "device",
		"deviceTypeCode":          "CSR1000V",
		"metroCode":               "SV",
		"throughput":              "100",
		"throughputUnit":          "Mbps",
		"tier":                    float64(1),
		"hostNamePrefix":          "test",
		"packageCode":             "SEC",
		"version":                 "9.0.1",
		"licenseMode":             "Sub",
		"licenseToken":            "sWf3df4gaAvbbexw45ga4f",
		"licenseFileId":           licenseFileIDs[0],
		"aclTemplateUuid":         "a624178c-6d59-4798-9a7f-2ddf2c7c5881",
		"accountNumber":           "123456",
		"notifications":           []any{"bla@bla.com"},
		"purchaseOrderNumber":     "1234567890",
		"termLength":              "1",
		"additionalBandwidth":     "50",
		"orderReference":          "12312121sddsf1231",
		"interfaceCount":          float64(10),
		"sshInterfaceId":          "5",
		"core":                    float64(2),
		"deviceManagementType":    "EQUINIX-CONFIGURED",
		"vendorConfig":            map[string]any{"key": "value"},
		"userPublicKey":           map[string]any{"username": "user", "keyName": "key"},
		"connectivity":            "INTERNET-ACCESS",
		"projectId":               "68ccfd49-39b1-478e-957a-67c72f719d7a",
		"diverseFromDeviceUuid":   "ed7891bd-15b4-4f72-ac56-d96cfdacddcc",
		"generateDefaultPassword": false,
	}, posts[0].Body, "create request body matches the API contract")

	// Creating a device waits for the device to be provisioned and licensed, and for its ACL template.
	assert.NotEmpty(t, api.RequestsFor(http.MethodGet, networkedge.DevicesPath+"/"+d.Id()), "device status was checked")
	assert.NotEmpty(t, api.RequestsFor(http.MethodGet, networkedge.DevicesPath+"/"+d.Id()+"/acl"), "ACL template status was checked")
}

func TestDevice_createWithSecondary(t *testing.T) {
	// given
	api := networkedge.NewMockAPI(t)
	licensePath := writeLicenseFile(t, "secondary.lic", "secondary license")
	d := schema.TestResourceDataRaw(t, deviceResource(t).Schema, map[string]any{
		"name":       "primary",
		"type_code":  "CSR1000V",
		"metro_code": "SV",
		"secondary_device": []any{map[string]any{
			"name":                 "device",
			"metro_code":           "SV",
			"hostname":             "SV5",
			"wan_interface_id":     "5",
			"license_token":        "sWf3df4gaAvbbexw45ga4f",
			"license_file":         licensePath,
			"acl_template_id":      "a624178c-6d59-4798-9a7f-2ddf2c7c5881",
			"account_number":       "123456",
			"notifications":        []any{"bla@bla.com"},
			"additional_bandwidth": 50,
			"vendor_configuration": map[string]any{"key": "value"},
			"ssh_key":              []any{map[string]any{"username": "user", "key_name": "testKey"}},
		}},
	})

	// when
	diags := deviceResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "create does not return errors: %v", diags)
	licenseFileIDs := api.LicenseFiles.IDs()
	require.Len(t, licenseFileIDs, 1, "API holds the secondary license file")

	posts := api.RequestsFor(http.MethodPost, networkedge.DevicesPath)
	require.Len(t, posts, 1, "one create request was sent")
	assert.Equal(t, map[string]any{
		"virtualDeviceName":   "device",
		"metroCode":           "SV",
		"hostNamePrefix":      "SV5",
		"sshInterfaceId":      "5",
		"licenseToken":        "sWf3df4gaAvbbexw45ga4f",
		"licenseFileId":       licenseFileIDs[0],
		"aclTemplateUuid":     "a624178c-6d59-4798-9a7f-2ddf2c7c5881",
		"accountNumber":       "123456",
		"notifications":       []any{"bla@bla.com"},
		"additionalBandwidth": "50",
		"vendorConfig":        map[string]any{"key": "value"},
		"userPublicKey":       map[string]any{"username": "user", "keyName": "testKey"},
	}, posts[0].Body["secondary"], "secondary device in the create request matches the API contract")
	assert.Len(t, api.Devices.IDs(), 2, "API holds the primary and the secondary device")
}

func TestDevice_read(t *testing.T) {
	// given
	const id = "0452fa68-8246-48b1-a1b2-817fb4baddcb"
	api := networkedge.NewMockAPI(t)
	api.Devices.Seed(id, map[string]any{
		"name":                    "device",
		"deviceTypeCode":          "CSR1000V",
		"projectId":               "68ccfd49-39b1-478e-957a-67c72f719d7a",
		"status":                  "PROVISIONED",
		"metroCode":               "SV",
		"throughput":              "100",
		"throughputUnit":          "Mbps",
		"hostName":                "test",
		"core":                    map[string]any{"core": 2, "tier": 1},
		"packageCode":             "SEC",
		"version":                 "9.0.1",
		"licenseType":             "BYOL",
		"licenseToken":            "sWf3df4gaAvbbexw45ga4f",
		"aclTemplateUuid":         "a624178c-6d59-4798-9a7f-2ddf2c7c5881",
		"accountNumber":           "123456",
		"notifications":           []any{"bla@bla.com"},
		"purchaseOrderNumber":     "1234567890",
		"termLength":              1,
		"additionalBandwidth":     50,
		"orderReference":          "12312121sddsf1231",
		"interfaceCount":          10,
		"sshInterfaceId":          "6",
		"deviceManagementType":    "SELF-CONFIGURED",
		"diverseFromDeviceUuid":   "68ccfd49-39b1-478e-957a-67c72f719d7a",
		"diverseFromDeviceName":   "diverseFromDeviceName",
		"generateDefaultPassword": false,
		"vendorConfig":            map[string]any{"key": "value"},
		"userPublicKey":           map[string]any{"username": "user", "publicKeyName": "key"},
		"asn":                     11222,
		"zoneCode":                "Zone2",
	})
	d := schema.TestResourceDataRaw(t, deviceResource(t).Schema, map[string]any{})
	d.SetId(id)

	// when
	diags := deviceResource(t).ReadContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "read does not return errors: %v", diags)
	assert.Equal(t, "device", d.Get("name"), "name matches")
	assert.Equal(t, "CSR1000V", d.Get("type_code"), "type code matches")
	assert.Equal(t, "SV", d.Get("metro_code"), "metro code matches")
	assert.Equal(t, 100, d.Get("throughput"), "throughput matches")
	assert.Equal(t, "Mbps", d.Get("throughput_unit"), "throughput unit matches")
	assert.Equal(t, "test", d.Get("hostname"), "hostname matches")
	assert.Equal(t, "SEC", d.Get("package_code"), "package code matches")
	assert.Equal(t, "9.0.1", d.Get("version"), "version matches")
	assert.Equal(t, true, d.Get("byol"), "BYOL matches")
	assert.Empty(t, d.Get("license_token"), "license token is not read from the API")
	assert.Equal(t, "a624178c-6d59-4798-9a7f-2ddf2c7c5881", d.Get("acl_template_id"), "ACL template ID matches")
	assert.Equal(t, "123456", d.Get("account_number"), "account number matches")
	assert.Equal(t, []string{"bla@bla.com"}, converters.SetToStringList(d.Get("notifications").(*schema.Set)), "notifications match")
	assert.Equal(t, "1234567890", d.Get("purchase_order_number"), "purchase order number matches")
	assert.Equal(t, 1, d.Get("term_length"), "term length matches")
	assert.Equal(t, 50, d.Get("additional_bandwidth"), "additional bandwidth matches")
	assert.Equal(t, "12312121sddsf1231", d.Get("order_reference"), "order reference matches")
	assert.Equal(t, 10, d.Get("interface_count"), "interface count matches")
	assert.Empty(t, d.Get("wan_interface_id"), "WAN interface ID is not read from the API")
	assert.Equal(t, 2, d.Get("core_count"), "core count matches")
	assert.Equal(t, 1, d.Get("tier"), "tier matches")
	assert.Equal(t, true, d.Get("self_managed"), "self managed matches")
	assert.Equal(t, map[string]any{"key": "value"}, d.Get("vendor_configuration"), "vendor configuration matches")
	keys := d.Get("ssh_key").(*schema.Set).List()
	require.Len(t, keys, 1, "one SSH key")
	assert.Equal(t, map[string]any{"username": "user", "key_name": "key"}, keys[0], "SSH key matches")
	assert.Equal(t, 11222, d.Get("asn"), "ASN matches")
	assert.Equal(t, "Zone2", d.Get("zone_code"), "zone code matches")
	assert.Equal(t, "68ccfd49-39b1-478e-957a-67c72f719d7a", d.Get("project_id"), "project ID matches")
	assert.Equal(t, "68ccfd49-39b1-478e-957a-67c72f719d7a", d.Get("diverse_device_id"), "diverse device ID matches")
	assert.Equal(t, "diverseFromDeviceName", d.Get("diverse_device_name"), "diverse device name matches")
	assert.Equal(t, false, d.Get("generate_default_password"), "generate default password matches")
}

func TestDevice_readSecondary(t *testing.T) {
	// given
	const (
		primaryID   = "0452fa68-8246-48b1-a1b2-817fb4baddcb"
		secondaryID = "c2a147a3-ff47-4a24-a6e5-d6d7ce6459f3"
	)
	api := networkedge.NewMockAPI(t)
	api.Devices.Seed(primaryID, map[string]any{"name": "primary", "status": "PROVISIONED", "redundantUuid": secondaryID})
	api.Devices.Seed(secondaryID, map[string]any{
		"name":                "device",
		"status":              "PROVISIONED",
		"licenseStatus":       "APPLIED",
		"metroCode":           "SV",
		"ibx":                 "SV5",
		"region":              "AMER",
		"hostName":            "test",
		"licenseFileId":       "d72dbe58-e596-4698-8b57-0a38e8077d25",
		"aclTemplateUuid":     "a624178c-6d59-4798-9a7f-2ddf2c7c5881",
		"sshIpAddress":        "1.1.1.1",
		"sshIpFqdn":           "test-1.1.1.1-SV.test.equinix.com",
		"accountNumber":       "123456",
		"notifications":       []any{"bla@bla.com"},
		"redundancyType":      "SECONDARY",
		"redundantUuid":       primaryID,
		"additionalBandwidth": 50,
		"projectId":           "68ccfd49-39b1-478e-957a-67c72f719d7a",
		"interfaces": []any{map[string]any{
			"id":                1,
			"name":              "GigabitEthernet1",
			"status":            "AVAILABLE",
			"operationalStatus": "UP",
			"macAddress":        "58-0A-C9-7A-DA-E9",
			"ipAddress":         "2.2.2.2",
			"assignedType":      "test-connection(AWS Direct Connect)",
			"type":              "DATA",
		}},
		"vendorConfig":  map[string]any{"key": "value"},
		"userPublicKey": map[string]any{"username": "user", "publicKeyName": "testKey"},
		"asn":           11222,
		"zoneCode":      "Zone2",
	})
	d := schema.TestResourceDataRaw(t, deviceResource(t).Schema, map[string]any{})
	d.SetId(primaryID)
	// A read keeps the license file path, license token and cloud-init file ID of the
	// secondary device from the state. The API does not return them.
	require.NoError(t, d.Set("secondary_device", []any{map[string]any{
		"license_file":       "/tmp/licenseFileSec",
		"license_token":      "sWf3df4gaAvbbexw45ga4f",
		"cloud_init_file_id": "59a29dac-73fb-4e09-8aba-215b30fed63e",
	}}))

	// when
	diags := deviceResource(t).ReadContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "read does not return errors: %v", diags)
	assert.Equal(t, secondaryID, d.Get("redundant_id"), "redundant ID matches")
	secondary := d.Get("secondary_device").([]any)
	require.Len(t, secondary, 1, "one secondary device")
	sec := secondary[0].(map[string]any)
	assert.Equal(t, secondaryID, sec["uuid"], "secondary UUID matches")
	assert.Equal(t, "device", sec["name"], "secondary name matches")
	assert.Equal(t, "PROVISIONED", sec["status"], "secondary status matches")
	assert.Equal(t, "APPLIED", sec["license_status"], "secondary license status matches")
	assert.Equal(t, "SV", sec["metro_code"], "secondary metro code matches")
	assert.Equal(t, "SV5", sec["ibx"], "secondary IBX matches")
	assert.Equal(t, "AMER", sec["region"], "secondary region matches")
	assert.Equal(t, "test", sec["hostname"], "secondary hostname matches")
	assert.Equal(t, "sWf3df4gaAvbbexw45ga4f", sec["license_token"], "secondary license token is kept from the state")
	assert.Equal(t, "d72dbe58-e596-4698-8b57-0a38e8077d25", sec["license_file_id"], "secondary license file ID matches")
	assert.Equal(t, "/tmp/licenseFileSec", sec["license_file"], "secondary license file path is kept from the state")
	assert.Equal(t, "59a29dac-73fb-4e09-8aba-215b30fed63e", sec["cloud_init_file_id"], "secondary cloud-init file ID is kept from the state")
	assert.Equal(t, "a624178c-6d59-4798-9a7f-2ddf2c7c5881", sec["acl_template_id"], "secondary ACL template ID matches")
	assert.Equal(t, "1.1.1.1", sec["ssh_ip_address"], "secondary SSH IP address matches")
	assert.Equal(t, "test-1.1.1.1-SV.test.equinix.com", sec["ssh_ip_fqdn"], "secondary SSH FQDN matches")
	assert.Equal(t, "123456", sec["account_number"], "secondary account number matches")
	assert.Equal(t, []string{"bla@bla.com"}, converters.SetToStringList(sec["notifications"].(*schema.Set)), "secondary notifications match")
	assert.Equal(t, "SECONDARY", sec["redundancy_type"], "secondary redundancy type matches")
	assert.Equal(t, primaryID, sec["redundant_id"], "secondary redundant ID matches")
	assert.Equal(t, 50, sec["additional_bandwidth"], "secondary additional bandwidth matches")
	assert.Equal(t, "68ccfd49-39b1-478e-957a-67c72f719d7a", sec["project_id"], "secondary project ID matches")
	assert.Equal(t, []any{map[string]any{
		"id":                 1,
		"name":               "GigabitEthernet1",
		"status":             "AVAILABLE",
		"operational_status": "UP",
		"mac_address":        "58-0A-C9-7A-DA-E9",
		"ip_address":         "2.2.2.2",
		"assigned_type":      "test-connection(AWS Direct Connect)",
		"type":               "DATA",
	}}, sec["interface"], "secondary interfaces match")
	assert.Equal(t, map[string]any{"key": "value"}, sec["vendor_configuration"], "secondary vendor configuration matches")
	assert.Equal(t, 11222, sec["asn"], "secondary ASN matches")
	assert.Equal(t, "Zone2", sec["zone_code"], "secondary zone code matches")
	assert.NotEmpty(t, sec["ssh_key"], "secondary SSH key is set")
}

func TestDevice_delete(t *testing.T) {
	// given
	const id = "0452fa68-8246-48b1-a1b2-817fb4baddcb"
	api := networkedge.NewMockAPI(t)
	api.Devices.Seed(id, map[string]any{"name": "device", "status": "PROVISIONED"})
	d := schema.TestResourceDataRaw(t, deviceResource(t).Schema, map[string]any{})
	d.SetId(id)

	// when
	diags := deviceResource(t).DeleteContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "delete does not return errors: %v", diags)
	deletes := api.RequestsFor(http.MethodDelete, networkedge.DevicesPath+"/"+id)
	require.Len(t, deletes, 1, "one delete request was sent")
	assert.Equal(t, "true", deletes[0].Query.Get("deleteRedundantDevice"), "the redundant device is deleted too")
	// Delete waits until the device is deprovisioned.
	assert.NotEmpty(t, api.RequestsFor(http.MethodGet, networkedge.DevicesPath+"/"+id), "device status was checked")
}

func TestDevice_update(t *testing.T) {
	// given
	const id = "0452fa68-8246-48b1-a1b2-817fb4baddcb"
	api := networkedge.NewMockAPI(t)
	api.Devices.Seed(id, map[string]any{
		"name":                "device",
		"status":              "PROVISIONED",
		"core":                map[string]any{"core": 2},
		"additionalBandwidth": 50,
	})
	r := deviceResource(t)
	state := &terraform.InstanceState{
		ID: id,
		Attributes: map[string]string{
			"id":                   id,
			"name":                 "device",
			"core_count":           "2",
			"additional_bandwidth": "50",
		},
	}
	config := terraform.NewResourceConfigRaw(map[string]any{
		"name":                 "device",
		"core_count":           4,
		"additional_bandwidth": 100,
	})
	diff, err := schema.InternalMap(r.Schema).Diff(context.Background(), state, config, nil, nil, false)
	require.NoError(t, err)
	d, err := schema.InternalMap(r.Schema).Data(state, diff)
	require.NoError(t, err)

	// when
	diags := r.UpdateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "update does not return errors: %v", diags)
	patches := api.RequestsFor(http.MethodPatch, networkedge.DevicesPath+"/"+id)
	require.Len(t, patches, 1, "one device update request was sent")
	assert.Equal(t, float64(4), patches[0].Body["core"], "core count in the device update request matches")
	puts := api.RequestsFor(http.MethodPut, networkedge.DevicesPath+"/"+id+"/additionalBandwidths")
	require.Len(t, puts, 1, "one additional bandwidth update request was sent")
	assert.Equal(t, map[string]any{"additionalBandwidth": float64(100)}, puts[0].Body, "additional bandwidth update request body matches the API contract")
	// Update waits until the additional bandwidth is provisioned.
	assert.NotEmpty(t, api.RequestsFor(http.MethodGet, networkedge.DevicesPath+"/"+id+"/additionalBandwidths"), "additional bandwidth status was checked")
}

func newDeviceCreateData(t *testing.T, extra map[string]any) *schema.ResourceData {
	t.Helper()
	raw := map[string]any{
		"name":       "device",
		"type_code":  "CSR1000V",
		"metro_code": "SV",
	}
	for k, v := range extra {
		raw[k] = v
	}
	return schema.TestResourceDataRaw(t, deviceResource(t).Schema, raw)
}

// newDeviceUpdateData returns resource data with a stored state and a changed configuration.
func newDeviceUpdateData(t *testing.T, id string, state map[string]string, config map[string]any) *schema.ResourceData {
	t.Helper()
	r := deviceResource(t)
	instanceState := &terraform.InstanceState{ID: id, Attributes: state}
	diff, err := schema.InternalMap(r.Schema).Diff(context.Background(), instanceState, terraform.NewResourceConfigRaw(config), nil, nil, false)
	require.NoError(t, err)
	d, err := schema.InternalMap(r.Schema).Data(instanceState, diff)
	require.NoError(t, err)
	return d
}

func TestDevice_createWaitsUntilProvisioned(t *testing.T) {
	t.Parallel()
	// given
	api := networkedge.NewMockAPI(t)
	api.Devices.Sequence("", "status", "PROVISIONING", "PROVISIONED")
	d := newDeviceCreateData(t, nil)

	// when
	diags := deviceResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "create does not return errors: %v", diags)
	assert.Equal(t, "PROVISIONED", d.Get("status"), "status matches")
	assert.GreaterOrEqual(t, len(api.RequestsFor(http.MethodGet, networkedge.DevicesPath+"/"+d.Id())), 3,
		"create kept checking the status until the device was provisioned")
}

func TestDevice_createWaitsUntilLicensed(t *testing.T) {
	t.Parallel()
	// given
	api := networkedge.NewMockAPI(t)
	// The first check belongs to the wait for provisioning. The next ones belong to the wait for the license.
	api.Devices.Sequence("", "licenseStatus", "APPLYING_LICENSE", "APPLYING_LICENSE", "APPLIED")
	d := newDeviceCreateData(t, nil)

	// when
	diags := deviceResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "create does not return errors: %v", diags)
	assert.Equal(t, "APPLIED", d.Get("license_status"), "license status matches")
	assert.GreaterOrEqual(t, len(api.RequestsFor(http.MethodGet, networkedge.DevicesPath+"/"+d.Id())), 4,
		"create kept checking the license status until the license was applied")
}

func TestDevice_createWaitsUntilACLTemplateProvisioned(t *testing.T) {
	t.Parallel()
	// given
	api := networkedge.NewMockAPI(t)
	api.Devices.Script("", "aclStatus", "PROVISIONING", "PROVISIONED")
	d := newDeviceCreateData(t, map[string]any{"acl_template_id": "a624178c-6d59-4798-9a7f-2ddf2c7c5881"})

	// when
	diags := deviceResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "create does not return errors: %v", diags)
	assert.Len(t, api.RequestsFor(http.MethodGet, networkedge.DevicesPath+"/"+d.Id()+"/acl"), 2,
		"create checked the ACL template status until it was provisioned")
}

func TestDevice_createFailsOnUnexpectedStatus(t *testing.T) {
	t.Parallel()
	// given
	api := networkedge.NewMockAPI(t)
	api.Devices.Sequence("", "status", "FAILED")
	d := newDeviceCreateData(t, nil)

	// when
	diags := deviceResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.True(t, diags.HasError(), "create returns an error")
	assert.Contains(t, diags[0].Summary, "error waiting for network device", "error comes from the wait")
}

func TestDevice_deleteWaitsUntilDeprovisioned(t *testing.T) {
	t.Parallel()
	// given
	const id = "0452fa68-8246-48b1-a1b2-817fb4baddcb"
	api := networkedge.NewMockAPI(t)
	api.Devices.Seed(id, map[string]any{"name": "device", "status": "PROVISIONED"})
	api.Devices.Sequence(id, "status", "DEPROVISIONING", "DEPROVISIONED")
	d := schema.TestResourceDataRaw(t, deviceResource(t).Schema, map[string]any{})
	d.SetId(id)

	// when
	diags := deviceResource(t).DeleteContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "delete does not return errors: %v", diags)
	assert.Len(t, api.RequestsFor(http.MethodGet, networkedge.DevicesPath+"/"+id), 2,
		"delete checked the status until the device was deprovisioned")
}

func TestDevice_updateWaitsUntilResourceUpgradeFinished(t *testing.T) {
	t.Parallel()
	// given
	const id = "0452fa68-8246-48b1-a1b2-817fb4baddcb"
	api := networkedge.NewMockAPI(t)
	api.Devices.Seed(id, map[string]any{"name": "device", "status": "PROVISIONED", "core": map[string]any{"core": 2}})
	api.Devices.Sequence(id, "status", "RESOURCE_UPGRADE_IN_PROGRESS", "PROVISIONED")
	d := newDeviceUpdateData(t, id,
		map[string]string{"id": id, "name": "device", "core_count": "2"},
		map[string]any{"name": "device", "core_count": 4})

	// when
	diags := deviceResource(t).UpdateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "update does not return errors: %v", diags)
	assert.Equal(t, "PROVISIONED", d.Get("status"), "status matches")
	assert.GreaterOrEqual(t, len(api.RequestsFor(http.MethodGet, networkedge.DevicesPath+"/"+id)), 3,
		"update kept checking the status until the resource upgrade finished")
}

func TestDevice_updateWaitsUntilAdditionalBandwidthProvisioned(t *testing.T) {
	t.Parallel()
	// given
	const id = "0452fa68-8246-48b1-a1b2-817fb4baddcb"
	api := networkedge.NewMockAPI(t)
	api.Devices.Seed(id, map[string]any{"name": "device", "status": "PROVISIONED", "additionalBandwidth": 50})
	api.Devices.Script(id, "additionalBandwidthStatus", "PROVISIONING", "PROVISIONED")
	d := newDeviceUpdateData(t, id,
		map[string]string{"id": id, "name": "device", "additional_bandwidth": "50"},
		map[string]any{"name": "device", "additional_bandwidth": 100})

	// when
	diags := deviceResource(t).UpdateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "update does not return errors: %v", diags)
	assert.Len(t, api.RequestsFor(http.MethodGet, networkedge.DevicesPath+"/"+id+"/additionalBandwidths"), 2,
		"update checked the additional bandwidth status until it was provisioned")
}
