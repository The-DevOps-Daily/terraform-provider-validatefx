package validators

import (
    "context"
    "regexp"
    "strings"

    frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SemVerRange returns a validator that ensures the string is a well-formed
// semantic version range expression like ">=1.0.0,<2.0.0".
func SemVerRange() frameworkvalidator.String { return semverRangeValidator{} }

type semverRangeValidator struct{}

var _ frameworkvalidator.String = (*semverRangeValidator)(nil)

// Accept comparators separated by commas. Each comparator has an operator
// (<=, >=, <, >, =) and a SemVer value (with optional leading v).
var (
    // reSemver is the same semantic version pattern used by the SemVer validator,
    // prefixed with optional 'v'.
    reSemver = regexp.MustCompile(`^v?` + semverPattern.String()[1:])
    reComparator = regexp.MustCompile(`^(<=|>=|<|>|=)\s*(.+)$`)
)

func (semverRangeValidator) Description(_ context.Context) string {
    return "value must be a valid SemVer range (e.g., >=1.0.0,<2.0.0)"
}

func (v semverRangeValidator) MarkdownDescription(ctx context.Context) string { return v.Description(ctx) }

func (semverRangeValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
    if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
        return
    }

    raw := strings.TrimSpace(req.ConfigValue.ValueString())
    if raw == "" {
        return
    }

    // Split by comma and require at least one comparator
    parts := strings.Split(raw, ",")
    if len(parts) == 0 {
        resp.Diagnostics.AddAttributeError(req.Path, "Invalid SemVer Range", "Range must contain at least one comparator")
        return
    }

    for _, part := range parts {
        p := strings.TrimSpace(part)
        if p == "" {
            resp.Diagnostics.AddAttributeError(req.Path, "Invalid SemVer Range", "Range must not contain empty comparators")
            return
        }
        m := reComparator.FindStringSubmatch(p)
        if m == nil {
            resp.Diagnostics.AddAttributeError(req.Path, "Invalid SemVer Range", "Each comparator must start with one of: <, <=, >, >=, =")
            return
        }
        ver := strings.TrimSpace(m[2])
        if !reSemver.MatchString(ver) {
            resp.Diagnostics.AddAttributeError(req.Path, "Invalid SemVer Range", "Comparator version must be a valid semantic version")
            return
        }
    }
}

