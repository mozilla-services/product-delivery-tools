package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentType(t *testing.T) {
	files := [][]string{
		[]string{"foo/bar/firefox-44.0a1.en-US.win32.mar", "application/octet-stream"},
		[]string{"foo/bar/firefox-44.0a1.en-US.win32.dmg", "application/x-apple-diskimage"},
		[]string{"foo/bar/firefox-44.0a1.en-US.win32.png", "image/png"},
		[]string{"foo/bar/firefox-44.0a1.en-US.win32.txt", "text/plain; charset=utf-8"},
		[]string{"foo/bar/firefox-44.0a1.en-US.win32.unknown", "application/octet-stream"},
	}

	for _, f := range files {
		assert.Equal(t, f[1], ContentType(f[0]))
	}
}
