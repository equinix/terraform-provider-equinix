package api

//Port describes Equinix Fabric's user port
type Port struct {
	UUID            *string `json:"uuid,omitempty"`
	Name            *string `json:"name,omitempty"`
	Region          *string `json:"region,omitempty"`
	IBX             *string `json:"ibx,omitempty"`
	MetroCode       *string `json:"metroCode,omitempty"`
	DevicePriority  *string `json:"devicePriority,omitempty"`
	Encapsulation   *string `json:"encapsulation,omitempty"`
	Buyout          *bool   `json:"buyout"`
	TotalBandwidth  *int64  `json:"totalBandwidth,omitempty"`
	ProvisionStatus *string `json:"provisionStatus,omitempty"`
}
