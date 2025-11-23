package validators

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var hostnameLabelRegex = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`)

var _ frameworkvalidator.String = Hostname()

// Hostname returns a schema.String validator enforcing RFC 1123 hostname rules.
func Hostname() frameworkvalidator.String {
	return hostnameValidator{}
}

type hostnameValidator struct{}

func (hostnameValidator) Description(_ context.Context) string {
	return "value must be a valid hostname"
}

func (v hostnameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (hostnameValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if value == "" {
		return
	}

	if !isRFC1123Hostname(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Hostname",
			fmt.Sprintf("Value %q is not a valid hostname", value),
		)
	}
}

func isRFC1123Hostname(value string) bool {
	if len(value) == 0 {
		return false
	}

	if strings.HasSuffix(value, ".") {
		value = value[:len(value)-1]
	}

	if len(value) == 0 || len(value) > 253 {
		return false
	}

	if strings.Contains(value, "..") {
		return false
	}

	labels := strings.Split(value, ".")

	for _, label := range labels {
		if !isHostnameLabel(label) {
			return false
		}
	}

	return true
}

func isHostnameLabel(label string) bool {
	if len(label) == 0 || len(label) > 63 {
		return false
	}

	if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
		return false
	}

	return hostnameLabelRegex.MatchString(label)
}
