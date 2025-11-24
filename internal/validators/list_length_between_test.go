package validators

import (
	"context"
	"testing"
)

func TestListLengthBetween(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		min         int
		max         int
		list        []string
		expectError bool
	}{
		// Valid lengths
		{
			name:        "length at minimum",
			min:         2,
			max:         5,
			list:        []string{"a", "b"},
			expectError: false,
		},
		{
			name:        "length at maximum",
			min:         2,
			max:         5,
			list:        []string{"a", "b", "c", "d", "e"},
			expectError: false,
		},
		{
			name:        "length in middle",
			min:         2,
			max:         5,
			list:        []string{"a", "b", "c"},
			expectError: false,
		},
		{
			name:        "exact length when min equals max",
			min:         3,
			max:         3,
			list:        []string{"a", "b", "c"},
			expectError: false,
		},
		{
			name:        "zero minimum with zero length",
			min:         0,
			max:         5,
			list:        []string{},
			expectError: false,
		},
		{
			name:        "single element list",
			min:         1,
			max:         1,
			list:        []string{"a"},
			expectError: false,
		},

		// Invalid lengths - too short
		{
			name:        "length below minimum",
			min:         2,
			max:         5,
			list:        []string{"a"},
			expectError: true,
		},
		{
			name:        "empty list with positive minimum",
			min:         1,
			max:         5,
			list:        []string{},
			expectError: true,
		},

		// Invalid lengths - too long
		{
			name:        "length above maximum",
			min:         2,
			max:         5,
			list:        []string{"a", "b", "c", "d", "e", "f"},
			expectError: true,
		},
		{
			name:        "length way above maximum",
			min:         1,
			max:         3,
			list:        []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			validator := NewListLengthBetween(tc.min, tc.max)
			err := validator.Validate(tc.list)

			if tc.expectError && err == nil {
				t.Fatalf("expected error for list length %d with min=%d, max=%d, got none", len(tc.list), tc.min, tc.max)
			}

			if !tc.expectError && err != nil {
				t.Fatalf("expected no error for list length %d with min=%d, max=%d, got: %v", len(tc.list), tc.min, tc.max, err)
			}
		})
	}
}

func TestListLengthBetweenDescription(t *testing.T) {
	t.Parallel()

	validator := NewListLengthBetween(2, 5)
	desc := validator.Description(context.Background())

	if desc == "" {
		t.Fatal("expected non-empty description")
	}

	markdownDesc := validator.MarkdownDescription(context.Background())
	if markdownDesc == "" {
		t.Fatal("expected non-empty markdown description")
	}
}
