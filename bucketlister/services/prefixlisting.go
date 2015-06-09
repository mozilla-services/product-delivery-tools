package services

import (
	"fmt"
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
	return filepath.Base(f.Name)
}

func (f *File) LastModifiedString() string {
	return f.LastModified.Format("02-Jan-2006 15:04")
}

func (f *File) SizeString() string {
	if f.Size < 1024 {
		return fmt.Sprintf("%d", f.Size)
	}

	if f.Size < 1024*1024 {
		return fmt.Sprintf("%dK", f.Size/1024)
	}

	return fmt.Sprintf("%dM", f.Size/(1024*1024))
}

// PrefixListing is a compact listing of an S3 prefix
type PrefixListing struct {
	Prefixes []string
	Files    []*File
}
