package tokenizer

import "testing"

func FuzzWordTokens(f *testing.F) {
	f.Add("Salam, d\u00fcnya!")
	f.Add("user@mail.az")
	f.Add("https://gov.az")
	f.Add("1.000.000,50")
	f.Add("")
	f.Add("\xff\xfe")
	f.Add("h h h h h h h h")
	f.Add(".user@domain.com")
	f.Fuzz(func(t *testing.T, s string) {
		tokens := WordTokens(s)
		verifyInvariants(t, s, tokens)
	})
}

func FuzzSentenceTokens(f *testing.F) {
	f.Add("Birinci. \u0130kinci.")
	f.Add("Prof. \u018eliyev g\u0259ldi.")
	f.Add("Ola bil\u0259r... B\u0259lk\u0259.")
	f.Add("")
	f.Add("Az.R. qanunu.")
	f.Fuzz(func(t *testing.T, s string) {
		tokens := SentenceTokens(s)
		verifyInvariants(t, s, tokens)
	})
}
