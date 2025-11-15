package validators

import (
	"math/rand"
	"testing"
)

// FuzzSetEqualsValidator validates robustness of SetEqualsValidator against arbitrary lists.
// We build two lists that should be equal as sets (with randomized order/duplicates),
// and a third that differs by at least one element, then check Validate results.
func FuzzSetEqualsValidator(f *testing.F) {
	f.Add(int64(1))
	f.Add(int64(42))
	f.Add(int64(2025))

	f.Fuzz(func(t *testing.T, seed int64) {
		t.Parallel()
		rng := rand.New(rand.NewSource(seed))
		base := []string{"a", "b", "c", "d", "e"}

		// pick a random subset of base as expected
		expected := make([]string, 0, len(base))
		for _, v := range base {
			if rng.Intn(2) == 1 {
				expected = append(expected, v)
			}
		}
		if len(expected) == 0 {
			expected = []string{"a"}
		}

		v := NewSetEquals(expected)

		// Build listA equal to expected as a set, with random duplicates/order
		listA := make([]string, 0, 10)
		for i := 0; i < 10; i++ {
			listA = append(listA, expected[rng.Intn(len(expected))])
		}

		// Build listB that differs by toggling one element
		listB := append([]string{}, listA...)
		// flip an element to something not in expected if possible
		pool := []string{"x", "y", "z"}
		listB[rng.Intn(len(listB))] = pool[rng.Intn(len(pool))]

		if err := v.Validate(listA); err != nil {
			t.Fatalf("expected listA to be equal set, got error: %v", err)
		}
		if err := v.Validate(listB); err == nil {
			t.Fatalf("expected listB to fail set equality validation")
		}
	})
}
