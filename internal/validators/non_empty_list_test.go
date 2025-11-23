package validators

import (
	"testing"
)

func TestNonEmptyListValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		values    []string
		wantError bool
	}{
		// Valid scenarios - non-empty lists
		{"single element", []string{"a"}, false},
		{"multiple elements", []string{"a", "b", "c"}, false},
		{"many elements", []string{"1", "2", "3", "4", "5"}, false},

		// Invalid scenarios - empty list
		{"empty list", []string{}, true},
		{"nil list", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := NewNonEmptyList()
			err := v.Validate(tt.values)

			if tt.wantError && err == nil {
				t.Errorf("expected error, got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestNonEmptyListValidator_NilValidator(t *testing.T) {
	t.Parallel()

	var v *NonEmptyListValidator
	err := v.Validate([]string{"a"})

	if err == nil {
		t.Error("expected error for nil validator")
	}
}
