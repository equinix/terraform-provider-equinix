package account

import (
	"context"
	"fmt"
	"strings"

	"github.com/equinix/ne-go"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/validaters"
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
}

var networkAccountDescriptions = map[string]string{
	"Name":      "Account name for filtering",
	"Number":    "Account unique number",
	"Status":    "Account status for filtering. Possible values are Active, Processing, Submitted, Staged",
	"UCMID":     "Account unique identifier",
	"MetroCode": "Account location metro cod",
}

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
				ValidateFunc: validaters.StringIsMetroCode(),
				Description:  networkAccountDescriptions["MetroCode"],
			},
		},
	}
}

func dataSourceNetworkAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*config.Config)
	var diags diag.Diagnostics
	metro := d.Get(networkAccountSchemaNames["MetroCode"]).(string)
	name := d.Get(networkAccountSchemaNames["Name"]).(string)
	status := d.Get(networkAccountSchemaNames["Status"]).(string)
	accounts, err := conf.NEClient.GetAccounts(metro)
	if err != nil {
		return diag.FromErr(err)
	}
	var filtered []ne.Account
	for _, account := range accounts {
		if name != "" && ne.StringValue(account.Name) != name {
			continue
		}
		if status != "" && !strings.EqualFold(ne.StringValue(account.Status), status) {
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

func updateNetworkAccountResource(account ne.Account, metroCode string, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%s-%s", metroCode, ne.StringValue(account.Name)))
	if err := d.Set(networkAccountSchemaNames["Name"], account.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(networkAccountSchemaNames["Number"], account.Number); err != nil {
		return fmt.Errorf("error reading Number: %s", err)
	}
	if err := d.Set(networkAccountSchemaNames["Status"], account.Status); err != nil {
		return fmt.Errorf("error reading Status: %s", err)
	}
	if err := d.Set(networkAccountSchemaNames["UCMID"], account.UCMID); err != nil {
		return fmt.Errorf("error reading UCMID: %s", err)
	}
	if err := d.Set(networkAccountSchemaNames["MetroCode"], metroCode); err != nil {
		return fmt.Errorf("error reading MetroCode: %s", err)
	}
	return nil
}
