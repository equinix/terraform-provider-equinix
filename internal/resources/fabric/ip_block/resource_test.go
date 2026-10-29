package ipblock_test

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const (
	prefixConflictErrorCode = "EQ-7100323"
	createPrefixAttempts    = 5
)

// attemptT lets a single resource.Test run fail without failing the parent test,
// so the caller can retry with a different prefix.
type attemptT struct {
	*testing.T
	mu     sync.Mutex
	failed bool
	msg    string
}

func (a *attemptT) record(msg string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.failed = true
	a.msg += msg + "\n"
}

func (a *attemptT) Error(args ...any)                 { a.record(fmt.Sprint(args...)) }
func (a *attemptT) Errorf(format string, args ...any) { a.record(fmt.Sprintf(format, args...)) }
func (a *attemptT) Fatal(args ...any)                 { a.record(fmt.Sprint(args...)); runtime.Goexit() }
func (a *attemptT) Fatalf(format string, args ...any) {
	a.record(fmt.Sprintf(format, args...))
	runtime.Goexit()
}
func (a *attemptT) Fail()    { a.record("Fail called") }
func (a *attemptT) FailNow() { a.record("FailNow called"); runtime.Goexit() }
func (a *attemptT) Failed() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.failed
}

// run executes fn on its own goroutine so that Fatal/FailNow only end this attempt.
func (a *attemptT) run(fn func(t *attemptT)) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn(a)
	}()
	<-done
}

func TestAccFabricIpBlockCreate_PFCR(t *testing.T) {
	testData := testinghelpers.GetFabricIpBlockTestData(t)
	pfcr, ok := testData["pfcr"]
	if !ok {
		t.Skip("skipping: no 'pfcr' key in TF_ACC_FABRIC_IP_BLOCK_TEST_DATA")
	}
	projectID := pfcr["project_id"]
	if projectID == "" {
		t.Skip("skipping: TF_ACC_FABRIC_IP_BLOCK_TEST_DATA missing project_id")
	}

	// A fixed prefix from the env data is used as-is; otherwise a random prefix from the
	// test pool is tried, up to createPrefixAttempts times if it is already taken.
	fixedPrefix := pfcr["prefix"]
	attempts := createPrefixAttempts
	if fixedPrefix != "" {
		attempts = 1
	}

	t.Parallel()
	var lastMsg string
	for i := 1; i <= attempts; i++ {
		prefix := fixedPrefix
		if prefix == "" {
			prefix = testinghelpers.RandomIpBlockPrefix()
		}

		at := &attemptT{T: t}
		at.run(func(at *attemptT) {
			resource.Test(at, ipBlockCreateTestCase(t, projectID, prefix))
		})
		if !at.Failed() {
			return
		}

		lastMsg = at.msg
		if !strings.Contains(lastMsg, prefixConflictErrorCode) {
			t.Fatalf("attempt %d with prefix %s failed: %s", i, prefix, lastMsg)
		}
		t.Logf("attempt %d/%d: prefix %s already in use, retrying", i, attempts, prefix)
	}
	t.Fatalf("no free prefix found after %d attempts, last error: %s", attempts, lastMsg)
}

func ipBlockCreateTestCase(t *testing.T, projectID, prefix string) resource.TestCase {
	return resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkIpBlockDeleted,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricIpBlockResourceConfig(projectID, prefix),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("equinix_fabric_ip_block.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_ip_block.test", tfjsonpath.New("state"), knownvalue.StringExact("ACTIVE")),
					statecheck.ExpectKnownValue("equinix_fabric_ip_block.test", tfjsonpath.New("href"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_ip_block.test", tfjsonpath.New("type"), knownvalue.StringExact("IPV4_IP_BLOCK")),
					statecheck.ExpectKnownValue("equinix_fabric_ip_block.test", tfjsonpath.New("prefix"), knownvalue.StringExact(prefix)),
					statecheck.ExpectKnownValue("equinix_fabric_ip_block.test", tfjsonpath.New("project"),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"project_id": knownvalue.StringExact(projectID),
						}),
					),
				},
			},
		},
	}
}

func testAccFabricIpBlockResourceConfig(projectID, prefix string) string {
	return fmt.Sprintf(`
resource "equinix_fabric_ip_block" "test" {
  type    = "IPV4_IP_BLOCK"
  prefix  = %q
  project = {
    project_id = %q
  }
}
`, prefix, projectID)
}

func checkIpBlockDeleted(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_ip_block" {
			continue
		}
		ipBlock, _, err := client.IPBlocksApi.GetIpBlock(ctx, rs.Primary.ID).Execute()
		if err == nil && ipBlock.GetState() == fabricv4.IPBLOCKSTATE_ACTIVE {
			return fmt.Errorf("IP block %s still exists and is ACTIVE", rs.Primary.ID)
		}
	}
	return nil
}
