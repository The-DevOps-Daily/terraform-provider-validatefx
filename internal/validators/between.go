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

	valid, boundsDiag, valueDiag := EvaluateBetween(value, v.min, v.max)
	if boundsDiag != nil {
		resp.Diagnostics.AddAttributeError(req.Path, boundsDiag.Summary, boundsDiag.Detail)
		return
	}

	if valueDiag != nil {
		resp.Diagnostics.AddAttributeError(req.Path, valueDiag.Summary, valueDiag.Detail)
		return
	}

	if !valid {
		return
	}
}

func EvaluateBetween(value, minRaw, maxRaw string) (bool, *BetweenDiagnostic, *BetweenDiagnostic) {
	value = strings.TrimSpace(value)
	minRaw = strings.TrimSpace(minRaw)
	maxRaw = strings.TrimSpace(maxRaw)

	min, minSet, diag := parseBound("Invalid Minimum", minRaw)
	if diag != nil {
		return false, diag, nil
	}

	max, maxSet, diag := parseBound("Invalid Maximum", maxRaw)
	if diag != nil {
		return false, diag, nil
	}

	if minSet && maxSet && min.Cmp(max) == 1 {
		return false, nil, &BetweenDiagnostic{
			Summary: "Invalid Range",
			Detail:  fmt.Sprintf("minimum %s cannot be greater than maximum %s", minRaw, maxRaw),
		}
	}

	num, ok := new(big.Float).SetString(value)
	if !ok {
		return false, nil, &BetweenDiagnostic{
			Summary: "Invalid Number",
			Detail:  fmt.Sprintf("Value %q is not a valid decimal", value),
		}
	}

	if minSet && num.Cmp(min) == -1 {
		return false, nil, &BetweenDiagnostic{
			Summary: "Value Too Small",
			Detail:  fmt.Sprintf("Value %q is less than minimum %s", value, min.Text('g', -1)),
		}
	}

	if maxSet && num.Cmp(max) == 1 {
		return false, nil, &BetweenDiagnostic{
			Summary: "Value Too Large",
			Detail:  fmt.Sprintf("Value %q is greater than maximum %s", value, max.Text('g', -1)),
		}
	}

	return true, nil, nil
}

func parseBound(summary, raw string) (*big.Float, bool, *BetweenDiagnostic) {
	if raw == "" {
		return nil, false, nil
	}

	num, ok := new(big.Float).SetString(raw)
	if !ok {
		return nil, false, &BetweenDiagnostic{
			Summary: summary,
			Detail:  fmt.Sprintf("%q is not a valid decimal", raw),
		}
	}

	return num, true, nil
}

// BetweenDiagnostic captures validation errors for invalid bounds or values.
type BetweenDiagnostic struct {
	Summary string
	Detail  string
}
