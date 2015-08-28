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

// SortMountedAt sorts a slice of bucketlisters by mountedAt
type SortMountedAt []*BucketLister

func (s SortMountedAt) Len() int { return len(s) }

func (s SortMountedAt) Less(i, j int) bool { return s[i].mountedAt < s[j].mountedAt }

func (s SortMountedAt) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// BucketLister is a directory listing service for S3
type BucketLister struct {
	Bucket     string
	mountedAt  string
	basePrefix string

	listers []*BucketLister

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
		listers:    []*BucketLister{},
	}
}

// AddBucketLister adds a lister to the root lister
//
// If a lister is attached, it will show up as a directory link
func (b *BucketLister) AddBucketLister(child *BucketLister) {
	if b.listers == nil {
		b.listers = []*BucketLister{child}
		return
	}
	b.listers = append(b.listers, child)
}

// Empty returns true if the bucket contains zero keys
func (b *BucketLister) Empty() (bool, error) {
	listParams := &s3.ListObjectsInput{
		Bucket:  aws.String(b.Bucket),
		MaxKeys: aws.Int64(1),
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

// Mount returns the mount point of this lister
func (b *BucketLister) Mount() string {
	return b.mountedAt
}

func (b *BucketLister) listerDirs(reqPath string) []string {
	dirs := make(map[string]bool)
	for _, lister := range b.listers {
		if strings.HasPrefix(lister.mountedAt, reqPath) {
			tmp := strings.TrimPrefix(lister.mountedAt, reqPath)
			idx := strings.Index(tmp, "/")
			if idx > 0 {
				tmp = tmp[0:idx]
			}
			dirs[path.Base(tmp)+"/"] = true
		}
	}

	res := make([]string, 0, len(dirs))
	for k := range dirs {
		res = append(res, k)
	}
	return res
}

func deduplicateSlice(slice []string) []string {
	keys := make(map[string]bool)
	for _, k := range slice {
		keys[k] = true
	}
	result := make([]string, 0, len(keys))
	for k := range keys {
		result = append(result, k)
	}
	return result
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

	listing.Prefixes = deduplicateSlice(append(listing.Prefixes, extraDirs...))

	sort.Strings(listing.Prefixes)

	for _, o := range objects {
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

	if reqPath != b.mountedAt && len(listing.Files) == 0 && len(listing.Prefixes) == 0 {
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
