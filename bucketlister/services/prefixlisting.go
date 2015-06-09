package services

import (
	"path/filepath"
	"time"
)

// File represents an object in a PrefixListing
type File struct {
	Name         string
	LastModified time.Time
	Size         int64
}

func (f *File) Base() string {
	return filepath.Base(f.Name) + "/"
}

// PrefixListing is a compact listing of an S3 prefix
type PrefixListing struct {
	Prefixes []string
	Files    []*File
}
