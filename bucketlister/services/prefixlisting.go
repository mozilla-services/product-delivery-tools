package services

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Prefix represents a prefix path
type Prefix string

func (p Prefix) Escaped() string {
	parts := strings.Split(string(p), "/")
	escapedParts := make([]string, len(parts))
	for i, p := range parts {
		escapedParts[i] = template.URLQueryEscaper(p)
	}
	return strings.Join(escapedParts, "/")
}

// File represents an object in a PrefixListing
type File struct {
	Name         string    `json:"name"`
	LastModified time.Time `json:"last_modified"`
	Size         int64     `json:"size"`
}

func (f *File) BaseEscaped() string {
	return template.URLQueryEscaper(f.Base())
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
	Prefixes []string `json:"prefixes"`
	Files    []*File  `json:"files"`
}

// PrefixStructs returns prefix objects
func (p *PrefixListing) PrefixStructs() []Prefix {
	tmp := make([]Prefix, len(p.Prefixes))
	for i, s := range p.Prefixes {
		tmp[i] = Prefix(s)
	}
	return tmp
}
