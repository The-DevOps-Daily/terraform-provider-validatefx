package validators

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ARN validates that a string is an AWS ARN (loosely, but safely).
// Pattern source: adapted to allow typical ARN segments: arn:partition:service:region:account-id:resource
// We avoid being overly strict on service-specific resource formats while ensuring the ARN skeleton is correct.
func ARN() validator.String { return arnValidator{} }

type arnValidator struct{}

var _ validator.String = (*arnValidator)(nil)

// Loose skeleton capture for service-aware validation: arn:partition:service:region:account:resource
var arnSkeleton = regexp.MustCompile(`^arn:([^:]+):([^:]+):([^:]*):([^:]*):(.+)$`)
var regionRe = regexp.MustCompile(`^[a-z]{2}-(gov-)?[a-z]+-\d$`)
var accountRe = regexp.MustCompile(`^\d{12}$`)

func (arnValidator) Description(_ context.Context) string {
	return "value must be a valid AWS ARN"
}

func (v arnValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (arnValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) { //nolint:cyclop
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	s := req.ConfigValue.ValueString()
	if s == "" {
		return
	}
	m := arnSkeleton.FindStringSubmatch(s)
	if m == nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", fmt.Sprintf("Value %q does not match ARN skeleton.", s))
		return
	}

	partition := m[1]
	service := m[2]
	region := m[3]
	account := m[4]
	resource := m[5]

	// Basic checks
	if resource == "" {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", "Resource component must be non-empty.")
		return
	}
	if strings.HasPrefix(resource, ":") {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", "Resource component must not start with a colon.")
		return
	}

	// Optional tighter partition check
	switch partition {
	case "aws", "aws-us-gov", "aws-cn":
		// ok
	default:
		// keep permissive; do not fail on custom partitions
	}

	switch service {
	case "s3":
		// S3 ARNs typically have empty region and account
		if region != "" {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", "S3 ARNs must have empty region.")
			return
		}
		if account != "" {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", "S3 ARNs must have empty account ID.")
			return
		}
		return
	case "iam":
		// IAM ARNs have empty region and 12-digit account; resource begins with known kinds
		if region != "" {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", "IAM ARNs must have empty region.")
			return
		}
		if !accountRe.MatchString(account) {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", "IAM ARNs must include a 12-digit account ID.")
			return
		}
		if !(strings.HasPrefix(resource, "user/") || strings.HasPrefix(resource, "role/") || strings.HasPrefix(resource, "group/") || strings.HasPrefix(resource, "policy/") || strings.HasPrefix(resource, "instance-profile/")) {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", "IAM resource must start with user/, role/, group/, policy/, or instance-profile/.")
			return
		}
		return
	case "lambda":
		// Lambda requires region and account; resource must start with function:
		if !regionRe.MatchString(region) {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", "Lambda ARNs must include a valid region.")
			return
		}
		if !accountRe.MatchString(account) {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", "Lambda ARNs must include a 12-digit account ID.")
			return
		}
		if !strings.HasPrefix(resource, "function:") {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", "Lambda resource must start with function:.")
			return
		}
		return
	default:
		// Generic rule: if account is present, it must be 12 digits; region if present should look like region
		if account != "" && !accountRe.MatchString(account) {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", "Account ID must be 12 digits when provided.")
			return
		}
		if region != "" && !regionRe.MatchString(region) {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid ARN", "Region must be a valid AWS region when provided.")
			return
		}
		return
	}
}
