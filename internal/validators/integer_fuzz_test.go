package validators

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzIntegerValidator(f *testing.F) {
	seeds := []string{"", "0", "-42", "+7", "3.14", "abc", "  10  "}
	for _, s := range seeds {
		f.Add(s)
	}

	v := Integer()
	re := regexp.MustCompile(`^[+-]?\d+$`)

	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("int"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		// Validator trims spaces; emulate that
		trimmed := s
		if len(s) > 0 && (s[0] == ' ' || s[len(s)-1] == ' ') {
			trimmed = strings.TrimSpace(s)
		}
		if trimmed == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty/space-only should not error")
			}
			return
		}
		expect := re.MatchString(trimmed)
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q (trimmed=%q): expect=%v diagErr=%v", s, trimmed, expect, resp.Diagnostics.HasError())
		}
	})
}
