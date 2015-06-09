package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// BucketLister is a directory listing service for S3
type BucketLister struct {
	Bucket     string
	mountedAt  string
	basePrefix string

	listers map[string][]*BucketLister

	AWSConfig *aws.Config
}

// NewBucketLister returns a *BucketLister
//
// prefix is the starting point for this lister
func NewBucketLister(bucket, prefix string, awsConfig *aws.Config) *BucketLister {
	trimmedPrefix := strings.Trim(prefix, "/")
	if trimmedPrefix != "" {
		trimmedPrefix += "/"
	}
	return &BucketLister{
		AWSConfig:  awsConfig,
		Bucket:     bucket,
		mountedAt:  "/" + trimmedPrefix,
		basePrefix: trimmedPrefix,
		listers:    make(map[string][]*BucketLister),
	}
}

// AddBucketLister adds a lister to the root lister
//
// If a lister is attached, it will show up as a directory link
func (b *BucketLister) AddBucketLister(mount string, child *BucketLister) {
	if b.listers[mount] == nil {
		b.listers[mount] = []*BucketLister{child}
		return
	}
	b.listers[mount] = append(b.listers[mount], child)
}

// Empty returns true if the bucket contains zero keys
func (b *BucketLister) Empty() (bool, error) {
	listParams := &s3.ListObjectsInput{
		Bucket:  aws.String(b.Bucket),
		MaxKeys: aws.Long(1),
	}

	s3Service := s3.New(b.AWSConfig)
	res, err := s3Service.ListObjects(listParams)
	if err != nil {
		return true, fmt.Errorf("listing %s err: %s", b.Bucket, err)
	}

	return len(res.Contents) <= 0, nil
}

func objectToListFileInfo(obj *s3.Object) *File {
	return &File{
		Name:         *obj.Key,
		LastModified: *obj.LastModified,
		Size:         *obj.Size,
	}
}

func (b *BucketLister) listerDirs(reqPath string) []string {
	dirs := []string{}
	for _, lister := range b.listers[reqPath] {
		if empty, err := lister.Empty(); empty {
			if err != nil {
				log.Println("Error checking empty: %s", err)
			}
			continue
		}
		dirs = append(dirs, path.Base(lister.mountedAt)+"/")
	}

	return dirs
}

func (b *BucketLister) listPrefix(reqPath, prefix string) (*PrefixListing, error) {
	s3Service := s3.New(b.AWSConfig)
	objects, prefixes, err := listObjects(s3Service, b.Bucket, prefix)
	if err != nil {
		return nil, err
	}

	extraDirs := b.listerDirs(reqPath)

	listing := &PrefixListing{
		Prefixes: make([]string, len(prefixes), len(prefixes)+len(extraDirs)),
		Files:    make([]*File, 0, len(objects)),
	}

	for i, p := range prefixes {
		listing.Prefixes[i] = path.Base(*p.Prefix) + "/"
	}

	listing.Prefixes = append(listing.Prefixes, extraDirs...)

	sort.Strings(listing.Prefixes)

	for _, o := range objects {
		if *o.Key == prefix {
			continue
		}
		o := objectToListFileInfo(o)
		o.Name = strings.TrimPrefix(o.Name, prefix)
		listing.Files = append(listing.Files, o)
	}

	return listing, nil
}

// ServeHTTP implements http.Handler
func (b *BucketLister) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	if !strings.HasSuffix(reqPath, "/") {
		reqPath += "/"
	}
	relPath := strings.TrimPrefix(reqPath, b.mountedAt)
	prefix := path.Join(b.basePrefix, relPath)
	if prefix != "" {
		prefix += "/"
	}

	listing, err := b.listPrefix(reqPath, prefix)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error."))
		log.Printf("Error %s", err)
		return
	}

	if len(listing.Files) == 0 && len(listing.Prefixes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	contentType := "text/html"
	body := new(bytes.Buffer)
	switch req.Header.Get("Accept") {
	case "application/json":
		contentType = "application/json"
		err := json.NewEncoder(body).Encode(listing)
		if err != nil {
			log.Printf("Error encoding JSON err: %s", err)
		}
	default:
		tmplParams := &listTemplateInput{
			Path:          reqPath,
			PrefixListing: listing,
		}
		err = listTemplate.Execute(body, tmplParams)
		if err != nil {
			log.Printf("Error executing template err: %s", err)
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error."))
		return
	}

	setExpiresIn(15*time.Minute, w)
	w.Header().Set("Vary", "Accept")
	w.Header().Set("Content-Type", contentType)
	w.Write(body.Bytes())
}
