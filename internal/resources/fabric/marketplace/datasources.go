package marketplace

import (
	"context"
	"net/http"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceFabricMarketplaceSubscription() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricMarketplaceSubscriptionRead,
		Schema:      fabricMarketplaceSubscriptionDataSourceSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch Marketplace Subscription detail for a given UUID",
	}
}

func dataSourceFabricMarketplaceSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return fabricMarketplaceSubscriptionRead(ctx, d, meta)
}

func fabricMarketplaceSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	subscription, resp, err := client.MarketplaceSubscriptionsApi.GetSubscriptionById(ctx, d.Id()).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return diag.Errorf("Fabric marketplace subscription %s not found", d.Id())
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(subscription.GetUuid())
	return setFabricMap(d, subscription)
}
