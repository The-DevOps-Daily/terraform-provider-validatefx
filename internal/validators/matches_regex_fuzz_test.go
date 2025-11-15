package validators

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzMatchesRegexValidator(f *testing.F) {
	seeds := []struct{ pattern, value string }{
		{"^[a-z]+$", "abc"},
		{"^[a-z]+$", "ABC"},
		{"^foo|bar$", "foobar"},
		{"(unclosed", "x"},
		{"^$", ""},
	}
	for _, s := range seeds {
		f.Add(s.pattern, s.value)
	}

	f.Fuzz(func(t *testing.T, pattern, value string) {
		t.Parallel()
		v := MatchesRegex(pattern)
		req := frameworkvalidator.StringRequest{Path: path.Root("re"), ConfigValue: types.StringValue(value)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if value == "" { // empty allowed
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty value should not error (pattern=%q)", pattern)
			}
			return
		}

		compiled, err := regexp.Compile(pattern)
		if err != nil {
			if !resp.Diagnostics.HasError() {
				t.Fatalf("invalid pattern should error")
			}
			return
		}

		expect := compiled.MatchString(value)
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for pattern=%q value=%q: match=%v diagErr=%v", pattern, value, expect, resp.Diagnostics.HasError())
		}
	})
}
