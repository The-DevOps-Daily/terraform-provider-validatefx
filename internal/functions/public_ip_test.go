package functions

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestPublicIPFunction(t *testing.T) {
	tests := []struct {
		name             string
		value            basetypes.StringValue
		excludeLinkLocal basetypes.BoolValue
		excludeReserved  basetypes.BoolValue
		expected         attr.Value
		expectError      bool
	}{
		{
			name:             "valid public IP",
			value:            basetypes.NewStringValue("8.8.8.8"),
			excludeLinkLocal: basetypes.NewBoolValue(false),
			excludeReserved:  basetypes.NewBoolValue(false),
			expected:         basetypes.NewBoolValue(true),
			expectError:      false,
		},
		{
			name:             "private IP",
			value:            basetypes.NewStringValue("192.168.1.1"),
			excludeLinkLocal: basetypes.NewBoolValue(false),
			excludeReserved:  basetypes.NewBoolValue(false),
			expected:         nil,
			expectError:      true,
		},
		{
			name:             "link-local with exclusion",
			value:            basetypes.NewStringValue("169.254.1.1"),
			excludeLinkLocal: basetypes.NewBoolValue(true),
			excludeReserved:  basetypes.NewBoolValue(false),
			expected:         nil,
			expectError:      true,
		},
		{
			name:             "public IP with both flags",
			value:            basetypes.NewStringValue("8.8.4.4"),
			excludeLinkLocal: basetypes.NewBoolValue(false),
			excludeReserved:  basetypes.NewBoolValue(true),
			expected:         basetypes.NewBoolValue(true),
			expectError:      false,
		},
		{
			name:             "reserved IP with exclusion",
			value:            basetypes.NewStringValue("255.255.255.255"),
			excludeLinkLocal: basetypes.NewBoolValue(false),
			excludeReserved:  basetypes.NewBoolValue(true),
			expected:         nil,
			expectError:      true,
		},
		{
			name:             "IPv6 public",
			value:            basetypes.NewStringValue("2001:4860:4860::8888"),
			excludeLinkLocal: basetypes.NewBoolValue(false),
			excludeReserved:  basetypes.NewBoolValue(false),
			expected:         basetypes.NewBoolValue(true),
			expectError:      false,
		},
		{
			name:             "IPv6 link-local with exclusion",
			value:            basetypes.NewStringValue("fe80::1"),
			excludeLinkLocal: basetypes.NewBoolValue(true),
			excludeReserved:  basetypes.NewBoolValue(false),
			expected:         nil,
			expectError:      true,
		},
		{
			name:             "null value",
			value:            basetypes.NewStringNull(),
			excludeLinkLocal: basetypes.NewBoolValue(false),
			excludeReserved:  basetypes.NewBoolValue(false),
			expected:         basetypes.NewBoolUnknown(),
			expectError:      false,
		},
		{
			name:             "unknown value",
			value:            basetypes.NewStringUnknown(),
			excludeLinkLocal: basetypes.NewBoolValue(false),
			excludeReserved:  basetypes.NewBoolValue(false),
			expected:         basetypes.NewBoolUnknown(),
			expectError:      false,
		},
		{
			name:             "unknown flag",
			value:            basetypes.NewStringValue("8.8.8.8"),
			excludeLinkLocal: basetypes.NewBoolUnknown(),
			excludeReserved:  basetypes.NewBoolValue(false),
			expected:         basetypes.NewBoolUnknown(),
			expectError:      false,
		},
		{
			name:             "invalid IP",
			value:            basetypes.NewStringValue("not-an-ip"),
			excludeLinkLocal: basetypes.NewBoolValue(false),
			excludeReserved:  basetypes.NewBoolValue(false),
			expected:         nil,
			expectError:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fn := NewPublicIPFunction()
			req := function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{tc.value, tc.excludeLinkLocal, tc.excludeReserved}),
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

func TestPublicIPFunctionMetadata(t *testing.T) {
	fn := NewPublicIPFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "public_ip" {
		t.Errorf("expected name 'public_ip', got %q", resp.Name)
	}
}

func TestPublicIPFunctionDefinition(t *testing.T) {
	fn := NewPublicIPFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Return == nil {
		t.Fatal("expected return definition, got nil")
	}
	if len(resp.Definition.Parameters) != 3 {
		t.Errorf("expected 3 parameters, got %d", len(resp.Definition.Parameters))
	}
}
