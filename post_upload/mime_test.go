package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentType(t *testing.T) {
	files := [][]string{
		[]string{"/foo/bar/firefox.mar", "application/octet-stream"},
		[]string{"/foo/bar/firefox.dmg", "application/x-apple-diskimage"},
		[]string{"/foo/bar/firefox.png", "image/png"},
		[]string{"/foo/bar/firefox.txt", "text/plain; charset=utf-8"},
		[]string{"/foo/bar/firefox.unknown", "application/octet-stream"},
	}

	for _, f := range files {
		assert.Equal(t, f[1], ContentType(f[0]))
	}
}
