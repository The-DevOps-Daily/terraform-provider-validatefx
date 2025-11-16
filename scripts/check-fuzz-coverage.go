package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    log.SetFlags(0)

    root := "internal/validators"
    entries, err := os.ReadDir(root)
    if err != nil {
        log.Fatalf("failed to read %s: %v", root, err)
    }

    var missing []string
    for _, e := range entries {
        name := e.Name()
        if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
            continue
        }
        base := strings.TrimSuffix(name, ".go")
        fuzz := filepath.Join(root, base+"_fuzz_test.go")
        if _, err := os.Stat(fuzz); err != nil {
            missing = append(missing, base)
        }
    }

    if len(missing) > 0 {
        log.Fatalf("fuzz coverage missing for validators: %s\nAdd *_fuzz_test.go under internal/validators for each.", strings.Join(missing, ", "))
    }
    fmt.Println("Fuzz coverage check passed: all validators have fuzz tests.")
}

