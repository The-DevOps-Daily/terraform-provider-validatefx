package functions

import "testing"

func TestPasswordStrengthFunc(t *testing.T) {
	ok, err := NewPasswordStrengthFunction().(*passwordStrengthFunction).CallTest("Abc@1234")
	if !ok || err != nil {
		t.Errorf("expected valid password, got error: %v", err)
	}

	ok, err = NewPasswordStrengthFunction().(*passwordStrengthFunction).CallTest("abc")
	if ok || err == nil {
		t.Errorf("expected invalid password but got no error")
	}
}

// helper to simulate Call
func (f *passwordStrengthFunction) CallTest(password string) (bool, error) {
	err := validators.PasswordStrength(password)
	return err == nil, err
}
