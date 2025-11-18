package functions

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestSubnetFunction(t *testing.T) {
	tests := []struct {
		name        string
		value       basetypes.StringValue
		expected    attr.Value
		expectError bool
	}{
		{
			name:        "valid IPv4 CIDR",
			value:       basetypes.NewStringValue("192.168.1.0/24"),
			expected:    basetypes.NewBoolValue(true),
			expectError: false,
		},
		{
			name:        "valid IPv6 CIDR",
			value:       basetypes.NewStringValue("2001:db8::/32"),
			expected:    basetypes.NewBoolValue(true),
			expectError: false,
		},
		{
			name:        "valid /32 subnet",
			value:       basetypes.NewStringValue("10.0.0.1/32"),
			expected:    basetypes.NewBoolValue(true),
			expectError: false,
		},
		{
			name:        "invalid missing mask",
			value:       basetypes.NewStringValue("192.168.1.0"),
			expected:    nil,
			expectError: true,
		},
		{
			name:        "invalid IP format",
			value:       basetypes.NewStringValue("not-a-cidr"),
			expected:    nil,
			expectError: true,
		},
		{
			name:        "invalid mask value",
			value:       basetypes.NewStringValue("192.168.1.0/33"),
			expected:    nil,
			expectError: true,
		},
		{
			name:        "null value",
			value:       basetypes.NewStringNull(),
			expected:    basetypes.NewBoolUnknown(),
			expectError: false,
		},
		{
			name:        "unknown value",
			value:       basetypes.NewStringUnknown(),
			expected:    basetypes.NewBoolUnknown(),
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fn := NewSubnetFunction()
			req := function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{tc.value}),
			}
			resp := &function.RunResponse{Result: function.NewResultData(basetypes.NewBoolNull())}

			fn.Run(context.Background(), req, resp)

			if tc.expectError {
				if resp.Error == nil {
					t.Fatalf("expected error, got none")
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %s", resp.Error)
			}

			boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok {
				t.Fatalf("unexpected result type %T", resp.Result.Value())
			}

			if diff := cmp.Diff(tc.expected, boolVal); diff != "" {
				t.Errorf("unexpected result (-expected +got): %s", diff)
			}
		})
	}
}
