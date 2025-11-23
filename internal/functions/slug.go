package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewSlugFunction exposes the slug validator as a Terraform function.
func NewSlugFunction() function.Function {
	return newStringValidationFunction(
		"slug",
		"Validate that a string is a valid slug.",
		"Returns true when the input string is a valid slug (lowercase letters, digits, and hyphens; no leading/trailing or consecutive hyphens).",
		validators.Slug(),
	)
}
