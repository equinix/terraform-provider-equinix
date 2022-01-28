package ne

import (
	"fmt"
	"net/url"

	"github.com/equinix/ne-go/internal/api"
	"github.com/equinix/rest-go"
)

//GetDeviceTypes retrieves list of devices types along with their details
func (c RestClient) GetDeviceTypes() ([]DeviceType, error) {
	path := "/ne/v1/deviceTypes"
	content, err := c.GetOffsetPaginated(path, &api.DeviceTypeResponse{},
		rest.DefaultOffsetPagingConfig())
	if err != nil {
		return nil, err
	}
	transformed := make([]DeviceType, len(content))
	for i := range content {
		transformed[i] = mapDeviceTypeAPIToDomain(content[i].(api.DeviceType))
	}
	return transformed, nil
}

//GetDeviceSoftwareVersions retrieves list of available software versions for a given device type
func (c RestClient) GetDeviceSoftwareVersions(deviceTypeCode string) ([]DeviceSoftwareVersion, error) {
	deviceType, err := c.getDeviceType(deviceTypeCode)
	if err != nil {
		return nil, err
	}
	return mapDeviceTypeAPIToDeviceSoftwareVersions(*deviceType), nil
}

//GetDevicePlatforms retrieves list of available platform configurations for a given device type
func (c RestClient) GetDevicePlatforms(deviceTypeCode string) ([]DevicePlatform, error) {
	deviceType, err := c.getDeviceType(deviceTypeCode)
	if err != nil {
		return nil, err
	}
	return mapDeviceTypeAPIToDevicePlatforms(*deviceType), nil
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Unexported package methods
//_______________________________________________________________________

func (c RestClient) getDeviceType(typeCode string) (*api.DeviceType, error) {
	path := "/ne/v1/deviceTypes"
	content, err := c.GetOffsetPaginated(path, &api.DeviceTypeResponse{},
		rest.DefaultOffsetPagingConfig().
			SetAdditionalParams(map[string]string{"deviceTypeCode": url.QueryEscape(typeCode)}))
	if err != nil {
		return nil, err
	}
	if len(content) < 1 {
		return nil, fmt.Errorf("device type query returned no results for given type code: %s", typeCode)
	}
	if len(content) > 1 {
		return nil, fmt.Errorf("device type query returned more than one result for a given type code: %s", typeCode)
	}
	devType := content[0].(api.DeviceType)
	return &devType, nil
}

func mapDeviceTypeAPIToDomain(apiDevice api.DeviceType) DeviceType {
	return DeviceType{
		Name:        apiDevice.Name,
		Code:        apiDevice.Code,
		Description: apiDevice.Description,
		Vendor:      apiDevice.Vendor,
		Category:    apiDevice.Category,
		MetroCodes:  mapDeviceTypeAvailableMetrosAPIToDomain(apiDevice.AvailableMetros),
	}
}

func mapDeviceTypeAvailableMetrosAPIToDomain(apiMetros []api.DeviceTypeAvailableMetro) []string {
	transformed := make([]string, len(apiMetros))
	for i := range apiMetros {
		transformed[i] = StringValue(apiMetros[i].Code)
	}
	return transformed
}

func mapDeviceTypeAPIToDeviceSoftwareVersions(apiType api.DeviceType) []DeviceSoftwareVersion {
	versionMap := make(map[string]*DeviceSoftwareVersion)
	for _, apiPkg := range apiType.SoftwarePackages {
		for _, apiVer := range apiPkg.VersionDetails {
			apiVerStr := StringValue(apiVer.Version)
			ver, ok := versionMap[apiVerStr]
			if !ok {
				ver = mapDeviceSoftwareVersionAPIToDomain(apiVer)
				versionMap[apiVerStr] = ver
			}
			ver.PackageCodes = append(ver.PackageCodes, StringValue(apiPkg.Code))
		}
	}
	transformed := make([]DeviceSoftwareVersion, 0, len(versionMap))
	for ver := range versionMap {
		transformed = append(transformed, *versionMap[ver])
	}
	return transformed
}

func mapDeviceSoftwareVersionAPIToDomain(apiVer api.DeviceTypeVersionDetails) *DeviceSoftwareVersion {
	return &DeviceSoftwareVersion{
		Version:          apiVer.Version,
		ImageName:        apiVer.ImageName,
		Date:             apiVer.Date,
		Status:           apiVer.Status,
		IsStable:         apiVer.IsStable,
		ReleaseNotesLink: apiVer.ReleaseNotesLink,
		PackageCodes:     []string{},
	}
}

func mapDeviceTypeAPIToDevicePlatforms(apiType api.DeviceType) []DevicePlatform {
	mgmtPlatforms := flattenMgmtType(apiType.DeviceManagementTypes.EquinixConfigured)
	mgmtPlatforms = append(mgmtPlatforms, flattenMgmtType(apiType.DeviceManagementTypes.SelfConfigured)...)
	pMap := make(map[string]*DevicePlatform)
	for i := range mgmtPlatforms {
		platform, ok := pMap[StringValue(mgmtPlatforms[i].Flavor)]
		if !ok {
			platform = &mgmtPlatforms[i]
			pMap[StringValue(mgmtPlatforms[i].Flavor)] = platform
		} else {
			platform.ManagementTypes = append(platform.ManagementTypes, mgmtPlatforms[i].ManagementTypes...)
		}
	}
	return devPlatfromMapToSlice(pMap)
}

func flattenMgmtType(mgmtType api.DeviceManagementType) []DevicePlatform {
	if !mgmtType.IsSupported {
		return []DevicePlatform{}
	}
	licPlatforms := mapLicenseOption(mgmtType.LicenseOptions.BYOL)
	licPlatforms = append(licPlatforms, mapLicenseOption(mgmtType.LicenseOptions.Sub)...)
	pMap := make(map[string]*DevicePlatform)
	for i := range licPlatforms {
		platform, ok := pMap[StringValue(licPlatforms[i].Flavor)]
		if !ok {
			platform = &licPlatforms[i]
			platform.ManagementTypes = append(platform.ManagementTypes, mgmtType.Type)
			pMap[StringValue(licPlatforms[i].Flavor)] = platform
		} else {
			platform.LicenseOptions = append(platform.LicenseOptions, licPlatforms[i].LicenseOptions...)
		}
	}
	return devPlatfromMapToSlice(pMap)
}

func mapLicenseOption(licOption api.DeviceLicenseOption) []DevicePlatform {
	if !BoolValue(licOption.IsSupported) {
		return []DevicePlatform{}
	}
	transformed := make([]DevicePlatform, len(licOption.Cores))
	for i := range licOption.Cores {
		transformed[i] = DevicePlatform{
			Flavor:          licOption.Cores[i].Flavor,
			CoreCount:       licOption.Cores[i].Core,
			Memory:          licOption.Cores[i].Memory,
			MemoryUnit:      licOption.Cores[i].Unit,
			PackageCodes:    mapPackageCodesAPIToDomain(licOption.Cores[i].PackageCodes),
			ManagementTypes: []string{},
			LicenseOptions:  []string{StringValue(licOption.Type)},
		}
	}
	return transformed
}

func mapPackageCodesAPIToDomain(apiCodes []api.DevicePackageCode) []string {
	transformed := make([]string, 0, len(apiCodes))
	for _, apiCode := range apiCodes {
		if !BoolValue(apiCode.IsSupported) {
			continue
		}
		transformed = append(transformed, StringValue(apiCode.PackageCode))
	}
	return transformed
}

func devPlatfromMapToSlice(devPlatformMap map[string]*DevicePlatform) []DevicePlatform {
	transformed := make([]DevicePlatform, 0, len(devPlatformMap))
	for _, v := range devPlatformMap {
		transformed = append(transformed, *v)
	}
	return transformed
}
