// Command jevmoji serves the jevmoji single-page app: type a group-chat
// message and get Jev-scored emoji reactions back.
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/amscotti/jevmoji/internal/server"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()

	// TYPESAFE_API_KEY comes from the environment. As a convenience, a .env
	// file fills in whatever is still unset — first from the working
	// directory, then from the directory holding the binary (so running the
	// binary from elsewhere still finds its config). The environment always
	// wins over both files.
	if os.Getenv("TYPESAFE_API_KEY") == "" {
		for _, path := range envFilePaths() {
			if err := loadEnvFile(path); err != nil && !errors.Is(err, os.ErrNotExist) {
				log.Fatalf("jevmoji: loading %s: %v", path, err)
			}
			if os.Getenv("TYPESAFE_API_KEY") != "" {
				break
			}
		}
	}
	if os.Getenv("TYPESAFE_API_KEY") == "" {
		log.Printf("jevmoji: warning: TYPESAFE_API_KEY is unset — the page serves, but /api/reactions will answer 502 until it is set")
	}

	srv, err := server.New(server.Config{APIKey: os.Getenv("TYPESAFE_API_KEY")})
	if err != nil {
		log.Fatalf("jevmoji: %v", err)
	}

	httpSrv := &http.Server{
		Addr:    *addr,
		Handler: srv.Routes(),
		// No WriteTimeout: a reactions request waits on the upstream Jev
		// call (8s per try, up to 3 tries plus backoff ≈ 25s worst case),
		// and killing slow writes would turn model latency into 500s.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("jevmoji listening on %s", *addr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("jevmoji: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("jevmoji: shutdown: %v", err)
	}
}

// envFilePaths lists the .env locations tried in order: the working
// directory first (invocation-specific), then the binary's own directory
// (installation config). loadEnvFile never overwrites, so the first file
// that sets a key wins and the environment beats both.
func envFilePaths() []string {
	paths := []string{".env"}
	if exe, err := os.Executable(); err == nil {
		if dir := filepath.Dir(exe); dir != "" && dir != "." {
			paths = append(paths, filepath.Join(dir, ".env"))
		}
	}
	return paths
}

// loadEnvFile parses a minimal KEY=VALUE .env file into the process
// environment. Variables already set in the environment win. Blank lines,
// '#' comments, and lines without '=' are ignored; a leading `export ` is
// stripped, and values wrapped in one matching pair of single or double
// quotes are unwrapped. There is deliberately no inline-comment stripping,
// escape processing, or multiline support — keep secrets simple.
func loadEnvFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = unquote(strings.TrimSpace(value))
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
	return nil
}

// unquote removes one matching pair of surrounding single or double quotes.
// Anything else — including mismatched or interior quotes — is left alone.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
