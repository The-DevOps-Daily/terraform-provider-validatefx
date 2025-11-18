package functions

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestURIFunction(t *testing.T) {
	tests := []struct {
		name        string
		value       basetypes.StringValue
		expected    attr.Value
		expectError bool
	}{
		{
			name:        "valid http URI",
			value:       basetypes.NewStringValue("http://example.com"),
			expected:    basetypes.NewBoolValue(true),
			expectError: false,
		},
		{
			name:        "valid https URI",
			value:       basetypes.NewStringValue("https://example.com/path"),
			expected:    basetypes.NewBoolValue(true),
			expectError: false,
		},
		{
			name:        "valid ftp URI",
			value:       basetypes.NewStringValue("ftp://files.example.com"),
			expected:    basetypes.NewBoolValue(true),
			expectError: false,
		},
		{
			name:        "valid URI with query",
			value:       basetypes.NewStringValue("https://example.com/search?q=test"),
			expected:    basetypes.NewBoolValue(true),
			expectError: false,
		},
		{
			name:        "invalid missing scheme",
			value:       basetypes.NewStringValue("example.com"),
			expected:    nil,
			expectError: true,
		},
		{
			name:        "invalid empty",
			value:       basetypes.NewStringValue(""),
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
			fn := NewURIFunction()
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
