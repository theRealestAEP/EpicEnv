package cmd

import (
	"strings"
	"testing"
)

func TestReadInputValueFromPipePreservesSingleLineWithoutTrailingNewline(t *testing.T) {
	expected := strings.Repeat("A", 1700)

	got, err := readInputValue(strings.NewReader(expected), 0, false)
	if err != nil {
		t.Fatalf("readInputValue returned error: %v", err)
	}

	if got != expected {
		t.Fatalf("value mismatch: got length %d, want length %d", len(got), len(expected))
	}
}
