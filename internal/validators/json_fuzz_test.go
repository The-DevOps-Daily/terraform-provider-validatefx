package validators

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzJSONValidator validates robustness against arbitrary inputs and checks that
// strings which unmarshal to JSON objects are accepted while others are rejected.
func FuzzJSONValidator(f *testing.F) {
	for _, s := range []string{
		"", "{}", "{\"a\":1}", "[]", "true", "null", "{]", "{\"nested\": {\"b\":2}}",
	} {
		f.Add(s)
	}

	v := JSON()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("json"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error: %v", resp.Diagnostics)
			}
			return
		}

		var decoded any
		err := json.Unmarshal([]byte(s), &decoded)
		ok := err == nil
		if ok {
			_, ok = decoded.(map[string]any)
		}

		if ok != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch: decoded-object=%v diagnostics.HasError=%v value=%q", ok, resp.Diagnostics.HasError(), s)
		}
	})
}
