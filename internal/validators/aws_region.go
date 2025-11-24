package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// AWSRegion validates that a string is a valid AWS region code.
func AWSRegion() validator.String { return awsRegionValidator{} }

type awsRegionValidator struct{}

var _ validator.String = (*awsRegionValidator)(nil)

// Valid AWS regions as of 2024
var validAWSRegions = map[string]bool{
	// US regions
	"us-east-1": true,
	"us-east-2": true,
	"us-west-1": true,
	"us-west-2": true,
	// US GovCloud
	"us-gov-east-1": true,
	"us-gov-west-1": true,
	// Canada
	"ca-central-1": true,
	// Europe
	"eu-central-1": true,
	"eu-central-2": true,
	"eu-west-1":    true,
	"eu-west-2":    true,
	"eu-west-3":    true,
	"eu-north-1":   true,
	"eu-south-1":   true,
	"eu-south-2":   true,
	// Asia Pacific
	"ap-east-1":      true,
	"ap-south-1":     true,
	"ap-south-2":     true,
	"ap-northeast-1": true,
	"ap-northeast-2": true,
	"ap-northeast-3": true,
	"ap-southeast-1": true,
	"ap-southeast-2": true,
	"ap-southeast-3": true,
	"ap-southeast-4": true,
	// South America
	"sa-east-1": true,
	// Middle East
	"me-central-1": true,
	"me-south-1":   true,
	// Africa
	"af-south-1": true,
	// China
	"cn-north-1":     true,
	"cn-northwest-1": true,
}

func (awsRegionValidator) Description(_ context.Context) string {
	return "value must be a valid AWS region code"
}

func (v awsRegionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (awsRegionValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if diag := validateStringInMap(value, validAWSRegions, req.Path, "Invalid AWS Region", "AWS region code"); diag != nil {
		resp.Diagnostics.Append(diag)
	}
}
