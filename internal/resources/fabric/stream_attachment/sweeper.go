package stream

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/sweep"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_stream", &resource.Sweeper{
		Name:         "equinix_fabric_stream",
		Dependencies: []string{},
		F:            testSweepStreams,
	})
}

func testSweepStreams(_ string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping Fabric Streams")
	ctx := context.Background()
	meta, err := sweep.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Streams: %s", err)
	}
	configLoadErr := meta.Load(ctx)
	if configLoadErr != nil {
		return fmt.Errorf("error loading configuration for sweeping Streams: %s", err)
	}
	fabric := meta.NewFabricClientForTesting(ctx)
	limit := int32(100)

	streams, _, err := fabric.StreamsApi.GetStreams(ctx).Limit(limit).Execute()
	if err != nil {
		return fmt.Errorf("error getting streams list for sweeping fabric streams: %s", err)
	}

	for _, stream := range streams.GetData() {
		if sweep.IsSweepableFabricTestResource(stream.GetName()) {
			log.Printf("[DEBUG] Deleting stream: %s", stream.GetName())
			_, resp, err := fabric.StreamsApi.DeleteStreamByUuid(ctx, stream.GetUuid()).Execute()
			if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err) != nil {
				errs = append(errs, fmt.Errorf("error deleting fabric stream: %s", err))
			}
		}
	}

	return errors.Join(errs...)
}
