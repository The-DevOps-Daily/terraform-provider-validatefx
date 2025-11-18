package functions

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestPortNumberFunction(t *testing.T) {
	tests := []struct {
		name        string
		value       basetypes.StringValue
		expected    attr.Value
		expectError bool
	}{
		{
			name:        "valid port 80",
			value:       basetypes.NewStringValue("80"),
			expected:    basetypes.NewBoolValue(true),
			expectError: false,
		},
		{
			name:        "valid port 443",
			value:       basetypes.NewStringValue("443"),
			expected:    basetypes.NewBoolValue(true),
			expectError: false,
		},
		{
			name:        "valid port 1",
			value:       basetypes.NewStringValue("1"),
			expected:    basetypes.NewBoolValue(true),
			expectError: false,
		},
		{
			name:        "valid port 65535",
			value:       basetypes.NewStringValue("65535"),
			expected:    basetypes.NewBoolValue(true),
			expectError: false,
		},
		{
			name:        "port 0 invalid",
			value:       basetypes.NewStringValue("0"),
			expected:    nil,
			expectError: true,
		},
		{
			name:        "port 65536 invalid",
			value:       basetypes.NewStringValue("65536"),
			expected:    nil,
			expectError: true,
		},
		{
			name:        "negative port",
			value:       basetypes.NewStringValue("-1"),
			expected:    nil,
			expectError: true,
		},
		{
			name:        "not a number",
			value:       basetypes.NewStringValue("abc"),
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
			fn := NewPortNumberFunction()
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
