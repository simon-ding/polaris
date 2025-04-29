//go:build !lib
package ui

import "embed"

//go:embed build/web/*
var Web embed.FS