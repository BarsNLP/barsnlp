package detect

import (
	"encoding/json"
	"flag"
	"os"
	"testing"
)

var updateGolden = flag.Bool("update", false, "regenerate golden test files")

// goldenCase represents a single golden test case for language detection.
type goldenCase struct {
	Name       string `json:"name"`
	Input      string `json:"input"`
	WantLang   string `json:"want_lang"`   // Language name: "Azerbaijani", "Russian", etc.
	WantScript string `json:"want_script"` // Script code: "Latn", "Cyrl", ""
	WantCode   string `json:"want_code"`   // ISO 639-1: "az", "ru", "en", "tr", ""
}

const goldenPath = "../data/golden/detect.json"

func TestGolden(t *testing.T) {
	if *updateGolden {
		updateGoldenFile(t)
		return
	}

	data, err := os.ReadFile(goldenPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("detect.json not found, run with -update to generate")
		}
		t.Fatalf("reading golden file: %v", err)
	}

	var cases []goldenCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("parsing golden file: %v", err)
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			got := Detect(tc.Input)

			if got.Lang.String() != tc.WantLang {
				t.Errorf("Lang: got %q, want %q", got.Lang.String(), tc.WantLang)
			}

			if got.Script.String() != tc.WantScript {
				t.Errorf("Script: got %q, want %q", got.Script.String(), tc.WantScript)
			}

			gotCode := Lang(tc.Input)
			if gotCode != tc.WantCode {
				t.Errorf("Lang code: got %q, want %q", gotCode, tc.WantCode)
			}

			if tc.WantLang != "Unknown" && got.Confidence <= 0 {
				t.Errorf("Confidence: expected > 0 for %q, got %f", tc.WantLang, got.Confidence)
			}
		})
	}
}

func updateGoldenFile(t *testing.T) {
	t.Helper()

	data, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("reading golden file for update: %v", err)
	}

	var cases []goldenCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("parsing golden file for update: %v", err)
	}

	for i := range cases {
		got := Detect(cases[i].Input)
		cases[i].WantLang = got.Lang.String()
		cases[i].WantScript = got.Script.String()
		cases[i].WantCode = Lang(cases[i].Input)
	}

	out, err := json.MarshalIndent(cases, "", "  ")
	if err != nil {
		t.Fatalf("marshaling golden data: %v", err)
	}

	out = append(out, '\n')

	if err := os.WriteFile(goldenPath, out, 0644); err != nil {
		t.Fatalf("writing golden file: %v", err)
	}

	t.Log("golden file updated, review with: git diff data/golden/detect.json")
}
