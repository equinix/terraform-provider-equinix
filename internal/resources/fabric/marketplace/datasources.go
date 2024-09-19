package marketplace

import (
	"context"
	"log"
	"strings"

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

func dataSourceFabricMarketplaceSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return fabricMarketplaceSubscriptionRead(ctx, d, meta)
}

func fabricMarketplaceSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	subscription, _, err := client.MarketplaceSubscriptionsApi.GetSubscriptionById(ctx, d.Id()).Execute()
	if err != nil {
		log.Printf("[WARN] Subscription %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(subscription.GetUuid())
	return setFabricMap(d, subscription)
}
