package networkedge

import (
	"net/http"

	"github.com/equinix/terraform-provider-equinix/internal/restmock"
)

// LicenseFilesPath is the path that device license files are uploaded to.
const LicenseFilesPath = DevicesPath + "/licenseFiles"

// addDeviceRoutes registers the device endpoints that need more than the generic collection behavior.
// Devices are provisioned, licensed and deprovisioned without delay.
func (m *MockAPI) addDeviceRoutes() {
	m.LicenseFiles = m.AddCollection(LicenseFilesPath, "metroCode", "deviceTypeCode", "licenseType", "deviceManagementType")
	m.LicenseFiles.CreatedIDProperty = "fileId"

	// A create request can contain a secondary device. The response names both devices.
	m.Handle("POST "+DevicesPath, func(req restmock.Request, _ func(string) string) (int, any) {
		primaryID := m.Devices.Add(newDevice(req.Body))
		resp := map[string]any{"uuid": primaryID}
		if secondary, ok := req.Body["secondary"].(map[string]any); ok {
			secondaryID := m.Devices.Add(newDevice(secondary))
			m.Devices.Merge(primaryID, map[string]any{"redundantUuid": secondaryID})
			m.Devices.Merge(secondaryID, map[string]any{"redundantUuid": primaryID})
			resp["secondaryUuid"] = secondaryID
		}
		return http.StatusOK, resp
	})

	m.Handle("PATCH "+DevicesPath+"/{id}", func(req restmock.Request, id func(string) string) (int, any) {
		fields := map[string]any{}
		if v, ok := req.Body["virtualDeviceName"]; ok {
			fields["name"] = v
		}
		if v, ok := req.Body["core"]; ok {
			fields["core"] = map[string]any{"core": v}
		}
		return updated(m.Devices.Merge(id("id"), fields))
	})

	m.Handle("DELETE "+DevicesPath+"/{id}", func(req restmock.Request, id func(string) string) (int, any) {
		deviceID := id("id")
		doc, ok := m.Devices.Get(deviceID)
		if !ok {
			return http.StatusNotFound, map[string]any{"errorCode": "EQ-MOCK-404", "errorMessage": "device does not exist"}
		}
		deprovisioned := map[string]any{"status": "DEPROVISIONED"}
		m.Devices.Merge(deviceID, deprovisioned)
		if secondaryID, ok := doc["redundantUuid"].(string); ok && req.Query.Get("deleteRedundantDevice") == "true" {
			m.Devices.Merge(secondaryID, deprovisioned)
		}
		return http.StatusNoContent, nil
	})

	m.Handle("GET "+DevicesPath+"/{id}/acl", func(_ restmock.Request, id func(string) string) (int, any) {
		if _, ok := m.Devices.Get(id("id")); !ok {
			return http.StatusNotFound, nil
		}
		return http.StatusOK, map[string]any{"status": m.scripted(id("id"), "aclStatus")}
	})

	m.Handle("PATCH "+DevicesPath+"/{id}/acl", func(req restmock.Request, id func(string) string) (int, any) {
		fields := map[string]any{}
		for _, k := range []string{"aclTemplateUuid", "mgmtAclTemplateUuid"} {
			if v, ok := req.Body[k]; ok {
				fields[k] = v
			}
		}
		return updated(m.Devices.Merge(id("id"), fields))
	})

	m.Handle("GET "+DevicesPath+"/{id}/additionalBandwidths", func(_ restmock.Request, id func(string) string) (int, any) {
		doc, ok := m.Devices.Get(id("id"))
		if !ok {
			return http.StatusNotFound, nil
		}
		return http.StatusOK, map[string]any{"status": m.scripted(id("id"), "additionalBandwidthStatus"), "additionalBandwidth": doc["additionalBandwidth"]}
	})

	m.Handle("PUT "+DevicesPath+"/{id}/additionalBandwidths", func(req restmock.Request, id func(string) string) (int, any) {
		return updated(m.Devices.Merge(id("id"), map[string]any{"additionalBandwidth": req.Body["additionalBandwidth"]}))
	})
}

// scripted returns the next scripted status of a device, or PROVISIONED if the test scripted none.
// Tests script a status with m.Devices.Script(id, "aclStatus", ...) or
// m.Devices.Script(id, "additionalBandwidthStatus", ...).
func (m *MockAPI) scripted(deviceID, key string) any {
	if v, ok := m.Devices.Next(deviceID, key); ok {
		return v
	}
	return "PROVISIONED"
}

// newDevice builds the stored document of a device that a create request describes.
func newDevice(req map[string]any) map[string]any {
	return map[string]any{
		"name":            req["virtualDeviceName"],
		"deviceTypeCode":  req["deviceTypeCode"],
		"metroCode":       req["metroCode"],
		"aclTemplateUuid": req["aclTemplateUuid"],
		"status":          "PROVISIONED",
		"licenseStatus":   "APPLIED",
	}
}

func updated(found bool) (int, any) {
	if !found {
		return http.StatusNotFound, map[string]any{"errorCode": "EQ-MOCK-404", "errorMessage": "device does not exist"}
	}
	return http.StatusNoContent, nil
}
