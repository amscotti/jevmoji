package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/amscotti/jevmoji/internal/jev"
	"github.com/amscotti/jevmoji/internal/suggest"
)

// fakeJudge records what the server sent and answers from a probability map.
type fakeJudge struct {
	calls     int
	state     any
	questions map[string]jev.Question
	probs     map[string]float64
	err       error
}

func (f *fakeJudge) SystemOne(_ context.Context, state any, questions map[string]jev.Question) (*jev.Response, error) {
	f.calls++
	f.state = state
	f.questions = questions
	if f.err != nil {
		return nil, f.err
	}
	answers := make(map[string]jev.Answer, len(questions))
	for id := range questions {
		p := f.probs[id]
		v := p
		answers[id] = jev.Answer{Type: jev.TypeNoul, Noul: &v}
	}
	return &jev.Response{
		Model:   "jev-latest",
		Answers: answers,
		Usage:   jev.Usage{InputTokens: 1200},
	}, nil
}

// stepClock hands out t0 on the first call and t0+187ms on the second.
func stepClock() func() time.Time {
	t0 := time.Unix(1000, 0)
	step := 0
	return func() time.Time {
		step++
		if step == 1 {
			return t0
		}
		return t0.Add(187 * time.Millisecond)
	}
}

func newTestServer(judge *fakeJudge) *httptest.Server {
	srv, err := New(Config{Judge: judge, Now: stepClock()})
	if err != nil {
		panic(err)
	}
	ts := httptest.NewServer(srv.Routes())
	return ts
}

func postForm(t *testing.T, endpoint, text string) *http.Response {
	t.Helper()
	resp, err := http.PostForm(endpoint, url.Values{"text": {text}})
	if err != nil {
		t.Fatalf("POST /api/reactions: %v", err)
	}
	return resp
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(body)
}

func TestHealth(t *testing.T) {
	ts := newTestServer(&fakeJudge{})
	t.Cleanup(ts.Close)

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusOK || body != "ok" {
		t.Errorf("health = %d %q", resp.StatusCode, body)
	}
}

func TestStaticAssetsAreServed(t *testing.T) {
	ts := newTestServer(&fakeJudge{})
	t.Cleanup(ts.Close)

	for _, path := range []string{"/static/app.css", "/static/ui.js"} {
		resp, err := http.Get(ts.URL + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		_ = readBody(t, resp)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("GET %s: status = %d, want 200", path, resp.StatusCode)
		}
	}
}

func TestSecurityHeaders(t *testing.T) {
	ts := newTestServer(&fakeJudge{})
	t.Cleanup(ts.Close)

	for _, path := range []string{"/", "/health", "/static/app.css"} {
		resp, err := http.Get(ts.URL + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		_ = readBody(t, resp)
		for header, want := range map[string]string{
			"X-Content-Type-Options": "nosniff",
			"X-Frame-Options":        "DENY",
			"Referrer-Policy":        "strict-origin-when-cross-origin",
		} {
			if got := resp.Header.Get(header); got != want {
				t.Errorf("GET %s: %s = %q, want %q", path, header, got, want)
			}
		}
	}

	resp := postForm(t, ts.URL+"/api/reactions", "hi")
	_ = readBody(t, resp)
	if got := resp.Header.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("POST /api/reactions: X-Content-Type-Options = %q", got)
	}
}

func TestPageRenders(t *testing.T) {
	ts := newTestServer(&fakeJudge{})
	t.Cleanup(ts.Close)

	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	for _, want := range []string{
		`id="message"`,
		"Words in.",
		"Emojis out.",
		"Your reactions will land here.",
		"Meeting fuel",
		"/static/ui.js",
		"https://unpkg.com/htmx.org@4.0.0/dist/htmx.min.js",
		`integrity="sha384-`,
		`crossorigin="anonymous"`,
		`rel="icon"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("page missing %q", want)
		}
	}
	for _, gone := range []string{
		"alpinejs",
		"/static/app.js",
		"x-data",
		"tray()",
	} {
		if strings.Contains(body, gone) {
			t.Errorf("page should not contain %q (read-only reactions)", gone)
		}
	}
}

func TestReactionsHTML(t *testing.T) {
	judge := &fakeJudge{probs: map[string]float64{"r00": 0.97, "r01": 0.93, "r38": 0.6}}
	ts := newTestServer(judge)
	t.Cleanup(ts.Close)

	resp := postForm(t, ts.URL+"/api/reactions", "You eat like a trash seagull")
	body := readBody(t, resp)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if !strings.Contains(body, "You eat like a trash seagull") {
		t.Error("message text missing from the card")
	}
	if got := strings.Count(body, "chosen-pill"); got != 3 {
		t.Errorf("chosen pills = %d, want 3", got)
	}
	// The reactions stay read-only. The only button copies the model's picks.
	for _, gone := range []string{"tray-pill", "pill-count", "x-on:click", "x-data", "message-head"} {
		if strings.Contains(body, gone) {
			t.Errorf("reactions partial should not contain %q (read-only)", gone)
		}
	}
	if got := strings.Count(body, "<button"); got != 1 || !strings.Contains(body, "data-copy-emojis=") {
		t.Error("expected only the copy-emojis button, with no reaction buttons")
	}
	for _, emoji := range []string{"🗑️", "🐦", "🍟"} {
		if !strings.Contains(body, emoji) {
			t.Errorf("body missing %s", emoji)
		}
	}
	if judge.calls != 1 {
		t.Errorf("judge calls = %d, want 1", judge.calls)
	}
}

func TestReactionsWithNoStrongMatches(t *testing.T) {
	judge := &fakeJudge{probs: map[string]float64{"r00": 0.49}}
	ts := newTestServer(judge)
	t.Cleanup(ts.Close)

	resp := postForm(t, ts.URL+"/api/reactions", "a message without strong matches")
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if !strings.Contains(body, "No strong matches this time.") {
		t.Error("missing explanation when no reaction clears the threshold")
	}
	if strings.Contains(body, "chosen-pill") || strings.Contains(body, "data-copy-emojis=") {
		t.Error("no-match results should not offer reactions or a copy button")
	}
}

func TestReactionsJSON(t *testing.T) {
	judge := &fakeJudge{probs: map[string]float64{"r00": 0.97, "r01": 0.93, "r38": 0.6}}
	ts := newTestServer(judge)
	t.Cleanup(ts.Close)

	resp := postForm(t, ts.URL+"/api/reactions?format=json", "You eat like a trash seagull")
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d body %s", resp.StatusCode, body)
	}

	var sug suggest.Suggestion
	if err := json.Unmarshal([]byte(body), &sug); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	if sug.State != "You eat like a trash seagull" {
		t.Errorf("state = %q", sug.State)
	}
	if len(sug.Reactions) != 200 {
		t.Fatalf("long list = %d reactions, want 200", len(sug.Reactions))
	}
	want := []string{"🗑️", "🐦", "🍟"}
	for i, emoji := range want {
		if sug.Chosen[i] != emoji {
			t.Fatalf("chosen = %v, want leading %v", sug.Chosen, want)
		}
	}
	if sug.Model != "jev-latest" || sug.Usage == nil || sug.Usage.InputTokens != 1200 {
		t.Errorf("model/usage = %q %+v", sug.Model, sug.Usage)
	}
	if sug.LatencyMS != 187 {
		t.Errorf("latency = %dms, want 187ms", sug.LatencyMS)
	}
	// long list is sorted best-first
	if sug.Reactions[0].Emoji != "🗑️" || sug.Reactions[1].Emoji != "🐦" {
		t.Errorf("sorted list = %v", sug.Reactions[:2])
	}
}

func TestFifthReactionInHTMLAndJSON(t *testing.T) {
	judge := &fakeJudge{probs: map[string]float64{
		"r00": 0.95,
		"r01": 0.94,
		"r02": 0.90,
		"r03": 0.80,
		"r04": 0.75,
		"r05": 0.74,
	}}
	ts := newTestServer(judge)
	t.Cleanup(ts.Close)

	body := readBody(t, postForm(t, ts.URL+"/api/reactions", "a message with several strong matches"))
	if got := strings.Count(body, "chosen-pill"); got != 5 {
		t.Errorf("chosen stickers = %d, want 5", got)
	}
	if !strings.Contains(body, `data-copy-emojis="🗑️ 🐦 🦝 🐀 🪳 "`) {
		t.Error("copy action should include all five chosen emojis")
	}

	bodyJSON := readBody(t, postForm(t, ts.URL+"/api/reactions?format=json", "a message with several strong matches"))
	var sug suggest.Suggestion
	if err := json.Unmarshal([]byte(bodyJSON), &sug); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(sug.Chosen) != 5 || len(sug.Reactions) != 200 {
		t.Fatalf("chosen = %d, reactions = %d, want 5 and 200", len(sug.Chosen), len(sug.Reactions))
	}
	if sug.Reactions[4].Probability != 0.75 || !sug.Reactions[4].Chosen || sug.Reactions[5].Chosen {
		t.Error("JSON chosen flags should match the conditional fifth sticker")
	}
}

func TestJSONBodyInput(t *testing.T) {
	judge := &fakeJudge{probs: map[string]float64{"r00": 0.9}}
	ts := newTestServer(judge)
	t.Cleanup(ts.Close)

	resp, err := http.Post(ts.URL+"/api/reactions?format=json", "application/json",
		strings.NewReader(`{"text":"  hello there  "}`))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	body := readBody(t, resp)
	var sug suggest.Suggestion
	if err := json.Unmarshal([]byte(body), &sug); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sug.State != "hello there" {
		t.Errorf("state = %q (want trimmed)", sug.State)
	}
}

func TestBadJSONBody(t *testing.T) {
	judge := &fakeJudge{}
	ts := newTestServer(judge)
	t.Cleanup(ts.Close)

	resp, err := http.Post(ts.URL+"/api/reactions", "application/json", strings.NewReader("{oops"))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	body := readBody(t, resp)
	if !strings.Contains(body, "expected a JSON body") {
		t.Errorf("error card = %q", body)
	}
	if judge.calls != 0 {
		t.Errorf("judge should not be called, calls = %d", judge.calls)
	}
}

func TestBlankMessageStaysIdle(t *testing.T) {
	judge := &fakeJudge{}
	ts := newTestServer(judge)
	t.Cleanup(ts.Close)

	for _, text := range []string{"", "   "} {
		resp := postForm(t, ts.URL+"/api/reactions", text)
		body := readBody(t, resp)
		if !strings.Contains(body, "Your reactions will land here.") {
			t.Errorf("text %q: idle card missing, body = %q", text, body)
		}
	}
	if judge.calls != 0 {
		t.Errorf("judge calls = %d, want 0", judge.calls)
	}
}

func TestShortMessageIsScored(t *testing.T) {
	judge := &fakeJudge{probs: map[string]float64{"r00": 0.9}}
	ts := newTestServer(judge)
	t.Cleanup(ts.Close)

	resp := postForm(t, ts.URL+"/api/reactions", "hi")
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if !strings.Contains(body, "emojis scored") {
		t.Errorf("short message should be scored, body = %q", body)
	}
	if judge.calls != 1 {
		t.Errorf("judge calls = %d, want 1", judge.calls)
	}
}

func TestJudgeErrorSurfacesCard(t *testing.T) {
	judge := &fakeJudge{err: errors.New("upstream down")}
	ts := newTestServer(judge)
	t.Cleanup(ts.Close)

	resp := postForm(t, ts.URL+"/api/reactions", "You eat like a trash seagull")
	body := readBody(t, resp)
	if !strings.Contains(body, "The emoji oracle is taking a breather.") {
		t.Errorf("error card missing, body = %q", body)
	}

	respJSON := postForm(t, ts.URL+"/api/reactions?format=json", "You eat like a trash seagull")
	bodyJSON := readBody(t, respJSON)
	if respJSON.StatusCode != http.StatusBadGateway {
		t.Errorf("json status = %d, want 502", respJSON.StatusCode)
	}
	var envelope map[string]string
	if err := json.Unmarshal([]byte(bodyJSON), &envelope); err != nil || envelope["error"] == "" {
		t.Errorf("json error envelope = %q (%v)", bodyJSON, err)
	}
}

func TestJudgeSeesStateAndFullQuestionSet(t *testing.T) {
	judge := &fakeJudge{}
	ts := newTestServer(judge)
	t.Cleanup(ts.Close)

	_ = readBody(t, postForm(t, ts.URL+"/api/reactions", "You eat like a trash seagull"))

	state, ok := judge.state.(map[string]any)
	if !ok {
		t.Fatalf("state = %T", judge.state)
	}
	if state["message"] != "You eat like a trash seagull" {
		t.Errorf("state message = %v", state["message"])
	}
	if state["context"] == "" || state["reactions"] == "" {
		t.Errorf("state = %v", state)
	}
	if len(judge.questions) != 200 {
		t.Fatalf("questions = %d, want 200", len(judge.questions))
	}
	for _, id := range []string{"r00", "r99", "r100", "r199"} {
		if _, ok := judge.questions[id]; !ok {
			t.Errorf("question %s missing", id)
		}
	}
}
