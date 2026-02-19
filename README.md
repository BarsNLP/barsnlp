# BarsNLP

[![CI](https://github.com/BarsNLP/barsnlp/actions/workflows/ci.yml/badge.svg)](https://github.com/BarsNLP/barsnlp/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/BarsNLP/barsnlp.svg)](https://pkg.go.dev/github.com/BarsNLP/barsnlp)
[![Go Version](https://img.shields.io/github/go-mod/go-version/BarsNLP/barsnlp)](https://github.com/BarsNLP/barsnlp/blob/main/go.mod)
[![License](https://img.shields.io/github/license/BarsNLP/barsnlp)](LICENSE)

NLP toolkit for Azerbaijani language. Pure Go, zero dependencies.

All packages are safe for concurrent use.

## Packages

| Package | Description |
|---------|-------------|
| [translit](#transliteration) | Latin / Cyrillic script conversion |
| [tokenizer](#tokenizer) | Word and sentence tokenization with byte offsets |
| [morph](#morphological-analysis) | Stem and suffix chain decomposition |
| [numtext](#number-to-text) | Number / text conversion ("123" &rarr; "yuz iyirmi uc") |
| [ner](#named-entity-recognition) | FIN, VOEN, phone, email, IBAN, plate, URL extraction |

## Install

```
go get github.com/BarsNLP/barsnlp
```

Requires Go 1.25.7 or later.

## Transliteration

Convert Azerbaijani text between Latin and Cyrillic scripts.

```go
translit.CyrillicToLatin("Азәрбајҹан")
// Azərbaycan

translit.LatinToCyrillic("Həyat gözəldir")
// Һәјат ҝөзәлдир
```

Contextual rules handle Cyrillic Г/г disambiguation automatically. Non-Azerbaijani characters (digits, punctuation, emoji) pass through unchanged.

## Tokenizer

Split Azerbaijani text into words and sentences with byte offsets.

```go
// Word tokenization
tokenizer.Words("Bakı'nın küçələri gözəldir.")
// [Bakı'nın küçələri gözəldir]

// Structured tokens with byte offsets
for _, t := range tokenizer.WordTokens("Salam, dünya!") {
    fmt.Printf("%s: %q\n", t.Type, t.Text)
}
// Word: "Salam"
// Punctuation: ","
// Space: " "
// Word: "dünya"
// Punctuation: "!"

// Sentence splitting
tokenizer.Sentences("Birinci cümlə. İkinci cümlə.")
// [Birinci cümlə.  İkinci cümlə.]
```

Handles URLs, emails, Azerbaijani abbreviations (Prof., Az.R.), thousand-separator dots (1.000.000), decimal commas (3,14), hyphens (sosial-iqtisadi), and apostrophe suffixes (Bakı'nın).

## Morphological Analysis

Decompose inflected Azerbaijani words into stem and suffix chain.

```go
// Extract stem from inflected word
morph.Stem("kitablarımızdan")
// kitab

// Full morphological analysis
for _, a := range morph.Analyze("kitablar") {
    fmt.Println(a)
}
// kitab[Plural:lar]
// kitabl[TenseAorist:ar]
// kitablar

// Batch stemming (pairs with tokenizer.Words)
morph.Stems([]string{"kitablarımızdan", "evlərdə", "gəlmişdir"})
// [kitab ev gəl]
```

Uses a table-driven morphotactic state machine with backtracking. Validates vowel harmony, consonant assimilation, and suffix ordering. Includes an embedded dictionary (~12K stems from Wiktionary) for stem validation.

## Number-to-Text

Convert between numbers and Azerbaijani text representations.

```go
// Cardinal number
numtext.Convert(123)
// yüz iyirmi üç

// Ordinal number with vowel-harmony suffix
numtext.ConvertOrdinal(5)
// beşinci

// Decimal: math mode
numtext.ConvertFloat("3.14", numtext.MathMode)
// üç tam yüzdə on dörd

// Decimal: digit-by-digit mode
numtext.ConvertFloat("3.14", numtext.DigitMode)
// üç vergül bir dörd

// Parse text back to number
n, _ := numtext.Parse("iki milyon üç yüz min doxsan beş")
fmt.Println(n)
// 2300095
```

Supports integers up to ±10^18, negative numbers, ordinals, and decimals with dot or comma separator. Parse is case-insensitive and accepts both canonical ("yüz") and explicit ("bir yüz") forms.

## Named Entity Recognition

Extract structured entities from Azerbaijani text: FIN, VOEN, phone numbers, emails, IBANs, license plates, and URLs.

```go
// Extract all entities with byte offsets
for _, e := range ner.Recognize("FIN: 5ARPXK2, tel +994501234567") {
    fmt.Printf("%s: %q (labeled=%v)\n", e.Type, e.Text, e.Labeled)
}
// FIN: "5ARPXK2" (labeled=true)
// Phone: "+994501234567" (labeled=false)

// Convenience functions return []string
ner.Phones("+994501234567 və 0551234567")
// [+994501234567 0551234567]

ner.Emails("info@gov.az")
// [info@gov.az]

ner.IBANs("AZ21NABZ00000000137010001944")
// [AZ21NABZ00000000137010001944]
```

FIN and VOEN patterns are ambiguous in isolation. When preceded by a keyword (e.g. "FIN:", "VOEN:"), `Entity.Labeled` is true, indicating higher confidence. Overlapping entities are resolved by preferring longer matches.

## Planned

- **spell** — spell checker (SymSpell algorithm)
- **datetime** — date/time parser ("5 mart 2026" &rarr; structured)
- **detect** — language detection (az/ru/en/tr)
- **normalize** — text normalization, diacritic restoration
- **keywords** — keyword extraction (TF-IDF / TextRank)
- **validate** — text validator (spelling + punctuation + layout)

## License

[MIT](LICENSE)
