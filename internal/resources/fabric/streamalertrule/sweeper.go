package streamalertrule

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

// AddTestSweeper registers the sweeper for Fabric Stream Alert Rules
func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_stream_alert_rule", &resource.Sweeper{
		Name:         "equinix_fabric_stream_alert_rule",
		Dependencies: []string{},
		F:            testSweepStreamAlertRules,
	})
}

func testSweepStreamAlertRules(_ string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping Fabric Stream Alert Rules")
	ctx := context.Background()
	meta, err := sweep.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Stream Alert Rules: %s", err)
	}
	configLoadErr := meta.Load(ctx)
	if configLoadErr != nil {
		return fmt.Errorf("error loading configuration for sweeping Stream Alert Rules: %s", configLoadErr)
	}
	fabric := meta.NewFabricClientForTesting(ctx)
	limit := int32(100)

	streams, _, err := fabric.StreamsApi.GetStreams(ctx).Limit(limit).Execute()
	if err != nil {
		return fmt.Errorf("error getting streams list for sweeping fabric streams: %s", err)
	}

	for _, stream := range streams.GetData() {
		if sweep.IsSweepableFabricTestResource(stream.GetName()) {
			alertRules, _, err := fabric.StreamAlertRulesApi.GetStreamAlertRules(ctx, stream.GetUuid()).Limit(limit).Execute()
			if err != nil {
				errs = append(errs, fmt.Errorf("error getting fabric stream subscriptions on stream %s: %s", stream.GetUuid(), err))
			}
			for _, alertRule := range alertRules.GetData() {
				log.Printf("[DEBUG] alert rule : %s", alertRule.GetName())
				if sweep.IsSweepableFabricTestResource(alertRule.GetName()) {
					log.Printf("[DEBUG] Deleting stream alert rule: %s", alertRule.GetName())
					_, resp, err := fabric.StreamAlertRulesApi.DeleteStreamAlertRuleByUuid(ctx, alertRule.GetUuid(), stream.GetUuid()).Execute()
					if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err) != nil {
						errs = append(errs, fmt.Errorf("error deleting fabric stream alert rule %s on stream %s: %s", alertRule.GetUuid(), stream.GetUuid(), err))
						continue
					}
				}
			}
			log.Printf("[DEBUG] Deleting stream: %s", stream.GetName())
			_, resp, err := fabric.StreamsApi.DeleteStreamByUuid(ctx, stream.GetUuid()).Execute()
			if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err) != nil {
				errs = append(errs, fmt.Errorf("error deleting fabric stream: %s", err))
			}
		}
	}
	return errors.Join(errs...)
}
