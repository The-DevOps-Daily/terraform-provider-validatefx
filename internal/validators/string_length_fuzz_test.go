package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzStringLengthValidator(f *testing.F) {
	seeds := []string{"", "ab", "abc", "abcdef", "こんにちは", "emoji🚀"}
	for _, s := range seeds {
		f.Add(s)
	}

	min, max := 3, 6
	v := NewStringLengthValidator(&min, &max)
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("len"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		// Expect valid if rune length within [3,6]
		l := len([]rune(s))
		expect := l >= 3 && l <= 6
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q (len=%d): expect=%v diagErr=%v", s, l, expect, resp.Diagnostics.HasError())
		}
	})
}
