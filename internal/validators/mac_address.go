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

func (macAddressValidator) MarkdownDescription(_ context.Context) string {
	return "value must be a valid MAC address"
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

	var separator rune
	if strings.Contains(trimmed, ":") {
		separator = ':'
	}
	if strings.Contains(trimmed, "-") {
		if separator != 0 {
			return "", fmt.Errorf("mixed separators are not permitted")
		}
		separator = '-'
	}

	if separator != 0 {
		parts := strings.Split(trimmed, string(separator))
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

	if len(trimmed) != 12 {
		return "", fmt.Errorf("expected exactly 12 hexadecimal characters for compact format")
	}

	if !isHexadecimal(trimmed) {
		return "", fmt.Errorf("value contains non-hexadecimal characters")
	}

	return strings.ToUpper(trimmed), nil
}

func isHexadecimal(value string) bool {
	for _, r := range value {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
			return false
		}
	}
	return true
}
