provider {
  name    = "validatefx"
  dir     = "."
  example = "examples/functions/assert/main.tf"
}

rendered_provider_name = "validatefx"
website_dir = "docs"
templates_dir = "templates"
examples_dir  = "examples"

schema_paths = ["internal/provider"]

functions = [
  { name = "assert",       example = "examples/functions/assert/main.tf" }
  { name = "email",        example = "examples/functions/email/main.tf" }
  { name = "uuid",         example = "examples/functions/uuid/main.tf" }
  { name = "base64",       example = "examples/functions/base64/main.tf" }
  { name = "credit_card",  example = "examples/functions/credit_card/main.tf" }
  { name = "domain",       example = "examples/functions/domain/main.tf" }
]
