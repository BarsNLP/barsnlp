package ner

import (
	"regexp"
	"sort"
	"strings"
)

// Compiled regexes for each entity type.
// Order matters: more specific patterns (IBAN, URL, Email) are matched first
// so they take priority over generic ones (FIN bare, VOEN bare) in overlap resolution.
var (
	// Phone: international format +994 XX XXX XX XX (spaces optional)
	rePhoneIntl = regexp.MustCompile(`\+994\s?(\d{2})\s?(\d{3})\s?(\d{2})\s?(\d{2})`)
	// Phone: local format 0XX XXX XX XX (spaces optional)
	rePhoneLocal = regexp.MustCompile(`\b0(\d{2})\s?(\d{3})\s?(\d{2})\s?(\d{2})\b`)

	// Email: standard pattern
	reEmail = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)

	// URL: http or https prefixed
	reURL = regexp.MustCompile(`https?://[^\s<>"{}|\\^` + "`" + ` ]+`)

	// IBAN: AZ + 2 digits + 4 uppercase letters + 20 alphanumeric chars = 28 total
	reIBAN = regexp.MustCompile(`\bAZ\d{2}[A-Z]{4}[A-Z0-9]{20}\b`)

	// LicensePlate: XX-YY-ZZZ format
	reLicensePlate = regexp.MustCompile(`\b\d{2}-[A-Z]{2}-\d{3}\b`)

	// FIN labeled: preceded by keyword "FIN" with optional colon/space
	reFINLabeled = regexp.MustCompile(`(?i)\bFIN[:\s]\s?([A-HJ-NP-Z0-9]{7})\b`)
	// FIN bare: 7 uppercase alphanumeric chars (no I, no O)
	reFINBare = regexp.MustCompile(`\b[A-HJ-NP-Z0-9]{7}\b`)

	// VOEN labeled: preceded by keyword "VOEN" or "VÖEN" with optional colon/space
	reVOENLabeled = regexp.MustCompile(`(?i)\bV[ÖO]EN[:\s]\s?(\d{10})\b`)
	// VOEN bare: exactly 10 digits
	reVOENBare = regexp.MustCompile(`\b\d{10}\b`)
)

// recognize is the internal implementation of Recognize.
func recognize(s string) []Entity {
	var all []Entity

	// High-specificity patterns first
	all = append(all, matchURL(s)...)
	all = append(all, matchEmail(s)...)
	all = append(all, matchIBAN(s)...)
	all = append(all, matchLicensePlate(s)...)
	all = append(all, matchPhone(s)...)

	// Ambiguous patterns last (FIN/VOEN labeled, then bare)
	all = append(all, matchFIN(s)...)
	all = append(all, matchVOEN(s)...)

	if len(all) == 0 {
		return nil
	}

	all = resolveOverlaps(all)
	sort.Slice(all, func(i, j int) bool {
		return all[i].Start < all[j].Start
	})
	return all
}

// matchPhone finds phone numbers in both international and local formats.
func matchPhone(s string) []Entity {
	var out []Entity
	for _, m := range rePhoneIntl.FindAllStringIndex(s, -1) {
		out = append(out, Entity{
			Text:  s[m[0]:m[1]],
			Start: m[0],
			End:   m[1],
			Type:  Phone,
		})
	}
	for _, m := range rePhoneLocal.FindAllStringIndex(s, -1) {
		out = append(out, Entity{
			Text:  s[m[0]:m[1]],
			Start: m[0],
			End:   m[1],
			Type:  Phone,
		})
	}
	return out
}

// matchEmail finds email addresses.
func matchEmail(s string) []Entity {
	var out []Entity
	for _, m := range reEmail.FindAllStringIndex(s, -1) {
		out = append(out, Entity{
			Text:  s[m[0]:m[1]],
			Start: m[0],
			End:   m[1],
			Type:  Email,
		})
	}
	return out
}

// matchURL finds HTTP/HTTPS URLs.
func matchURL(s string) []Entity {
	var out []Entity
	for _, m := range reURL.FindAllStringIndex(s, -1) {
		text := s[m[0]:m[1]]
		// Trim trailing punctuation that is likely not part of the URL.
		text = strings.TrimRight(text, ".,;:!?)]}>")
		end := m[0] + len(text)
		out = append(out, Entity{
			Text:  text,
			Start: m[0],
			End:   end,
			Type:  URL,
		})
	}
	return out
}

// matchIBAN finds Azerbaijani IBAN numbers.
func matchIBAN(s string) []Entity {
	var out []Entity
	for _, m := range reIBAN.FindAllStringIndex(s, -1) {
		out = append(out, Entity{
			Text:  s[m[0]:m[1]],
			Start: m[0],
			End:   m[1],
			Type:  IBAN,
		})
	}
	return out
}

// matchLicensePlate finds Azerbaijani license plates.
func matchLicensePlate(s string) []Entity {
	var out []Entity
	for _, m := range reLicensePlate.FindAllStringIndex(s, -1) {
		out = append(out, Entity{
			Text:  s[m[0]:m[1]],
			Start: m[0],
			End:   m[1],
			Type:  LicensePlate,
		})
	}
	return out
}

// matchFIN finds FIN codes. Labeled matches (preceded by "FIN" keyword) take
// priority over bare matches at the same position.
func matchFIN(s string) []Entity {
	var out []Entity

	// Labeled matches: "FIN: XXXXXXX" or "FIN XXXXXXX"
	for _, sub := range reFINLabeled.FindAllStringSubmatchIndex(s, -1) {
		// sub[2]:sub[3] is the capture group (the 7-char code)
		out = append(out, Entity{
			Text:    s[sub[2]:sub[3]],
			Start:   sub[2],
			End:     sub[3],
			Type:    FIN,
			Labeled: true,
		})
	}

	// Bare matches: any 7-char [A-HJ-NP-Z0-9] with word boundaries
	for _, m := range reFINBare.FindAllStringIndex(s, -1) {
		out = append(out, Entity{
			Text:  s[m[0]:m[1]],
			Start: m[0],
			End:   m[1],
			Type:  FIN,
		})
	}

	return out
}

// matchVOEN finds VOEN codes. Labeled matches take priority over bare ones.
func matchVOEN(s string) []Entity {
	var out []Entity

	// Labeled matches: "VOEN: 1234567890" or "VÖEN 1234567890"
	for _, sub := range reVOENLabeled.FindAllStringSubmatchIndex(s, -1) {
		out = append(out, Entity{
			Text:    s[sub[2]:sub[3]],
			Start:   sub[2],
			End:     sub[3],
			Type:    VOEN,
			Labeled: true,
		})
	}

	// Bare matches: any 10-digit sequence with word boundaries
	for _, m := range reVOENBare.FindAllStringIndex(s, -1) {
		out = append(out, Entity{
			Text:  s[m[0]:m[1]],
			Start: m[0],
			End:   m[1],
			Type:  VOEN,
		})
	}

	return out
}

// resolveOverlaps removes overlapping entities. When two entities overlap:
//   - The longer (more specific) match wins.
//   - If equal length, labeled wins over unlabeled.
//   - If still tied, the first one encountered wins.
func resolveOverlaps(entities []Entity) []Entity {
	if len(entities) <= 1 {
		return entities
	}

	// Sort by start offset, then by length descending, then labeled first.
	sort.Slice(entities, func(i, j int) bool {
		if entities[i].Start != entities[j].Start {
			return entities[i].Start < entities[j].Start
		}
		li := entities[i].End - entities[i].Start
		lj := entities[j].End - entities[j].Start
		if li != lj {
			return li > lj
		}
		if entities[i].Labeled != entities[j].Labeled {
			return entities[i].Labeled
		}
		return false
	})

	result := make([]Entity, 0, len(entities))
	maxEnd := 0

	for _, e := range entities {
		if e.Start >= maxEnd {
			result = append(result, e)
			maxEnd = e.End
		} else if e.End > maxEnd {
			// Partial overlap: keep the one that extends further only if
			// we haven't already committed a result covering this range.
			// Since we sorted by start then by length desc, the first entity
			// at a given start is already the best. Skip partial overlaps.
			continue
		}
	}

	return result
}
