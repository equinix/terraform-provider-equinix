package equinix

import (
	"context"
	"fmt"
	v4 "github.com/equinix/terraform-provider-equinix/internal/apis/fabric/v4"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"sort"
	"time"
)

//variables should be snake case
var ecxFabricConnectionSchemaNames = map[string]string{
	"Uuid":           "uuid",
	"Name":           "name",
	"Type_":          "type",
	"Href":           "Href",
	"Description":    "description",
	"State":          "state",
	"Change":         "Change",
	"Operation":      "operation",
	"Order":          "order",
	"Notifications":  "Notifications",
	"Account":        "account",
	"ChangeLog":      "changeLog",
	"Bandwidth":      "bandwidth",
	"Redundancy":     "redundancy",
	"IsRemote":       "isRemote",
	"Direction":      "direction",
	"ASide":          "aSide",
	"ZSide":          "zSide",
	"AdditionalInfo": "additionalInfo",
}

var ecxFabricConnectionSchemaDescription = map[string]string{
	"Uuid":           "Unique identifier of the connection, will be part of created response",
	"Name":           "Connection name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
	"Type_":          "Defines the connection type like VG_VC, EVPL_VC, EPL_VC, EC_VC, GW_VC, ACCESS_EPL_VC ",
	"Href":           "Connection URI information",
	"Description":    "Customer-provided connection description",
	"State":          "Connection overall state",
	"Change":         "Represents latest change request and its state information",
	"Operation":      "Connection type-specific operational data",
	"Order":          "Order related to this connection information",
	"Notifications":  "Preferences for notifications on connection configuration or status changes",
	"Account":        "Customer account information that is associated with this connection",
	"ChangeLog":      "Captures connection lifecycle change information",
	"Bandwidth":      "Connection bandwidth in Mbps",
	"Redundancy":     "Redundancy Information",
	"IsRemote":       "Connection property derived from access point locations",
	"Direction":      "Connection directionality from the requester point of view",
	"ASide":          "Requester or Customer side connection configuration object of the multi-segment connection",
	"ZSide":          "Destination or Provider side connection configuration object of the multi-segment connection",
	"AdditionalInfo": "Connection additional information",
}

var ecxFabricConnectionAdditionalInfoSchemaNames = map[string]string{
	"Name":  "name",
	"Value": "value",
}

var ecxFabricConnectionAdditionalInfoDescriptions = map[string]string{
	"Name":  "Additional information key",
	"Value": "Additional information value",
}

var ecxFabricConnectionChangeSchemaNames = map[string]string{
	"Uuid":             "uuid",
	"Type_":            "type",
	"Status":           "status",
	"CreationDateTime": "creationDateTime",
	"UpdatedDateTime":  "updatedDateTime",
	"Information":      "information",
	"Data":             "data",
}

var ecxFabricConnectionChangeDescription = map[string]string{
	"Uuid":             "Unique identifier of the change",
	"Type_":            "Type of change",
	"Status":           "Current outcome of the change flow",
	"CreationDateTime": "Time change request received",
	"UpdatedDateTime":  "Record last updated",
	"Information":      "Additional information",
	"Data":             "Change operation data",
}

var ecxFabricConnectionChangeOperationSchemaNames = map[string]string{
	"Op":    "Op",
	"Path":  "Path",
	"Value": "Value",
}

var ecxFabricConnectionChangeOperationDescription = map[string]string{
	"Op":    "operation name",
	"Path":  "path inside document leading to updated parameter",
	"Value": "new value for updated parameter",
}

var ecxFabricConnectionOperationalSchemaNames = map[string]string{
	"ProviderStatus":    "providerStatus",
	"EquinixStatus":     "equinixStatus",
	"OperationalStatus": "operationalStatus",
	"Errors":            "errors",
	"OpStatusChangedAt": "opStatusChangedAt",
}

var ecxFabricConnectionOperationalDescription = map[string]string{
	"ProviderStatus":    "Connection provider readiness status",
	"EquinixStatus":     "Connection status",
	"OperationalStatus": "Connection operational status",
	"Errors":            "Errors occurred",
	"OpStatusChangedAt": "When connection transitioned into current operational status",
}

var ecxFabricConnectionOrderSchemaNames = map[string]string{
	"PurchaseOrderNumber": "purchaseOrderNumber",
	"BillingTier":         "billingTier",
	"OrderId":             "orderId",
	"OrderNumber":         "orderNumber",
}

var ecxFabricConnectionOrderDescription = map[string]string{
	"PurchaseOrderNumber": "Purchase order number",
	"BillingTier":         "Billing tier for connection bandwidth",
	"OrderId":             "Order Identification",
	"OrderNumber":         "Order Reference Number",
}

var ecxFabricConnectionNotificationSchemaNames = map[string]string{
	"Type_":        "type",
	"SendInterval": "sendInterval",
	"Emails":       "emails",
}

var ecxFabricConnectionNotificationDescription = map[string]string{
	"Type_":        "Notification Type",
	"SendInterval": "send interval",
	"Emails":       "Array of contact emails",
}

var ecxFabricConnectionAccountSchemaNames = map[string]string{
	"AccountNumber":          "accountNumber",
	"AccountName":            "accountName",
	"OrgId":                  "orgId",
	"OrganizationName":       "organizationName",
	"GlobalOrgId":            "globalOrgId",
	"GlobalOrganizationName": "globalOrganizationName",
	"UcmId":                  "ucmId",
	"GlobalCustId":           "globalCustId",
}

var ecxFabricConnectionAccountDescription = map[string]string{
	"AccountNumber":          "Account Number",
	"AccountName":            "Account Name",
	"OrgId":                  "Customer organization identifier",
	"OrganizationName":       "Customer organization name",
	"GlobalOrgId":            "Global organization identifier",
	"GlobalOrganizationName": "Global organization name",
	"UcmId":                  "System unique identifier ",
	"GlobalCustId":           "Global Customer organization identifier",
}

var ecxFabricConnectionChangeLogSchemaNames = map[string]string{
	"CreatedBy":         "createdBy",
	"CreatedByFullName": "createdByFullName",
	"CreatedByEmail":    "createdByEmail",
	"CreatedDateTime":   "createdDateTime",
	"UpdatedBy":         "updatedBy",
	"UpdatedByFullName": "updatedByFullName",
	"UpdatedByEmail":    "updatedByEmail",
	"UpdatedDateTime":   "updatedDateTime",
	"DeletedBy":         "deletedBy",
	"DeletedByFullName": "deletedByFullName",
	"DeletedByEmail":    "deletedByEmail",
	"DeletedDateTime":   "deletedDateTime",
}

var ecxFabricConnectionChangeLogSchemaDescription = map[string]string{
	"CreatedBy":         "Created by User Key",
	"CreatedByFullName": "Created by User Full Name",
	"CreatedByEmail":    "Created by User Email Address",
	"CreatedDateTime":   "Created by Date and Time",
	"UpdatedBy":         "Updated by User Key",
	"UpdatedByFullName": "Updated by User Full Name",
	"UpdatedByEmail":    "Updated by User Email Address",
	"UpdatedDateTime":   "Updated by Date and Time",
	"DeletedBy":         "Deleted by User Key",
	"DeletedByFullName": "Deleted by User Full Name",
	"DeletedByEmail":    "Deleted by User Email Address",
	"DeletedDateTime":   "Deleted by Date and Time",
}

var ecxFabricConnectionRedundancySchemaNames = map[string]string{
	"Group":    "group",
	"Priority": "priority",
}

var ecxFabricConnectionRedundancyDescription = map[string]string{
	"Group":    "Redundancy group identifier",
	"Priority": "Priority type",
}

var ecxFabricConnectionConnectionSideSchemaNames = map[string]string{
	"Invitation":     "invitation",
	"ServiceToken":   "serviceToken",
	"AccessPoint":    "accessPoint",
	"CompanyProfile": "companyProfile",
	"Nat":            "nat",
	"AdditionalInfo": "additionalInfo",
}

var ecxFabricConnectionConnectionSideDescription = map[string]string{
	"Invitation":     "Invitation based on connection request",
	"ServiceToken":   "For service token based connections, Service tokens authorize users to access protected resources and services. Resource owners can distribute the tokens to trusted partners and vendors, allowing selected third parties to work directly with Equinix network assets",
	"AccessPoint":    "Point of access details",
	"CompanyProfile": "Company Profile",
	"Nat":            "Network Address Translation type",
	"AdditionalInfo": "side additional details",
}

var ecxFabricConnectionInvitationSchemaNames = map[string]string{
	"Type_":   "type",
	"Uuid":    "uuid",
	"State":   "state",
	"Message": "message",
	"Email":   "email",
	"Expiry":  "expiry",
}

var ecxFabricConnectionInvitationDescription = map[string]string{
	"Type_":   "Invitation type",
	"Uuid":    "Equinix-assigned invitation identifier",
	"State":   "Invitation status as it is today",
	"Message": "Invitation message",
	"Email":   "Invitation recipient",
	"Expiry":  "Invitation expiry time",
}

var ecxFabricConnectionServiceTokenSchemaNames = map[string]string{
	"Type_":              "type",
	"Href":               "href",
	"Uuid":               "uuid",
	"Description":        "description",
	"Expiry":             "expiry",
	"ExpirationDateTime": "expirationDateTime",
	"Connection":         "connection",
	"State":              "state",
	"Notifications":      "notifications",
	"Account":            "account",
	"Changelog":          "changelog",
}

var ecxFabricConnectionServiceTokenDescription = map[string]string{
	"Type_":              "token type",
	"Href":               "An absolute URL that is the subject of the link's context",
	"Uuid":               "Equinix-assigned service token identifier",
	"Description":        "Service token description",
	"Expiry":             "Lifespan (in days) of the service token",
	"ExpirationDateTime": "Expiration date and time of the service token",
	"Connection":         "connection",
	"State":              "state",
	"Notifications":      "Service token related notifications",
	"Account":            "account",
	"Changelog":          "changelog",
}

var ecxFabricConnectionServiceTokenConnectionSchemaNames = map[string]string{
	"Type_":                 "type",
	"AllowRemoteConnection": "allowRemoteConnection",
	"BandwidthLimit":        "bandwidthLimit",
	"SupportedBandwidths":   "supportedBandwidths",
	"ASide":                 "aSide",
	"ZSide":                 "zSide",
}

var ecxFabricConnectionServiceTokenConnectionDescription = map[string]string{
	"Type_":                 "Type of Connection",
	"AllowRemoteConnection": "Authorization to connect remotely",
	"BandwidthLimit":        "Connection bandwidth limit in Mbps",
	"SupportedBandwidths":   "List of permitted bandwidths.",
	"ASide":                 "aSide",
	"ZSide":                 "zSide",
}

var ecxFabricConnectionServiceTokenSideSchemaNames = map[string]string{
	"AccessPointSelectors": "accessPointSelectors",
}

var ecxFabricConnectionServiceTokenSideDescription = map[string]string{
	"AccessPointSelectors": "Access Point Selectors",
}

var ecxFabricConnectionServiceTokenAccessPointSelectorSchemaNames = map[string]string{
	"Type_":        "type",
	"Port":         "port",
	"LinkProtocol": "linkProtocol",
}

var ecxFabricConnectionServiceTokenAccessPointSelectorDescription = map[string]string{
	"Type_":        " Selector Type",
	"Port":         "port",
	"LinkProtocol": "linkProtocol",
}

var ecxFabricConnectionServiceTokenAccessPointSelectorLinkProtocolSchemaNames = map[string]string{
	"Type_": "type",
}

var ecxFabricConnectionServiceTokenAccessPointSelectorLinkProtocolDescription = map[string]string{
	"Type_": "type",
}

var ecxFabricConnectionServiceTokenAccessPointSelectorSimplifiedMetadataEntitySchemaNames = map[string]string{
	"Href":  "href",
	"Uuid":  "uuid",
	"Type_": "type",
}

var ecxFabricConnectionServiceTokenAccessPointSelectorSimplifiedMetadataEntityDescription = map[string]string{
	"Href":  "url to entity",
	"Uuid":  "Equinix assigned Identifier",
	"Type_": "Type of Port",
}

var ecxFabricConnectionAccessPointSchemaNames = map[string]string{
	"Type_":                "type",
	"Account":              "account",
	"Location":             "location",
	"Port":                 "port",
	"Profile":              "profile",
	"Gateway":              "gateway",
	"LinkProtocol":         "linkProtocol",
	"VirtualDevice":        "virtualDevice",
	"Interface_":           "interface",
	"SellerRegion":         "sellerRegion",
	"PeeringType":          "peeringType",
	"AuthenticationKey":    "authenticationKey",
	"RoutingProtocols":     "routingProtocols",
	"AdditionalInfo":       "additionalInfo",
	"ProviderConnectionId": "providerConnectionId",
}

var ecxFabricConnectionAccessPointDescription = map[string]string{
	"Type_":                "Access point type",
	"Account":              "account",
	"Location":             "Access point location",
	"Port":                 "Port access point information",
	"Profile":              "Service Profile ",
	"Gateway":              "Gateway access point information",
	"LinkProtocol":         "Connection link protocol",
	"VirtualDevice":        "Virtual device",
	"Interface_":           "Virtual device interface",
	"SellerRegion":         "Access point seller region",
	"PeeringType":          "Peering Type",
	"AuthenticationKey":    "Access point authentication key",
	"RoutingProtocols":     "Access point routing protocols configuration",
	"AdditionalInfo":       "Access point additional Information",
	"ProviderConnectionId": "Provider assigned Connection Id",
}

var ecxFabricConnectionSideAccessPointLocationSchemaNames = map[string]string{
	"Href":      "href",
	"Region":    "region",
	"MetroName": "metroName",
	"MetroCode": "metroCode",
	"Ibx":       "ibx",
}

var ecxFabricConnectionSideAccessPointLocationDescription = map[string]string{
	"Href":      "href",
	"Region":    "Access point region",
	"MetroName": "Access point metro name",
	"MetroCode": "Access point metro code",
	"Ibx":       "ibx",
}

var ecxFabricConnectionSideAccessPointSimplifiedLinkProtocolSchemaNames = map[string]string{
	"Type_":    "type",
	"VlanTag":  "vlanTag",
	"VlanSTag": "vlanSTag",
	"VlanCTag": "vlanCTag",
	"Unit":     "unit",
	"Vni":      "vni",
	"IntUnit":  "intUnit",
}

var ecxFabricConnectionSideAccessPointSimplifiedLinkProtocolDescription = map[string]string{
	"Type_":    "type",
	"VlanTag":  "vlanTag value specified for DOT1Q connections",
	"VlanSTag": "vlanSTag value specified for QINQ connections",
	"VlanCTag": "vlanCTag value specified for QINQ connections",
	"Unit":     "unit",
	"Vni":      "vni",
	"IntUnit":  "intUnit",
}

var ecxFabricConnectionSideAccessPointVirtualDeviceSchemaNames = map[string]string{
	"Href":  "href",
	"Uuid":  "uuid",
	"Name":  "name",
	"Type_": "type",
}

var ecxFabricConnectionSideAccessPointVirtualDeviceDescription = map[string]string{
	"Href":  "Virtual Device URI",
	"Uuid":  "Equinix-assigned Virtual Device identifier",
	"Name":  "Customer-assigned Virtual Device name",
	"Type_": "Virtual Device type",
}

var ecxFabricConnectionSideAccessPointVirtualDeviceInterfaceSchemaNames = map[string]string{
	"Uuid":  "uuid",
	"Id":    "id",
	"Type_": "type",
}

var ecxFabricConnectionSideAccessPointVirtualDeviceInterfaceDescription = map[string]string{
	"Uuid":  "Equinix-assigned Virtual Device Interface identifier",
	"Id":    "Interface id",
	"Type_": "Virtual Device Interface type",
}

var ecxFabricConnectionServiceProfileSchemaNames = map[string]string{
	"Href":                   "href",
	"Type_":                  "type",
	"Name":                   "name",
	"Uuid":                   "uuid",
	"Description":            "description",
	"Notifications":          "notifications",
	"Tags":                   "tags",
	"Visibility":             "visibility",
	"AllowedEmails":          "allowedEmails",
	"AccessPointTypeConfigs": "accessPointTypeConfigs",
	"CustomFields":           "customFields",
	"MarketingInfo":          "marketingInfo",
	"Ports":                  "ports",
	"VirtualDevices":         "virtualDevices",
	"Metros":                 "metros",
	"SelfProfile":            "selfProfile",
}

var ecxFabricConnectionServiceProfileDescription = map[string]string{
	"Href":                   "Service Profile URI response attribute",
	"Type_":                  "Service profile type",
	"Name":                   "Customer-assigned service profile name",
	"Uuid":                   "Equinix assigned service profile identifier",
	"Description":            "User-provided service description",
	"Notifications":          "Recipients of notifications on service profile change",
	"Tags":                   "Tags",
	"Visibility":             "Visibility of the service profile",
	"AllowedEmails":          "User Emails that are allowed to access this service profile",
	"AccessPointTypeConfigs": "Access Point Type ",
	"CustomFields":           "Custom Fields",
	"MarketingInfo":          "Marketing Info",
	"Ports":                  "Ports",
	"VirtualDevices":         "Virtual Devices",
	"Metros":                 "Derived response attribute",
	"SelfProfile":            "response attribute indicates whether the profile belongs to the same organization as the api-invoker",
}

var ecxFabricConnectionServiceProfileAccessTypeSchemaNames = map[string]string{
	"Type_": "type",
	"Uuid":  "uuid",
}

var ecxFabricConnectionServiceProfileAccessTypeDescription = map[string]string{
	"Type_": "type",
	"Uuid":  "uuid",
}

var ecxFabricConnectionServiceProfileCustomFieldSchemaNames = map[string]string{
	"Label":          "label",
	"Description":    "description",
	"Required":       "required",
	"DataType":       "dataType",
	"Options":        "options",
	"CaptureInEmail": "captureInEmail",
}

var ecxFabricConnectionServiceProfileCustomFieldDescription = map[string]string{
	"Label":          "label",
	"Description":    "description",
	"Required":       "required",
	"DataType":       "dataType",
	"Options":        "options",
	"CaptureInEmail": "capture this field as a part of email notification",
}

var ecxFabricConnectionServiceProfileMarketingInfoSchemaNames = map[string]string{
	"Logo":         "logo",
	"Promotion":    "promotion",
	"ProcessSteps": "processSteps",
}

var ecxFabricConnectionServiceProfileMarketingInfoDescription = map[string]string{
	"Logo":         "Logo file name",
	"Promotion":    "Profile promotion on marketplace",
	"ProcessSteps": "processSteps",
}

var ecxFabricConnectionServiceProfileMarketingProcessStepsSchemaNames = map[string]string{
	"Title":       "title",
	"SubTitle":    "subTitle",
	"Description": "description",
}

var ecxFabricConnectionServiceProfileMarketingProcessStepsDescription = map[string]string{
	"Title":       "Service profile custom step title",
	"SubTitle":    "Service profile custom step sub title",
	"Description": "Service profile custom step description",
}

var ecxFabricConnectionServiceProfileAccessPointPortsColoSchemaNames = map[string]string{
	"SellerRegion":            "sellerRegion",
	"SellerRegionDescription": "sellerRegionDescription",
	"CrossConnectId":          "crossConnectId",
	"Type_":                   "type",
	"Uuid":                    "uuid",
	"Location":                "location",
}

var ecxFabricConnectionServiceProfileAccessPointPortsColoDescription = map[string]string{
	"SellerRegion":            "Seller Region",
	"SellerRegionDescription": "Seller Region Description",
	"CrossConnectId":          " Cross ConnectId",
	"Type_":                   " Type",
	"Uuid":                    "uuid",
	"Location":                "Location",
}

var ecxFabricConnectionServiceProfileAccessPointPortsVDSchemaNames = map[string]string{
	"Type_":         "type",
	"Uuid":          "uuid",
	"Location":      "location",
	"InterfaceUuid": "interfaceUuid",
}

var ecxFabricConnectionServiceProfileAccessPointPortsVDDescription = map[string]string{
	"Type_":         "type",
	"Uuid":          "uuid",
	"Location":      "location",
	"InterfaceUuid": "interfaceUuid",
}

var ecxFabricConnectionServiceProfileServiceMetrosSchemaNames = map[string]string{
	"Code":          "code",
	"Name":          "name",
	"Ibxs":          "ibxs",
	"InTrail":       "inTrail",
	"DisplayName":   "displayName",
	"SellerRegions": "sellerRegions",
}

var ecxFabricConnectionServiceProfileServiceMetrosDescription = map[string]string{
	"Code":          "Metro code",
	"Name":          "Metro name",
	"Ibxs":          "Allowed ibxes in the metro",
	"InTrail":       "inTrail",
	"DisplayName":   "Service metro display name",
	"SellerRegions": "Seller Regions",
}

func resourceFabricConnection() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		//Read:   resourceFabricConnectionRead,
		CreateContext: resourceFabricConnectionCreate,
		//Update: resourceFabricConnectionUpdate,
		//Delete: resourceFabricConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the connection, will be part of created response",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Customer-provided connection description",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Defines the connection type like VG_VC, EVPL_VC, EPL_VC, EC_VC, GW_VC, ACCESS_EPL_VC",
			},
			"notifications": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: stringIsEmailAddress(),
				},
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "sate",
			},
			"redundancy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "sate",
			},
			/*"a_side": {
				Type:        schema.,
				Required:    true,
				Description: "sate",
			},
			"z_side": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "sate",
			},*/
		},
	}
}

func resourceFabricConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabric

	conType := v4.ConnectionType(d.Get("type").(string))

	notifications := []v4.SimplifiedNotification{{
		Type_:  "Email",
		Emails: expandSetToStringList(d.Get("notifications").(*schema.Set)),
	}}

	priority := v4.ConnectionPriority(d.Get("redundancy").(string))

	red := v4.ConnectionRedundancy{
		Priority: &priority,
	}

	createRequest := v4.ConnectionPostRequest{
		Name:          d.Get("name").(string),
		Type_:         &conType,
		Notifications: notifications,
		Bandwidth:     int32(d.Get("bandwidth").(int)),
		Redundancy:    &red,
	}

	conn, _, err := client.ConnectionsApi.CreateConnection(ctx, createRequest)

	fmt.Println(err)
	if err != nil {
		return diag.FromErr(err)
	}
	uuid := conn.Uuid
	d.SetId(uuid)

	return resourceFabricConnectionRead(ctx, d, meta)
}

func resourceFabricConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabric

	conn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, d.Id(), nil)

	if err != nil {
		log.Printf("[WARN] Connection %s not found ", d.Id())
		d.SetId("")
		return nil
	}

	m := map[string]interface{}{
		"name":        conn.Name,
		"description": conn.Description,
	}

	return setFabricMap(d, m)
}

func setFabricMap(d *schema.ResourceData, m map[string]interface{}) diag.Diagnostics {
	errs := &multierror.Error{}
	var diags diag.Diagnostics
	for key, v := range m {
		var err error
		if f, ok := v.(setFn); ok {
			err = f(d, key)
		} else {
			err = d.Set(key, v)
		}

		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	sort.Sort(errs)

	return diags
}

/*func resourceFabricConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metal

	sPtr := func(s string) *string { return &s }
	iPtr := func(i int) *int { return &i }

	updateRequest := &packngo.VRFUpdateRequest{}
	if d.HasChange("name") {
		updateRequest.Name = sPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateRequest.Description = sPtr(d.Get("description").(string))
	}
	if d.HasChange("local_asn") {
		updateRequest.LocalASN = iPtr(d.Get("local_asn").(int))
	}
	if d.HasChange("ip_ranges") {
		ipRanges := expandSetToStringList(d.Get("ip_ranges").(*schema.Set))
		updateRequest.IPRanges = &ipRanges
	}

	_, _, err := client.VRFs.Update(d.Id(), updateRequest)
	if err != nil {
		return friendlyError(err)
	}

	return resourceMetalVRFRead(d, meta)
}
func resourceFabricConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metal

	resp, err := client.VRFs.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) == nil {
		d.SetId("")
	}

	return friendlyError(err)
}*/
