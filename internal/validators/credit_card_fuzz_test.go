package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzCreditCardValidator validates robustness of Luhn-based credit card validator.
func FuzzCreditCardValidator(f *testing.F) {
	seeds := []string{
		"", "4111 1111 1111 1111", // Visa test number
		"5500-0000-0000-0004", // MasterCard test pattern
		"340000000000009",     // Amex-like length
		"1234567890123456",    // invalid Luhn
		"0000 0000 0000 0000", // all zeros -> invalid
		"4242 4242 4242 4242",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	v := CreditCard()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("credit_card"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
			return
		}
		// Simple oracle using internal helper
		expect := isValidCreditCard(s)
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q: luhn-ok=%v diagErr=%v", s, expect, resp.Diagnostics.HasError())
		}
	})
}
