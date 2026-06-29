package cmd

import "testing"

func TestFormatEnvfileValueRoundTripsThroughImportParser(t *testing.T) {
	values := []string{
		"contains spaces",
		"contains 'single quotes'",
		`contains "quotes"`,
		"contains $dollars",
		"contains # hash",
		`contains \ backslashes`,
		"contains literal #personal marker",
		`all together: spaces "quotes" $dollars #hash \slashes #personal`,
		`fallback: 'single' "double" C:\tmp\path #hash #personal`,
		"fallback with `backticks`: 'single' \"double\" #personal",
		`trailing: 'single' and backslash\`,
	}

	for _, want := range values {
		key, got, personal, ok := parseEnvfileLine("KEY=" + formatEnvfileValue(want))
		if !ok {
			t.Fatalf("parseEnvfileLine did not parse formatted value %q", want)
		}
		if key != "KEY" {
			t.Fatalf("key mismatch: got %q, want KEY", key)
		}
		if personal {
			t.Fatalf("formatted shared value %q was marked personal", want)
		}
		if got != want {
			t.Fatalf("value mismatch: got %q, want %q", got, want)
		}
	}
}

func TestFormatEnvfileValueUsesDotenvQuotedLiterals(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"hello world", "'hello world'"},
		{`contains "quotes"`, `'contains "quotes"'`},
		{"pa$word", "'pa$word'"},
		{"abc # def", "'abc # def'"},
		{`C:\tmp\path`, `'C:\tmp\path'`},
		{"foo #personal", "'foo #personal'"},
		{`has 'single' only`, `"has 'single' only"`},
		{`has 'single' and C:\tmp\path`, "`has 'single' and C:\\tmp\\path`"},
		{`has 'single' and trailing\`, "`has 'single' and trailing\\`"},
		{`has 'single' and "double"`, "`has 'single' and \"double\"`"},
		{"has 'single', \"double\", and `backticks`", "`has 'single', \"double\", and `backticks``"},
	}

	for _, tt := range tests {
		got := formatEnvfileValue(tt.value)
		if got != tt.want {
			t.Fatalf("formatEnvfileValue(%q) = %q, want %q", tt.value, got, tt.want)
		}
	}
}

func TestQuotedPersonalMarkerStaysInSharedValue(t *testing.T) {
	want := "foo #personal"

	_, got, personal, ok := parseEnvfileLine("KEY=" + formatEnvfileValue(want))
	if !ok {
		t.Fatal("parseEnvfileLine did not parse formatted value")
	}
	if personal {
		t.Fatal("quoted shared value was marked personal")
	}
	if got != want {
		t.Fatalf("value mismatch: got %q, want %q", got, want)
	}
}

func TestTrailingPersonalAnnotationMarksPersonalValue(t *testing.T) {
	want := "foo #personal"

	_, got, personal, ok := parseEnvfileLine("KEY=" + formatEnvfileValue(want) + " #personal")
	if !ok {
		t.Fatal("parseEnvfileLine did not parse formatted personal value")
	}
	if !personal {
		t.Fatal("trailing #personal annotation was not detected")
	}
	if got != want {
		t.Fatalf("value mismatch: got %q, want %q", got, want)
	}
}
