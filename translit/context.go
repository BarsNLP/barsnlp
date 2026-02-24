package translit

import (
	"strings"
	"unicode"
)

// containsGje reports whether s contains Ҝ or ҝ anywhere.
// If present, the text uses Soviet orthography where Ҝ=G and Г=Q unambiguously.
func containsGje(s string) bool {
	return strings.ContainsAny(s, "Ҝҝ")
}

// resolveG returns the Latin rune for Cyrillic Г/г based on context.
//
// When hasGje is true (text contains Ҝ/ҝ), Г always maps to Q.
// Otherwise, lookahead skips non-letters to find the next letter:
//   - front vowel (ә, е, и, ө, ү) → G
//   - back vowel (а, о, у, ы), consonant, or end of string → Q
func resolveG(upper bool, rest string, hasGje bool) rune {
	if hasGje {
		if upper {
			return 'Q'
		}
		return 'q'
	}

	for _, r := range rest {
		if !unicode.IsLetter(r) {
			continue
		}
		if frontVowels[r] {
			if upper {
				return 'G'
			}
			return 'g'
		}
		// Back vowel or consonant.
		if upper {
			return 'Q'
		}
		return 'q'
	}

	// No letter found after Г (end of string).
	if upper {
		return 'Q'
	}
	return 'q'
}
