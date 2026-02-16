package translit

import (
	"fmt"
	"strings"
	"testing"
)

func TestCyrillicToLatin(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Individual letters (lowercase)
		{"а→a", "а", "a"},
		{"б→b", "б", "b"},
		{"в→v", "в", "v"},
		{"ғ→ğ", "ғ", "ğ"},
		{"д→d", "д", "d"},
		{"е→e", "е", "e"},
		{"ә→ə", "ә", "ə"},
		{"ж→j", "ж", "j"},
		{"з→z", "з", "z"},
		{"и→i", "и", "i"},
		{"ј→y", "ј", "y"},
		{"к→k", "к", "k"},
		{"ҝ→g", "ҝ", "g"},
		{"л→l", "л", "l"},
		{"м→m", "м", "m"},
		{"н→n", "н", "n"},
		{"о→o", "о", "o"},
		{"ө→ö", "ө", "ö"},
		{"п→p", "п", "p"},
		{"р→r", "р", "r"},
		{"с→s", "с", "s"},
		{"т→t", "т", "t"},
		{"у→u", "у", "u"},
		{"ү→ü", "ү", "ü"},
		{"ф→f", "ф", "f"},
		{"х→x", "х", "x"},
		{"һ→h", "һ", "h"},
		{"ч→ç", "ч", "ç"},
		{"ҹ→c", "ҹ", "c"},
		{"ш→ş", "ш", "ş"},
		{"ы→ı", "ы", "ı"},

		// Individual letters (uppercase)
		{"А→A", "А", "A"},
		{"Б→B", "Б", "B"},
		{"В→V", "В", "V"},
		{"Ғ→Ğ", "Ғ", "Ğ"},
		{"Д→D", "Д", "D"},
		{"Е→E", "Е", "E"},
		{"Ә→Ə", "Ә", "Ə"},
		{"Ж→J", "Ж", "J"},
		{"З→Z", "З", "Z"},
		{"И→İ", "И", "İ"},
		{"Ј→Y", "Ј", "Y"},
		{"К→K", "К", "K"},
		{"Ҝ→G", "Ҝ", "G"},
		{"Л→L", "Л", "L"},
		{"М→M", "М", "M"},
		{"Н→N", "Н", "N"},
		{"О→O", "О", "O"},
		{"Ө→Ö", "Ө", "Ö"},
		{"П→P", "П", "P"},
		{"Р→R", "Р", "R"},
		{"С→S", "С", "S"},
		{"Т→T", "Т", "T"},
		{"У→U", "У", "U"},
		{"Ү→Ü", "Ү", "Ü"},
		{"Ф→F", "Ф", "F"},
		{"Х→X", "Х", "X"},
		{"Һ→H", "Һ", "H"},
		{"Ч→Ç", "Ч", "Ç"},
		{"Ҹ→C", "Ҹ", "C"},
		{"Ш→Ş", "Ш", "Ş"},
		{"Ы→I", "Ы", "I"},

		// Compatibility: Й maps same as Ј
		{"й→y", "й", "y"},
		{"Й→Y", "Й", "Y"},

		// Full words from spec
		{"Азәрбајҹан", "Азәрбајҹан", "Azərbaycan"},
		{"Бакы шәһәри", "Бакы шәһәри", "Bakı şəhəri"},

		// Contextual Г/г: with Ҝ present → Г always Q
		{"Ҝәнҹә (has Ҝ)", "Ҝәнҹә", "Gəncə"},
		{"Ҝ forces Г→Q", "Ҝәнҹә Гала", "Gəncə Qala"},

		// Contextual Г/г: without Ҝ, before back vowel → Q
		{"Гала→Qala", "Гала", "Qala"},
		{"Гол→Qol", "Гол", "Qol"},
		{"Гурд→Qurd", "Гурд", "Qurd"},
		{"Гырмызы→Qırmızı", "Гырмызы", "Qırmızı"},

		// Contextual Г/г: without Ҝ, before front vowel → G
		{"Гәнҹ→Gənc", "Гәнҹ", "Gənc"},
		{"Гөзәл→Gözəl", "Гөзәл", "Gözəl"},
		{"Гүл→Gül", "Гүл", "Gül"},

		// Contextual Г/г: before consonant → Q
		{"Гран→Qran", "Гран", "Qran"},

		// Contextual Г/г: end of string → Q
		{"trailing Г", "баг", "baq"},

		// Contextual Г/г: lookahead skips non-letters
		{"Г.ә→G.ə", "Г.ә", "G.ə"},
		{"Г-а→Q-a", "Г-а", "Q-a"},
		{"Г 123 ә→G 123 ə", "Г 123 ә", "G 123 ə"},
		{"Г.→Q.", "Г.", "Q."},

		// Soft/hard signs removed
		{"soft sign removed", "Письмо", "Pismo"},
		{"hard sign removed", "объект", "obekt"},

		// Case preservation
		{"mixed case", "ГаЛа", "QaLa"},

		// Empty string
		{"empty", "", ""},

		// Non-Azerbaijani passthrough
		{"ASCII passthrough", "Hello 123!", "Hello 123!"},
		{"digits", "12345", "12345"},
		{"punctuation", ".,;:!?", ".,;:!?"},

		// Emoji passthrough
		{"emoji", "Бакы🏙️", "Bakı🏙️"},

		// Mixed content
		{"mixed Az and non-Az", "Бакы city 2024", "Bakı city 2024"},

		// Unicode edge cases
		{"CJK passthrough", "Бакы中文", "Bakı中文"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CyrillicToLatin(tt.input)
			if got != tt.want {
				t.Errorf("CyrillicToLatin(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestLatinToCyrillic(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Individual letters (lowercase)
		{"a→а", "a", "а"},
		{"b→б", "b", "б"},
		{"c→ҹ", "c", "ҹ"},
		{"ç→ч", "ç", "ч"},
		{"d→д", "d", "д"},
		{"e→е", "e", "е"},
		{"ə→ә", "ə", "ә"},
		{"f→ф", "f", "ф"},
		{"g→ҝ", "g", "ҝ"},
		{"ğ→ғ", "ğ", "ғ"},
		{"h→һ", "h", "һ"},
		{"ı→ы", "ı", "ы"},
		{"i→и", "i", "и"},
		{"j→ж", "j", "ж"},
		{"k→к", "k", "к"},
		{"l→л", "l", "л"},
		{"m→м", "m", "м"},
		{"n→н", "n", "н"},
		{"o→о", "o", "о"},
		{"ö→ө", "ö", "ө"},
		{"p→п", "p", "п"},
		{"q→г", "q", "г"},
		{"r→р", "r", "р"},
		{"s→с", "s", "с"},
		{"ş→ш", "ş", "ш"},
		{"t→т", "t", "т"},
		{"u→у", "u", "у"},
		{"ü→ү", "ü", "ү"},
		{"v→в", "v", "в"},
		{"x→х", "x", "х"},
		{"y→ј", "y", "ј"},
		{"z→з", "z", "з"},

		// Individual letters (uppercase)
		{"A→А", "A", "А"},
		{"B→Б", "B", "Б"},
		{"C→Ҹ", "C", "Ҹ"},
		{"Ç→Ч", "Ç", "Ч"},
		{"D→Д", "D", "Д"},
		{"E→Е", "E", "Е"},
		{"Ə→Ә", "Ə", "Ә"},
		{"F→Ф", "F", "Ф"},
		{"G→Ҝ", "G", "Ҝ"},
		{"Ğ→Ғ", "Ğ", "Ғ"},
		{"H→Һ", "H", "Һ"},
		{"I→Ы", "I", "Ы"},
		{"İ→И", "İ", "И"},
		{"J→Ж", "J", "Ж"},
		{"K→К", "K", "К"},
		{"L→Л", "L", "Л"},
		{"M→М", "M", "М"},
		{"N→Н", "N", "Н"},
		{"O→О", "O", "О"},
		{"Ö→Ө", "Ö", "Ө"},
		{"P→П", "P", "П"},
		{"Q→Г", "Q", "Г"},
		{"R→Р", "R", "Р"},
		{"S→С", "S", "С"},
		{"Ş→Ш", "Ş", "Ш"},
		{"T→Т", "T", "Т"},
		{"U→У", "U", "У"},
		{"Ü→Ү", "Ü", "Ү"},
		{"V→В", "V", "В"},
		{"X→Х", "X", "Х"},
		{"Y→Ј", "Y", "Ј"},
		{"Z→З", "Z", "З"},

		// Full words from spec
		{"Azərbaycan", "Azərbaycan", "Азәрбајҹан"},
		{"Həyat gözəldir", "Həyat gözəldir", "Һәјат ҝөзәлдир"},

		// Empty string
		{"empty", "", ""},

		// Passthrough
		{"ASCII passthrough", "Hello 123!", "Һелло 123!"},
		{"digits only", "12345", "12345"},
		{"emoji", "Bakı🏙️", "Бакы🏙️"},

		// Mixed content
		{"mixed", "Bakı city 2024", "Бакы ҹитј 2024"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LatinToCyrillic(tt.input)
			if got != tt.want {
				t.Errorf("LatinToCyrillic(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRoundTripForward(t *testing.T) {
	// CyrillicToLatin(LatinToCyrillic(s)) should return the original for standard Latin input.
	inputs := []string{
		"Azərbaycan",
		"Bakı",
		"Həyat gözəldir",
		"Qala",
		"Gəncə",
		"",
	}
	for _, s := range inputs {
		t.Run(s, func(t *testing.T) {
			got := CyrillicToLatin(LatinToCyrillic(s))
			if got != s {
				t.Errorf("round-trip failed: %q → LatinToCyrillic → CyrillicToLatin → %q", s, got)
			}
		})
	}
}

func TestRoundTripReverseLossy(t *testing.T) {
	// LatinToCyrillic(CyrillicToLatin(s)) is lossy when input contains Ь or Ъ.
	input := "Письмо"
	cyr := LatinToCyrillic(CyrillicToLatin(input))
	if cyr == input {
		t.Errorf("expected lossy round-trip for %q but got exact match", input)
	}
}

func TestArabicStubs(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string) string
		in   string
	}{
		{"ArabicToLatin", ArabicToLatin, "مرحبا"},
		{"LatinToArabic", LatinToArabic, "salam"},
		{"ArabicToLatin empty", ArabicToLatin, ""},
		{"LatinToArabic empty", LatinToArabic, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(tt.in); got != tt.in {
				t.Errorf("%s(%q) = %q, want %q (stub should return input)", tt.name, tt.in, got, tt.in)
			}
		})
	}
}

func TestLargeInput(t *testing.T) {
	// 1MB+ input should complete without panic.
	chunk := "Азәрбајҹан Бакы шәһәри Гала "
	input := strings.Repeat(chunk, 40000) // ~1.1MB
	got := CyrillicToLatin(input)
	if len(got) == 0 {
		t.Error("expected non-empty output for large input")
	}
}

func TestMalformedUTF8(t *testing.T) {
	// Invalid UTF-8 bytes produce U+FFFD via Go's range; should not panic.
	input := "Бакы\xff\xfeшәһәр"
	got := CyrillicToLatin(input)
	if len(got) == 0 {
		t.Error("expected non-empty output for malformed UTF-8 input")
	}
}

// Benchmarks

func BenchmarkCyrillicToLatin(b *testing.B) {
	input := strings.Repeat("Азәрбајҹан Бакы шәһәри Гала Гәнҹ ", 1000)
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for b.Loop() {
		CyrillicToLatin(input)
	}
}

func BenchmarkLatinToCyrillic(b *testing.B) {
	input := strings.Repeat("Azərbaycan Bakı şəhəri Qala Gənc ", 1000)
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for b.Loop() {
		LatinToCyrillic(input)
	}
}

// Examples

func ExampleCyrillicToLatin() {
	fmt.Println(CyrillicToLatin("Азәрбајҹан"))
	fmt.Println(CyrillicToLatin("Бакы шәһәри"))
	// Output:
	// Azərbaycan
	// Bakı şəhəri
}

func ExampleLatinToCyrillic() {
	fmt.Println(LatinToCyrillic("Azərbaycan"))
	fmt.Println(LatinToCyrillic("Həyat gözəldir"))
	// Output:
	// Азәрбајҹан
	// Һәјат ҝөзәлдир
}
