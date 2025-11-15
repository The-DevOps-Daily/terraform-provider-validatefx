package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzUsernameValidator(f *testing.F) {
	seeds := []string{"", "ab", "abc", "user_name", "User-Name", "123", "a_very_long_username_exceeding_limit"}
	for _, s := range seeds {
		f.Add(s)
	}

	v := DefaultUsernameValidator()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("username"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		// Expect valid when length is within default bounds and only letters/digits/underscore
		validChars := isValidUsername(s)
		l := len([]rune(s))
		expect := validChars && l >= defaultUsernameMinLength && l <= defaultUsernameMaxLength

		if s == "" {
			expect = false
		}

		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q: expect=%v diagErr=%v", s, expect, resp.Diagnostics.HasError())
		}
	})
}
