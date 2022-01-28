//Package ecx implements Equinix Fabric client
package ecx

const (
	//ConnectionStatusNotAvailable indicates that request to create connection was not sent
	//to the provider. Applicable for provider status only
	ConnectionStatusNotAvailable = "NOT_AVAILABLE"
	//ConnectionStatusPendingApproval indicates that connection awaits provider's approval.
	ConnectionStatusPendingApproval = "PENDING_APPROVAL"
	//ConnectionStatusPendingAutoApproval indicates that connection is in process of
	//automatic approval
	ConnectionStatusPendingAutoApproval = "PENDING_AUTO_APPROVAL"
	//ConnectionStatusProvisioning indicates that connection is in creation process
	ConnectionStatusProvisioning = "PROVISIONING"
	//ConnectionStatusRejected indicates that provider has rejected the connection
	ConnectionStatusRejected = "REJECTED"
	//ConnectionStatusPendingBGPPeering indicates that connection was approved by provider and
	//awaits for BGP peering configuration on provider side
	ConnectionStatusPendingBGPPeering = "PENDING_BGP_PEERING"
	//ConnectionStatusPendingProviderVlan indicates that connection awaits for provider approval
	//and vlan assignment
	ConnectionStatusPendingProviderVlan = "PENDING_PROVIDER_VLAN"
	//ConnectionStatusProvisioned indicates that connection is created successfully
	ConnectionStatusProvisioned = "PROVISIONED"
	//ConnectionStatusAvailable indicates that connection is established.
	//Applicable for provider status only
	ConnectionStatusAvailable = "AVAILABLE"
	//ConnectionStatusPendingDelete indicates that connection is in deletion process and awaits
	//for providers approval to be removed
	ConnectionStatusPendingDelete = "PENDING_DELETE"
	//ConnectionStatusDeprovisioning indicates that connection is being removed
	ConnectionStatusDeprovisioning = "DEPROVISIONING"
	//ConnectionStatusDeprovisioned indicates that connection is removed
	ConnectionStatusDeprovisioned = "DEPROVISIONED"
	//ConnectionStatusDeleted indicates that connection was administratively deleted
	ConnectionStatusDeleted = "DELETED"
)

//Client describes operations provided by Equinix Fabric client module
type Client interface {
	GetUserPorts() ([]Port, error)

	GetL2OutgoingConnections(statuses []string) ([]L2Connection, error)
	GetL2Connection(uuid string) (*L2Connection, error)
	CreateL2Connection(conn L2Connection) (*string, error)
	CreateL2RedundantConnection(priConn L2Connection, secConn L2Connection) (*string, *string, error)
	NewL2ConnectionUpdateRequest(uuid string) L2ConnectionUpdateRequest
	DeleteL2Connection(uuid string) error
	ConfirmL2Connection(uuid string, confirmConn L2ConnectionToConfirm) (*L2ConnectionConfirmation, error)

	GetL2SellerProfiles() ([]L2ServiceProfile, error)
	GetL2ServiceProfile(uuid string) (*L2ServiceProfile, error)
	CreateL2ServiceProfile(sp L2ServiceProfile) (*string, error)
	UpdateL2ServiceProfile(sp L2ServiceProfile) error
	DeleteL2ServiceProfile(uuid string) error
}

//L2ConnectionUpdateRequest describes composite request to update given Layer2 connection
type L2ConnectionUpdateRequest interface {
	WithName(name string) L2ConnectionUpdateRequest
	WithBandwidth(speed int, speedUnit string) L2ConnectionUpdateRequest
	WithSpeed(speed int) L2ConnectionUpdateRequest
	WithSpeedUnit(speedUnit string) L2ConnectionUpdateRequest
	Execute() error
}

//Error describes Equinix Fabric error that occurs during API call processing
type Error struct {
	//ErrorCode is short error identifier
	ErrorCode string
	//ErrorMessage is textual description of an error
	ErrorMessage string
}

//L2Connection describes layer 2 connection managed by Equinix Fabric
type L2Connection struct {
	UUID                *string
	Name                *string
	ProfileUUID         *string
	Speed               *int
	SpeedUnit           *string
	Status              *string
	ProviderStatus      *string
	Notifications       []string
	PurchaseOrderNumber *string
	PortUUID            *string
	DeviceUUID          *string
	DeviceInterfaceID   *int
	VlanSTag            *int
	VlanCTag            *int
	NamedTag            *string
	AdditionalInfo      []L2ConnectionAdditionalInfo
	ZSidePortUUID       *string
	ZSideVlanSTag       *int
	ZSideVlanCTag       *int
	SellerRegion        *string
	SellerMetroCode     *string
	AuthorizationKey    *string
	RedundantUUID       *string
	RedundancyType      *string
	Actions             []L2ConnectionAction
}

//L2ConnectionAdditionalInfo additional info object used in L2 connections
type L2ConnectionAdditionalInfo struct {
	Name  *string
	Value *string
}

//L2ConnectionAction describes pending actions to complete connection provisioning
type L2ConnectionAction struct {
	Type         *string
	Message      *string
	OperationID  *string
	RequiredData []L2ConnectionActionData
}

//L2ConnectionActionData describes data required for a given to complete
type L2ConnectionActionData struct {
	Key               *string
	Label             *string
	Value             *string
	IsEditable        *bool
	ValidationPattern *string
}

//L2ConnectionToConfirm accepts the hosted connection in the seller side
type L2ConnectionToConfirm struct {
	AccessKey *string
	SecretKey *string
}

//L2ConnectionConfirmation describes a connection confirmed
type L2ConnectionConfirmation struct {
	PrimaryConnectionID *string
	Message             *string
}

//L2ServiceProfile describes layer 2 service profile managed by Equinix Fabric
type L2ServiceProfile struct {
	UUID                                *string
	State                               *string
	AlertPercentage                     *float64
	AllowCustomSpeed                    *bool
	AllowOverSubscription               *bool
	APIAvailable                        *bool
	AuthKeyLabel                        *string
	ConnectionNameLabel                 *string
	CTagLabel                           *string
	EnableAutoGenerateServiceKey        *bool
	EquinixManagedPortAndVlan           *bool
	Features                            L2ServiceProfileFeatures
	IntegrationID                       *string
	Name                                *string
	OnBandwidthThresholdNotification    []string
	OnProfileApprovalRejectNotification []string
	OnVcApprovalRejectionNotification   []string
	OverSubscription                    *string
	Ports                               []L2ServiceProfilePort
	Private                             *bool
	PrivateUserEmails                   []string
	RequiredRedundancy                  *bool
	SpeedBands                          []L2ServiceProfileSpeedBand
	SpeedFromAPI                        *bool
	TagType                             *string
	VlanSameAsPrimary                   *bool
	Description                         *string
	Metros                              []L2SellerProfileMetro
	AdditionalInfos                     []L2SellerProfileAdditionalInfo
	Encapsulation                       *string
	GlobalOrganization                  *string
	OrganizationName                    *string
}

//L2ServiceProfilePort describes port used in L2 service profile
type L2ServiceProfilePort struct {
	ID        *string
	MetroCode *string
}

//L2ServiceProfileSpeedBand describes speed / bandwidth used in L2 service profile
type L2ServiceProfileSpeedBand struct {
	Speed     *int
	SpeedUnit *string
}

//L2ServiceProfileFeatures describes features used in L2 service profile
type L2ServiceProfileFeatures struct {
	CloudReach  *bool
	TestProfile *bool
}

//Port describes Equinix Fabric's user port
type Port struct {
	UUID          *string
	Name          *string
	Region        *string
	IBX           *string
	MetroCode     *string
	Priority      *string
	Encapsulation *string
	Buyout        *bool
	Bandwidth     *string
	Status        *string
}

//L2SellerProfileMetro describes details of a metro in which service provices is present
type L2SellerProfileMetro struct {
	Code    *string
	Name    *string
	IBXes   []string
	Regions map[string]string
}

//L2SellerProfileAdditionalInfo describes additional information that might be provided by service buyer when using given seller profile
type L2SellerProfileAdditionalInfo struct {
	Name             *string
	Description      *string
	DataType         *string
	IsMandatory      *bool
	IsCaptureInEmail *bool
}
