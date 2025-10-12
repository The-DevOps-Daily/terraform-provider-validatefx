package validatefx

import "testing"

func TestIsJSON(t *testing.T) {
    validJSON := `{"name": "Krishna", "age": 20}`
    invalidJSON := `{"name": "Krishna", "age": 20` // missing closing }

    if err := IsJSON(validJSON); err != nil {
        t.Errorf("Expected valid JSON, got error: %v", err)
    }

    if err := IsJSON(invalidJSON); err == nil {
        t.Errorf("Expected error for invalid JSON, got nil")
    }
}
