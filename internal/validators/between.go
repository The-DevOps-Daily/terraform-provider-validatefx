package validators

import (
	"context"
	"fmt"
	"math/big"
	"strings"

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

	min, minSet, err := parseDecimal(v.min)
	if err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Minimum", err.Error())
		return
	}

	max, maxSet, err := parseDecimal(v.max)
	if err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Maximum", err.Error())
		return
	}

	if minSet && maxSet && min.Cmp(max) == 1 {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Range", fmt.Sprintf("minimum %s cannot be greater than maximum %s", v.min, v.max))
		return
	}

	num, ok := new(big.Float).SetString(value)
	if !ok {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Number", fmt.Sprintf("Value %q is not a valid decimal", value))
		return
	}

	if minSet && num.Cmp(min) == -1 {
		resp.Diagnostics.AddAttributeError(req.Path, "Value Too Small", fmt.Sprintf("Value %q is less than minimum %s", value, v.min))
		return
	}

	if maxSet && num.Cmp(max) == 1 {
		resp.Diagnostics.AddAttributeError(req.Path, "Value Too Large", fmt.Sprintf("Value %q is greater than maximum %s", value, v.max))
		return
	}
}

func parseDecimal(input string) (*big.Float, bool, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil, false, nil
	}

	num, ok := new(big.Float).SetString(trimmed)
	if !ok {
		return nil, false, fmt.Errorf("%q is not a valid decimal", input)
	}

	return num, true, nil
}
