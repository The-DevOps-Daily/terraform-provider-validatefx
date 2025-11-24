package validators

import (
	"context"
	"testing"
)

func TestDependentValue(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		condition      string
		dependent      string
		expectError    bool
	}{
		// Valid cases
		{
			name:        "both empty",
			condition:   "",
			dependent:   "",
			expectError: false,
		},
		{
			name:        "both set",
			condition:   "value1",
			dependent:   "value2",
			expectError: false,
		},
		{
			name:        "condition empty, dependent set",
			condition:   "",
			dependent:   "value",
			expectError: false,
		},
		{
			name:        "both set with different values",
			condition:   "enable_feature",
			dependent:   "feature_config",
			expectError: false,
		},

		// Invalid cases
		{
			name:        "condition set, dependent empty",
			condition:   "value",
			dependent:   "",
			expectError: true,
		},
		{
			name:        "condition true, dependent empty",
			condition:   "true",
			dependent:   "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			validator := NewDependentValue()
			err := validator.Validate(tc.condition, tc.dependent)

			if tc.expectError && err == nil {
				t.Fatalf("expected error for condition=%q, dependent=%q, got none", tc.condition, tc.dependent)
			}

			if !tc.expectError && err != nil {
				t.Fatalf("expected no error for condition=%q, dependent=%q, got: %v", tc.condition, tc.dependent, err)
			}
		})
	}
}

func TestDependentValueDescription(t *testing.T) {
	t.Parallel()

	validator := NewDependentValue()
	desc := validator.Description(context.Background())

	if desc == "" {
		t.Fatal("expected non-empty description")
	}

	markdownDesc := validator.MarkdownDescription(context.Background())
	if markdownDesc == "" {
		t.Fatal("expected non-empty markdown description")
	}
}
