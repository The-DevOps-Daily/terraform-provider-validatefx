package validators

import (
	"context"
	"testing"
)

func TestListUnique(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		list        []string
		expectError bool
	}{
		// Valid - all unique
		{
			name:        "all unique elements",
			list:        []string{"a", "b", "c"},
			expectError: false,
		},
		{
			name:        "single element",
			list:        []string{"a"},
			expectError: false,
		},
		{
			name:        "empty list",
			list:        []string{},
			expectError: false,
		},
		{
			name:        "many unique elements",
			list:        []string{"apple", "banana", "cherry", "date", "elderberry"},
			expectError: false,
		},
		{
			name:        "numbers as strings",
			list:        []string{"1", "2", "3", "4", "5"},
			expectError: false,
		},

		// Invalid - contains duplicates
		{
			name:        "simple duplicate",
			list:        []string{"a", "b", "a"},
			expectError: true,
		},
		{
			name:        "multiple duplicates",
			list:        []string{"a", "b", "a", "b"},
			expectError: true,
		},
		{
			name:        "duplicate at end",
			list:        []string{"a", "b", "c", "a"},
			expectError: true,
		},
		{
			name:        "consecutive duplicates",
			list:        []string{"a", "a"},
			expectError: true,
		},
		{
			name:        "all same elements",
			list:        []string{"a", "a", "a", "a"},
			expectError: true,
		},
		{
			name:        "duplicate numbers",
			list:        []string{"1", "2", "3", "1"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			validator := NewListUnique()
			err := validator.Validate(tc.list)

			if tc.expectError && err == nil {
				t.Fatalf("expected error for list %v, got none", tc.list)
			}

			if !tc.expectError && err != nil {
				t.Fatalf("expected no error for list %v, got: %v", tc.list, err)
			}
		})
	}
}

func TestListUniqueDescription(t *testing.T) {
	t.Parallel()

	validator := NewListUnique()
	desc := validator.Description(context.Background())

	if desc == "" {
		t.Fatal("expected non-empty description")
	}

	markdownDesc := validator.MarkdownDescription(context.Background())
	if markdownDesc == "" {
		t.Fatal("expected non-empty markdown description")
	}
}
