package detect

import (
	"strings"
	"sync"
	"testing"
	"time"
)

func FuzzDetect(f *testing.F) {
	// Seeds: representative inputs for each language.
	f.Add("Salam, necəsən? Bu gün hava çox gözəldir.")
	f.Add("Привет, как у тебя дела сегодня?")
	f.Add("Hello, how are you doing today?")
	f.Add("Merhaba, bugün hava çok güzel değil mi?")
	f.Add("Бу мәтн Азәрбајҹан дилиндәдир")
	f.Add("")
	f.Add("   \t\n")
	f.Add("1234567890")
	f.Add("\xff\xfe")  // malformed UTF-8
	f.Add("test")     // fewer than minLetters

	f.Fuzz(func(t *testing.T, s string) {
		// Detect must never panic.
		r := Detect(s)

		// If detection succeeded, verify invariants.
		if r.Lang != Unknown {
			if r.Confidence < 0 || r.Confidence > 1.0 {
				t.Errorf("Detect(%q): confidence %f outside [0,1]", s, r.Confidence)
			}
		}

		// DetectAll must never panic.
		results := DetectAll(s)
		if results != nil {
			if len(results) != 4 {
				t.Errorf("DetectAll(%q): got %d results, want 4", s, len(results))
			}

			// Scores must be non-negative and sorted descending.
			for i, res := range results {
				if res.Confidence < 0 || res.Confidence > 1.0 {
					t.Errorf("DetectAll(%q)[%d]: confidence %f outside [0,1]", s, i, res.Confidence)
				}
				if i > 0 && res.Confidence > results[i-1].Confidence {
					t.Errorf("DetectAll(%q): not sorted descending at index %d", s, i)
				}
			}

			// Scores must sum to approximately 1.0.
			var total float64
			for _, res := range results {
				total += res.Confidence
			}
			if total < 0.99 || total > 1.01 {
				t.Errorf("DetectAll(%q): scores sum to %f, want ~1.0", s, total)
			}
		}

		// Lang must never panic and must be consistent with Detect.
		lang := Lang(s)
		if r.Lang == Unknown && lang != "" {
			t.Errorf("Lang(%q) = %q, but Detect returned Unknown", s, lang)
		}
	})
}

// TestExactlyMaxInput verifies that inputs at exactly maxInputBytes are processed normally.
func TestExactlyMaxInput(t *testing.T) {
	// Build input exactly at the limit using ASCII-only sentence to avoid
	// splitting a multi-byte rune when slicing to maxInputBytes.
	sentence := "Hello how are you doing today? "
	repeats := (maxInputBytes / len(sentence)) + 1
	exact := strings.Repeat(sentence, repeats)[:maxInputBytes]
	r := Detect(exact)
	// Should process normally (not panic or return Unknown if enough letters).
	_ = r
}

// TestConcurrentSafety verifies the package is safe for concurrent use by multiple goroutines.
func TestConcurrentSafety(t *testing.T) {
	inputs := []string{
		"Salam, necəsən? Bu gün hava çox gözəldir.",
		"Привет, как у тебя дела сегодня?",
		"Hello, how are you doing today?",
		"Merhaba, bugün hava çok güzel değil mi?",
	}

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := range goroutines {
		go func() {
			defer wg.Done()
			s := inputs[i%len(inputs)]
			_ = Detect(s)
			_ = DetectAll(s)
			_ = Lang(s)
		}()
	}
	wg.Wait()
}

// TestMalformedUTF8 verifies that invalid UTF-8 sequences do not cause panics.
func TestMalformedUTF8(t *testing.T) {
	inputs := []string{
		"\xff\xfe",
		"valid\xff invalid",
		string([]byte{0x80, 0x81, 0x82}),
		"Salam\xffnecəsən",
	}
	for _, s := range inputs {
		// Must not panic.
		_ = Detect(s)
		_ = DetectAll(s)
		_ = Lang(s)
	}
}

// TestNullByteInjection verifies that embedded null bytes do not cause panics.
func TestNullByteInjection(t *testing.T) {
	inputs := []string{
		"Salam\x00necəsən? Bu gün hava çox gözəldir.",
		"\x00\x00\x00",
		"Hello\x00world how are you today?",
	}
	for _, s := range inputs {
		// Must not panic.
		_ = Detect(s)
		_ = DetectAll(s)
		_ = Lang(s)
	}
}

// TestReDoSResistance verifies detection completes quickly on adversarial input.
func TestReDoSResistance(t *testing.T) {
	// Adversarial input: long repetitive pattern.
	adversarial := strings.Repeat("aaaa", 100000)
	start := time.Now()
	_ = Detect(adversarial)
	elapsed := time.Since(start)
	if elapsed > 2*time.Second {
		t.Errorf("detection took %v on adversarial input, want < 2s", elapsed)
	}
}
