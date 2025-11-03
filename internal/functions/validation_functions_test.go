package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func callFunction(t *testing.T, fn function.Function, args ...attr.Value) *function.RunResponse {
	t.Helper()

	resp := &function.RunResponse{}
	fn.Run(context.Background(), function.RunRequest{
		Arguments: function.NewArgumentsData(args),
	}, resp)
	return resp
}

func boolResult(t *testing.T, resp *function.RunResponse) basetypes.BoolValue {
	t.Helper()

	value := resp.Result.Value()
	boolVal, ok := value.(basetypes.BoolValue)
	if !ok {
		t.Fatalf("unexpected result type %T", value)
	}

	return boolVal
}

func TestStringValidationFunctions(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		factory func() function.Function
		valid   string
		invalid string
	}{
		"base64": {
			factory: NewBase64Function,
			valid:   "U29sdmVkIQ==",
			invalid: "not-base64",
		},
		"credit_card": {
			factory: NewCreditCardFunction,
			valid:   "4532015112830366",
			invalid: "4532015112830367",
		},
		"domain": {
			factory: NewDomainFunction,
			valid:   "example.com",
			invalid: "invalid..domain",
		},
		"email": {
			factory: NewEmailFunction,
			valid:   "alice@example.com",
			invalid: "bad-email",
		},
		"hostname": {
			factory: NewHostnameFunction,
			valid:   "service.internal",
			invalid: "bad_name",
		},
		"ip": {
			factory: NewIPFunction,
			valid:   "127.0.0.1",
			invalid: "999.999.999.999",
		},
		"json": {
			factory: NewJSONFunction,
			valid:   "{\"key\":\"value\"}",
			invalid: "invalid-json",
		},
		"mac_address": {
			factory: NewMACAddressFunction,
			valid:   "00:1A:2B:3C:4D:5E",
			invalid: "AA:BB:CC:DD:EE",
		},
		"phone": {
			factory: NewPhoneFunction,
			valid:   "+14155552671",
			invalid: "14155552671",
		},
		"url": {
			factory: NewURLFunction,
			valid:   "https://example.com",
			invalid: "ftp://example.com",
		},
		"uuid": {
			factory: NewUUIDFunction,
			valid:   "d9428888-122b-11e1-b85c-61cd3cbb3210",
			invalid: "not-a-uuid",
		},
		"cidr": {
			factory: NewCIDRFunction,
			valid:   "10.0.0.0/24",
			invalid: "10.0.0.0/33",
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fn := tc.factory()

			resp := callFunction(t, fn, basetypes.NewStringValue(tc.valid))
			if resp.Error != nil {
				t.Fatalf("unexpected error: %s", resp.Error)
			}

			result := boolResult(t, resp)
			if !result.ValueBool() {
				t.Fatalf("expected true result for %s", name)
			}

			invalidResp := callFunction(t, fn, basetypes.NewStringValue(tc.invalid))
			if invalidResp.Error == nil {
				t.Fatalf("expected error for invalid value %q", tc.invalid)
			}

			nullResp := callFunction(t, fn, basetypes.NewStringNull())
			if nullResp.Error != nil {
				t.Fatalf("unexpected error for null input: %s", nullResp.Error)
			}

			nullResult := boolResult(t, nullResp)
			if !nullResult.IsUnknown() {
				t.Fatalf("expected unknown result for null input")
			}
		})
	}
}

func TestBetweenFunction(t *testing.T) {
	t.Parallel()

	fn := NewBetweenFunction()

	resp := callFunction(t, fn,
		basetypes.NewStringValue("7.5"),
		basetypes.NewStringValue("5"),
		basetypes.NewStringValue("10"),
	)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error)
	}

	result := boolResult(t, resp)
	if !result.ValueBool() {
		t.Fatalf("expected true result")
	}

	smallResp := callFunction(t, fn,
		basetypes.NewStringValue("2"),
		basetypes.NewStringValue("5"),
		basetypes.NewStringValue("10"),
	)
	if smallResp.Error == nil {
		t.Fatalf("expected error for out-of-range value")
	}

	boundsResp := callFunction(t, fn,
		basetypes.NewStringValue("7"),
		basetypes.NewStringValue("10"),
		basetypes.NewStringValue("5"),
	)
	if boundsResp.Error == nil {
		t.Fatalf("expected error for invalid bounds")
	}

	unknownResp := callFunction(t, fn,
		basetypes.NewStringUnknown(),
		basetypes.NewStringValue("5"),
		basetypes.NewStringValue("10"),
	)
	if unknownResp.Error != nil {
		t.Fatalf("unexpected error for unknown value: %s", unknownResp.Error)
	}

	unknownResult := boolResult(t, unknownResp)
	if !unknownResult.IsUnknown() {
		t.Fatalf("expected unknown result for unknown value")
	}
}

func TestStringLengthFunction(t *testing.T) {
	t.Parallel()

	fn := NewStringLengthFunction()

	resp := callFunction(t, fn,
		basetypes.NewStringValue("hello"),
		basetypes.NewInt64Value(3),
		basetypes.NewInt64Value(10),
	)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error)
	}

	result := boolResult(t, resp)
	if !result.ValueBool() {
		t.Fatalf("expected true result")
	}

	longResp := callFunction(t, fn,
		basetypes.NewStringValue("this string is too long"),
		basetypes.NewInt64Value(3),
		basetypes.NewInt64Value(10),
	)
	if longResp.Error == nil {
		t.Fatalf("expected error for string exceeding max length")
	}

	unknownResp := callFunction(t, fn,
		basetypes.NewStringUnknown(),
		basetypes.NewInt64Null(),
		basetypes.NewInt64Null(),
	)
	if unknownResp.Error != nil {
		t.Fatalf("unexpected error for unknown value: %s", unknownResp.Error)
	}

	unknownResult := boolResult(t, unknownResp)
	if !unknownResult.IsUnknown() {
		t.Fatalf("expected unknown result for unknown value")
	}
}

func TestMatchesRegexFunction(t *testing.T) {
	t.Parallel()

	fn := NewMatchesRegexFunction()

	resp := callFunction(t, fn,
		basetypes.NewStringValue("user_123"),
		basetypes.NewStringValue("^[a-z0-9_]+$"),
	)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error)
	}

	if !boolResult(t, resp).ValueBool() {
		t.Fatalf("expected true result")
	}

	mismatch := callFunction(t, fn,
		basetypes.NewStringValue("Invalid-User"),
		basetypes.NewStringValue("^[a-z0-9_]+$"),
	)
	if mismatch.Error == nil {
		t.Fatalf("expected error for pattern mismatch")
	}

	compileErr := callFunction(t, fn,
		basetypes.NewStringValue("value"),
		basetypes.NewStringValue("("),
	)
	if compileErr.Error == nil {
		t.Fatalf("expected error for invalid regex pattern")
	}

	unknown := callFunction(t, fn,
		basetypes.NewStringUnknown(),
		basetypes.NewStringValue(".*"),
	)
	if unknown.Error != nil {
		t.Fatalf("unexpected error for unknown value: %s", unknown.Error)
	}

	if !boolResult(t, unknown).IsUnknown() {
		t.Fatalf("expected unknown result for unknown value")
	}
}

func TestDateTimeFunction(t *testing.T) {
	t.Parallel()

	fn := NewDateTimeFunction()

	nullLayouts := basetypes.NewListNull(basetypes.StringType{})

	resp := callFunction(t, fn,
		basetypes.NewStringValue("2025-11-02T15:04:05Z"),
		nullLayouts,
	)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error)
	}

	if !boolResult(t, resp).ValueBool() {
		t.Fatalf("expected true result")
	}

	customLayout := basetypes.NewListValueMust(
		basetypes.StringType{},
		[]attr.Value{basetypes.NewStringValue("2006-01-02 15:04:05")},
	)
	customResp := callFunction(t, fn,
		basetypes.NewStringValue("2025-11-02 15:04:05"),
		customLayout,
	)
	if customResp.Error != nil {
		t.Fatalf("unexpected error with custom layouts: %s", customResp.Error)
	}

	if !boolResult(t, customResp).ValueBool() {
		t.Fatalf("expected true result with custom layout")
	}

	invalidResp := callFunction(t, fn,
		basetypes.NewStringValue("2025-13-02T15:04:05Z"),
		nullLayouts,
	)
	if invalidResp.Error == nil {
		t.Fatalf("expected error for invalid datetime")
	}

	unknownResp := callFunction(t, fn,
		basetypes.NewStringUnknown(),
		nullLayouts,
	)
	if unknownResp.Error != nil {
		t.Fatalf("unexpected error for unknown value: %s", unknownResp.Error)
	}

	if !boolResult(t, unknownResp).IsUnknown() {
		t.Fatalf("expected unknown result for unknown value")
	}
}

func TestAssertFunction(t *testing.T) {
	t.Parallel()

	fn := NewAssertFunction()

	success := callFunction(t, fn,
		basetypes.NewBoolValue(true),
		basetypes.NewStringValue("message"),
	)
	if success.Error != nil {
		t.Fatalf("unexpected error: %s", success.Error)
	}

	if !boolResult(t, success).ValueBool() {
		t.Fatalf("expected true result")
	}

	failure := callFunction(t, fn,
		basetypes.NewBoolValue(false),
		basetypes.NewStringValue("boom"),
	)
	if failure.Error == nil {
		t.Fatalf("expected error for false condition")
	}

	unknown := callFunction(t, fn,
		basetypes.NewBoolUnknown(),
		basetypes.NewStringValue("ignored"),
	)
	if unknown.Error != nil {
		t.Fatalf("unexpected error for unknown condition: %s", unknown.Error)
	}

	if !boolResult(t, unknown).IsUnknown() {
		t.Fatalf("expected unknown result for unknown condition")
	}
}
