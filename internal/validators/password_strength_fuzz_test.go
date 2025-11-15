package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzPasswordStrengthValidator ensures no panics and basic expectation:
// strong-looking strings should pass, weak-looking ones should fail.
func FuzzPasswordStrengthValidator(f *testing.F) {
	seeds := []string{"", "short", "NoNumber!", "noupper1!", "NOLOWER1!", "Valid123!", "Another$Good9"}
	for _, s := range seeds {
		f.Add(s)
	}

	v := PasswordStrengthValidator()

	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("password"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
			return
		}

		// Heuristic expectation matching validator definition
		strong := len(s) >= 8
		hasUpper := false
		hasLower := false
		hasNum := false
		hasSpecial := false
		for _, r := range s {
			switch {
			case r >= 'A' && r <= 'Z':
				hasUpper = true
			case r >= 'a' && r <= 'z':
				hasLower = true
			case r >= '0' && r <= '9':
				hasNum = true
			case (r >= 33 && r <= 47) || (r >= 58 && r <= 64) || (r >= 91 && r <= 96) || (r >= 123 && r <= 126):
				hasSpecial = true
			}
		}
		strong = strong && hasUpper && hasLower && hasNum && hasSpecial

		if strong != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q: strong=%v diagErr=%v", s, strong, resp.Diagnostics.HasError())
		}
	})
}
