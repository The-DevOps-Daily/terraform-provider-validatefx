package validators

import (
	"context"
	"fmt"
	"strings"
	"time"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const defaultDateTimeLayout = time.RFC3339Nano

var _ frameworkvalidator.String = DateTime(nil)

// DateTime returns a schema.String validator enforcing ISO 8601 / RFC 3339 datetimes.
// Optional layouts may be provided to extend accepted formats.
func DateTime(layouts []string) frameworkvalidator.String {
	return &dateTimeValidator{layouts: normalizeLayouts(layouts)}
}

type dateTimeValidator struct {
	layouts []string
}

func (v *dateTimeValidator) Description(_ context.Context) string {
	return "value must be a valid datetime string"
}

func (v *dateTimeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *dateTimeValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := strings.TrimSpace(req.ConfigValue.ValueString())
	if value == "" {
		return
	}

	layouts := v.layouts
	if len(layouts) == 0 {
		layouts = []string{defaultDateTimeLayout}
	}

	if err := validateAgainstLayouts(value, layouts); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, err.Summary, err.Detail)
	}
}

func normalizeLayouts(layouts []string) []string {
	if len(layouts) == 0 {
		return nil
	}

	uniq := make([]string, 0, len(layouts))
	seen := make(map[string]struct{})
	for _, layout := range layouts {
		layout = strings.TrimSpace(layout)
		if layout == "" {
			continue
		}
		if _, ok := seen[layout]; ok {
			continue
		}
		uniq = append(uniq, layout)
		seen[layout] = struct{}{}
	}

	return uniq
}

func validateAgainstLayouts(value string, layouts []string) *dateTimeError {
	var firstErr *dateTimeError

	for _, layout := range layouts {
		if layout == "" {
			continue
		}

		if layout == defaultDateTimeLayout {
			if err := validateRFC3339(value); err == nil {
				return nil
			} else if firstErr == nil {
				firstErr = err
			}
			continue
		}

		if _, err := time.Parse(layout, value); err == nil {
			return nil
		} else if firstErr == nil {
			firstErr = &dateTimeError{
				Summary: "Invalid Datetime",
				Detail:  fmt.Sprintf("Value %q does not match layout %q (%s)", value, layout, err.Error()),
			}
		}
	}

	if firstErr != nil {
		return firstErr
	}

	return &dateTimeError{
		Summary: "Invalid Datetime",
		Detail:  fmt.Sprintf("Value %q is not a valid datetime.", value),
	}
}

func validateRFC3339(value string) *dateTimeError {
	if _, err := time.Parse(defaultDateTimeLayout, value); err == nil {
		return nil
	}

	segments := strings.SplitN(value, "T", 2)
	if len(segments) != 2 {
		return &dateTimeError{
			Summary: "Invalid Datetime",
			Detail:  fmt.Sprintf("Value %q must contain a 'T' between date and time components", value),
		}
	}

	datePart, timePart := segments[0], segments[1]

	if _, err := time.Parse("2006-01-02", datePart); err != nil {
		return &dateTimeError{
			Summary: "Invalid Datetime",
			Detail:  fmt.Sprintf("invalid date component (%s)", err.Error()),
		}
	}

	timeLayouts := []string{
		"15:04:05Z07:00",
		"15:04:05.999999999Z07:00",
	}

	for _, layout := range timeLayouts {
		if _, err := time.Parse(layout, timePart); err == nil {
			return nil
		}
	}

	return &dateTimeError{
		Summary: "Invalid Datetime",
		Detail:  "invalid time component",
	}
}

type dateTimeError struct {
	Summary string
	Detail  string
}
