# Following example uses a boolean variable that may eventually be set to you false when you update your equinix_metal_plans filter criteria because you need a device plan with a new feature.
variable "ignore_plans_metros_changes" {
  type = bool
  description = "If set to true, it will ignore plans or metros changes"
  default = false
}

data "equinix_metal_plans" "example" {
  // new search criteria
}

resource "equinix_metal_device" "example" {
  // required device arguments

  lifecycle {
    ignore_changes = var.ignore_plans_metros_changes ? [
        plan,
        metro,
    ] : []
  }
}
