// Package web holds the frontend assets embedded into the jevmoji binary.
package web

import "embed"

// FS contains the page templates and static assets (CSS, JS, font).
//
//go:embed templates static
var FS embed.FS
