package ecx

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/equinix/ecx-go/v2/internal/api"
	"github.com/equinix/rest-go"
)

//GetL2SellerProfiles operations retrieves available layer2 seller service profiles
func (c RestClient) GetL2SellerProfiles() ([]L2ServiceProfile, error) {
	path := "/ecx/v3/l2/serviceprofiles/services"
	content, err := c.GetPaginated(path, &api.L2SellerProfilesResponse{}, rest.DefaultPagingConfig().
		SetSizeParamName("pageSize").
		SetPageParamName("pageNumber").
		SetFirstPageNumber(0))
	if err != nil {
		return nil, err
	}
	transformed := make([]L2ServiceProfile, len(content))
	for i := range content {
		transformed[i] = *mapL2ServiceProfileAPIToDomain(content[i].(api.L2ServiceProfile))
	}
	return transformed, nil
}

//GetL2ServiceProfile operation retrieves layer 2 servie profile with a given UUID
func (c RestClient) GetL2ServiceProfile(uuid string) (*L2ServiceProfile, error) {
	path := "/ecx/v3/l2/serviceprofiles/" + url.PathEscape(uuid)
	respBody := api.L2ServiceProfile{}
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, http.MethodGet, path); err != nil {
		return nil, err
	}
	return mapL2ServiceProfileAPIToDomain(respBody), nil
}

//CreateL2ServiceProfile operation creates layer 2 service profile with a given profile structure.
//Upon successful creation, connection structure with assigned UUID will be returned
func (c RestClient) CreateL2ServiceProfile(l2profile L2ServiceProfile) (*string, error) {
	path := "/ecx/v3/l2/serviceprofiles"
	reqBody := mapL2ServiceProfileDomainToAPI(l2profile)
	respBody := api.CreateL2ServiceProfileResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)
	if err := c.Execute(req, http.MethodPost, path); err != nil {
		return nil, err
	}
	return respBody.UUID, nil
}

//UpdateL2ServiceProfile operation updates layer 2 service profile by replacing existing profile with a given profile structure.
//Target profile structure needs to have UUID defined
func (c RestClient) UpdateL2ServiceProfile(sp L2ServiceProfile) error {
	if StringValue(sp.UUID) == "" {
		return fmt.Errorf("target profile structure needs to have UUID defined")
	}
	path := "/ecx/v3/l2/serviceprofiles"
	reqBody := mapL2ServiceProfileDomainToAPI(sp)
	respBody := api.CreateL2ServiceProfileResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)
	if err := c.Execute(req, http.MethodPut, path); err != nil {
		return err
	}
	return nil
}

//DeleteL2ServiceProfile deletes layer 2 service profile with a given UUID
func (c RestClient) DeleteL2ServiceProfile(uuid string) error {
	path := "/ecx/v3/l2/serviceprofiles/" + url.PathEscape(uuid)
	respBody := api.L2ServiceProfileDeleteResponse{}
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, http.MethodDelete, path); err != nil {
		return err
	}
	return nil
}

func mapL2ServiceProfileDomainToAPI(l2profile L2ServiceProfile) api.L2ServiceProfile {
	return api.L2ServiceProfile{
		UUID:                                l2profile.UUID,
		State:                               l2profile.State,
		AlertPercentage:                     l2profile.AlertPercentage,
		AllowCustomSpeed:                    l2profile.AllowCustomSpeed,
		AllowOverSubscription:               l2profile.AllowOverSubscription,
		APIAvailable:                        l2profile.APIAvailable,
		AuthKeyLabel:                        l2profile.AuthKeyLabel,
		ConnectionNameLabel:                 l2profile.ConnectionNameLabel,
		CTagLabel:                           l2profile.CTagLabel,
		EnableAutoGenerateServiceKey:        l2profile.EnableAutoGenerateServiceKey,
		EquinixManagedPortAndVlan:           l2profile.EquinixManagedPortAndVlan,
		Features:                            mapFeaturesDomainToAPI(l2profile.Features),
		IntegrationID:                       l2profile.IntegrationID,
		Name:                                l2profile.Name,
		OnBandwidthThresholdNotification:    l2profile.OnBandwidthThresholdNotification,
		OnProfileApprovalRejectNotification: l2profile.OnProfileApprovalRejectNotification,
		OnVcApprovalRejectionNotification:   l2profile.OnVcApprovalRejectionNotification,
		OverSubscription:                    l2profile.OverSubscription,
		Ports:                               mapPortsDomainToAPI(l2profile.Ports),
		Private:                             l2profile.Private,
		PrivateUserEmails:                   l2profile.PrivateUserEmails,
		RequiredRedundancy:                  l2profile.RequiredRedundancy,
		SpeedBands:                          mapSpeedBandsDomainToAPI(l2profile.SpeedBands),
		SpeedFromAPI:                        l2profile.SpeedFromAPI,
		TagType:                             l2profile.TagType,
		VlanSameAsPrimary:                   l2profile.VlanSameAsPrimary,
		Description:                         l2profile.Description,
	}
}

func mapL2ServiceProfileAPIToDomain(apiProfile api.L2ServiceProfile) *L2ServiceProfile {
	return &L2ServiceProfile{
		UUID:                                apiProfile.UUID,
		State:                               apiProfile.State,
		AlertPercentage:                     apiProfile.AlertPercentage,
		AllowCustomSpeed:                    apiProfile.AllowCustomSpeed,
		AllowOverSubscription:               apiProfile.AllowOverSubscription,
		APIAvailable:                        apiProfile.APIAvailable,
		AuthKeyLabel:                        apiProfile.AuthKeyLabel,
		ConnectionNameLabel:                 apiProfile.ConnectionNameLabel,
		CTagLabel:                           apiProfile.CTagLabel,
		EnableAutoGenerateServiceKey:        apiProfile.EnableAutoGenerateServiceKey,
		EquinixManagedPortAndVlan:           apiProfile.EquinixManagedPortAndVlan,
		Features:                            mapFeaturesAPIToDomain(apiProfile.Features),
		IntegrationID:                       apiProfile.IntegrationID,
		Name:                                apiProfile.Name,
		OnBandwidthThresholdNotification:    apiProfile.OnBandwidthThresholdNotification,
		OnProfileApprovalRejectNotification: apiProfile.OnProfileApprovalRejectNotification,
		OnVcApprovalRejectionNotification:   apiProfile.OnVcApprovalRejectionNotification,
		OverSubscription:                    apiProfile.OverSubscription,
		Ports:                               mapPortsAPIToDomain(apiProfile.Ports),
		Private:                             apiProfile.Private,
		PrivateUserEmails:                   apiProfile.PrivateUserEmails,
		RequiredRedundancy:                  apiProfile.RequiredRedundancy,
		SpeedBands:                          mapSpeedBandsAPIToDomain(apiProfile.SpeedBands),
		SpeedFromAPI:                        apiProfile.SpeedFromAPI,
		TagType:                             apiProfile.TagType,
		VlanSameAsPrimary:                   apiProfile.VlanSameAsPrimary,
		Description:                         apiProfile.Description,
		Metros:                              mapL2SellerProfileMetrosAPIToDomain(apiProfile.Metros),
		AdditionalInfos:                     mapL2SellerProfileAdditionalInfosAPIToDomain(apiProfile.AdditionalInfos),
		Encapsulation:                       apiProfile.ProfileEncapsulation,
		OrganizationName:                    apiProfile.OrganizationName,
		GlobalOrganization:                  apiProfile.GlobalOrganization,
	}
}

func mapFeaturesDomainToAPI(features L2ServiceProfileFeatures) api.L2ServiceProfileFeatures {
	return api.L2ServiceProfileFeatures{
		CloudReach:  features.CloudReach,
		TestProfile: features.TestProfile,
	}
}

func mapFeaturesAPIToDomain(apiFeatures api.L2ServiceProfileFeatures) L2ServiceProfileFeatures {
	return L2ServiceProfileFeatures{
		CloudReach:  apiFeatures.CloudReach,
		TestProfile: apiFeatures.TestProfile,
	}
}

func mapPortsDomainToAPI(ports []L2ServiceProfilePort) []api.L2ServiceProfilePort {
	apiPorts := make([]api.L2ServiceProfilePort, len(ports))
	for i, v := range ports {
		apiPorts[i] = api.L2ServiceProfilePort{
			ID:        v.ID,
			MetroCode: v.MetroCode,
		}
	}
	return apiPorts
}

func mapPortsAPIToDomain(apiPorts []api.L2ServiceProfilePort) []L2ServiceProfilePort {
	ports := make([]L2ServiceProfilePort, len(apiPorts))
	for i, v := range apiPorts {
		ports[i] = L2ServiceProfilePort{
			ID:        v.ID,
			MetroCode: v.MetroCode,
		}
	}
	return ports
}

func mapSpeedBandsDomainToAPI(bands []L2ServiceProfileSpeedBand) []api.L2ServiceProfileSpeedBand {
	apiBands := make([]api.L2ServiceProfileSpeedBand, len(bands))
	for i, v := range bands {
		apiBands[i] = api.L2ServiceProfileSpeedBand{
			Speed:     v.Speed,
			SpeedUnit: v.SpeedUnit,
		}
	}
	return apiBands
}

func mapSpeedBandsAPIToDomain(apiBands []api.L2ServiceProfileSpeedBand) []L2ServiceProfileSpeedBand {
	bands := make([]L2ServiceProfileSpeedBand, len(apiBands))
	for i, v := range apiBands {
		bands[i] = L2ServiceProfileSpeedBand{
			Speed:     v.Speed,
			SpeedUnit: v.SpeedUnit,
		}
	}
	return bands
}

func mapL2SellerProfileMetrosAPIToDomain(apiMetros []api.L2SellerProfileMetro) []L2SellerProfileMetro {
	transformed := make([]L2SellerProfileMetro, len(apiMetros))
	for i := range apiMetros {
		transformed[i] = L2SellerProfileMetro{
			Code:    apiMetros[i].Code,
			Name:    apiMetros[i].Name,
			IBXes:   apiMetros[i].IBXs,
			Regions: apiMetros[i].Regions}
	}
	return transformed
}

func mapL2SellerProfileAdditionalInfosAPIToDomain(apiInfos []api.L2SellerProfileAdditionalInfo) []L2SellerProfileAdditionalInfo {
	transformed := make([]L2SellerProfileAdditionalInfo, len(apiInfos))
	for i := range apiInfos {
		transformed[i] = L2SellerProfileAdditionalInfo{
			Name:             apiInfos[i].Name,
			Description:      apiInfos[i].Description,
			DataType:         apiInfos[i].DataType,
			IsMandatory:      apiInfos[i].Mandatory,
			IsCaptureInEmail: apiInfos[i].CaptureInEmail}
	}
	return transformed
}
