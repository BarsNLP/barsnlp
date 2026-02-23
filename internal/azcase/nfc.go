package azcase

import "strings"

// ComposeNFC replaces known NFD decomposed sequences for the 6 Azerbaijani
// letters with diacritics: ö, ü, ç, ş, ğ, İ.
// This is NOT full Unicode NFC — only Azerbaijani-specific pairs.
// For full NFC, preprocess with golang.org/x/text/unicode/norm externally.
func ComposeNFC(s string) string {
	// Fast path: scan for combining marks U+0306, U+0307, U+0308, U+0327.
	hasCombiner := false
	for _, r := range s {
		if r == 0x0306 || r == 0x0307 || r == 0x0308 || r == 0x0327 {
			hasCombiner = true
			break
		}
	}
	if !hasCombiner {
		return s
	}

	// Slow path: compose known Azerbaijani decomposed pairs.
	// Lowercase
	s = strings.ReplaceAll(s, "o\u0308", "\u00f6") // o + diaeresis -> ö
	s = strings.ReplaceAll(s, "u\u0308", "\u00fc") // u + diaeresis -> ü
	s = strings.ReplaceAll(s, "c\u0327", "\u00e7") // c + cedilla   -> ç
	s = strings.ReplaceAll(s, "s\u0327", "\u015f") // s + cedilla   -> ş
	s = strings.ReplaceAll(s, "g\u0306", "\u011f") // g + breve     -> ğ
	// Uppercase
	s = strings.ReplaceAll(s, "O\u0308", "\u00d6") // O + diaeresis -> Ö
	s = strings.ReplaceAll(s, "U\u0308", "\u00dc") // U + diaeresis -> Ü
	s = strings.ReplaceAll(s, "C\u0327", "\u00c7") // C + cedilla   -> Ç
	s = strings.ReplaceAll(s, "S\u0327", "\u015e") // S + cedilla   -> Ş
	s = strings.ReplaceAll(s, "G\u0306", "\u011e") // G + breve     -> Ğ
	s = strings.ReplaceAll(s, "I\u0307", "\u0130") // I + dot above -> İ
	return s
}
