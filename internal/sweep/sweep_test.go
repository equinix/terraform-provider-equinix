package sweep_test

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/project"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/ssh_key"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vrf"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	// Register legacy sweepers for resources balled up under equinix/
	equinix.AddMetalTestSweepers()

	// Register sweepers for individual resource packages
	addTestSweepers()

	resource.TestMain(m)
}

func addTestSweepers() {
	ssh_key.AddTestSweeper()
	project.AddTestSweeper()
	vrf.AddTestSweeper()
}
