package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzSemVerValidator checks regex-based semver validator for robustness.
func FuzzSemVerValidator(f *testing.F) {
	seeds := []string{"", "1.0.0", "v1.2.3", "1.0.0-alpha.1", "1.0.0+build.1", "01.0.0", "1.2", "not-a-version"}
	for _, s := range seeds {
		f.Add(s)
	}
	v := SemVer()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("semver"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
			return
		}
		// Expectation based on semverPattern
		expect := semverPattern.MatchString(s)
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q: match=%v diagErr=%v", s, expect, resp.Diagnostics.HasError())
		}
	})
}
