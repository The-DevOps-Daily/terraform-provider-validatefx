package validators

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	validLabelRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]*[a-zA-Z0-9])?$`)
)

var _ frameworkvalidator.String = Domain()

// Domain returns a schema.String validator which enforces valid domain format.
// The validator checks for proper domain format according to RFC 1123 and RFC 952.
func Domain() frameworkvalidator.String {
	return domainValidator{}
}

type domainValidator struct{}

func (domainValidator) Description(_ context.Context) string {
	return "value must be a valid domain name"
}

func (domainValidator) MarkdownDescription(_ context.Context) string {
	return "value must be a valid domain name"
}

func (domainValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if value == "" {
		return
	}

	if !isValidDomain(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Domain",
			fmt.Sprintf("Value %q is not a valid domain name", value),
		)
	}
}

// isValidDomain validates a domain name according to RFC 1123 and RFC 952
func isValidDomain(domain string) bool {
	// Basic length check
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// Remove trailing dot if present (FQDN)
	if strings.HasSuffix(domain, ".") {
		domain = domain[:len(domain)-1]
	}

	// Check if domain is empty after removing trailing dot
	if len(domain) == 0 {
		return false
	}

	// Split domain into labels
	labels := strings.Split(domain, ".")

	// Must have at least two labels (e.g., example.com)
	if len(labels) < 2 {
		return false
	}

	// Validate each label
	for _, label := range labels {
		if !isValidLabel(label) {
			return false
		}
	}

	// The last label (TLD) must not be all numeric
	lastLabel := labels[len(labels)-1]
	if isAllNumeric(lastLabel) {
		return false
	}

	return true
}

// isValidLabel validates a single domain label
func isValidLabel(label string) bool {
	// Label length must be 1-63 characters
	if len(label) == 0 || len(label) > 63 {
		return false
	}

	// Label cannot start or end with a hyphen
	if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
		return false
	}

	// Label can only contain alphanumeric characters and hyphens
	return validLabelRegex.MatchString(label)
}

// isAllNumeric checks if a string contains only numeric characters
func isAllNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
