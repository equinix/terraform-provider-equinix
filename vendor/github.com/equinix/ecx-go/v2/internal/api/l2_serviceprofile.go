package api

//L2ServiceProfile l2 service profile resource used in get, post and put operations
type L2ServiceProfile struct {
	UUID                                *string                         `json:"uuid,omitempty"`
	State                               *string                         `json:"state,omitempty"`
	AlertPercentage                     *float64                        `json:"alertPercentage,omitempty"`
	AllowCustomSpeed                    *bool                           `json:"allowCustomSpeed"`
	AllowOverSubscription               *bool                           `json:"allowOverSubscription"`
	APIAvailable                        *bool                           `json:"apiAvailable"`
	AuthKeyLabel                        *string                         `json:"authKeyLabel,omitempty"`
	ConnectionNameLabel                 *string                         `json:"connectionNameLabel,omitempty"`
	CTagLabel                           *string                         `json:"ctagLabel,omitempty"`
	EnableAutoGenerateServiceKey        *bool                           `json:"enableAutoGenerateServiceKey"`
	EquinixManagedPortAndVlan           *bool                           `json:"equinixManagedPortAndVlan"`
	Features                            L2ServiceProfileFeatures        `json:"features,omitempty"`
	IntegrationID                       *string                         `json:"integrationId,omitempty"`
	Name                                *string                         `json:"name,omitempty"`
	OnBandwidthThresholdNotification    []string                        `json:"onBandwidthThresholdNotification,omitempty"`
	OnProfileApprovalRejectNotification []string                        `json:"onProfileApprovalRejectNotification,omitempty"`
	OnVcApprovalRejectionNotification   []string                        `json:"onVcApprovalRejectionNotification,omitempty"`
	OverSubscription                    *string                         `json:"overSubscription,omitempty"`
	Ports                               []L2ServiceProfilePort          `json:"ports,omitempty"`
	Private                             *bool                           `json:"private"`
	PrivateUserEmails                   []string                        `json:"privateUserEmails,omitempty"`
	RequiredRedundancy                  *bool                           `json:"requiredRedundancy"`
	SpeedBands                          []L2ServiceProfileSpeedBand     `json:"speedBands,omitempty"`
	SpeedFromAPI                        *bool                           `json:"speedFromAPI"`
	TagType                             *string                         `json:"tagType,omitempty"`
	VlanSameAsPrimary                   *bool                           `json:"vlanSameAsPrimary"`
	Description                         *string                         `json:"description,omitempty"`
	Metros                              []L2SellerProfileMetro          `json:"metros,omitempty"`
	AdditionalInfos                     []L2SellerProfileAdditionalInfo `json:"additionalBuyerInfo,omitempty"`
	ProfileEncapsulation                *string                         `json:"profileEncapsulation,omitempty"`
	GlobalOrganization                  *string                         `json:"globalOrganization,omitempty"`
	OrganizationName                    *string                         `json:"organizationName,omitempty"`
}

//L2ServiceProfileDeleteResponse delete l2 service profile response
type L2ServiceProfileDeleteResponse struct {
	Message *string `json:"message,omitempty"`
	Status  *string `json:"status,omitempty"`
}

//L2ServiceProfilePort port used in L2 service profile
type L2ServiceProfilePort struct {
	ID        *string `json:"id,omitempty"`
	MetroCode *string `json:"metroCode,omitempty"`
}

//L2ServiceProfileSpeedBand speed band used in L2 service profile
type L2ServiceProfileSpeedBand struct {
	Speed     *int    `json:"speed,omitempty"`
	SpeedUnit *string `json:"unit,omitempty"`
}

//L2ServiceProfileFeatures feature used in L2 service profile
type L2ServiceProfileFeatures struct {
	CloudReach  *bool `json:"cloudReach"`
	TestProfile *bool `json:"testProfile"`
}

//CreateL2ServiceProfileResponse post l2 service profile response
type CreateL2ServiceProfileResponse struct {
	UUID *string `json:"uuid,omitempty"`
}

//L2SellerProfilesResponse response with list of l2 selller profiles
type L2SellerProfilesResponse struct {
	IsLastPage  *bool              `json:"isLastPage"`
	IsFirstPage *bool              `json:"isFirstPage"`
	TotalCount  *int               `json:"totalCount,omitempty"`
	PageSize    *int               `json:"PageSize,omitempty"`
	Content     []L2ServiceProfile `json:"content,omitempty"`
}

//L2SellerProfileMetro describces details of a metro in which service provider is present
type L2SellerProfileMetro struct {
	Code    *string           `json:"code,omitempty"`
	Name    *string           `json:"name,omitempty"`
	IBXs    []string          `json:"ibxs,omitempty"`
	Regions map[string]string `json:"sellerRegions,omitempty"`
}

//L2SellerProfileAdditionalInfo describces additional information that might be provided by service buyer when using given seller profil
type L2SellerProfileAdditionalInfo struct {
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	Mandatory      *bool   `json:"mandatory"`
	DataType       *string `json:"datatype,omitempty"`
	CaptureInEmail *bool   `json:"captureInEmail"`
}
