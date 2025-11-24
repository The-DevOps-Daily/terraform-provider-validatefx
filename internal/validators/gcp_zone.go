package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// GCPZone validates that a string is a valid GCP zone.
func GCPZone() validator.String { return gcpZoneValidator{} }

type gcpZoneValidator struct{}

var _ validator.String = (*gcpZoneValidator)(nil)

// Valid GCP zones as of 2024
// Zones are region + zone suffix (a, b, c, d, e, f)
var validGCPZones = map[string]bool{
	// us-central1
	"us-central1-a": true,
	"us-central1-b": true,
	"us-central1-c": true,
	"us-central1-f": true,
	// us-east1
	"us-east1-b": true,
	"us-east1-c": true,
	"us-east1-d": true,
	// us-east4
	"us-east4-a": true,
	"us-east4-b": true,
	"us-east4-c": true,
	// us-east5
	"us-east5-a": true,
	"us-east5-b": true,
	"us-east5-c": true,
	// us-south1
	"us-south1-a": true,
	"us-south1-b": true,
	"us-south1-c": true,
	// us-west1
	"us-west1-a": true,
	"us-west1-b": true,
	"us-west1-c": true,
	// us-west2
	"us-west2-a": true,
	"us-west2-b": true,
	"us-west2-c": true,
	// us-west3
	"us-west3-a": true,
	"us-west3-b": true,
	"us-west3-c": true,
	// us-west4
	"us-west4-a": true,
	"us-west4-b": true,
	"us-west4-c": true,
	// northamerica-northeast1
	"northamerica-northeast1-a": true,
	"northamerica-northeast1-b": true,
	"northamerica-northeast1-c": true,
	// northamerica-northeast2
	"northamerica-northeast2-a": true,
	"northamerica-northeast2-b": true,
	"northamerica-northeast2-c": true,
	// southamerica-east1
	"southamerica-east1-a": true,
	"southamerica-east1-b": true,
	"southamerica-east1-c": true,
	// southamerica-west1
	"southamerica-west1-a": true,
	"southamerica-west1-b": true,
	"southamerica-west1-c": true,
	// europe-central2
	"europe-central2-a": true,
	"europe-central2-b": true,
	"europe-central2-c": true,
	// europe-north1
	"europe-north1-a": true,
	"europe-north1-b": true,
	"europe-north1-c": true,
	// europe-southwest1
	"europe-southwest1-a": true,
	"europe-southwest1-b": true,
	"europe-southwest1-c": true,
	// europe-west1
	"europe-west1-b": true,
	"europe-west1-c": true,
	"europe-west1-d": true,
	// europe-west2
	"europe-west2-a": true,
	"europe-west2-b": true,
	"europe-west2-c": true,
	// europe-west3
	"europe-west3-a": true,
	"europe-west3-b": true,
	"europe-west3-c": true,
	// europe-west4
	"europe-west4-a": true,
	"europe-west4-b": true,
	"europe-west4-c": true,
	// europe-west6
	"europe-west6-a": true,
	"europe-west6-b": true,
	"europe-west6-c": true,
	// europe-west8
	"europe-west8-a": true,
	"europe-west8-b": true,
	"europe-west8-c": true,
	// europe-west9
	"europe-west9-a": true,
	"europe-west9-b": true,
	"europe-west9-c": true,
	// europe-west10
	"europe-west10-a": true,
	"europe-west10-b": true,
	"europe-west10-c": true,
	// europe-west12
	"europe-west12-a": true,
	"europe-west12-b": true,
	"europe-west12-c": true,
	// asia-east1
	"asia-east1-a": true,
	"asia-east1-b": true,
	"asia-east1-c": true,
	// asia-east2
	"asia-east2-a": true,
	"asia-east2-b": true,
	"asia-east2-c": true,
	// asia-northeast1
	"asia-northeast1-a": true,
	"asia-northeast1-b": true,
	"asia-northeast1-c": true,
	// asia-northeast2
	"asia-northeast2-a": true,
	"asia-northeast2-b": true,
	"asia-northeast2-c": true,
	// asia-northeast3
	"asia-northeast3-a": true,
	"asia-northeast3-b": true,
	"asia-northeast3-c": true,
	// asia-south1
	"asia-south1-a": true,
	"asia-south1-b": true,
	"asia-south1-c": true,
	// asia-south2
	"asia-south2-a": true,
	"asia-south2-b": true,
	"asia-south2-c": true,
	// asia-southeast1
	"asia-southeast1-a": true,
	"asia-southeast1-b": true,
	"asia-southeast1-c": true,
	// asia-southeast2
	"asia-southeast2-a": true,
	"asia-southeast2-b": true,
	"asia-southeast2-c": true,
	// australia-southeast1
	"australia-southeast1-a": true,
	"australia-southeast1-b": true,
	"australia-southeast1-c": true,
	// australia-southeast2
	"australia-southeast2-a": true,
	"australia-southeast2-b": true,
	"australia-southeast2-c": true,
	// me-central1
	"me-central1-a": true,
	"me-central1-b": true,
	"me-central1-c": true,
	// me-central2
	"me-central2-a": true,
	"me-central2-b": true,
	"me-central2-c": true,
	// me-west1
	"me-west1-a": true,
	"me-west1-b": true,
	"me-west1-c": true,
	// africa-south1
	"africa-south1-a": true,
	"africa-south1-b": true,
	"africa-south1-c": true,
}

func (gcpZoneValidator) Description(_ context.Context) string {
	return "value must be a valid GCP zone"
}

func (v gcpZoneValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (gcpZoneValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if diag := validateStringInMap(value, validGCPZones, req.Path, "Invalid GCP Zone", "GCP zone"); diag != nil {
		resp.Diagnostics.Append(diag)
	}
}
