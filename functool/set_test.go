package functool

import (
	"testing"
)

// NewStringSetWithValue
func TestNewStringSetWithValue(t *testing.T) {
	names := []string{"zhang", "wang", "sun"}
	ss := NewStringSetWithValue(names)
	ss.Remove("zhang")
	ss.Remove("sun")
	ss.Remove("wang")
	ss.Add("li")

	if ss.Len() != 1 {
		t.Fatalf("Unexpected set length %d != 1", ss.Len())
	}
}
