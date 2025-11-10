package validators

import "testing"

func TestSetEqualsValidator(t *testing.T) {
	t.Parallel()

	v := NewSetEquals([]string{"a", "b", "c"})

	tests := []struct {
		name    string
		input   []string
		wantErr bool
	}{
		{name: "same order", input: []string{"a", "b", "c"}, wantErr: false},
		{name: "different order", input: []string{"c", "b", "a"}, wantErr: false},
		{name: "with duplicates", input: []string{"a", "a", "b", "c", "c"}, wantErr: false},
		{name: "missing element", input: []string{"a", "b"}, wantErr: true},
		{name: "extra element", input: []string{"a", "b", "c", "d"}, wantErr: true},
		{name: "completely different", input: []string{"x", "y"}, wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := v.Validate(tc.input)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
