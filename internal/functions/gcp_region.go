package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewGCPRegionFunction returns a Terraform function that validates GCP region codes.
func NewGCPRegionFunction() function.Function {
	return newStringValidationFunction(
		"gcp_region",
		"Validate that a string is a valid GCP region.",
		"Returns true when the input value is a valid GCP region (e.g., us-central1, europe-west1).",
		validators.GCPRegion(),
	)
}
