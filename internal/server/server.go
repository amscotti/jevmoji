// Package server wires the jevmoji HTTP surface: the single page, the
// reactions endpoint, and the embedded static assets.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/amscotti/jevmoji/internal/catalog"
	"github.com/amscotti/jevmoji/internal/jev"
	"github.com/amscotti/jevmoji/internal/suggest"
	"github.com/amscotti/jevmoji/web"
)

// Message bounds: blank input stays idle (no point asking the jury about
// nothing), short greetings included — even a "hi" deserves a reaction.
// The upper bound keeps the state small enough to stay cheap.
const (
	minMessageRunes = 1
	maxMessageRunes = 800
)

// Judge evaluates the emoji questions against a message state. *jev.Client
// implements it; tests bring their own.
type Judge interface {
	SystemOne(ctx context.Context, state any, questions map[string]jev.Question) (*jev.Response, error)
}

// Config carries the process-level settings the server needs. Judge and Now
// are optional seams for tests.
type Config struct {
	APIKey string
	Judge  Judge
	Now    func() time.Time
}

// Server holds the HTTP handlers.
type Server struct {
	judge Judge
	tmpl  *template.Template
	now   func() time.Time
}

// pageView is what both the full page and the results partial render.
type pageView struct {
	Suggestion *suggest.Suggestion
	Error      string
}

// New builds a Server from cfg, parsing the embedded templates.
func New(cfg Config) (*Server, error) {
	funcs := template.FuncMap{
		"pct": func(p float64) string { return fmt.Sprintf("%.0f%%", p*100) },
	}
	tmpl, err := template.New("jevmoji").Funcs(funcs).ParseFS(web.FS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("server: parsing templates: %w", err)
	}
	judge := cfg.Judge
	if judge == nil {
		judge = jev.NewClient(cfg.APIKey)
	}
	now := cfg.Now
	if now == nil {
		now = time.Now
	}
	return &Server{judge: judge, tmpl: tmpl, now: now}, nil
}

// Routes returns the full route table.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", s.handlePage)
	mux.HandleFunc("POST /api/reactions", s.handleReactions)
	// /health (not /healthz): Cloud Run's Google Front End reserves paths
	// ending in "z" and intercepts them with its own 404 before the request
	// reaches the container.
	mux.HandleFunc("GET /health", s.handleHealth)

	static, err := fs.Sub(web.FS, "static")
	if err != nil {
		panic("server: embedded static fs: " + err.Error())
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))

	return securityHeaders(mux)
}

// securityHeaders sets the baseline headers for every response. The page
// pulls pinned scripts and styles from CDNs, so there is no CSP allowlist
// here — but nothing on this server should be sniffed, framed, or sent as a
// referrer to another origin.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handlePage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, "page", pageView{}); err != nil {
		log.Printf("jevmoji: rendering page: %v", err)
	}
}

// handleReactions scores every catalog emoji as a reply to the message and
// renders the results partial (or the JSON envelope with ?format=json).
func (s *Server) handleReactions(w http.ResponseWriter, r *http.Request) {
	text, err := extractText(r)
	if err != nil {
		s.render(w, r, http.StatusBadRequest, pageView{Error: err.Error()})
		return
	}
	if utf8.RuneCountInString(text) < minMessageRunes {
		s.render(w, r, http.StatusOK, pageView{})
		return
	}

	start := s.now()
	resp, err := s.judge.SystemOne(r.Context(), catalog.State(text), catalog.Questions())
	elapsed := s.now().Sub(start)
	if err != nil {
		log.Printf("jevmoji: systemone: %v", err)
		s.render(w, r, http.StatusBadGateway, pageView{Error: friendlyError(err)})
		return
	}

	sug, err := suggest.Build(text, catalog.Catalog, resp.Answers)
	if err != nil {
		log.Printf("jevmoji: building suggestion: %v", err)
		s.render(w, r, http.StatusBadGateway,
			pageView{Error: "the model answered, but not with what we expected"})
		return
	}
	sug.Model = resp.Model
	usage := resp.Usage
	sug.Usage = &usage
	sug.LatencyMS = elapsed.Milliseconds()
	s.render(w, r, http.StatusOK, pageView{Suggestion: &sug})
}

// extractText pulls the message out of a JSON body {"text": "..."} or the
// form field text=..., and normalizes it (trim, bounded length).
func extractText(r *http.Request) (string, error) {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		var body struct {
			Text string `json:"text"`
		}
		if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&body); err != nil {
			return "", errors.New(`expected a JSON body like {"text": "..."}`)
		}
		return normalize(body.Text), nil
	}
	if err := r.ParseForm(); err != nil {
		return "", errors.New("could not read the form")
	}
	return normalize(r.FormValue("text")), nil
}

func normalize(text string) string {
	text = strings.TrimSpace(text)
	if utf8.RuneCountInString(text) > maxMessageRunes {
		runes := []rune(text)
		text = string(runes[:maxMessageRunes])
	}
	return text
}

// friendlyError turns transport and API failures into something worth
// reading on the page.
func friendlyError(err error) string {
	var apiErr *jev.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusUnauthorized:
			return "the TypeSafe API key was rejected (401) — check TYPESAFE_API_KEY"
		case http.StatusUnprocessableEntity:
			return "the TypeSafe API called the request malformed (422)"
		case http.StatusTooManyRequests, 529:
			return "the TypeSafe API is rate limiting or overloaded — try again in a moment"
		}
		if apiErr.Message != "" {
			return "the TypeSafe API said: " + apiErr.Message
		}
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return "the emoji jury timed out"
	}
	return "could not reach the TypeSafe API"
}

// render writes the results partial as HTML, or the JSON envelope when the
// caller asked for JSON (?format=json or Accept: application/json).
func (s *Server) render(w http.ResponseWriter, r *http.Request, status int, view pageView) {
	if wantsJSON(r) {
		renderJSON(w, status, view)
		return
	}
	// HTMX 4 swaps 4xx/5xx responses by default, so HTML keeps the honest
	// status code while the partial carries the error state.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := s.tmpl.ExecuteTemplate(w, "results", view); err != nil {
		log.Printf("jevmoji: rendering results: %v", err)
	}
}

func renderJSON(w http.ResponseWriter, status int, view pageView) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if view.Error != "" {
		_ = json.NewEncoder(w).Encode(map[string]string{"error": view.Error})
		return
	}
	if view.Suggestion == nil {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"state":     "",
			"reactions": []suggest.Reaction{},
			"chosen":    []string{},
		})
		return
	}
	_ = json.NewEncoder(w).Encode(view.Suggestion)
}

func wantsJSON(r *http.Request) bool {
	if r.URL.Query().Get("format") == "json" {
		return true
	}
	return strings.Contains(r.Header.Get("Accept"), "application/json")
}
