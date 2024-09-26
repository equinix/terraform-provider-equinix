package equinix

import (
	"fmt"
	"strings"
	"time"
)

// list of plans and metros and os used as filter criteria to find available hardware to run tests
var (
	preferable_plans  = []string{"x1.small.x86", "t1.small.x86", "c2.medium.x86", "c3.small.x86", "c3.medium.x86", "m3.small.x86"}
	preferable_metros = []string{"ch", "ny", "sv", "ty", "am"}
	preferable_os     = []string{"ubuntu_20_04"}
)

// Deprecated: use the identical TestDeviceTerminationTime from internal/acceptance instead
func testDeviceTerminationTime() string {
	return time.Now().UTC().Add(60 * time.Minute).Format(time.RFC3339)
}

// This function should be used to find available plans in all test where a metal_device resource is needed.
//
// TODO consider adding a datasource for equinix_metal_operating_system and making the local.os conditional
//
//	https://github.com/equinix/terraform-provider-equinix/pull/220#discussion_r915418418equinix_metal_operating_system
//	https://github.com/equinix/terraform-provider-equinix/discussions/221
func confAccMetalDevice_base(plans, metros, os []string) string {
	return fmt.Sprintf(`
data "equinix_metal_plans" "test" {
    sort {
        attribute = "id"
        direction = "asc"
    }

    filter {
        attribute = "name"
        values    = [%s]
    }
    filter {
        attribute = "available_in_metros"
        values    = [%s]
    }
    filter {
        attribute = "deployment_types"
        values    = ["on_demand", "spot_market"]
    }
}

// Select a metal plan randomly and lock it in
// so that we don't pick a different one for
// every subsequent terraform plan
resource "random_integer" "plan_idx" {
  min = 0
  max = length(data.equinix_metal_plans.test.plans) - 1
}

resource "terraform_data" "plan" {
  input = data.equinix_metal_plans.test.plans[random_integer.plan_idx.result]

  lifecycle {
	ignore_changes = ["input"]
  }
}

resource "terraform_data" "facilities" {
  input = sort(tolist(setsubtract(terraform_data.plan.output.available_in, ["nrt1", "dfw2", "ewr1", "ams1", "sjc1", "ld7", "sy4", "ny6"])))

  lifecycle {
    ignore_changes = ["input"]
  }
}

// Select a metal facility randomly and lock it in
// so that we don't pick a different one for
// every subsequent terraform plan
resource "random_integer" "facility_idx" {
  min = 0
  max = length(local.facilities) - 1
}

resource "terraform_data" "facility" {
  input = local.facilities[random_integer.facility_idx.result]

  lifecycle {
	ignore_changes = ["input"]
  }
}

// Select a metal metro randomly and lock it in
// so that we don't pick a different one for
// every subsequent terraform plan
resource "random_integer" "metro_idx" {
  min = 0
  max = length(local.metros) - 1
}

resource "terraform_data" "metro" {
  input = local.metros[random_integer.metro_idx.result]

  lifecycle {
	ignore_changes = ["input"]
  }
}

locals {
    // Select a random plan
    plan              = terraform_data.plan.output.slug

    // Select a random facility from the facilities in which the selected plan is available, excluding decommed facilities
    facilities             = terraform_data.facilities.output
    facility               = terraform_data.facility.output

    // Select a random metro from the metros in which the selected plan is available
    metros             = sort(tolist(terraform_data.plan.output.available_in_metros))
    metro              = terraform_data.metro.output

    os = [%s][0]
}
`, fmt.Sprintf("\"%s\"", strings.Join(plans[:], `","`)), fmt.Sprintf("\"%s\"", strings.Join(metros[:], `","`)), fmt.Sprintf("\"%s\"", strings.Join(os[:], `","`)))
}
