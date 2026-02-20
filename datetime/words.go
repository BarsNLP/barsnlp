// Data tables for Azerbaijani date/time parsing.
package datetime

import (
	"regexp"
	"time"
)

// Input validation limits.
const (
	maxInputBytes = 1 << 20 // 1 MiB
	maxResults    = 10000
)

// Numeric validation bounds.
const (
	minDay      = 1
	maxDay      = 31
	minMonth    = 1
	maxMonth    = 12
	minHour     = 0
	maxHour     = 23
	maxMinute   = 59
	maxSecond   = 59
	minYear     = 1
	maxYear     = 9999
	daysPerWeek = 7
)

// Maximum gap in bytes between adjacent date and time spans for merging
// into a single TypeDateTime result (allows for ", saat " or " " between them).
const maxMergeGap = 20

// months maps all Azerbaijani month name forms (bare + 5 noun cases) to time.Month.
// 72 entries total: 12 bare + 12×5 inflected.
var months = map[string]time.Month{
	// Yanvar (January) — back vowel, unrounded
	"yanvar":    time.January,
	"yanvarda":  time.January,
	"yanvardan": time.January,
	"yanvarın":  time.January,
	"yanvarı":   time.January,
	"yanvara":   time.January,

	// Fevral (February) — back vowel, unrounded
	"fevral":    time.February,
	"fevralda":  time.February,
	"fevraldan": time.February,
	"fevralın":  time.February,
	"fevralı":   time.February,
	"fevrala":   time.February,

	// Mart (March) — back vowel, unrounded
	"mart":    time.March,
	"martda":  time.March,
	"martdan": time.March,
	"martın":  time.March,
	"martı":   time.March,
	"marta":   time.March,

	// Aprel (April) — front vowel, unrounded
	"aprel":    time.April,
	"apreldə":  time.April,
	"apreldən": time.April,
	"aprelin":  time.April,
	"apreli":   time.April,
	"aprelə":   time.April,

	// May (May) — back vowel, unrounded
	"may":    time.May,
	"mayda":  time.May,
	"maydan": time.May,
	"mayın":  time.May,
	"mayı":   time.May,
	"maya":   time.May,

	// İyun (June) — back vowel, rounded
	"iyun":    time.June,
	"iyunda":  time.June,
	"iyundan": time.June,
	"iyunun":  time.June,
	"iyunu":   time.June,
	"iyuna":   time.June,

	// İyul (July) — back vowel, rounded
	"iyul":    time.July,
	"iyulda":  time.July,
	"iyuldan": time.July,
	"iyulun":  time.July,
	"iyulu":   time.July,
	"iyula":   time.July,

	// Avqust (August) — back vowel, rounded
	"avqust":    time.August,
	"avqustda":  time.August,
	"avqustdan": time.August,
	"avqustun":  time.August,
	"avqustu":   time.August,
	"avqusta":   time.August,

	// Sentyabr (September) — back vowel, unrounded
	"sentyabr":    time.September,
	"sentyabrda":  time.September,
	"sentyabrdan": time.September,
	"sentyabrın":  time.September,
	"sentyabrı":   time.September,
	"sentyabra":   time.September,

	// Oktyabr (October) — back vowel, unrounded
	"oktyabr":    time.October,
	"oktyabrda":  time.October,
	"oktyabrdan": time.October,
	"oktyabrın":  time.October,
	"oktyabrı":   time.October,
	"oktyabra":   time.October,

	// Noyabr (November) — back vowel, unrounded
	"noyabr":    time.November,
	"noyabrda":  time.November,
	"noyabrdan": time.November,
	"noyabrın":  time.November,
	"noyabrı":   time.November,
	"noyabra":   time.November,

	// Dekabr (December) — back vowel, unrounded
	"dekabr":    time.December,
	"dekabrda":  time.December,
	"dekabrdan": time.December,
	"dekabrın":  time.December,
	"dekabrı":   time.December,
	"dekabra":   time.December,
}

// genitiveMonths identifies month forms that expect a following possessive day number.
// Used to detect patterns like "martın 15-i", "fevralın 3-ü".
var genitiveMonths = map[string]bool{
	"yanvarın":   true,
	"fevralın":   true,
	"martın":     true,
	"aprelin":    true,
	"mayın":      true,
	"iyunun":     true,
	"iyulun":     true,
	"avqustun":   true,
	"sentyabrın": true,
	"oktyabrın":  true,
	"noyabrın":   true,
	"dekabrın":   true,
}

// weekdayEntry holds a weekday name and its time.Weekday value.
type weekdayEntry struct {
	name    string
	weekday time.Weekday
}

// weekdays lists Azerbaijani weekday names ordered by length descending.
// Multi-word names must be checked before their single-word prefixes to prevent
// "bazar" from matching before "bazar ertəsi".
var weekdays = []weekdayEntry{
	{"çərşənbə axşamı", time.Tuesday},
	{"bazar ertəsi", time.Monday},
	{"cümə axşamı", time.Thursday},
	{"çərşənbə", time.Wednesday},
	{"şənbə", time.Saturday},
	{"cümə", time.Friday},
	{"bazar", time.Sunday},
}

// dayOffsets maps single-word and two-word relative date keywords to day offsets from ref.
var dayOffsets = map[string]int{
	"bu gün":   0,
	"bugün":    0,
	"sabah":    1,
	"birigün":  2,
	"dünən":    -1,
	"srağagün": -2,
}

// periodPrefix maps modifier words to their period offset.
// -1 = previous, 0 = current, +1 = next.
var periodPrefix = map[string]int{
	"keçən": -1,
	"bu":    0,
	"gələn": 1,
}

// periodKind identifies what kind of period a unit word represents.
type periodKind int

const (
	periodWeek periodKind = iota
	periodMonth
	periodYear
)

// periodUnits maps period unit words to their kind.
var periodUnits = map[string]periodKind{
	"həftə": periodWeek,
	"ay":    periodMonth,
	"il":    periodYear,
}

// dirKind represents past or future in quantity-direction expressions.
type dirKind int

const (
	dirBefore dirKind = iota
	dirAfter
)

// directionWords maps Azerbaijani direction keywords to their meaning.
var directionWords = map[string]dirKind{
	"əvvəl": dirBefore,
	"öncə":  dirBefore,
	"sonra": dirAfter,
}

// qtyUnit identifies the time unit in quantity-direction expressions.
type qtyUnit int

const (
	qtyDay qtyUnit = iota
	qtyWeek
	qtyMonth
	qtyYear
	qtyHour
	qtyMinute
	qtySecond
)

// quantityUnits maps unit words to their type.
var quantityUnits = map[string]qtyUnit{
	"gün":    qtyDay,
	"həftə":  qtyWeek,
	"ay":     qtyMonth,
	"il":     qtyYear,
	"saat":   qtyHour,
	"dəqiqə": qtyMinute,
	"saniyə": qtySecond,
}

// timeShift represents AM/PM disambiguation for time-of-day words.
type timeShift int

const (
	shiftAM timeShift = iota // Keep hour as-is (morning/midday)
	shiftPM                  // Add 12 to hour (evening/night)
)

// timeOfDayWords maps time-of-day words to their disambiguation effect.
// These only fire when combined with an explicit hour, not as standalone results.
var timeOfDayWords = map[string]timeShift{
	"səhər":  shiftAM,
	"gündüz": shiftAM,
	"axşam":  shiftPM,
	"gecə":   shiftPM,
}

// bridgeWord is the possessive compound connector "ayının"
// in formal date patterns like "mart ayının 15-i".
const bridgeWord = "ayının"

// Compiled regex patterns for numeric date/time formats.
var (
	// ISO 8601: YYYY-MM-DD
	reISO = regexp.MustCompile(`\b(\d{4})-(\d{2})-(\d{2})\b`)

	// Dot-separated: DD.MM.YYYY (Azerbaijani convention)
	reDot = regexp.MustCompile(`\b(\d{1,2})\.(\d{1,2})\.(\d{4})\b`)

	// Slash-separated: DD/MM/YYYY
	reSlash = regexp.MustCompile(`\b(\d{1,2})/(\d{1,2})/(\d{4})\b`)

	// Time: HH:MM or HH:MM:SS (24-hour format)
	reTime = regexp.MustCompile(`\b(\d{1,2}):(\d{2})(?::(\d{2}))?\b`)

	// Ordinal or possessive suffix on a digit: 5-ci, 1-inci, 5-nci, 15-i, 3-ü.
	// Group 1: the digit(s).
	// Covers all 12 ordinal variants plus possessive forms.
	// Uses ^/$ instead of \b because Go's \b is ASCII-only and fails on ü/ı suffixes.
	// This regex is only applied to pre-split words via parseOrdinalWord.
	reOrdinalDay = regexp.MustCompile(`^(\d{1,2})[-.]?(?:(?:[iıuü])?nc[iıuü]|c[iıuü]|[iıuü])$`)
)
