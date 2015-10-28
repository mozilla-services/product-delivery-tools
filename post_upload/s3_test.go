package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyCacheControl(t *testing.T) {
	cases := [][]string{
		{"pub/firefox/nightly/latest-trunk/firefox-44.0a1.en-US.win32.installer.exe", "max-age=3600"},
		{"pub/firefox/releases/41.0.2/win32/en-US/Firefox%20Setup%2041.0.2.exe", ""},
	}

	for _, c := range cases {
		res := keyCacheControl(c[0])
		if c[1] == "" {
			assert.Nil(t, res)
		} else {
			assert.Equal(t, c[1], *res)
		}
	}
}
