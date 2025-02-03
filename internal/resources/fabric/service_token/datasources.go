package service_token

import (
	"context"
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema:      dataSourceBaseSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch service token for a given UUID

Additional documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/Fabric/service%20tokens/Fabric-Service-Tokens.htm
* API: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/ConnectUsingServiceToken.htm`,
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceRead(ctx, d, meta)
}

func DataSourceSearch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSearch,
		Schema:      dataSourceSearchSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch service token for a given search data set

Additional documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/Fabric/service%20tokens/Fabric-Service-Tokens.htm
* API: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/ConnectUsingServiceToken.htm`,
	}
}

func dataSourceSearch(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	searchRequest := buildSearchRequest(d)

	serviceTokens, _, err := client.ServiceTokensApi.SearchServiceTokens(ctx).ServiceTokenSearchRequest(searchRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	if len(serviceTokens.Data) < 1 {
		return diag.FromErr(fmt.Errorf("no records are found for the route filter search criteria provided - %d , please change the search criteria", len(serviceTokens.Data)))
	}

	d.SetId(serviceTokens.Data[0].GetUuid())
	return setServiceTokensData(d, serviceTokens)
}
