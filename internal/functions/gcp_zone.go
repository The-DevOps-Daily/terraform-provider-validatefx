package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewGCPZoneFunction returns a Terraform function that validates GCP zone codes.
func NewGCPZoneFunction() function.Function {
	return newStringValidationFunction(
		"gcp_zone",
		"Validate that a string is a valid GCP zone.",
		"Returns true when the input value is a valid GCP zone (e.g., us-central1-a, europe-west1-b).",
		validators.GCPZone(),
	)
}
