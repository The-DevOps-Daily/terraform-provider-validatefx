terraform {
  required_providers {
    validatefx = {
      source  = "The-DevOps-Daily/validatefx"
      version = "0.0.1"
    }
  }
}

provider "validatefx" {}

locals {
  filenames = [
    "app.log",
    "app.txt",
    "app.csv",
  ]

  results = [
    for name in local.filenames : {
      value = name
      valid = provider::validatefx::has_suffix(name, [".log", ".txt"])
    }
  ]
}

output "has_suffix_example" {
  value = local.results
}
