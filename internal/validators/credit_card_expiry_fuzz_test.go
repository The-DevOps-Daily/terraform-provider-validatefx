package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzCreditCardExpiryValidator validates robustness of credit card expiry validator.
func FuzzCreditCardExpiryValidator(f *testing.F) {
	seeds := []string{
		"",
		"12/25",
		"01/2025",
		"12/99",
		"00/25",
		"13/25",
		"1/25",
		"01/5",
		"0125",
		"01-25",
		"AB/CD",
		"12/2000",
		"06/30",
		"12/2099",
		"99/99",
		"01/025",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	v := CreditCardExpiry()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{
			Path:        path.Root("expiry"),
			ConfigValue: types.StringValue(s),
		}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		// Empty strings should not error
		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty string should not error")
			}
			return
		}

		// The validator should never panic and always return a valid response
		// We don't check correctness here, just that it handles all inputs gracefully
	})
}
