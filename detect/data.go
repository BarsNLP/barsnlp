package detect

import (
	"math"
	"unicode"
)

// --- Character-set maps ---

// azLatinUnique contains characters unique to Azerbaijani Latin script that
// do not appear in Turkish or English. The schwa (ə/Ə) is the single strongest
// discriminator between Azerbaijani and Turkish Latin texts.
var azLatinUnique = map[rune]bool{
	'ə': true, 'Ə': true, // U+0259, U+018F — schwa, exclusive to Azerbaijani Latin
}

// azCyrillicUnique contains characters unique to Azerbaijani Cyrillic that
// are absent from Russian Cyrillic. Presence of any of these characters is a
// strong signal that the text is Azerbaijani Cyrillic.
var azCyrillicUnique = map[rune]bool{
	'ә': true, 'Ә': true, // U+04D9, U+04D8 — schwa (Azerbaijani ə in Cyrillic)
	'ғ': true, 'Ғ': true, // U+0493, U+0492 — Ğ equivalent
	'ҹ': true, 'Ҹ': true, // U+04B9, U+04B8 — C equivalent
	'ҝ': true, 'Ҝ': true, // U+049D, U+049C — G equivalent
	'ө': true, 'Ө': true, // U+04E9, U+04E8 — Ö equivalent
	'ү': true, 'Ү': true, // U+04AF, U+04AE — Ü equivalent
	'һ': true, 'Һ': true, // U+04BB, U+04BA — H equivalent
	'ј': true, 'Ј': true, // U+0458, U+0408 — Y/Je equivalent
}

// ruCyrillicUnique contains characters present in Russian Cyrillic but absent
// from Azerbaijani Cyrillic. These are strong negative signals for Azerbaijani
// and positive signals for Russian.
var ruCyrillicUnique = map[rune]bool{
	'ы': true, 'Ы': true, // U+044B, U+042B — back unrounded vowel, Russian only
	'э': true, 'Э': true, // U+044D, U+042D — open e, Russian only
	'щ': true, 'Щ': true, // U+0449, U+0429 — shcha, Russian only
}

// trAzSharedSpecial contains special Latin characters shared between Turkish
// and Azerbaijani Latin scripts. Their presence indicates a Turkic language
// but does not distinguish Azerbaijani from Turkish on its own.
var trAzSharedSpecial = map[rune]bool{
	'ğ': true, 'Ğ': true, // U+011F, U+011E — soft g
	'ş': true, 'Ş': true, // U+015F, U+015E — sh sound
	'ç': true, 'Ç': true, // U+00E7, U+00C7 — ch sound
	'ö': true, 'Ö': true, // U+00F6, U+00D6 — front rounded o
	'ü': true, 'Ü': true, // U+00FC, U+00DC — front rounded u
	'ı': true, 'İ': true, // U+0131, U+0130 — dotless i / dotted I
}

// azLatinXQ contains Latin letters that are common in Azerbaijani texts but
// rare in Turkish. Q (representing the uvular stop) and X (representing the
// velar fricative) appear frequently in native Azerbaijani vocabulary.
var azLatinXQ = map[rune]bool{
	'x': true, 'X': true, // U+0078, U+0058 — velar fricative, frequent in Azerbaijani
	'q': true, 'Q': true, // U+0071, U+0051 — uvular stop, frequent in Azerbaijani
}

// isCyrillic reports whether r is a Cyrillic letter (Unicode block U+0400..U+04FF).
func isCyrillic(r rune) bool {
	return r >= 0x0400 && r <= 0x04FF
}

// toLowerTurkic returns the Turkic-aware lowercase of r.
// In Azerbaijani and Turkish, I (U+0049) maps to ı (U+0131, dotless)
// and İ (U+0130, dotted) maps to i (U+0069). All other runes use
// standard Unicode lowercasing.
func toLowerTurkic(r rune) rune {
	switch r {
	case 'İ': // U+0130 — dotted capital I → standard lowercase i
		return 'i'
	case 'I': // U+0049 — ASCII capital I → dotless ı in Turkic
		return 'ı'
	default:
		return unicode.ToLower(r)
	}
}

// --- Trigram profiles ---

// trigramSize is the number of consecutive runes in a single trigram.
const trigramSize = 3

// azLatnTrigrams is the top-50 character trigram frequency profile for
// Azerbaijani Latin script, derived from a representative corpus.
// Values are normalized relative frequencies.
var azLatnTrigrams = map[string]float64{
	"lar": 0.012697,
	"lər": 0.010075,
	"dir": 0.006809,
	"arı": 0.006206,
	"əri": 0.005485,
	"ilə": 0.005140,
	"nda": 0.005035,
	"dən": 0.004908,
	"dır": 0.004820,
	"bir": 0.004702,
	"ara": 0.004600,
	"dan": 0.004578,
	"rin": 0.004330,
	"ini": 0.004303,
	"ndə": 0.004287,
	"ind": 0.004253,
	"anı": 0.004096,
	"ələ": 0.003683,
	"edi": 0.003669,
	"nla": 0.003651,
	"ını": 0.003557,
	"ası": 0.003553,
	"lan": 0.003423,
	"əsi": 0.003417,
	"ınd": 0.003275,
	"adı": 0.003207,
	"rın": 0.003150,
	"ala": 0.003137,
	"nın": 0.003100,
	"əni": 0.003041,
	"rdi": 0.002986,
	"alı": 0.002975,
	"idi": 0.002945,
	"dil": 0.002934,
	"iri": 0.002905,
	"miş": 0.002783,
	"əli": 0.002776,
	"ili": 0.002720,
	"nin": 0.002706,
	"əti": 0.002646,
	"ayı": 0.002609,
	"olu": 0.002569,
	"ərə": 0.002536,
	"mış": 0.002500,
	"rdı": 0.002430,
	"sın": 0.002406,
	"ada": 0.002405,
	"mən": 0.002397,
	"şdı": 0.002394,
	"inə": 0.002387,
}

// trTrigrams is the top-50 character trigram frequency profile for Turkish,
// derived from standard Turkish corpus statistics.
// Turkish is agglutinative like Azerbaijani but uses front vowel 'e' instead
// of schwa 'ə', and has distinctive suffixes such as -yor (present continuous),
// -rak (-arak converb), -mak (infinitive), and participial forms with -dığ/-lığ/-ığı.
var trTrigrams = map[string]float64{
	"lar": 0.013200,
	"ler": 0.011500,
	"bir": 0.005100,
	"ile": 0.004900,
	"nda": 0.004800,
	"dan": 0.004600,
	"ını": 0.004400,
	"rin": 0.004300,
	"ara": 0.004200,
	"ini": 0.004100,
	"anı": 0.004000,
	"lan": 0.003900,
	"ind": 0.003800,
	"ala": 0.003700,
	"nin": 0.003600,
	"eri": 0.003500,
	"ili": 0.003400,
	"ası": 0.003300,
	"olu": 0.003200,
	"edi": 0.003100,
	"idi": 0.003000,
	"ınd": 0.002900,
	"arı": 0.002800,
	"alı": 0.002700,
	"dir": 0.002600,
	"sin": 0.002500,
	"yor": 0.002400,
	"ıyo": 0.002300,
	"nde": 0.002200,
	"den": 0.002100,
	"yan": 0.002000,
	"yen": 0.001900,
	"ter": 0.001800,
	"esi": 0.001700,
	"ine": 0.001600,
	"lma": 0.001500,
	"aya": 0.001400,
	"ard": 0.001300,
	"lik": 0.001200,
	"rak": 0.001100,
	"mak": 0.001000,
	"ken": 0.000950,
	"aki": 0.000900,
	"eki": 0.000850,
	"dığ": 0.000800,
	"lığ": 0.000750,
	"ığı": 0.000700,
	"tır": 0.000650,
	"dır": 0.000600,
	"rdi": 0.000550,
}

// --- Trigram functions ---

// extractTrigrams builds a frequency map of character trigrams from s,
// considering only letter runes. Non-letter runes are skipped and do not
// interrupt trigram boundaries — letters are collected into a contiguous
// sequence before sliding the trigram window.
// Letters are lowercased using Turkic-aware rules (İ→i, I→ı) to match
// the lowercase trigram profiles.
func extractTrigrams(s string) map[string]float64 {
	letters := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters = append(letters, toLowerTurkic(r))
		}
	}

	counts := make(map[string]float64)
	limit := len(letters) - trigramSize + 1
	for i := range limit {
		trigram := string(letters[i : i+trigramSize])
		counts[trigram]++
	}

	total := float64(limit)
	if total <= 0 {
		return counts
	}
	for k := range counts {
		counts[k] /= total
	}
	return counts
}

// trigramScore computes the cosine similarity between the trigram profile of s
// (considering only letter runes) and the reference profile.
// Returns a value in [0.0, 1.0].
// Returns 0.0 when s contains fewer than 3 letter runes.
func trigramScore(s string, profile map[string]float64) float64 {
	input := extractTrigrams(s)
	if len(input) == 0 {
		return 0.0
	}

	var dot, normInput, normProfile float64

	for trigram, inputFreq := range input {
		normInput += inputFreq * inputFreq
		if profileFreq, ok := profile[trigram]; ok {
			dot += inputFreq * profileFreq
		}
	}

	for _, profileFreq := range profile {
		normProfile += profileFreq * profileFreq
	}

	denom := math.Sqrt(normInput) * math.Sqrt(normProfile)
	if denom == 0 {
		return 0.0
	}
	return dot / denom
}
