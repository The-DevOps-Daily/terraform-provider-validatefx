package validators

import (
	"context"
	"net/mail"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzEmailValidator fuzzes the Email() validator to ensure it never panics
// and broadly agrees with net/mail.ParseAddress semantics. Empty, null, and
// unknown are allowed by design and should not produce diagnostics.
func FuzzEmailValidator(f *testing.F) {
	// Seed with a few interesting cases
	for _, s := range []string{
		"", "a@b.com", "user.name+tag+sorting@example.com", "not-an-email", "missing-at.example.com",
		"A@b.c", "plainaddress", "@nouser", "name@localhost", "foo@bar..com",
	} {
		f.Add(s)
	}

	v := Email()

	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{
			Path:        path.Root("email"),
			ConfigValue: types.StringValue(s),
		}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		// Empty strings are allowed (treated as no-op)
		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty string should not error, got %v", resp.Diagnostics)
			}
			return
		}

		// Cross-check with net/mail.ParseAddress
		_, err := mail.ParseAddress(s)
		hasErr := err != nil
		if hasErr != resp.Diagnostics.HasError() {
			t.Fatalf("mismatch: ParseAddress error=%v, diagnostics.HasError=%v (value=%q)", hasErr, resp.Diagnostics.HasError(), s)
		}
	})
}
