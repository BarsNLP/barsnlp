// Unexported conversion functions for Azerbaijani number-to-text conversion.
package numtext

import (
	"strconv"
	"strings"
)

const (
	growGroup   = 32   // estimated bytes for a 0-999 group
	growFloat   = 128  // estimated bytes for a decimal conversion
	maxDenomFD  = 3    // max fractional digits with a named denominator
	asciiRuneHi = 0x80 // upper bound for single-byte UTF-8 runes
)

// convert converts an int64 to Azerbaijani cardinal text.
// Returns "" if abs(n) exceeds maxAbs.
func convert(n int64) string {
	if n > maxAbs || n < -maxAbs {
		return ""
	}
	if n == 0 {
		return wordZero
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var parts []string

	for _, mag := range magnitudes {
		count := n / mag.value
		if count > 0 {
			// "bir min" → "min" (omit "bir" before "min" only)
			if mag.value == 1_000 && count == 1 {
				parts = append(parts, mag.word)
			} else {
				parts = append(parts, convertGroup(count)+" "+mag.word)
			}
			n %= mag.value
		}
	}

	if n > 0 {
		parts = append(parts, convertGroup(n))
	}

	result := strings.Join(parts, " ")

	if negative {
		return wordNegative + " " + result
	}
	return result
}

// convertGroup converts a number in [0, 999] to Azerbaijani text.
// Returns "" for 0; callers handle the zero case themselves.
func convertGroup(n int64) string {
	if n == 0 {
		return ""
	}

	var b strings.Builder
	b.Grow(growGroup)

	h := n / hundred
	if h == 1 {
		b.WriteString(wordHundred)
	} else if h > 1 {
		b.WriteString(ones[h])
		b.WriteByte(' ')
		b.WriteString(wordHundred)
	}

	r := n % hundred
	t := r / 10
	o := r % 10

	if t > 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(tens[t])
	}

	if o > 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(ones[o])
	}

	return b.String()
}

// convertOrdinal converts an int64 to Azerbaijani ordinal text.
// Returns "" if abs(n) exceeds maxAbs.
func convertOrdinal(n int64) string {
	if n > maxAbs || n < -maxAbs {
		return ""
	}

	negative := n < 0
	absN := n
	if negative {
		absN = -n
	}

	cardinal := convert(absN)
	if cardinal == "" {
		return ""
	}

	lv := lastVowel(cardinal)
	if lv == 0 {
		return ""
	}

	var suffix string
	lastRune, _ := lastRuneOf(cardinal)
	if isVowel(lastRune) {
		suffix = ordinalShort[lv]
	} else {
		suffix = ordinalFull[lv]
	}

	result := cardinal + suffix
	if negative {
		return wordNegative + " " + result
	}
	return result
}

// convertFloat converts a decimal number string to Azerbaijani text using
// the given Mode (MathMode or DigitMode).
func convertFloat(s string, mode Mode) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	negative := false
	switch s[0] {
	case '-':
		negative = true
		s = s[1:]
	case '+':
		s = s[1:]
	}

	sepIdx := strings.IndexAny(s, ".,")

	if sepIdx == -1 {
		// No decimal separator; treat as plain integer.
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return ""
		}
		if negative {
			val = -val
		}
		return convert(val)
	}

	wholePart := s[:sepIdx]
	fracPart := s[sepIdx+1:]

	if (wholePart != "" && !allDigits(wholePart)) || !allDigits(fracPart) || fracPart == "" {
		return ""
	}

	var wholeVal int64
	if wholePart == "" {
		wholeVal = 0
	} else {
		var err error
		wholeVal, err = strconv.ParseInt(wholePart, 10, 64)
		if err != nil {
			return ""
		}
	}

	wholeText := convert(wholeVal)
	if wholeText == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(growFloat)

	if negative {
		b.WriteString(wordNegative)
		b.WriteByte(' ')
	}
	b.WriteString(wholeText)

	switch mode {
	case MathMode:
		fracDigits := len(fracPart)

		// Parse fractional part as integer (leading zeros are significant for
		// the denominator but not for the numerator value).
		numeratorVal, err := strconv.ParseInt(fracPart, 10, 64)
		if err != nil {
			return ""
		}
		numeratorText := convert(numeratorVal)
		if numeratorText == "" {
			return ""
		}

		var denomWord string
		if fracDigits <= maxDenomFD {
			denomWord = denominators[fracDigits]
		} else {
			// Compose denominator for fracDigits > 3: convert(10^fracDigits) + locative suffix.
			denomBase := powerOf10Text(fracDigits)
			if denomBase == "" {
				return ""
			}
			denomWord = denomBase + locativeSuffix(denomBase)
		}

		b.WriteByte(' ')
		b.WriteString(wordExact)
		b.WriteByte(' ')
		b.WriteString(denomWord)
		b.WriteByte(' ')
		b.WriteString(numeratorText)

	case DigitMode:
		b.WriteByte(' ')
		b.WriteString(wordComma)
		for _, ch := range fracPart {
			d := int(ch - '0')
			b.WriteByte(' ')
			b.WriteString(ones[d])
		}
	}

	return b.String()
}

// powerOf10Text returns the Azerbaijani text for 10^exp.
// Used for composing denominators beyond 3 fractional digits.
func powerOf10Text(exp int) string {
	// Build 10^exp as int64 when possible, otherwise it exceeds maxAbs.
	var val int64 = 1
	for range exp {
		val *= 10
		if val > maxAbs {
			return ""
		}
	}
	return convert(val)
}

// lastVowel scans s backwards and returns the last rune that is an Azerbaijani vowel.
// Returns 0 if no vowel is found.
func lastVowel(s string) rune {
	for i := len(s); i > 0; {
		r, size := lastRuneAt(s, i)
		i -= size
		if isVowel(r) {
			return r
		}
	}
	return 0
}

// isVowel reports whether r is an Azerbaijani vowel.
func isVowel(r rune) bool {
	return strings.ContainsRune(azVowels, r)
}

// locativeSuffix returns the Azerbaijani locative case suffix ("da" or "də")
// based on vowel harmony of the last vowel in s.
// Back vowels (a, ı, o, u) → "da"; front vowels (e, ə, i, ö, ü) → "də".
func locativeSuffix(s string) string {
	lv := lastVowel(s)
	switch lv {
	case 'a', 'ı', 'o', 'u':
		return "da"
	default:
		return "də"
	}
}

// lastRuneOf returns the last rune in s and its byte size.
func lastRuneOf(s string) (rune, int) {
	return lastRuneAt(s, len(s))
}

// lastRuneAt decodes the rune that ends at byte position end in s.
func lastRuneAt(s string, end int) (rune, int) {
	// Walk back up to 4 bytes to find a valid UTF-8 start byte.
	for size := 1; size <= 4 && size <= end; size++ {
		r, n := decodeRuneAt(s, end-size)
		if n == size {
			return r, size
		}
	}
	return rune(s[end-1]), 1
}

// decodeRuneAt decodes the rune at byte position i in s.
func decodeRuneAt(s string, i int) (rune, int) {
	// Use a small slice trick to leverage the runtime UTF-8 decoder.
	r, size := rune(s[i]), 1
	if r < asciiRuneHi {
		return r, 1
	}
	// Multi-byte: range over a slice starting at i.
	for _, rv := range s[i:] {
		return rv, len(string(rv))
	}
	return r, size
}

// allDigits reports whether s consists entirely of ASCII digit characters.
// An empty string returns false.
func allDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
