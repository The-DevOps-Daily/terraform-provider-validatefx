package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewResourceNameFunction exposes the resource name validator as a Terraform function.
func NewResourceNameFunction() function.Function {
	return newStringValidationFunction(
		"resource_name",
		"Validate that a string is a valid Terraform resource name.",
		"Returns true when the input string matches Terraform resource naming conventions (lowercase letters, digits, underscores, and hyphens; must start with letter or underscore).",
		validators.ResourceName(),
	)
}
