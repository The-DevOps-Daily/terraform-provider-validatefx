package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// GCPRegion validates that a string is a valid GCP region.
func GCPRegion() validator.String { return gcpRegionValidator{} }

type gcpRegionValidator struct{}

var _ validator.String = (*gcpRegionValidator)(nil)

// Valid GCP regions as of 2024
var validGCPRegions = map[string]bool{
	// Americas
	"us-central1":             true,
	"us-east1":                true,
	"us-east4":                true,
	"us-east5":                true,
	"us-south1":               true,
	"us-west1":                true,
	"us-west2":                true,
	"us-west3":                true,
	"us-west4":                true,
	"northamerica-northeast1": true,
	"northamerica-northeast2": true,
	"southamerica-east1":      true,
	"southamerica-west1":      true,
	// Europe
	"europe-central2":   true,
	"europe-north1":     true,
	"europe-southwest1": true,
	"europe-west1":      true,
	"europe-west2":      true,
	"europe-west3":      true,
	"europe-west4":      true,
	"europe-west6":      true,
	"europe-west8":      true,
	"europe-west9":      true,
	"europe-west10":     true,
	"europe-west12":     true,
	// Asia Pacific
	"asia-east1":      true,
	"asia-east2":      true,
	"asia-northeast1": true,
	"asia-northeast2": true,
	"asia-northeast3": true,
	"asia-south1":     true,
	"asia-south2":     true,
	"asia-southeast1": true,
	"asia-southeast2": true,
	// Australia
	"australia-southeast1": true,
	"australia-southeast2": true,
	// Middle East
	"me-central1": true,
	"me-central2": true,
	"me-west1":    true,
	// Africa
	"africa-south1": true,
}

func (gcpRegionValidator) Description(_ context.Context) string {
	return "value must be a valid GCP region"
}

func (v gcpRegionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (gcpRegionValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if diag := validateStringInMap(value, validGCPRegions, req.Path, "Invalid GCP Region", "GCP region"); diag != nil {
		resp.Diagnostics.Append(diag)
	}
}
