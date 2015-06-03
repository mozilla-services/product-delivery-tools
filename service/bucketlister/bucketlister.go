package bucketlister

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"sort"
	"strings"

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

// New returns a *BucketLister
//
// prefix is the starting point for this lister
func New(bucket, prefix string, awsConfig *aws.Config) *BucketLister {
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

func objectToListFileInfo(obj *s3.Object) *listFileInfo {
	size := *obj.Size
	sizeStr := ""
	if size < 1024 {
		sizeStr = fmt.Sprintf("%d B", size)
	} else {
		sizeStr = fmt.Sprintf("%d KB", size/1024)
	}

	return &listFileInfo{
		Key:          *obj.Key,
		LastModified: (*obj.LastModified).Format("02-Jan-2006 15:04"),
		Size:         sizeStr,
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

	prefixes := []*s3.CommonPrefix{}
	objects := []*s3.Object{}

	s3Service := s3.New(b.AWSConfig)

	objects, prefixes, err := listObjects(s3Service, b.Bucket, prefix)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error."))
		log.Printf("Error %s", err)
		return
	}

	if len(b.listers[reqPath]) == 0 && len(objects) == 0 && len(prefixes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	tmplParams := &listTemplateInput{
		Path:        reqPath,
		Directories: make([]string, len(prefixes)),
		Files:       make([]*listFileInfo, 0, len(objects)),
	}

	for i, p := range prefixes {
		tmplParams.Directories[i] = path.Base(*p.Prefix) + "/"
	}

	extraDirs := b.listerDirs(reqPath)
	if len(extraDirs) > 0 {
		tmplParams.Directories = append(tmplParams.Directories, extraDirs...)
		sort.Strings(tmplParams.Directories)
	}

	for _, o := range objects {
		if *o.Key == prefix {
			continue
		}
		o := objectToListFileInfo(o)
		o.Key = strings.TrimPrefix(o.Key, prefix)
		tmplParams.Files = append(tmplParams.Files, o)
	}

	w.Header().Set("Content-Type", "text/html")
	err = listTemplate.Execute(w, tmplParams)
	if err != nil {
		log.Printf("Error executing template err: %s", err)
	}
}
