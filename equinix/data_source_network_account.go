package equinix

import (
	"context"
	"fmt"
	"strings"

	"github.com/equinix/ne-go"
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

func dataSourceNetworkAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkAccountRead,
		Schema: map[string]*schema.Schema{
			networkAccountSchemaNames["Name"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			networkAccountSchemaNames["Number"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			networkAccountSchemaNames["Status"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Processing", "Submitted", "Staged"}, true),
			},
			networkAccountSchemaNames["UCMID"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			networkAccountSchemaNames["MetroCode"]: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: stringIsMetroCode(),
			},
		},
	}
}

func dataSourceNetworkAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*Config)
	var diags diag.Diagnostics
	metro := d.Get(networkAccountSchemaNames["MetroCode"]).(string)
	name := d.Get(networkAccountSchemaNames["Name"]).(string)
	status := d.Get(networkAccountSchemaNames["Status"]).(string)
	accounts, err := conf.ne.GetAccounts(metro)
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
