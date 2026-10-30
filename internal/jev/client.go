package jev

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Endpoint is the TypeSafe System One evaluation endpoint.
const Endpoint = "https://api.typesafe.ai/v1/systemone"

// Model is the model alias sent with every request. jev-latest tracks the
// newest stable Jev release (Jev 1.14 / jev-1.14.0 as of October 2026); the
// response's model field reports the exact version that answered.
const Model = "jev-latest"

// Client talks to the TypeSafe System One API.
type Client struct {
	apiKey   string
	baseURL  string
	http     *http.Client
	backoff  time.Duration
	attempts int
}

// Option customizes a Client.
type Option func(*Client)

// WithHTTPClient replaces the internal HTTP client (8s timeout by default).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.http = hc }
}

// WithBaseURL points the client at another host (tests).
func WithBaseURL(u string) Option {
	return func(c *Client) { c.baseURL = u }
}

// WithBackoff sets the base retry delay for 429/529 (default 250ms). A
// server-sent Retry-After hint can extend an individual delay (capped); see
// maxRetryAfter.
func WithBackoff(d time.Duration) Option {
	return func(c *Client) { c.backoff = d }
}

// WithAttempts sets how many total tries SystemOne makes (default 3).
// Values below 1 are clamped to 1.
func WithAttempts(n int) Option {
	return func(c *Client) {
		if n < 1 {
			n = 1
		}
		c.attempts = n
	}
}

// NewClient builds a client for the given API key.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:   apiKey,
		baseURL:  Endpoint,
		http:     &http.Client{Timeout: 8 * time.Second},
		backoff:  250 * time.Millisecond,
		attempts: 3,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type request struct {
	State     json.RawMessage     `json:"state"`
	Model     string              `json:"model"`
	Questions map[string]Question `json:"questions"`
}

// SystemOne evaluates the questions against state in a single call. All
// questions are answered in parallel server-side; the response carries one
// typed answer per question key.
func (c *Client) SystemOne(ctx context.Context, state any, questions map[string]Question) (*Response, error) {
	if len(questions) == 0 {
		return nil, fmt.Errorf("jev: at least one question is required")
	}
	for key, q := range questions {
		if err := q.Validate(); err != nil {
			return nil, fmt.Errorf("jev: question %q: %w", key, err)
		}
	}

	// state is documented as string | object | array — nothing else.
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("jev: encoding state: %w", err)
	}
	if len(stateJSON) == 0 || !strings.ContainsRune(`"{[`, rune(stateJSON[0])) {
		return nil, fmt.Errorf("jev: state must be a string, object, or array")
	}

	payload, err := json.Marshal(request{State: stateJSON, Model: Model, Questions: questions})
	if err != nil {
		return nil, fmt.Errorf("jev: encoding request: %w", err)
	}

	var lastErr error
	var retryAfter time.Duration
	attempts := max(c.attempts, 1)
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			delay := c.backoff * time.Duration(1<<(attempt-1))
			if retryAfter > delay {
				delay = retryAfter
			}
			if err := sleepCtx(ctx, delay); err != nil {
				return nil, err
			}
		}
		var retryable bool
		var resp *Response
		resp, retryAfter, retryable, lastErr = c.do(ctx, payload)
		if lastErr == nil {
			return resp, nil
		}
		if !retryable {
			return nil, lastErr
		}
	}
	return nil, fmt.Errorf("jev: giving up after %d attempts: %w", attempts, lastErr)
}

// do performs one attempt. It reports the server's Retry-After hint (if any)
// alongside whether a failure is worth retrying (429 rate limit and 529
// overloaded per the API reference).
func (c *Client) do(ctx context.Context, payload []byte) (resp *Response, retryAfter time.Duration, retryable bool, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(payload))
	if err != nil {
		return nil, 0, false, fmt.Errorf("jev: building request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	httpResp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, true, fmt.Errorf("jev: request failed: %w", err)
	}
	defer func() { _ = httpResp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(httpResp.Body, 4<<20))
	if err != nil {
		return nil, 0, true, fmt.Errorf("jev: reading response: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		retryable := httpResp.StatusCode == http.StatusTooManyRequests || httpResp.StatusCode == 529
		var hint time.Duration
		if retryable {
			hint = cappedRetryAfter(httpResp.Header)
		}
		return nil, hint, retryable,
			&APIError{StatusCode: httpResp.StatusCode, Message: extractMessage(body), Raw: string(body)}
	}

	var decoded Response
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, 0, false, fmt.Errorf("jev: decoding response: %w", err)
	}
	if err := decoded.Validate(); err != nil {
		return nil, 0, false, err
	}
	return &decoded, 0, false, nil
}

// extractMessage pulls a human-readable message out of the common error
// envelopes; the API reference does not fix an error body schema.
func extractMessage(body []byte) string {
	var envelope struct {
		Message string `json:"message"`
		Detail  string `json:"detail"`
		Error   any    `json:"error"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return strings.TrimSpace(string(body))
	}
	switch e := envelope.Error.(type) {
	case string:
		return e
	case map[string]any:
		if msg, ok := e["message"].(string); ok {
			return msg
		}
	}
	if envelope.Message != "" {
		return envelope.Message
	}
	return envelope.Detail
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return fmt.Errorf("jev: %w", ctx.Err())
	case <-timer.C:
		return nil
	}
}

// maxRetryAfter caps a server-sent Retry-After hint so one hostile or buggy
// header cannot stall a request indefinitely. Context cancellation still
// applies while waiting.
const maxRetryAfter = 30 * time.Second

// cappedRetryAfter parses the Retry-After hint, bounded by maxRetryAfter.
func cappedRetryAfter(header http.Header) time.Duration {
	if hint := RetryAfterHint(header); hint > 0 {
		return min(hint, maxRetryAfter)
	}
	return 0
}

// RetryAfterHint returns the seconds a 429/529 response asked for, or 0.
// Kept small and separate so tests can assert the parse without a round trip.
func RetryAfterHint(header http.Header) time.Duration {
	v := header.Get("Retry-After")
	if v == "" {
		return 0
	}
	secs, err := strconv.Atoi(v)
	if err != nil || secs <= 0 {
		return 0
	}
	return time.Duration(secs) * time.Second
}
