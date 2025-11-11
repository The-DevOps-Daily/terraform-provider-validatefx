package validators

import "testing"

func TestPasswordStrength(t *testing.T) {
	cases := []struct {
		password string
		wantErr  bool
	}{
		{"Abc@1234", false},
		{"abc123", true},
		{"ABC123@", true},
		{"Abcdefgh", true},
		{"Ab1!", true},
	}

	for _, c := range cases {
		err := PasswordStrength(c.password)
		if (err != nil) != c.wantErr {
			t.Errorf("PasswordStrength(%q) = %v, wantErr %v", c.password, err, c.wantErr)
		}
	}
}
