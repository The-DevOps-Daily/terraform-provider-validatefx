package validators

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ frameworkvalidator.String = Between("", "")

// Between returns a validator ensuring the value falls within an inclusive numeric range.
func Between(minStr, maxStr string) frameworkvalidator.String {
	return &betweenValidator{
		min: strings.TrimSpace(minStr),
		max: strings.TrimSpace(maxStr),
	}
}

type betweenValidator struct {
	min string
	max string
}

func (v *betweenValidator) Description(_ context.Context) string {
	return "value must be a number between the configured minimum and maximum"
}

func (v *betweenValidator) MarkdownDescription(_ context.Context) string {
	return "value must be a number between the configured minimum and maximum"
}

func (v *betweenValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := strings.TrimSpace(req.ConfigValue.ValueString())
	if value == "" {
		return
	}

	min, minSet := parseBound(resp, req.Path, "Invalid Minimum", v.min)
	if resp.Diagnostics.HasError() {
		return
	}

	max, maxSet := parseBound(resp, req.Path, "Invalid Maximum", v.max)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateBounds(min, max, minSet, maxSet, v.min, v.max, req.Path, resp); err != nil {
		return
	}

	if err := validateValueWithin(value, min, max, minSet, maxSet, req.Path, resp); err != nil {
		return
	}
}

func parseBound(resp *frameworkvalidator.StringResponse, p path.Path, summary, bound string) (*big.Float, bool) {
	trimmed := strings.TrimSpace(bound)
	if trimmed == "" {
		return nil, false
	}

	num, ok := new(big.Float).SetString(trimmed)
	if !ok {
		resp.Diagnostics.AddAttributeError(p, summary, fmt.Sprintf("%q is not a valid decimal", bound))
		return nil, false
	}

	return num, true
}

func validateBounds(min, max *big.Float, minSet, maxSet bool, minRaw, maxRaw string, p path.Path, resp *frameworkvalidator.StringResponse) error {
	if minSet && maxSet && min.Cmp(max) == 1 {
		resp.Diagnostics.AddAttributeError(p, "Invalid Range", fmt.Sprintf("minimum %s cannot be greater than maximum %s", minRaw, maxRaw))
		return fmt.Errorf("invalid range")
	}
	return nil
}

func validateValueWithin(value string, min, max *big.Float, minSet, maxSet bool, p path.Path, resp *frameworkvalidator.StringResponse) error {
	num, ok := new(big.Float).SetString(value)
	if !ok {
		resp.Diagnostics.AddAttributeError(p, "Invalid Number", fmt.Sprintf("Value %q is not a valid decimal", value))
		return fmt.Errorf("invalid number")
	}

	if minSet && num.Cmp(min) == -1 {
		resp.Diagnostics.AddAttributeError(p, "Value Too Small", fmt.Sprintf("Value %q is less than minimum %s", value, min.Text('g', -1)))
		return fmt.Errorf("too small")
	}

	if maxSet && num.Cmp(max) == 1 {
		resp.Diagnostics.AddAttributeError(p, "Value Too Large", fmt.Sprintf("Value %q is greater than maximum %s", value, max.Text('g', -1)))
		return fmt.Errorf("too large")
	}

	return nil
}
