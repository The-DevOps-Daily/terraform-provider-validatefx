package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewAWSRegionFunction returns a Terraform function that validates AWS region codes.
func NewAWSRegionFunction() function.Function {
	return newStringValidationFunction(
		"aws_region",
		"Validate that a string is a valid AWS region code.",
		"Returns true when the input value is a valid AWS region code (e.g., us-east-1, eu-west-1).",
		validators.AWSRegion(),
	)
}
