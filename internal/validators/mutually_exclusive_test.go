package validators

import (
	"testing"
)

func TestMutuallyExclusiveValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		values    []string
		wantError bool
	}{
		// Valid scenarios - exactly one value set
		{"first set only", []string{"value1", "", ""}, false},
		{"second set only", []string{"", "value2", ""}, false},
		{"third set only", []string{"", "", "value3"}, false},
		{"middle set only", []string{"", "value", "", ""}, false},

		// Invalid scenarios - zero or multiple values set
		{"none set", []string{"", "", ""}, true},
		{"all empty", []string{}, true},
		{"two set", []string{"value1", "value2", ""}, true},
		{"all set", []string{"value1", "value2", "value3"}, true},
		{"multiple non-empty", []string{"a", "b", "", "c"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := NewMutuallyExclusive()
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

func TestMutuallyExclusiveValidator_NilValidator(t *testing.T) {
	t.Parallel()

	var v *MutuallyExclusiveValidator
	err := v.Validate([]string{"a"})

	if err == nil {
		t.Error("expected error for nil validator")
	}
}
