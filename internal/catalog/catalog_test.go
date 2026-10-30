package catalog

import (
	"strings"
	"testing"
)

func TestCatalogShape(t *testing.T) {
	if len(Catalog) != 200 {
		t.Fatalf("catalog length = %d, want 200", len(Catalog))
	}
	seenEmoji := map[string]bool{}
	seenLabel := map[string]bool{}
	seenID := map[string]bool{}
	for i, e := range Catalog {
		for name, v := range map[string]string{"emoji": e.Emoji, "name": e.Name, "label": e.Label, "hint": e.Hint} {
			if strings.TrimSpace(v) == "" {
				t.Errorf("entry %d: empty %s", i, name)
			}
		}
		if seenEmoji[e.Emoji] {
			t.Errorf("duplicate emoji %s", e.Emoji)
		}
		if seenLabel[e.Label] {
			t.Errorf("duplicate label %q", e.Label)
		}
		seenEmoji[e.Emoji] = true
		seenLabel[e.Label] = true
		id := ID(i)
		if seenID[id] {
			t.Errorf("duplicate question ID %q", id)
		}
		seenID[id] = true
	}
	for i, want := range map[int]string{0: "r00", 9: "r09", 99: "r99", 100: "r100", 199: "r199"} {
		if got := ID(i); got != want {
			t.Errorf("ID(%d) = %q, want %q", i, got, want)
		}
	}
	for _, emoji := range []string{"👋", "👍", "✅", "👏", "🙏", "🙌", "❤️", "😂", "🤣", "🎉", "🥳", "👀", "🤔"} {
		if !seenEmoji[emoji] {
			t.Errorf("missing everyday reaction %s", emoji)
		}
	}
}

func TestQuestions(t *testing.T) {
	questions := Questions()
	if len(questions) != len(Catalog) {
		t.Fatalf("questions = %d, want %d", len(questions), len(Catalog))
	}
	for id, q := range questions {
		if err := q.Validate(); err != nil {
			t.Errorf("question %s: %v", id, err)
		}
		if !strings.Contains(q.Instructions, "reply to this message") {
			t.Errorf("question %s instructions = %q", id, q.Instructions)
		}
	}
	// The trash emoji is the first catalog entry, framed with its joke hint.
	q := questions[ID(0)]
	if !strings.Contains(q.Instructions, "🗑️") || !strings.Contains(q.Instructions, "trash") {
		t.Errorf("trash question = %q", q.Instructions)
	}
}

func TestState(t *testing.T) {
	state := State("You eat like a trash seagull")
	if state["message"] != "You eat like a trash seagull" {
		t.Errorf("state = %v", state)
	}
	if state["context"] == "" || state["reactions"] == "" {
		t.Errorf("state = %v, want message + context + reactions", state)
	}
}
