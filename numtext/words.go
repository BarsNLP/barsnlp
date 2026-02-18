// Word tables for Azerbaijani number-to-text conversion.
package numtext

const (
	maxAbs  int64 = 1_000_000_000_000_000_000
	hundred int64 = 100

	wordNegative = "mənfi"
	wordHundred  = "yüz"
	wordExact    = "tam"
	wordComma    = "vergül"
	wordZero     = "sıfır"
)

var ones = [10]string{
	"sıfır",
	"bir",
	"iki",
	"üç",
	"dörd",
	"beş",
	"altı",
	"yeddi",
	"səkkiz",
	"doqquz",
}

// tens is indexed by tens digit (1–9); index 0 is unused.
var tens = [10]string{
	"",
	"on",
	"iyirmi",
	"otuz",
	"qırx",
	"əlli",
	"altmış",
	"yetmiş",
	"səksən",
	"doxsan",
}

type magnitude struct {
	value int64
	word  string
}

// magnitudes lists named powers of ten from largest to smallest.
// yüz (100) is handled separately within group conversion and is not listed here.
var magnitudes = []magnitude{
	{value: 1_000_000_000_000_000_000, word: "kvintilyon"},
	{value: 1_000_000_000_000_000, word: "kvadrilyon"},
	{value: 1_000_000_000_000, word: "trilyon"},
	{value: 1_000_000_000, word: "milyard"},
	{value: 1_000_000, word: "milyon"},
	{value: 1_000, word: "min"},
}

// ordinalFull maps the last vowel of a cardinal that ends in a consonant
// to the full ordinal suffix.
var ordinalFull = map[rune]string{
	'a': "ıncı",
	'ı': "ıncı",
	'e': "inci",
	'ə': "inci",
	'i': "inci",
	'o': "uncu",
	'u': "uncu",
	'ö': "üncü",
	'ü': "üncü",
}

// ordinalShort maps the last vowel of a cardinal that ends in a vowel
// to the short ordinal suffix (drops the initial vowel).
var ordinalShort = map[rune]string{
	'a': "ncı",
	'ı': "ncı",
	'e': "nci",
	'ə': "nci",
	'i': "nci",
	'o': "ncu",
	'u': "ncu",
	'ö': "ncü",
	'ü': "ncü",
}

// azVowels contains all Azerbaijani vowel characters for quick membership testing.
const azVowels = "aeəıioöuü"

// denominators maps the number of fractional digits (1–3) to the Azerbaijani
// denominator word used in math-mode decimal reading.
// Denominators for more than 3 digits are composed programmatically.
var denominators = map[int]string{
	1: "onda",
	2: "yüzdə",
	3: "mində",
}
