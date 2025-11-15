package validators

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzUUIDValidator cross-checks the validator against google/uuid parsing
// and ensures no panics on random inputs.
func FuzzUUIDValidator(f *testing.F) {
	for i := 0; i < 5; i++ {
		f.Add(uuid.NewString())
	}
	for _, s := range []string{"", "123", "not-a-uuid", "550e8400-e29b-41d4-a716-446655440000x"} {
		f.Add(s)
	}

	v := UUID()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("uuid"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
			return
		}

		_, err := uuid.Parse(s)
		expectValid := err == nil
		if expectValid != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q: parse-ok=%v diagErr=%v", s, expectValid, resp.Diagnostics.HasError())
		}
	})
}
