package sweep_test

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/sweep/services"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	// Register sweepers for individual resource packages
	services.AddFabricTestSweepers()
	services.AddMetalTestSweepers()

	resource.TestMain(m)
}
