package validators

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Simple E.164-like expectation: leading + and 8-15 digits
var e164Like = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)

func FuzzPhoneValidator(f *testing.F) {
	seeds := []string{"", "+12025550123", "+442071838750", "+8613800138000", "12345", "+000", "+1(202)555-0123"}
	for _, s := range seeds {
		f.Add(s)
	}
	v := Phone()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("phone"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
			return
		}
		expect := e164Like.MatchString(s)
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q: e164like=%v diagErr=%v", s, expect, resp.Diagnostics.HasError())
		}
	})
}
