package web

import "embed"

// Assets contains the public and admin Vite build outputs.
//
//go:embed all:dist
var Assets embed.FS
