package network

import (
	"context"
	"errors"
	"fmt"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"net/http"

	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"log"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_network", &resource.Sweeper{
		Name:         "equinix_fabric_connection",
		Dependencies: []string{},
		F:            testSweepNetworks,
	})
}

func testSweepNetworks(region string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping Fabric Networks")
	ctx := context.Background()
	meta, err := sweep.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Networks: %s", err)
	}
	err = meta.Load(ctx)
	if err != nil {
		log.Printf("Error loading meta: %v", err)
		return err
	}
	fabric := meta.NewFabricClientForTesting()

	name := fabricv4.NETWORKSEARCHFIELDNAME_NAME
	likeOperator := fabricv4.NETWORKFILTEROPERATOR_LIKE
	limit := int32(100)
	offset := int32(10)

	networkSearchRequest := fabricv4.NetworkSearchRequest{
		Filter: &fabricv4.NetworkFilter{
			And: []fabricv4.NetworkFilter{
				{
					Property: &name,
					Operator: &likeOperator,
					Values:   []string{"%_PFCR"},
				},
				{
					Property: &name,
					Operator: &likeOperator,
					Values:   []string{"%_PFNV"},
				},
			},
		},
		Pagination: &fabricv4.PaginationRequest{
			Offset: &offset,
			Limit:  &limit,
		},
	}
	networks, _, err := fabric.NetworksApi.SearchNetworks(ctx).NetworkSearchRequest(networkSearchRequest).Execute()
	if err != nil {
		return fmt.Errorf("error searching networks: %s", err)
	}

	for _, network := range networks.Data {
		if sweep.IsSweepableFabricTestResource(network.GetName()) {
			log.Printf("[DEBUG] Deleting Networks: %s", network.GetName())
			_, resp, err := fabric.NetworksApi.DeleteNetworkByUuid(ctx, network.GetUuid()).Execute()
			if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err) != nil {
				errs = append(errs, fmt.Errorf("error deleting network: %s", err))
			}
		}
	}
	return errors.Join(errs...)
}
