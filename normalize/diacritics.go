package normalize

import (
	"strings"
	"unicode"

	"github.com/az-ai-labs/az-lang-nlp/morph"
)

// maxSubstitutablePositions caps the variant generation to avoid
// combinatorial explosion. 2^10 = 1024 candidates max.
const maxSubstitutablePositions = 10

// maxWordBytes is the maximum byte length for a word to attempt
// diacritic restoration on. Matches the morph package limit.
const maxWordBytes = 256

// asciiToDiacritic maps lowercase ASCII characters to their possible
// Azerbaijani diacritic equivalents. Applied AFTER Turkic-aware lowercasing.
//
// The 'i' -> 'ı' mapping is intentionally excluded:
// after azLower, 'i' means confirmed dotted-i (from lowercase 'i' or uppercase 'İ'),
// and 'ı' means confirmed dotless-i (from uppercase 'I').
// Both are already correct after Turkic lowering.
var asciiToDiacritic = [128]rune{
	'e': '\u0259', // ə
	'o': '\u00f6', // ö
	'u': '\u00fc', // ü
	'g': '\u011f', // ğ
	'c': '\u00e7', // ç
	's': '\u015f', // ş
}

// hasDiacriticAlt reports whether the rune has a diacritic alternative.
func hasDiacriticAlt(r rune) bool {
	return r < 128 && asciiToDiacritic[r] != 0
}

// restoreWord attempts to restore diacritics on a single word.
// Returns the original word unchanged if:
//   - the word is empty or too long
//   - it has no substitutable characters
//   - it is already a known dictionary stem
//   - it has too many substitutable positions
//   - zero or multiple dictionary matches (ambiguous/unknown)
func restoreWord(word string) string {
	if word == "" || len(word) > maxWordBytes {
		return word
	}

	lowered := toLower(word)

	// Find substitutable positions in the lowered form.
	runes := []rune(lowered)
	var positions []int
	for i, r := range runes {
		if hasDiacriticAlt(r) {
			positions = append(positions, i)
		}
	}

	if len(positions) == 0 {
		return word
	}

	// If the ASCII form is already a known stem, do not modify it.
	// This prevents changing valid words like "ac" (hungry) to "aç" (open).
	if morph.IsKnownStem(lowered) {
		return word
	}

	if len(positions) > maxSubstitutablePositions {
		return word
	}

	// Generate variants lazily and check against dictionary.
	// Short-circuit on second match (ambiguous).
	totalVariants := 1 << len(positions)
	var matchRunes []rune
	matchCount := 0

	candidate := make([]rune, len(runes))
	for mask := 1; mask < totalVariants; mask++ {
		copy(candidate, runes)
		for bit, pos := range positions {
			if mask&(1<<bit) != 0 {
				candidate[pos] = asciiToDiacritic[runes[pos]]
			}
		}

		if morph.IsKnownStem(string(candidate)) {
			matchCount++
			if matchCount == 1 {
				matchRunes = make([]rune, len(candidate))
				copy(matchRunes, candidate)
			} else {
				// Ambiguous: two or more matches.
				return word
			}
		}
	}

	if matchCount != 1 {
		return word
	}

	// Restore the original case pattern onto the matched runes.
	return restoreCase(word, matchRunes)
}

// restoreCase applies the case pattern from the original word to the
// restored runes. Original uppercase positions become uppercase in the output.
func restoreCase(original string, restored []rune) string {
	origRunes := []rune(original)
	if len(origRunes) != len(restored) {
		return original
	}

	var b strings.Builder
	b.Grow(len(original) + len(origRunes)) // diacritics may use more bytes
	for i, r := range restored {
		if unicode.IsUpper(origRunes[i]) {
			b.WriteRune(azUpper(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// azLower returns the Azerbaijani-aware lowercase form of r.
// Handles the dotted/dotless I distinction:
//   - I (U+0049) -> ı (U+0131, dotless small i)
//   - İ (U+0130, dotted capital I) -> i (U+0069)
func azLower(r rune) rune {
	switch r {
	case 'I':
		return '\u0131' // I -> ı
	case '\u0130':
		return 'i' // İ -> i
	default:
		return unicode.ToLower(r)
	}
}

// azUpper returns the Azerbaijani-aware uppercase form of r.
// Handles the dotted/dotless I distinction:
//   - i (U+0069) -> İ (U+0130, dotted capital I)
//   - ı (U+0131, dotless small i) -> I (U+0049)
func azUpper(r rune) rune {
	switch r {
	case 'i':
		return '\u0130' // i -> İ
	case '\u0131':
		return 'I' // ı -> I
	default:
		return unicode.ToUpper(r)
	}
}

// toLower returns s with Azerbaijani-aware lowercasing applied to every rune.
func toLower(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		b.WriteRune(azLower(r))
	}
	return b.String()
}
