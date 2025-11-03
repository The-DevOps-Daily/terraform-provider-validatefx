terraform {
  required_providers {
    validatefx = {
      source  = "The-DevOps-Daily/validatefx"
      version = ">= 0.1.0"
    }
  }
}

provider "validatefx" {}

locals {
  samples = [
    {
      label = "within"
      value = "7.5"
    },
    {
      label = "too-low"
      value = "2"
    },
    {
      label = "too-high"
      value = "11"
    },
    {
      label = "not-number"
      value = "abc"
    },
  ]

  results = [
    for sample in local.samples : {
      label = sample.label
      value = sample.value
      valid = provider::validatefx::between(sample.value, "5", "10")
    }
  ]
}

output "between_validation" {
  value = local.results
}
