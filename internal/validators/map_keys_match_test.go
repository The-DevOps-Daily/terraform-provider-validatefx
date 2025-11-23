package validators

import (
	"testing"
)

func TestMapKeysMatchValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		allowedKeys  []string
		requiredKeys []string
		inputKeys    []string
		wantError    bool
	}{
		// Valid scenarios
		{"all keys allowed", []string{"a", "b", "c"}, nil, []string{"a", "b"}, false},
		{"exact match", []string{"a", "b"}, nil, []string{"a", "b"}, false},
		{"required keys present", []string{"a", "b", "c"}, []string{"a"}, []string{"a", "b"}, false},
		{"only required keys", nil, []string{"a", "b"}, []string{"a", "b", "c"}, false},
		{"empty allowed means all allowed", []string{}, []string{"a"}, []string{"a", "b", "c"}, false},

		// Invalid scenarios
		{"disallowed key", []string{"a", "b"}, nil, []string{"a", "c"}, true},
		{"missing required key", []string{"a", "b", "c"}, []string{"a", "b"}, []string{"a"}, true},
		{"all keys missing", []string{"a", "b"}, []string{"a", "b"}, []string{}, true},
		{"disallowed and missing required", []string{"a", "b"}, []string{"a"}, []string{"c"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := NewMapKeysMatch(tt.allowedKeys, tt.requiredKeys)
			err := v.Validate(tt.inputKeys)

			if tt.wantError && err == nil {
				t.Errorf("expected error, got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMapKeysMatchValidator_NilValidator(t *testing.T) {
	t.Parallel()

	var v *MapKeysMatchValidator
	err := v.Validate([]string{"a"})

	if err == nil {
		t.Error("expected error for nil validator")
	}
}
