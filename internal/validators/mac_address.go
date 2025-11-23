package validators

import (
	"context"
	"fmt"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ frameworkvalidator.String = MACAddress()

// MACAddress returns a schema.String validator that ensures the value is a valid MAC address.
func MACAddress() frameworkvalidator.String {
	return macAddressValidator{}
}

type macAddressValidator struct{}

func (macAddressValidator) Description(_ context.Context) string {
	return "value must be a valid MAC address"
}

func (v macAddressValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (macAddressValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value == "" {
		return
	}

	if _, err := normalizeMACAddress(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid MAC Address",
			fmt.Sprintf("Value %q is not a valid MAC address: %s", value, err.Error()),
		)
	}
}

// normalizeMACAddress removes separators and uppercases hexadecimal digits for validation.
func normalizeMACAddress(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("value is empty after trimming whitespace")
	}

	separator, err := detectMACSeparator(trimmed)
	if err != nil {
		return "", err
	}

	if separator != 0 {
		return normalizeSeparatedMAC(trimmed, separator)
	}

	return normalizeCompactMAC(trimmed)
}

func isHexadecimal(value string) bool {
	for _, r := range value {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
			return false
		}
	}
	return true
}

func detectMACSeparator(value string) (rune, error) {
	hasColon := strings.Contains(value, ":")
	hasDash := strings.Contains(value, "-")

	switch {
	case hasColon && hasDash:
		return 0, fmt.Errorf("mixed separators are not permitted")
	case hasColon:
		return ':', nil
	case hasDash:
		return '-', nil
	default:
		return 0, nil
	}
}

func normalizeSeparatedMAC(value string, separator rune) (string, error) {
	parts := strings.Split(value, string(separator))
	if len(parts) != 6 {
		return "", fmt.Errorf("expected exactly 6 octets separated by %q", string(separator))
	}

	for i, part := range parts {
		if len(part) != 2 {
			return "", fmt.Errorf("octet %d must be exactly 2 hexadecimal characters", i+1)
		}
		if !isHexadecimal(part) {
			return "", fmt.Errorf("octet %d contains non-hexadecimal characters", i+1)
		}
	}

	return strings.ToUpper(strings.Join(parts, "")), nil
}

func normalizeCompactMAC(value string) (string, error) {
	if len(value) != 12 {
		return "", fmt.Errorf("expected exactly 12 hexadecimal characters for compact format")
	}

	if !isHexadecimal(value) {
		return "", fmt.Errorf("value contains non-hexadecimal characters")
	}

	return strings.ToUpper(value), nil
}
