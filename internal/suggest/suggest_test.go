package suggest

import (
	"math"
	"slices"
	"testing"

	"github.com/amscotti/jevmoji/internal/catalog"
	"github.com/amscotti/jevmoji/internal/jev"
)

func noul(p float64) jev.Answer {
	v := p
	return jev.Answer{Type: jev.TypeNoul, Noul: &v}
}

// answersFor builds the full answer map with one probability per entry.
func answersFor(probs ...float64) map[string]jev.Answer {
	answers := make(map[string]jev.Answer, len(probs))
	for i, p := range probs {
		answers[catalog.ID(i)] = noul(p)
	}
	return answers
}

func TestBuildSortsAndChooses(t *testing.T) {
	entries := []catalog.Entry{
		{Emoji: "a", Name: "a", Label: "a"},
		{Emoji: "b", Name: "b", Label: "b"},
		{Emoji: "c", Name: "c", Label: "c"},
		{Emoji: "d", Name: "d", Label: "d"},
		{Emoji: "e", Name: "e", Label: "e"},
		{Emoji: "f", Name: "f", Label: "f"},
	}
	sug, err := Build("hello", entries, answersFor(0.9, 0.6, 0.7, 0.55, 0.8, 0.1))
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	// Five candidates clear 50%, but the fifth is only at 55%, so it does
	// not qualify for the extra slot. The 0.1 candidate stays unchosen.
	wantOrder := []string{"a", "e", "c", "b", "d", "f"}
	for i, label := range wantOrder {
		if sug.Reactions[i].Label != label {
			t.Fatalf("reactions[%d] = %s, want %s (order %v)", i, sug.Reactions[i].Label, label, sug.Reactions)
		}
	}
	wantChosen := []string{"a", "e", "c", "b"}
	if len(sug.Chosen) != len(wantChosen) {
		t.Fatalf("chosen = %v, want %v", sug.Chosen, wantChosen)
	}
	for i, emoji := range wantChosen {
		if sug.Chosen[i] != emoji {
			t.Errorf("chosen[%d] = %s, want %s", i, sug.Chosen[i], emoji)
		}
	}
	for _, r := range sug.Reactions {
		want := r.Probability >= ChosenThreshold && r.Label != "d"
		if r.Chosen != want {
			t.Errorf("reaction %s chosen = %v, want %v", r.Label, r.Chosen, want)
		}
	}
	if sug.State != "hello" {
		t.Errorf("state = %q", sug.State)
	}
}

func TestBuildConditionalFifthReaction(t *testing.T) {
	entries := []catalog.Entry{
		{Emoji: "a", Name: "a", Label: "a"},
		{Emoji: "b", Name: "b", Label: "b"},
		{Emoji: "c", Name: "c", Label: "c"},
		{Emoji: "d", Name: "d", Label: "d"},
		{Emoji: "e", Name: "e", Label: "e"},
		{Emoji: "f", Name: "f", Label: "f"},
	}
	cases := []struct {
		name  string
		probs []float64
		want  []string
	}{
		{
			name:  "base threshold is inclusive",
			probs: []float64{0.50, 0.49, 0.40, 0.30, 0.20, 0.10},
			want:  []string{"a"},
		},
		{
			name:  "fewer matches are not padded",
			probs: []float64{0.90, 0.80, 0.49, 0.20, 0.10, 0},
			want:  []string{"a", "b"},
		},
		{
			name:  "fifth just below threshold is excluded",
			probs: []float64{0.91, 0.88, 0.82, 0.80, math.Nextafter(0.75, 0), 0.73},
			want:  []string{"a", "b", "c", "d"},
		},
		{
			name:  "fifth at threshold is included",
			probs: []float64{0.91, 0.88, 0.82, 0.80, 0.75, 0.74},
			want:  []string{"a", "b", "c", "d", "e"},
		},
		{
			name:  "fifth just above threshold is included",
			probs: []float64{0.91, 0.88, 0.82, 0.80, math.Nextafter(0.75, 1), 0.74},
			want:  []string{"a", "b", "c", "d", "e"},
		},
		{
			name:  "strong matches are sorted and capped at five",
			probs: []float64{0.83, 0.92, 0.75, 0.87, 0.96, 0.89},
			want:  []string{"e", "b", "f", "d", "a"},
		},
		{
			name:  "ties keep catalog order at the cap",
			probs: []float64{0.80, 0.80, 0.80, 0.80, 0.80, 0.80},
			want:  []string{"a", "b", "c", "d", "e"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sug, err := Build("hello", entries, answersFor(tc.probs...))
			if err != nil {
				t.Fatalf("Build: %v", err)
			}
			if !slices.Equal(sug.Chosen, tc.want) {
				t.Errorf("chosen = %v, want %v", sug.Chosen, tc.want)
			}
			for _, reaction := range sug.Reactions {
				if want := slices.Contains(tc.want, reaction.Emoji); reaction.Chosen != want {
					t.Errorf("reaction %s chosen = %v, want %v", reaction.Emoji, reaction.Chosen, want)
				}
			}
		})
	}
}

func TestBuildBelowThresholdNothingChosen(t *testing.T) {
	entries := []catalog.Entry{{Emoji: "a", Name: "a", Label: "a"}, {Emoji: "b", Name: "b", Label: "b"}}
	sug, err := Build("x", entries, answersFor(0.49, 0.2))
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if len(sug.Chosen) != 0 {
		t.Errorf("chosen = %v, want empty", sug.Chosen)
	}
	for _, r := range sug.Reactions {
		if r.Chosen {
			t.Errorf("reaction %s must not be chosen", r.Label)
		}
	}
}

func TestBuildTiesBreakOnCatalogOrder(t *testing.T) {
	entries := []catalog.Entry{{Emoji: "a", Name: "a", Label: "a"}, {Emoji: "b", Name: "b", Label: "b"}}
	sug, err := Build("x", entries, answersFor(0.8, 0.8))
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if sug.Reactions[0].Label != "a" || sug.Reactions[1].Label != "b" {
		t.Errorf("tie order = %v, want catalog order", sug.Reactions)
	}
	if sug.Chosen[0] != "a" || sug.Chosen[1] != "b" {
		t.Errorf("chosen = %v", sug.Chosen)
	}
}

func TestBuildRejectsMissingAndWrongAnswers(t *testing.T) {
	entries := []catalog.Entry{{Emoji: "a", Name: "a", Label: "a"}, {Emoji: "b", Name: "b", Label: "b"}}

	if _, err := Build("x", entries, answersFor(0.5)); err == nil {
		t.Error("missing answer must error")
	}

	one := float64(1)
	bad := answersFor(0.5, 0.5)
	bad[catalog.ID(1)] = jev.Answer{Type: jev.TypeChoice, Probabilities: map[string]float64{"a": 1}, Confidence: &one}
	if _, err := Build("x", entries, bad); err == nil {
		t.Error("non-noul answer must error")
	}
}

func TestBuildRejectsOutOfRangeProbabilities(t *testing.T) {
	entries := []catalog.Entry{{Emoji: "a", Name: "a", Label: "a"}, {Emoji: "b", Name: "b", Label: "b"}}
	bad := []float64{-0.1, 1.7, math.NaN()}
	for _, p := range bad {
		answers := answersFor(0.5, 0.5)
		v := p
		answers[catalog.ID(1)] = jev.Answer{Type: jev.TypeNoul, Noul: &v}
		if _, err := Build("x", entries, answers); err == nil {
			t.Errorf("probability %v must error", p)
		}
	}
}

func TestBuildTrashSeagull(t *testing.T) {
	// A message about a trash seagull should surface a trash-themed and a
	// bird-themed emoji.
	probs := make([]float64, len(catalog.Catalog))
	for i := range probs {
		probs[i] = 0.1
	}
	probs[0] = 0.97 // 🗑️ trash
	probs[1] = 0.93 // 🐦 seagull
	probs[38] = 0.6 // 🍟 fries

	sug, err := Build("You eat like a trash seagull", catalog.Catalog, answersFor(probs...))
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	want := []string{"🗑️", "🐦", "🍟"}
	for i, emoji := range want {
		if sug.Chosen[i] != emoji {
			t.Fatalf("chosen = %v, want leading %v", sug.Chosen, want)
		}
	}
	if len(sug.Reactions) != len(catalog.Catalog) {
		t.Errorf("long list = %d reactions, want %d", len(sug.Reactions), len(catalog.Catalog))
	}
}
