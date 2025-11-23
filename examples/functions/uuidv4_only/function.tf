terraform {
  required_providers {
    validatefx = {
      source = "the-devops-daily/validatefx"
    }
  }
}

locals {
  # Valid UUIDv4
  valid_v4_lowercase = provider::validatefx::uuidv4_only("550e8400-e29b-41d4-a716-446655440000")
  valid_v4_uppercase = provider::validatefx::uuidv4_only("550E8400-E29B-41D4-A716-446655440000")
  valid_v4_mixed     = provider::validatefx::uuidv4_only("f47ac10b-58cc-4372-a567-0e02b2c3d479")
  valid_v4_random    = provider::validatefx::uuidv4_only("123e4567-e89b-42d3-a456-426614174000")
}

output "uuidv4_only_checks" {
  value = {
    lowercase = local.valid_v4_lowercase
    uppercase = local.valid_v4_uppercase
    mixed     = local.valid_v4_mixed
    random    = local.valid_v4_random
  }
}
