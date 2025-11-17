package validators

import (
	"context"
	"testing"

	frameworkdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIPRangeSizeValidator_Valid(t *testing.T) {
	t.Parallel()

	v := NewIPRangeSizeValidator(8, 28)

	cases := []string{
		"10.0.0.0/8",
		"192.168.1.0/24",
		"2001:db8::/28",
	}

	for _, c := range cases {
		req := frameworkvalidator.StringRequest{Path: path.Root("cidr"), ConfigValue: types.StringValue(c)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for %s, got %v", c, resp.Diagnostics)
		}
	}
}

func TestIPRangeSizeValidator_OutOfRange(t *testing.T) {
	t.Parallel()

	v := NewIPRangeSizeValidator(8, 24)

	cases := []string{
		"10.0.0.0/7",  // too broad
		"10.0.0.0/25", // too narrow
	}

	for _, c := range cases {
		req := frameworkvalidator.StringRequest{Path: path.Root("cidr"), ConfigValue: types.StringValue(c)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected diagnostics for %s", c)
		}
		if resp.Diagnostics[0].Severity() != frameworkdiag.SeverityError {
			t.Fatalf("expected error severity for %s", c)
		}
	}
}

func TestIPRangeSizeValidator_InvalidCIDR(t *testing.T) {
	t.Parallel()

	v := NewIPRangeSizeValidator(8, 24)

	cases := []string{"not-a-cidr", "192.168.0.1", "2001:db8::/129"}
	for _, c := range cases {
		req := frameworkvalidator.StringRequest{Path: path.Root("cidr"), ConfigValue: types.StringValue(c)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected diagnostics for %s", c)
		}
	}
}

func TestIPRangeSizeValidator_NullUnknown(t *testing.T) {
	t.Parallel()

	v := NewIPRangeSizeValidator(8, 24)

	cases := []frameworkvalidator.StringRequest{
		{Path: path.Root("cidr"), ConfigValue: types.StringNull()},
		{Path: path.Root("cidr"), ConfigValue: types.StringUnknown()},
		{Path: path.Root("cidr"), ConfigValue: types.StringValue("")},
	}
	for _, req := range cases {
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for %v", req.ConfigValue)
		}
	}
}
