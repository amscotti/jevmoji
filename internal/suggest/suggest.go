// Package suggest turns the model's per-emoji probabilities into the data the
// page renders: a long list of scored reactions and the short "chosen" list
// thrown onto the message.
package suggest

import (
	"fmt"
	"math"
	"sort"

	"github.com/amscotti/jevmoji/internal/catalog"
	"github.com/amscotti/jevmoji/internal/jev"
)

// The strongest four reactions can qualify at 50%. A fifth must reach 75%,
// so an extra sticker is reserved for a stronger match.
const (
	ChosenThreshold      = 0.50
	DefaultChosenCap     = 4
	ExtraChosenThreshold = 0.75
	ChosenCap            = 5
)

// Reaction is one scored emoji in the long list.
type Reaction struct {
	Emoji       string  `json:"emoji"`
	Name        string  `json:"name"`
	Label       string  `json:"label"`
	Probability float64 `json:"probability"`
	Chosen      bool    `json:"chosen"`
}

// Suggestion is the app's answer to a message: the state it scored, every
// reaction with its probability (long list, best first), and the few emojis
// chosen to sit on the message.
type Suggestion struct {
	State     string     `json:"state"`
	Reactions []Reaction `json:"reactions"`
	Chosen    []string   `json:"chosen"`
	Model     string     `json:"model,omitempty"`
	Usage     *jev.Usage `json:"usage,omitempty"`
	LatencyMS int64      `json:"latency_ms,omitempty"`
}

// Build validates the answers against the catalog and assembles the sorted
// reaction list. Every catalog entry must have a noul answer under its id.
func Build(state string, entries []catalog.Entry, answers map[string]jev.Answer) (Suggestion, error) {
	reactions := make([]Reaction, 0, len(entries))
	for i, e := range entries {
		id := catalog.ID(i)
		ans, ok := answers[id]
		if !ok {
			return Suggestion{}, fmt.Errorf("suggest: missing answer for %s (%s)", id, e.Emoji)
		}
		if ans.Type != jev.TypeNoul || ans.Noul == nil {
			return Suggestion{}, fmt.Errorf("suggest: answer for %s (%s) is not a noul answer", id, e.Emoji)
		}
		if math.IsNaN(*ans.Noul) || *ans.Noul < 0 || *ans.Noul > 1 {
			return Suggestion{}, fmt.Errorf("suggest: answer for %s (%s) has probability %v, want [0,1]", id, e.Emoji, *ans.Noul)
		}
		reactions = append(reactions, Reaction{
			Emoji:       e.Emoji,
			Name:        e.Name,
			Label:       e.Label,
			Probability: *ans.Noul,
		})
	}

	// Stable sort keeps catalog order as the tie-break.
	sort.SliceStable(reactions, func(a, b int) bool {
		return reactions[a].Probability > reactions[b].Probability
	})

	suggestion := Suggestion{
		State:     state,
		Reactions: reactions,
		Chosen:    make([]string, 0, ChosenCap),
	}
	for i := range reactions {
		if len(suggestion.Chosen) >= ChosenCap {
			break
		}
		threshold := ChosenThreshold
		if len(suggestion.Chosen) >= DefaultChosenCap {
			threshold = ExtraChosenThreshold
		}
		if reactions[i].Probability < threshold {
			break // Later reactions have no higher probability.
		}
		reactions[i].Chosen = true
		suggestion.Chosen = append(suggestion.Chosen, reactions[i].Emoji)
	}
	return suggestion, nil
}
