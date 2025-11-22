package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewAzureLocationFunction returns a Terraform function that validates Azure location codes.
func NewAzureLocationFunction() function.Function {
	return newStringValidationFunction(
		"azure_location",
		"Validate that a string is a valid Azure location.",
		"Returns true when the input value is a valid Azure location (e.g., eastus, westeurope).",
		validators.AzureLocation(),
	)
}
