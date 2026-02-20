package datetime

import (
	"encoding/json"
	"flag"
	"os"
	"testing"
	"time"
)

var updateGolden = flag.Bool("update", false, "regenerate golden test files")

// goldenCase represents a single golden test case.
type goldenCase struct {
	Name    string   `json:"name"`
	Input   string   `json:"input"`
	Ref     string   `json:"ref"`
	Results []Result `json:"results"`
}

const goldenPath = "../data/golden/datetime.json"

func TestGolden(t *testing.T) {
	if *updateGolden {
		updateGoldenFile(t)
		return
	}

	data, err := os.ReadFile(goldenPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("datetime.json not found, run with -update to generate")
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

			ref, err := time.Parse(time.RFC3339, tc.Ref)
			if err != nil {
				t.Fatalf("parsing ref time: %v", err)
			}

			got := Extract(tc.Input, ref)

			// Verify offset invariant for every result.
			for _, r := range got {
				if tc.Input[r.Start:r.End] != r.Text {
					t.Errorf("invariant broken: s[%d:%d]=%q != Text=%q",
						r.Start, r.End, tc.Input[r.Start:r.End], r.Text)
				}
			}

			compareResults(t, tc.Results, got)
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
		ref, err := time.Parse(time.RFC3339, cases[i].Ref)
		if err != nil {
			t.Fatalf("case %q: parsing ref time: %v", cases[i].Name, err)
		}
		cases[i].Results = Extract(cases[i].Input, ref)
	}

	out, err := json.MarshalIndent(cases, "", "  ")
	if err != nil {
		t.Fatalf("marshaling golden data: %v", err)
	}

	out = append(out, '\n')

	if err := os.WriteFile(goldenPath, out, 0644); err != nil {
		t.Fatalf("writing golden file: %v", err)
	}

	t.Log("golden file updated, review with: git diff data/golden/datetime.json")
}
