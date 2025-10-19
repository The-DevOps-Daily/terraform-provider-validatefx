package validatefx

import (
    "encoding/json"
    "fmt"
)

// IsJSON checks if the input string is valid JSON
func IsJSON(input string) error {
    var js interface{}
    if err := json.Unmarshal([]byte(input), &js); err != nil {
        return fmt.Errorf("invalid JSON: %v", err)
    }
    return nil
}
