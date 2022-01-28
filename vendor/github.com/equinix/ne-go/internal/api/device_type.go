package api

//DeviceTypeResponse describes response for Network Edge device types query
type DeviceTypeResponse struct {
	Pagination Pagination   `json:"pagination,omitempty"`
	Data       []DeviceType `json:"data,omitempty"`
}

//DeviceType describes Network Edge device type
type DeviceType struct {
	Code                  *string                     `json:"deviceTypeCode,omitempty"`
	Name                  *string                     `json:"name,omitempty"`
	Description           *string                     `json:"description,omitempty"`
	Vendor                *string                     `json:"vendor,omitempty"`
	Category              *string                     `json:"category,omitempty"`
	AvailableMetros       []DeviceTypeAvailableMetro  `json:"availableMetros,omitempty"`
	SoftwarePackages      []DeviceTypeSoftwarePackage `json:"softwarePackages,omitempty"`
	DeviceManagementTypes DeviceManagementTypes       `json:"deviceManagementTypes,omitempty"`
}

//DeviceTypeAvailableMetro describes metro in which network edge device is available
type DeviceTypeAvailableMetro struct {
	Code        *string `json:"metroCode,omitempty"`
	Description *string `json:"metroDescription,omitempty"`
	Region      *string `json:"region,omitempty"`
}

//DeviceTypeSoftwarePackage describes device software package details
type DeviceTypeSoftwarePackage struct {
	Code           *string                    `json:"packageCode,omitempty"`
	Name           *string                    `json:"name,omitempty"`
	VersionDetails []DeviceTypeVersionDetails `json:"versionDetails,omitempty"`
}

//DeviceTypeVersionDetails describes device software version details
type DeviceTypeVersionDetails struct {
	Version                   *string  `json:"version,omitempty"`
	ImageName                 *string  `json:"imageName,omitempty"`
	Date                      *string  `json:"versionDate,omitempty"`
	Status                    *string  `json:"status,omitempty"`
	IsStable                  *bool    `json:"stableVersion,omitempty"`
	AllowedUpgradableVersions []string `json:"allowedUpgradableVersions,omitempty"`
	IsUpgradeAllowed          *bool    `json:"upgradeAllowed,omitempty"`
	ReleaseNotesLink          *string  `json:"releaseNotesLink,omitempty"`
}

//DeviceManagementTypes describes device management types
//offered for a given device type
type DeviceManagementTypes struct {
	EquinixConfigured DeviceManagementType `json:"EQUINIX-CONFIGURED,omitempty"`
	SelfConfigured    DeviceManagementType `json:"SELF-CONFIGURED,omitempty"`
}

//DeviceManagementType describes details of a given device management type
type DeviceManagementType struct {
	Type           string               `json:"type,omitempty"`
	LicenseOptions DeviceLicenseOptions `json:"licenseOptions,omitempty"`
	IsSupported    bool                 `json:"supported,omitempty"`
}

//DeviceLicenseOptions describes licensing options offered for a given
//device management type
type DeviceLicenseOptions struct {
	Sub  DeviceLicenseOption `json:"SUB,omitempty"`
	BYOL DeviceLicenseOption `json:"BYOL,omitempty"`
}

//DeviceLicenseOption describes details of a given licesing option
type DeviceLicenseOption struct {
	Type        *string      `json:"type,omitempty"`
	Name        *string      `json:"name,omitempty"`
	Cores       []DeviceCore `json:"cores,omitempty"`
	IsSupported *bool        `json:"supported,omitempty"`
}

//DeviceCore describes CPU and memory configurations supported by given
//licesing option within given management type
type DeviceCore struct {
	Core         *int                `json:"core,omitempty"`
	Memory       *int                `json:"memory,omitempty"`
	Unit         *string             `json:"unit,omitempty"`
	Flavor       *string             `json:"flavor,omitempty"`
	PackageCodes []DevicePackageCode `json:"packageCodes,omitempty"`
	IsSupported  *bool               `json:"supported,omitempty"`
}

//DevicePackageCode describes device software package code supported
//by given CPU and memory configuration
type DevicePackageCode struct {
	PackageCode *string `json:"packageCode,omitempty"`
	IsSupported *bool   `json:"supported,omitempty"`
}
