package utils

import "testing"

func TestToAPISeason(t *testing.T) {
	tests := map[string]string{
		"23/24":    "20232024",
		"20232024": "20232024",
		"99/00":    "19992000",
	}

	for input, expected := range tests {
		actual, err := ToAPISeason(input)
		if err != nil {
			t.Fatalf("ToAPISeason(%q) returned error: %v", input, err)
		}

		if actual != expected {
			t.Fatalf("ToAPISeason(%q) = %q, want %q", input, actual, expected)
		}
	}
}

func TestToDisplaySeason(t *testing.T) {
	tests := map[string]string{
		"23/24":    "23/24",
		"20232024": "23/24",
	}

	for input, expected := range tests {
		actual, err := ToDisplaySeason(input)
		if err != nil {
			t.Fatalf("ToDisplaySeason(%q) returned error: %v", input, err)
		}

		if actual != expected {
			t.Fatalf("ToDisplaySeason(%q) = %q, want %q", input, actual, expected)
		}
	}
}
