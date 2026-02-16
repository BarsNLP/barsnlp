# BarsNLP

NLP toolkit for Azerbaijani language. Pure Go, zero dependencies.

## Install

```
go get github.com/BarsNLP/barsnlp
```

## Transliteration

Convert Azerbaijani text between Latin and Cyrillic scripts.

```go
package main

import (
	"fmt"
	"github.com/BarsNLP/barsnlp/translit"
)

func main() {
	fmt.Println(translit.CyrillicToLatin("Азәрбајҹан"))
	// Azərbaycan

	fmt.Println(translit.LatinToCyrillic("Həyat gözəldir"))
	// Һәјат ҝөзәлдир
}
```

Contextual rules handle Cyrillic Г/г disambiguation automatically. Non-Azerbaijani characters (digits, punctuation, emoji) pass through unchanged.

All functions are safe for concurrent use.

## License

MIT
