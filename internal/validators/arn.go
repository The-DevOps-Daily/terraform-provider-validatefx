package validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ARN validates that a string is an AWS ARN (loosely, but safely).
// Pattern source: adapted to allow typical ARN segments: arn:partition:service:region:account-id:resource
// We avoid being overly strict on service-specific resource formats while ensuring the ARN skeleton is correct.
func ARN() validator.String { return arnValidator{} }

type arnValidator struct{}

var _ validator.String = (*arnValidator)(nil)

// Accept empty region/account segments, but enforce skeleton and resource presence.
var arnRe = regexp.MustCompile(`^arn:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:(?:[0-9]{12})?:.+`)

func (arnValidator) Description(_ context.Context) string {
	return "value must be a valid AWS ARN"
}

func (v arnValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (arnValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	s := req.ConfigValue.ValueString()
	if s == "" {
		return
	}
	if !arnRe.MatchString(s) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid ARN",
			fmt.Sprintf("Value %q is not a valid AWS ARN.", s),
		)
	}
}
