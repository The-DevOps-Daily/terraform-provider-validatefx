package validators

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzBase64Validator(f *testing.F) {
	seeds := []string{"", "U29mdHdhcmU=", "SGVsbG8gV29ybGQh", "not*base64", "Zm9vYmFy", "===="}
	for _, s := range seeds {
		f.Add(s)
	}

	v := Base64Validator()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("b64"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
			return
		}
		_, err := base64.StdEncoding.DecodeString(s)
		expect := err == nil
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch: decode-ok=%v diagErr=%v for %q", expect, resp.Diagnostics.HasError(), s)
		}
	})
}
