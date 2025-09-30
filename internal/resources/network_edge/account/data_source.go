// Package account provides the network_account data source
package account

import (
	"context"
	"fmt"
	"strings"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_validation "github.com/equinix/terraform-provider-equinix/internal/validation"

	"github.com/equinix/equinix-sdk-go/services/networkedgev1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var networkAccountSchemaNames = map[string]string{
	"Name":      "name",
	"Number":    "number",
	"Status":    "status",
	"UCMID":     "ucm_id",
	"MetroCode": "metro_code",
	"ProjectID": "project_id",
}

var networkAccountDescriptions = map[string]string{
	"Name":      "Account name for filtering",
	"Number":    "Account unique number",
	"Status":    "Account status for filtering. Possible values are Active, Processing, Submitted, Staged",
	"UCMID":     "Account unique identifier",
	"MetroCode": "Account location metro cod",
	"ProjectID": "The unique identifier of Project Resource to which billing account is scoped to",
}

// DataSource creates a new Terraform data source for retrieving account data
func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkAccountRead,
		Description: "Use this data source to get number and identifier of Equinix Network Edge billing account in a given metro location",
		Schema: map[string]*schema.Schema{
			networkAccountSchemaNames["Name"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  networkAccountDescriptions["Name"],
			},
			networkAccountSchemaNames["Number"]: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: networkAccountDescriptions["Number"],
			},
			networkAccountSchemaNames["Status"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Processing", "Submitted", "Staged"}, true),
				Description:  networkAccountDescriptions["Status"],
			},
			networkAccountSchemaNames["UCMID"]: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: networkAccountDescriptions["UCMID"],
			},
			networkAccountSchemaNames["MetroCode"]: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: equinix_validation.StringIsMetroCode,
				Description:  networkAccountDescriptions["MetroCode"],
			},
			networkAccountSchemaNames["ProjectID"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
				Description:  networkAccountDescriptions["ProjectID"],
			},
		},
	}
}

func dataSourceNetworkAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).NewNetworkEdgeClientForSDK(ctx, d)
	var diags diag.Diagnostics
	metro := d.Get(networkAccountSchemaNames["MetroCode"]).(string)
	name := d.Get(networkAccountSchemaNames["Name"]).(string)
	status := d.Get(networkAccountSchemaNames["Status"]).(string)
	projectID := d.Get(networkAccountSchemaNames["ProjectID"]).(string)
	accounts, _, err := client.SetupApi.GetAccountsWithStatusUsingGET(ctx, metro).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	var filtered []networkedgev1.MetroAccountResponse
	for _, account := range accounts.Accounts {
		if name != "" && account.GetAccountName() != name {
			continue
		}
		// TODO: API spec doesn't declare `projectId` for the MetroAccountResponse schema
		if projectID != "" && account.AdditionalProperties["projectId"] != projectID {
			continue
		}
		if status != "" && !strings.EqualFold(account.GetAccountStatus(), status) {
			continue
		}
		filtered = append(filtered, account)
	}
	if len(filtered) < 1 {
		return diag.Errorf("network account query returned no results, please change your search criteria")
	}
	if len(filtered) > 1 {
		return diag.Errorf("network account query returned more than one result, please try more specific search criteria")
	}
	if err := updateNetworkAccountResource(filtered[0], metro, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateNetworkAccountResource(account networkedgev1.MetroAccountResponse, metroCode string, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%s-%s", metroCode, account.GetAccountName()))
	if err := d.Set(networkAccountSchemaNames["Name"], account.GetAccountName()); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(networkAccountSchemaNames["Number"], account.GetAccountNumber()); err != nil {
		return fmt.Errorf("error reading Number: %s", err)
	}
	if err := d.Set(networkAccountSchemaNames["Status"], account.GetAccountStatus()); err != nil {
		return fmt.Errorf("error reading Status: %s", err)
	}
	if err := d.Set(networkAccountSchemaNames["UCMID"], account.GetAccountUcmId()); err != nil {
		return fmt.Errorf("error reading UCMID: %s", err)
	}
	if err := d.Set(networkAccountSchemaNames["MetroCode"], metroCode); err != nil {
		return fmt.Errorf("error reading MetroCode: %s", err)
	}
	// TODO: API spec doesn't declare `projectId` for the MetroAccountResponse schema
	if err := d.Set(networkAccountSchemaNames["ProjectID"], account.AdditionalProperties["projectId"]); err != nil {
		return fmt.Errorf("error reading ProjectID: %s", err)
	}
	return nil
}
